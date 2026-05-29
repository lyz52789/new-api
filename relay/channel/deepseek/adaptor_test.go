package deepseek

import (
	"testing"

	"github.com/QuantumNous/new-api/dto"
	relaycommon "github.com/QuantumNous/new-api/relay/common"
)

func TestApplyDeepSeekV4OpenAIThinkingSuffix(t *testing.T) {
	request := &dto.GeneralOpenAIRequest{Model: "deepseek-v4-flash-none"}
	info := &relaycommon.RelayInfo{ChannelMeta: &relaycommon.ChannelMeta{UpstreamModelName: "deepseek-v4-flash-none"}}

	if err := applyDeepSeekV4OpenAIThinkingSuffix(info, request); err != nil {
		t.Fatalf("applyDeepSeekV4OpenAIThinkingSuffix returned error: %v", err)
	}

	if request.Model != "deepseek-v4-flash" {
		t.Fatalf("request.Model = %q, want deepseek-v4-flash", request.Model)
	}
	if string(request.THINKING) != `{"type":"disabled"}` {
		t.Fatalf("request.THINKING = %s, want disabled thinking", request.THINKING)
	}
	if request.ReasoningEffort != "" {
		t.Fatalf("request.ReasoningEffort = %q, want empty", request.ReasoningEffort)
	}
	if info.ChannelMeta.UpstreamModelName != "deepseek-v4-flash" {
		t.Fatalf("UpstreamModelName = %q, want deepseek-v4-flash", info.ChannelMeta.UpstreamModelName)
	}
}

func TestApplyDeepSeekV4ClaudeThinkingSuffix(t *testing.T) {
	request := &dto.ClaudeRequest{Model: "deepseek-v4-pro-max"}
	info := &relaycommon.RelayInfo{ChannelMeta: &relaycommon.ChannelMeta{UpstreamModelName: "deepseek-v4-pro-max"}}

	if err := applyDeepSeekV4ClaudeThinkingSuffix(info, request); err != nil {
		t.Fatalf("applyDeepSeekV4ClaudeThinkingSuffix returned error: %v", err)
	}

	if request.Model != "deepseek-v4-pro" {
		t.Fatalf("request.Model = %q, want deepseek-v4-pro", request.Model)
	}
	if request.Thinking == nil || request.Thinking.Type != "enabled" {
		t.Fatalf("request.Thinking = %#v, want enabled thinking", request.Thinking)
	}
	if string(request.OutputConfig) != `{"effort":"max"}` {
		t.Fatalf("request.OutputConfig = %s, want max effort", request.OutputConfig)
	}
	if info.ReasoningEffort != "max" {
		t.Fatalf("ReasoningEffort = %q, want max", info.ReasoningEffort)
	}
	if info.ChannelMeta.UpstreamModelName != "deepseek-v4-pro" {
		t.Fatalf("UpstreamModelName = %q, want deepseek-v4-pro", info.ChannelMeta.UpstreamModelName)
	}
}
