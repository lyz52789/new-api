package doubao

import "strings"

var ModelList = []string{
	"doubao-seedance-1-0-pro-250528",
	"doubao-seedance-1-0-lite-t2v",
	"doubao-seedance-1-0-lite-i2v",
	"doubao-seedance-1-5-pro-251215",
	"doubao-seedance-2-0-260128",
	"doubao-seedance-2-0-fast-260128",
}

var ChannelName = "doubao-video"

// videoPriceKey 价格表的键：输出分辨率档（is1080p/is4k 均为 false 即 480p/720p 基准档）、输入是否含视频。
type videoPriceKey struct {
	is1080p  bool
	is4k     bool
	hasVideo bool
}

// videoPriceTable stores BytePlus overseas official USD/M token prices for
// Seedance 2.0. Only the relative ratio is used at runtime:
// actual price / base price becomes OtherRatio. Generic aliases use the
// 480p/720p no-video-input price as the base. Fixed-resolution public aliases
// use their own no-video-input price as the base, so marketplace display price
// and runtime unit price stay aligned.
var videoPriceTable = map[string]map[videoPriceKey]float64{
	"doubao-seedance-2-0-260128": seedance20OverseasPrices,
	"Seedance-2.0-海外版":           seedance20OverseasPrices,
	"Seedance-2.0-480P-海外版":      seedance20OverseasPrices,
	"Seedance-2.0-720P-海外版":      seedance20OverseasPrices,
	"Seedance-2.0-1080P-海外版":     seedance20OverseasPrices,
	"Seedance-2.0-4K-海外版":        seedance20OverseasPrices,
}

var seedance20OverseasPrices = map[videoPriceKey]float64{
	{hasVideo: false}:                7.0,
	{hasVideo: true}:                 4.3,
	{is1080p: true, hasVideo: false}: 7.7,
	{is1080p: true, hasVideo: true}:  4.7,
	{is4k: true, hasVideo: false}:    4.0,
	{is4k: true, hasVideo: true}:     2.4,
}

var seedance15ProOverseasPrices = map[bool]float64{
	true:  2.4, // audio video, online inference, USD/M tokens
	false: 1.2, // silent video, online inference, USD/M tokens
}

// GetVideoInputRatio 返回指定模型在给定输出分辨率/是否含视频输入下，相对基准价的计费倍率。
// 第二个返回值表示该模型是否配置了价格表；倍率为 1.0 时调用方可忽略该 OtherRatio。
func GetVideoInputRatio(modelName, resolution string, hasVideo bool) (float64, bool) {
	prices, ok := videoPriceTable[modelName]
	if !ok {
		return 0, false
	}
	baseKey := videoPriceKey{}
	actualResolution := NormalizeVideoResolution(resolution)
	if forcedResolution, ok := ForcedSeedanceResolutionForModel(modelName); ok {
		actualResolution = forcedResolution
		baseKey = videoPriceKeyForResolution(forcedResolution, false)
	}
	base := prices[baseKey]
	if base <= 0 {
		return 0, false
	}
	price, ok := prices[videoPriceKeyForResolution(actualResolution, hasVideo)]
	if !ok {
		// 未配置的组合按基准价计费即可。
		return 1.0, true
	}
	return price / base, true
}

func GetSeedance15AudioRatio(modelName string, generateAudio bool) (float64, bool) {
	name := strings.ToLower(modelName)
	if !strings.Contains(name, "seedance-1.5-pro") &&
		!strings.Contains(name, "seedance-1-5-pro") &&
		!strings.Contains(name, "seedance15") {
		return 0, false
	}
	base := seedance15ProOverseasPrices[true]
	price := seedance15ProOverseasPrices[generateAudio]
	if base <= 0 || price <= 0 {
		return 0, false
	}
	return price / base, true
}

func videoPriceKeyForResolution(resolution string, hasVideo bool) videoPriceKey {
	res := NormalizeVideoResolution(resolution)
	return videoPriceKey{is1080p: res == "1080p", is4k: res == "4k", hasVideo: hasVideo}
}

func NormalizeVideoResolution(resolution string) string {
	res := strings.ToLower(strings.TrimSpace(resolution))
	res = strings.ReplaceAll(res, "_", "")
	res = strings.ReplaceAll(res, "-", "")
	res = strings.ReplaceAll(res, " ", "")
	switch res {
	case "480", "480p":
		return "480p"
	case "720", "720p":
		return "720p"
	case "1080", "1080p":
		return "1080p"
	case "4k", "uhd":
		return "4k"
	default:
		return res
	}
}

// ForcedSeedanceResolutionForModel returns the resolution encoded in a public
// Seedance 2.0 alias. Generic aliases intentionally return false.
func ForcedSeedanceResolutionForModel(modelName string) (string, bool) {
	name := strings.ToLower(modelName)
	if !strings.Contains(name, "seedance-2.0") && !strings.Contains(name, "seedance-2-0") {
		return "", false
	}
	switch {
	case strings.Contains(name, "4k"):
		return "4k", true
	case strings.Contains(name, "1080p"), strings.Contains(name, "1080"):
		return "1080p", true
	case strings.Contains(name, "720p"), strings.Contains(name, "720"):
		return "720p", true
	case strings.Contains(name, "480p"), strings.Contains(name, "480"):
		return "480p", true
	default:
		return "", false
	}
}
