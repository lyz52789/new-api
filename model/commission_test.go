package model

import (
	"testing"
	"time"

	"github.com/QuantumNous/new-api/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func insertAffiliateUser(t *testing.T, id int, username string, inviterId int) {
	t.Helper()
	user := &User{
		Id:        id,
		Username:  username,
		Status:    common.UserStatusEnabled,
		AffCode:   username + "_aff",
		InviterId: inviterId,
	}
	require.NoError(t, DB.Create(user).Error)
}

func insertAffiliateTopUp(t *testing.T, tradeNo string, userID int, paymentMethod string, money float64, status string) *TopUp {
	t.Helper()
	topUp := &TopUp{
		UserId:        userID,
		Amount:        10,
		Money:         money,
		TradeNo:       tradeNo,
		PaymentMethod: paymentMethod,
		Status:        status,
		CreateTime:    time.Now().Unix(),
	}
	require.NoError(t, topUp.Insert())
	return topUp
}

func countCommissionRecords(t *testing.T) int64 {
	t.Helper()
	var count int64
	require.NoError(t, DB.Model(&CommissionRecord{}).Count(&count).Error)
	return count
}

func TestCreateAffiliateCommissionForTopUp_CreatesPendingCommission(t *testing.T) {
	truncateTables(t)

	insertAffiliateUser(t, 501, "aff_inviter", 0)
	insertAffiliateUser(t, 502, "aff_invitee", 501)
	topUp := insertAffiliateTopUp(t, "affiliate-commission", 502, PaymentMethodWaffo, 123.45, common.TopUpStatusSuccess)

	record, err := CreateAffiliateCommissionForTopUp(topUp)
	require.NoError(t, err)
	require.NotNil(t, record)

	assert.Equal(t, 501, record.InviterId)
	assert.Equal(t, 502, record.InviteeId)
	assert.Equal(t, int64(12345), record.TopupAmount)
	assert.Equal(t, int64(1235), record.CommissionAmount)
	assert.Equal(t, CommissionStatusPending, record.Status)
	assert.GreaterOrEqual(t, record.SettleAt, common.GetTimestamp()+AffiliateCommissionSettleDelaySeconds-1)

	again, err := CreateAffiliateCommissionForTopUp(topUp)
	require.NoError(t, err)
	require.NotNil(t, again)
	assert.Equal(t, int64(1), countCommissionRecords(t))
}

func TestRechargeWaffo_CreatesAffiliateCommission(t *testing.T) {
	truncateTables(t)

	insertAffiliateUser(t, 601, "waffo_inviter", 0)
	insertAffiliateUser(t, 602, "waffo_invitee", 601)
	insertAffiliateTopUp(t, "affiliate-waffo", 602, PaymentMethodWaffo, 50, common.TopUpStatusPending)

	err := RechargeWaffo("affiliate-waffo", "127.0.0.1")
	require.NoError(t, err)

	var record CommissionRecord
	require.NoError(t, DB.Where("topup_trade_no = ?", "affiliate-waffo").First(&record).Error)
	assert.Equal(t, int64(5000), record.TopupAmount)
	assert.Equal(t, int64(500), record.CommissionAmount)
	assert.Equal(t, CommissionStatusPending, record.Status)

	err = RechargeWaffo("affiliate-waffo", "127.0.0.1")
	require.NoError(t, err)
	assert.Equal(t, int64(1), countCommissionRecords(t))
}

func TestManualCompleteTopUp_DoesNotCreateAffiliateCommission(t *testing.T) {
	truncateTables(t)

	insertAffiliateUser(t, 701, "manual_inviter", 0)
	insertAffiliateUser(t, 702, "manual_invitee", 701)
	insertAffiliateTopUp(t, "affiliate-manual", 702, PaymentMethodWaffo, 80, common.TopUpStatusPending)

	err := ManualCompleteTopUp("affiliate-manual", "127.0.0.1")
	require.NoError(t, err)
	assert.Equal(t, int64(0), countCommissionRecords(t))
}

func TestAffiliateDashboardCountsInviteRelations(t *testing.T) {
	truncateTables(t)

	insertAffiliateUser(t, 801, "dashboard_inviter", 0)
	insertAffiliateUser(t, 802, "dashboard_invitee_one", 801)
	insertAffiliateUser(t, 803, "dashboard_invitee_two", 801)

	dashboard, err := GetAffiliateDashboard(801)
	require.NoError(t, err)
	require.NotNil(t, dashboard)
	assert.Equal(t, 2, dashboard.InviteCount)
}
