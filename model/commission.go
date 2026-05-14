package model

import (
	"errors"

	"github.com/QuantumNous/new-api/common"
	appmetrics "github.com/QuantumNous/new-api/metrics"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

const (
	CommissionStatusPending = iota
	CommissionStatusAvailable
	CommissionStatusWithdrawn
	CommissionStatusCanceled
)

const (
	AffiliateCommissionRate               = 0.10
	AffiliateCommissionSettleDelaySeconds = 21 * 24 * 60 * 60
)

type CommissionRecord struct {
	Id               int64  `json:"id" gorm:"primaryKey;autoIncrement"`
	InviterId        int    `json:"inviter_id" gorm:"index;not null"`
	InviteeId        int    `json:"invitee_id" gorm:"index;not null"`
	TopupTradeNo     string `json:"topup_trade_no" gorm:"type:varchar(255);uniqueIndex;not null"`
	PaymentMethod    string `json:"payment_method" gorm:"type:varchar(50)"`
	TopupAmount      int64  `json:"topup_amount" gorm:"not null"`
	CommissionAmount int64  `json:"commission_amount" gorm:"not null"`
	Status           int    `json:"status" gorm:"index;not null;default:0"`
	SettleAt         int64  `json:"settle_at" gorm:"index;not null"`
	CreatedAt        int64  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt        int64  `json:"updated_at" gorm:"autoUpdateTime"`
}

func (CommissionRecord) TableName() string {
	return "commission_records"
}

type AffiliateDashboard struct {
	AffCode             string `json:"aff_code"`
	InviteCount         int    `json:"invite_count"`
	CommissionTotal     int64  `json:"commission_total"`
	CommissionPending   int64  `json:"commission_pending"`
	CommissionAvailable int64  `json:"commission_available"`
}

func moneyYuanToCents(money float64) int64 {
	return decimal.NewFromFloat(money).Mul(decimal.NewFromInt(100)).Round(0).IntPart()
}

func CalculateAffiliateCommissionCents(topupAmountCents int64) int64 {
	if topupAmountCents <= 0 {
		return 0
	}
	return decimal.NewFromInt(topupAmountCents).Mul(decimal.NewFromFloat(AffiliateCommissionRate)).Round(0).IntPart()
}

func CreateAffiliateCommissionForTopUp(topUp *TopUp) (*CommissionRecord, error) {
	if topUp == nil || topUp.UserId <= 0 || topUp.TradeNo == "" || topUp.Status != common.TopUpStatusSuccess {
		return nil, nil
	}

	var existing CommissionRecord
	if err := DB.Where("topup_trade_no = ?", topUp.TradeNo).First(&existing).Error; err == nil {
		return &existing, nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	var invitee User
	if err := DB.Where("id = ?", topUp.UserId).First(&invitee).Error; err != nil {
		return nil, err
	}
	if invitee.InviterId <= 0 || invitee.InviterId == invitee.Id {
		return nil, nil
	}

	topupAmountCents := moneyYuanToCents(topUp.Money)
	commissionAmount := CalculateAffiliateCommissionCents(topupAmountCents)
	if commissionAmount <= 0 {
		return nil, nil
	}

	record := &CommissionRecord{
		InviterId:        invitee.InviterId,
		InviteeId:        invitee.Id,
		TopupTradeNo:     topUp.TradeNo,
		PaymentMethod:    topUp.PaymentMethod,
		TopupAmount:      topupAmountCents,
		CommissionAmount: commissionAmount,
		Status:           CommissionStatusPending,
		SettleAt:         common.GetTimestamp() + AffiliateCommissionSettleDelaySeconds,
	}
	if err := DB.Create(record).Error; err != nil {
		return nil, err
	}
	appmetrics.RecordAffiliateCommissionCreated(record.PaymentMethod, "pending", record.CommissionAmount)
	return record, nil
}

func GetAffiliateDashboard(userId int) (*AffiliateDashboard, error) {
	user, err := GetUserById(userId, true)
	if err != nil {
		return nil, err
	}

	dashboard := &AffiliateDashboard{
		AffCode:     user.AffCode,
		InviteCount: user.AffCount,
	}
	if dashboard.AffCode == "" {
		dashboard.AffCode = common.GetRandomString(4)
		user.AffCode = dashboard.AffCode
		if err := user.Update(false); err != nil {
			return nil, err
		}
	}

	var inviteCount int64
	if err := DB.Model(&User{}).Where("inviter_id = ?", userId).Count(&inviteCount).Error; err != nil {
		return nil, err
	}
	dashboard.InviteCount = int(inviteCount)

	if err := DB.Model(&CommissionRecord{}).
		Where("inviter_id = ? AND status <> ?", userId, CommissionStatusCanceled).
		Select("COALESCE(SUM(commission_amount), 0)").
		Scan(&dashboard.CommissionTotal).Error; err != nil {
		return nil, err
	}
	if err := DB.Model(&CommissionRecord{}).
		Where("inviter_id = ? AND status = ?", userId, CommissionStatusPending).
		Select("COALESCE(SUM(commission_amount), 0)").
		Scan(&dashboard.CommissionPending).Error; err != nil {
		return nil, err
	}
	if err := DB.Model(&CommissionRecord{}).
		Where("inviter_id = ? AND status = ?", userId, CommissionStatusAvailable).
		Select("COALESCE(SUM(commission_amount), 0)").
		Scan(&dashboard.CommissionAvailable).Error; err != nil {
		return nil, err
	}

	return dashboard, nil
}

func GetAffiliateCommissionRecords(inviterId int, pageInfo *common.PageInfo) (records []*CommissionRecord, total int64, err error) {
	query := DB.Model(&CommissionRecord{}).Where("inviter_id = ?", inviterId)
	if err = query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err = query.Order("id desc").Limit(pageInfo.GetPageSize()).Offset(pageInfo.GetStartIdx()).Find(&records).Error
	return records, total, err
}
