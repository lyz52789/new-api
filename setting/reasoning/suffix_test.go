package reasoning

import "testing"

func TestParseDeepSeekV4ThinkingSuffix(t *testing.T) {
	tests := []struct {
		name         string
		model        string
		baseModel    string
		thinkingType string
		effort       string
		ok           bool
	}{
		{
			name:         "disables thinking with none suffix",
			model:        "deepseek-v4-flash-none",
			baseModel:    "deepseek-v4-flash",
			thinkingType: "disabled",
			ok:           true,
		},
		{
			name:         "enables max thinking with max suffix",
			model:        "deepseek-v4-pro-max",
			baseModel:    "deepseek-v4-pro",
			thinkingType: "enabled",
			effort:       "max",
			ok:           true,
		},
		{
			name:      "ignores base v4 model without suffix",
			model:     "deepseek-v4-flash",
			baseModel: "deepseek-v4-flash",
			ok:        false,
		},
		{
			name:      "ignores non deepseek v4 model",
			model:     "deepseek-chat-none",
			baseModel: "deepseek-chat-none",
			ok:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			baseModel, thinkingType, effort, ok := ParseDeepSeekV4ThinkingSuffix(tt.model)
			if baseModel != tt.baseModel || thinkingType != tt.thinkingType || effort != tt.effort || ok != tt.ok {
				t.Fatalf("ParseDeepSeekV4ThinkingSuffix(%q) = (%q, %q, %q, %v), want (%q, %q, %q, %v)", tt.model, baseModel, thinkingType, effort, ok, tt.baseModel, tt.thinkingType, tt.effort, tt.ok)
			}
		})
	}
}
