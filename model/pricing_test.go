package model

import "testing"

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
