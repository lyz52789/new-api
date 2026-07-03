package common

import (
	"testing"

	"github.com/QuantumNous/new-api/constant"
)

func TestDoubaoVideoUsesOpenAIVideoEndpoint(t *testing.T) {
	endpointTypes := GetEndpointTypesByChannelType(constant.ChannelTypeDoubaoVideo, "Seedance-2.0-720P-海外版")
	if len(endpointTypes) != 1 || endpointTypes[0] != constant.EndpointTypeOpenAIVideo {
		t.Fatalf("endpoint types = %v, want [%s]", endpointTypes, constant.EndpointTypeOpenAIVideo)
	}

	endpointInfo, ok := GetDefaultEndpointInfo(endpointTypes[0])
	if !ok {
		t.Fatalf("default endpoint info missing for %s", endpointTypes[0])
	}
	if endpointInfo.Path != "/v1/videos" || endpointInfo.Method != "POST" {
		t.Fatalf("endpoint info = %+v, want POST /v1/videos", endpointInfo)
	}
}
