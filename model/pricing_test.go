package model

import (
	"strings"
	"testing"

	"github.com/QuantumNous/new-api/common"
)

func TestAttachVideoPricingForSeedance20Fast(t *testing.T) {
	pricing := &Pricing{ModelName: "Seedance-2.0-fast-海外版", ModelRatio: 22.484}
	attachVideoPricing(pricing)

	if pricing.VideoPricing == nil {
		t.Fatal("expected video pricing")
	}
	if len(pricing.VideoPricing.Rows) != 4 {
		t.Fatalf("rows = %d, want 4", len(pricing.VideoPricing.Rows))
	}
	for _, row := range pricing.VideoPricing.Rows {
		if row.Resolution != "480p" && row.Resolution != "720p" {
			t.Fatalf("unexpected fast resolution %q", row.Resolution)
		}
	}
	if pricing.VideoPricing.Rows[0].SaleRMBPerMTokens != 44.968 {
		t.Fatalf("sale RMB/M = %v, want 44.968", pricing.VideoPricing.Rows[0].SaleRMBPerMTokens)
	}
	if pricing.VideoPricing.Rows[2].SaleRMBPerMTokens != 26.499 {
		t.Fatalf("video input sale RMB/M = %v, want 26.499", pricing.VideoPricing.Rows[2].SaleRMBPerMTokens)
	}
}

func TestAttachVideoPricingUsesConfiguredModelRatio(t *testing.T) {
	oldPrice := &Pricing{ModelName: "Seedance-2.0-fast-海外版", ModelRatio: 44.968}
	attachVideoPricing(oldPrice)
	if oldPrice.VideoPricing.Rows[0].SaleRMBPerMTokens != 89.936 {
		t.Fatalf("old sale RMB/M = %v, want 89.936", oldPrice.VideoPricing.Rows[0].SaleRMBPerMTokens)
	}

	newPrice := &Pricing{ModelName: "Seedance-2.0-fast-海外版", ModelRatio: 22.484}
	attachVideoPricing(newPrice)
	if newPrice.VideoPricing.Rows[0].SaleRMBPerMTokens != 44.968 {
		t.Fatalf("new sale RMB/M = %v, want 44.968", newPrice.VideoPricing.Rows[0].SaleRMBPerMTokens)
	}
}

func TestAttachVideoPricingAppliesOfficialScenarioRatios(t *testing.T) {
	pricing := &Pricing{ModelName: "Seedance-2.0-海外版", ModelRatio: 28.105}
	attachVideoPricing(pricing)

	if pricing.VideoPricing == nil {
		t.Fatal("expected video pricing")
	}
	if pricing.VideoPricing.Rows[2].Resolution != "1080p" {
		t.Fatalf("row[2] resolution = %q, want 1080p", pricing.VideoPricing.Rows[2].Resolution)
	}
	if pricing.VideoPricing.Rows[2].SaleRMBPerMTokens != 61.831 {
		t.Fatalf("1080p sale RMB/M = %v, want 61.831", pricing.VideoPricing.Rows[2].SaleRMBPerMTokens)
	}
}

func TestAttachVideoPricingShowsSeedance15PerSecondReference(t *testing.T) {
	pricing := &Pricing{ModelName: "Seedance-1.5-pro-海外版", ModelRatio: 19.272}
	attachVideoPricing(pricing)

	if pricing.VideoPricing == nil {
		t.Fatal("expected video pricing")
	}
	row := pricing.VideoPricing.Rows[0]
	if row.Resolution != "480p" || row.Scenario != "audio" {
		t.Fatalf("row[0] = %s/%s, want 480p/audio", row.Resolution, row.Scenario)
	}
	if row.SaleRMBPerVideo != 1.9272 {
		t.Fatalf("sale RMB/video = %v, want 1.9272", row.SaleRMBPerVideo)
	}
	if row.SaleRMBPerSecond != 0.3854 {
		t.Fatalf("sale RMB/second = %v, want 0.3854", row.SaleRMBPerSecond)
	}
}

func TestVideoPricingJSONDoesNotExposeCostBasis(t *testing.T) {
	pricing := &Pricing{ModelName: "Seedance-2.0-fast-海外版", ModelRatio: 22.484}
	attachVideoPricing(pricing)

	payload, err := common.Marshal(pricing)
	if err != nil {
		t.Fatalf("Marshal returned error: %v", err)
	}
	body := string(payload)
	for _, forbidden := range []string{
		"official_usd",
		"usd_cny_rate",
		"markup",
	} {
		if strings.Contains(body, forbidden) {
			t.Fatalf("public pricing JSON leaked %q: %s", forbidden, body)
		}
	}
	if !strings.Contains(body, "sale_rmb_per_m_tokens") {
		t.Fatalf("public pricing JSON should still include sale prices: %s", body)
	}
}
