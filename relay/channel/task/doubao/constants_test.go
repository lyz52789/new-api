package doubao

import "testing"

func TestForcedSeedanceResolutionForModel(t *testing.T) {
	tests := []struct {
		name   string
		want   string
		forced bool
	}{
		{name: "Seedance-2.0-海外版"},
		{name: "Seedance-2.0-720P-海外版", want: "720p", forced: true},
		{name: "Seedance-2.0-1080P-海外版", want: "1080p", forced: true},
		{name: "Seedance-2.0-4K-海外版", want: "4k", forced: true},
		{name: "doubao-seedance-2-0-260128"},
		{name: "Seedance-2.0-fast-海外版"},
	}
	for _, tt := range tests {
		got, forced := ForcedSeedanceResolutionForModel(tt.name)
		if got != tt.want || forced != tt.forced {
			t.Fatalf("ForcedSeedanceResolutionForModel(%q) = (%q, %v), want (%q, %v)", tt.name, got, forced, tt.want, tt.forced)
		}
	}
}

func TestGetVideoInputRatioUsesResolutionUnitPrice(t *testing.T) {
	tests := []struct {
		name       string
		resolution string
		hasVideo   bool
		want       float64
	}{
		{name: "base", resolution: "720P", want: 1},
		{name: "1080p", resolution: "1080P", want: 7.7 / 7.0},
		{name: "4k", resolution: "4K", want: 4.0 / 7.0},
		{name: "4k video input", resolution: "4K", hasVideo: true, want: 2.4 / 7.0},
	}
	for _, tt := range tests {
		got, ok := GetVideoInputRatio("Seedance-2.0-海外版", tt.resolution, tt.hasVideo)
		if !ok {
			t.Fatalf("%s: expected configured price table", tt.name)
		}
		if got != tt.want {
			t.Fatalf("%s: got %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestGetVideoInputRatioSupportsResolutionAliases(t *testing.T) {
	got, ok := GetVideoInputRatio("Seedance-2.0-1080P-海外版", "1080P", false)
	if !ok {
		t.Fatal("expected configured price table for 1080P alias")
	}
	if got != 1.0 {
		t.Fatalf("got %v, want %v", got, 1.0)
	}

	got, ok = GetVideoInputRatio("Seedance-2.0-1080P-海外版", "1080P", true)
	if !ok {
		t.Fatal("expected configured price table for 1080P alias")
	}
	if got != 4.7/7.7 {
		t.Fatalf("got %v, want %v", got, 4.7/7.7)
	}
}

func TestGetVideoInputRatioSupportsFastModel(t *testing.T) {
	got, ok := GetVideoInputRatio("Seedance-2.0-fast-海外版", "720P", false)
	if !ok {
		t.Fatal("expected configured price table for Seedance 2.0 Fast")
	}
	if got != 1.0 {
		t.Fatalf("got %v, want 1", got)
	}

	got, ok = GetVideoInputRatio("Seedance-2.0-fast-海外版", "720P", true)
	if !ok {
		t.Fatal("expected configured price table for Seedance 2.0 Fast")
	}
	if got != 3.3/5.6 {
		t.Fatalf("got %v, want %v", got, 3.3/5.6)
	}
}

func TestGetSeedance15AudioRatio(t *testing.T) {
	got, ok := GetSeedance15AudioRatio("Seedance-1.5-pro-海外版", true)
	if !ok {
		t.Fatal("expected Seedance 1.5 Pro price table")
	}
	if got != 1.0 {
		t.Fatalf("audio ratio = %v, want 1", got)
	}

	got, ok = GetSeedance15AudioRatio("Seedance-1.5-pro-海外版", false)
	if !ok {
		t.Fatal("expected Seedance 1.5 Pro price table")
	}
	if got != 1.2/2.4 {
		t.Fatalf("silent ratio = %v, want %v", got, 1.2/2.4)
	}
}
