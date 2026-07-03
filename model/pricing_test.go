package model

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestAttachVideoPricingForSeedance20Fast(t *testing.T) {
	pricing := &Pricing{ModelName: "Seedance-2.0-fast-海外版"}
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
	if pricing.VideoPricing.Rows[0].SaleRMBPerMTokens != 89.936 {
		t.Fatalf("sale RMB/M = %v, want 89.936", pricing.VideoPricing.Rows[0].SaleRMBPerMTokens)
	}
}

func TestVideoPricingJSONDoesNotExposeCostBasis(t *testing.T) {
	pricing := &Pricing{ModelName: "Seedance-2.0-fast-海外版"}
	attachVideoPricing(pricing)

	payload, err := json.Marshal(pricing)
	if err != nil {
		t.Fatalf("Marshal returned error: %v", err)
	}
	body := string(payload)
	for _, forbidden := range []string{
		"official_usd",
		"usd_cny_rate",
		"markup",
		"7.3",
		"2.2",
	} {
		if strings.Contains(body, forbidden) {
			t.Fatalf("public pricing JSON leaked %q: %s", forbidden, body)
		}
	}
	if !strings.Contains(body, "sale_rmb_per_m_tokens") {
		t.Fatalf("public pricing JSON should still include sale prices: %s", body)
	}
}
