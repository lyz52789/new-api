package doubao

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/QuantumNous/new-api/model"
	relaycommon "github.com/QuantumNous/new-api/relay/common"
	"github.com/gin-gonic/gin"
)

func TestForceSeedanceAliasResolution(t *testing.T) {
	req := relaycommon.TaskSubmitReq{Model: "Seedance-2.0-1080P-海外版"}
	if err := forceSeedanceAliasResolution(&req); err != nil {
		t.Fatalf("forceSeedanceAliasResolution returned error: %v", err)
	}
	if got := req.Metadata["resolution"]; got != "1080p" {
		t.Fatalf("resolution = %v, want 1080p", got)
	}
}

func TestForceSeedanceAliasResolutionRejectsConflict(t *testing.T) {
	req := relaycommon.TaskSubmitReq{
		Model:    "Seedance-2.0-720P-海外版",
		Metadata: map[string]interface{}{"resolution": "4K"},
	}
	err := forceSeedanceAliasResolution(&req)
	if err == nil {
		t.Fatal("expected conflict error")
	}
	if !strings.Contains(err.Error(), "requires resolution=720p") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestForceSeedance480PRejects4KConflict(t *testing.T) {
	req := relaycommon.TaskSubmitReq{
		Model:    "Seedance-2.0-480P-海外版",
		Metadata: map[string]interface{}{"resolution": "4K"},
	}
	err := forceSeedanceAliasResolution(&req)
	if err == nil {
		t.Fatal("expected conflict error")
	}
	if !strings.Contains(err.Error(), "requires resolution=480p") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestEstimateBillingUsesForced4KAliasBasePrice(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Set("task_request", relaycommon.TaskSubmitReq{
		Model: "Seedance-2.0-4K-海外版",
		Metadata: map[string]interface{}{
			"resolution": "4k",
			"content": []interface{}{
				map[string]interface{}{"type": "video_url", "video_url": map[string]interface{}{"url": "https://example.com/input.mp4"}},
			},
		},
	})

	ratios := (&TaskAdaptor{}).EstimateBilling(c, &relaycommon.RelayInfo{OriginModelName: "Seedance-2.0-4K-海外版"})
	got, ok := ratios["seedance_unit_price"]
	if !ok {
		t.Fatal("expected seedance_unit_price ratio")
	}
	if got != 2.4/4.0 {
		t.Fatalf("seedance_unit_price = %v, want %v", got, 2.4/4.0)
	}
}

func TestEstimateBillingUsesSeedance15SilentUnitPrice(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Set("task_request", relaycommon.TaskSubmitReq{
		Model: "Seedance-1.5-pro-海外版",
		Metadata: map[string]interface{}{
			"resolution":     "720p",
			"generate_audio": false,
		},
	})

	ratios := (&TaskAdaptor{}).EstimateBilling(c, &relaycommon.RelayInfo{OriginModelName: "Seedance-1.5-pro-海外版"})
	got, ok := ratios["seedance_unit_price"]
	if !ok {
		t.Fatal("expected seedance_unit_price ratio")
	}
	if got != 1.2/2.4 {
		t.Fatalf("seedance_unit_price = %v, want %v", got, 1.2/2.4)
	}
}

func TestParseTaskResultFallsBackToCompletionTokens(t *testing.T) {
	adaptor := &TaskAdaptor{}
	result, err := adaptor.ParseTaskResult([]byte(`{
		"id": "task-1",
		"status": "succeeded",
		"content": {"video_url": "https://example.com/video.mp4"},
		"usage": {"completion_tokens": 1234}
	}`))
	if err != nil {
		t.Fatalf("ParseTaskResult returned error: %v", err)
	}
	if result.Status != model.TaskStatusSuccess {
		t.Fatalf("status = %v, want success", result.Status)
	}
	if result.CompletionTokens != 1234 {
		t.Fatalf("completion tokens = %d, want 1234", result.CompletionTokens)
	}
	if result.TotalTokens != 1234 {
		t.Fatalf("total tokens = %d, want completion token fallback 1234", result.TotalTokens)
	}
}

func TestParseTaskResultPrefersTotalTokens(t *testing.T) {
	adaptor := &TaskAdaptor{}
	result, err := adaptor.ParseTaskResult([]byte(`{
		"id": "task-1",
		"status": "succeeded",
		"usage": {"completion_tokens": 1234, "total_tokens": 5678}
	}`))
	if err != nil {
		t.Fatalf("ParseTaskResult returned error: %v", err)
	}
	if result.TotalTokens != 5678 {
		t.Fatalf("total tokens = %d, want 5678", result.TotalTokens)
	}
}
