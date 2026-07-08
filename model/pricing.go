package model

import (
	"fmt"
	"math"
	"strings"

	"sync"
	"time"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/constant"
	"github.com/QuantumNous/new-api/setting/ratio_setting"
	"github.com/QuantumNous/new-api/types"
)

type Pricing struct {
	ModelName              string                  `json:"model_name"`
	Description            string                  `json:"description,omitempty"`
	Icon                   string                  `json:"icon,omitempty"`
	Tags                   string                  `json:"tags,omitempty"`
	VendorID               int                     `json:"vendor_id,omitempty"`
	QuotaType              int                     `json:"quota_type"`
	ModelRatio             float64                 `json:"model_ratio"`
	ModelPrice             float64                 `json:"model_price"`
	OwnerBy                string                  `json:"owner_by"`
	CompletionRatio        float64                 `json:"completion_ratio"`
	CacheRatio             *float64                `json:"cache_ratio,omitempty"`
	CreateCacheRatio       *float64                `json:"create_cache_ratio,omitempty"`
	ImageRatio             *float64                `json:"image_ratio,omitempty"`
	AudioRatio             *float64                `json:"audio_ratio,omitempty"`
	AudioCompletionRatio   *float64                `json:"audio_completion_ratio,omitempty"`
	EnableGroup            []string                `json:"enable_groups"`
	SupportedEndpointTypes []constant.EndpointType `json:"supported_endpoint_types"`
	VideoPricing           *VideoPricing           `json:"video_pricing,omitempty"`
	PricingVersion         string                  `json:"pricing_version,omitempty"`
}

type VideoPricing struct {
	BillingMode       string            `json:"billing_mode"`
	Currency          string            `json:"currency"`
	GroupRatioApplied bool              `json:"group_ratio_applied"`
	Formula           string            `json:"formula"`
	Rows              []VideoPricingRow `json:"rows"`
}

type VideoPricingRow struct {
	Resolution             string  `json:"resolution"`
	Scenario               string  `json:"scenario"`
	ScenarioLabel          string  `json:"scenario_label"`
	OfficialUSDPerMTokens  float64 `json:"-"`
	SaleRMBPerMTokens      float64 `json:"sale_rmb_per_m_tokens,omitempty"`
	OfficialUSDPerVideo    float64 `json:"-"`
	SaleRMBPerVideo        float64 `json:"sale_rmb_per_video,omitempty"`
	OfficialUSDPerSecond   float64 `json:"-"`
	SaleRMBPerSecond       float64 `json:"sale_rmb_per_second,omitempty"`
	OfficialUSDPerVideoMin float64 `json:"-"`
	OfficialUSDPerVideoMax float64 `json:"-"`
	SaleRMBPerVideoMin     float64 `json:"sale_rmb_per_video_min,omitempty"`
	SaleRMBPerVideoMax     float64 `json:"sale_rmb_per_video_max,omitempty"`
}

type PricingVendor struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Icon        string `json:"icon,omitempty"`
}

var (
	pricingMap           []Pricing
	vendorsList          []PricingVendor
	supportedEndpointMap map[string]common.EndpointInfo
	lastGetPricingTime   time.Time
	updatePricingLock    sync.Mutex

	// 缓存映射：模型名 -> 启用分组 / 计费类型
	modelEnableGroups     = make(map[string][]string)
	modelQuotaTypeMap     = make(map[string]int)
	modelEnableGroupsLock = sync.RWMutex{}
)

var (
	modelSupportEndpointTypes = make(map[string][]constant.EndpointType)
	modelSupportEndpointsLock = sync.RWMutex{}
)

func GetPricing() []Pricing {
	if time.Since(lastGetPricingTime) > time.Minute*1 || len(pricingMap) == 0 {
		updatePricingLock.Lock()
		defer updatePricingLock.Unlock()
		// Double check after acquiring the lock
		if time.Since(lastGetPricingTime) > time.Minute*1 || len(pricingMap) == 0 {
			modelSupportEndpointsLock.Lock()
			defer modelSupportEndpointsLock.Unlock()
			updatePricing()
		}
	}
	return pricingMap
}

// GetVendors 返回当前定价接口使用到的供应商信息
func GetVendors() []PricingVendor {
	if time.Since(lastGetPricingTime) > time.Minute*1 || len(pricingMap) == 0 {
		// 保证先刷新一次
		GetPricing()
	}
	return vendorsList
}

func GetModelSupportEndpointTypes(model string) []constant.EndpointType {
	if model == "" {
		return make([]constant.EndpointType, 0)
	}
	modelSupportEndpointsLock.RLock()
	defer modelSupportEndpointsLock.RUnlock()
	if endpoints, ok := modelSupportEndpointTypes[model]; ok {
		return endpoints
	}
	return make([]constant.EndpointType, 0)
}

func updatePricing() {
	//modelRatios := common.GetModelRatios()
	enableAbilities, err := GetAllEnableAbilityWithChannels()
	if err != nil {
		common.SysLog(fmt.Sprintf("GetAllEnableAbilityWithChannels error: %v", err))
		return
	}
	// 预加载模型元数据与供应商一次，避免循环查询
	var allMeta []Model
	_ = DB.Find(&allMeta).Error
	metaMap := make(map[string]*Model)
	prefixList := make([]*Model, 0)
	suffixList := make([]*Model, 0)
	containsList := make([]*Model, 0)
	for i := range allMeta {
		m := &allMeta[i]
		if m.NameRule == NameRuleExact {
			metaMap[m.ModelName] = m
		} else {
			switch m.NameRule {
			case NameRulePrefix:
				prefixList = append(prefixList, m)
			case NameRuleSuffix:
				suffixList = append(suffixList, m)
			case NameRuleContains:
				containsList = append(containsList, m)
			}
		}
	}

	// 将非精确规则模型匹配到 metaMap
	for _, m := range prefixList {
		for _, pricingModel := range enableAbilities {
			if strings.HasPrefix(pricingModel.Model, m.ModelName) {
				if _, exists := metaMap[pricingModel.Model]; !exists {
					metaMap[pricingModel.Model] = m
				}
			}
		}
	}
	for _, m := range suffixList {
		for _, pricingModel := range enableAbilities {
			if strings.HasSuffix(pricingModel.Model, m.ModelName) {
				if _, exists := metaMap[pricingModel.Model]; !exists {
					metaMap[pricingModel.Model] = m
				}
			}
		}
	}
	for _, m := range containsList {
		for _, pricingModel := range enableAbilities {
			if strings.Contains(pricingModel.Model, m.ModelName) {
				if _, exists := metaMap[pricingModel.Model]; !exists {
					metaMap[pricingModel.Model] = m
				}
			}
		}
	}

	// 预加载供应商
	var vendors []Vendor
	_ = DB.Find(&vendors).Error
	vendorMap := make(map[int]*Vendor)
	for i := range vendors {
		vendorMap[vendors[i].Id] = &vendors[i]
	}

	// 初始化默认供应商映射
	initDefaultVendorMapping(metaMap, vendorMap, enableAbilities)

	// 构建对前端友好的供应商列表
	vendorsList = make([]PricingVendor, 0, len(vendorMap))
	for _, v := range vendorMap {
		vendorsList = append(vendorsList, PricingVendor{
			ID:          v.Id,
			Name:        v.Name,
			Description: v.Description,
			Icon:        v.Icon,
		})
	}

	modelGroupsMap := make(map[string]*types.Set[string])

	for _, ability := range enableAbilities {
		groups, ok := modelGroupsMap[ability.Model]
		if !ok {
			groups = types.NewSet[string]()
			modelGroupsMap[ability.Model] = groups
		}
		groups.Add(ability.Group)
	}

	//这里使用切片而不是Set，因为一个模型可能支持多个端点类型，并且第一个端点是优先使用端点
	modelSupportEndpointsStr := make(map[string][]string)

	// 先根据已有能力填充原生端点
	for _, ability := range enableAbilities {
		endpoints := modelSupportEndpointsStr[ability.Model]
		channelTypes := common.GetEndpointTypesByChannelType(ability.ChannelType, ability.Model)
		for _, channelType := range channelTypes {
			if !common.StringsContains(endpoints, string(channelType)) {
				endpoints = append(endpoints, string(channelType))
			}
		}
		modelSupportEndpointsStr[ability.Model] = endpoints
	}

	// 再补充模型自定义端点：若配置有效则替换默认端点，不做合并
	for modelName, meta := range metaMap {
		if strings.TrimSpace(meta.Endpoints) == "" {
			continue
		}
		var raw map[string]interface{}
		if err := common.Unmarshal([]byte(meta.Endpoints), &raw); err == nil {
			endpoints := make([]string, 0, len(raw))
			for k, v := range raw {
				switch v.(type) {
				case string, map[string]interface{}:
					if !common.StringsContains(endpoints, k) {
						endpoints = append(endpoints, k)
					}
				}
			}
			if len(endpoints) > 0 {
				modelSupportEndpointsStr[modelName] = endpoints
			}
		}
	}

	modelSupportEndpointTypes = make(map[string][]constant.EndpointType)
	for model, endpoints := range modelSupportEndpointsStr {
		supportedEndpoints := make([]constant.EndpointType, 0)
		for _, endpointStr := range endpoints {
			endpointType := constant.EndpointType(endpointStr)
			supportedEndpoints = append(supportedEndpoints, endpointType)
		}
		modelSupportEndpointTypes[model] = supportedEndpoints
	}

	// 构建全局 supportedEndpointMap（默认 + 自定义覆盖）
	supportedEndpointMap = make(map[string]common.EndpointInfo)
	// 1. 默认端点
	for _, endpoints := range modelSupportEndpointTypes {
		for _, et := range endpoints {
			if info, ok := common.GetDefaultEndpointInfo(et); ok {
				if _, exists := supportedEndpointMap[string(et)]; !exists {
					supportedEndpointMap[string(et)] = info
				}
			}
		}
	}
	// 2. 自定义端点（models 表）覆盖默认
	for _, meta := range metaMap {
		if strings.TrimSpace(meta.Endpoints) == "" {
			continue
		}
		var raw map[string]interface{}
		if err := common.Unmarshal([]byte(meta.Endpoints), &raw); err == nil {
			for k, v := range raw {
				switch val := v.(type) {
				case string:
					supportedEndpointMap[k] = common.EndpointInfo{Path: val, Method: "POST"}
				case map[string]interface{}:
					ep := common.EndpointInfo{Method: "POST"}
					if p, ok := val["path"].(string); ok {
						ep.Path = p
					}
					if m, ok := val["method"].(string); ok {
						ep.Method = strings.ToUpper(m)
					}
					supportedEndpointMap[k] = ep
				default:
					// ignore unsupported types
				}
			}
		}
	}

	pricingMap = make([]Pricing, 0)
	for model, groups := range modelGroupsMap {
		pricing := Pricing{
			ModelName:              model,
			EnableGroup:            groups.Items(),
			SupportedEndpointTypes: modelSupportEndpointTypes[model],
		}

		// 补充模型元数据（描述、标签、供应商、状态）
		if meta, ok := metaMap[model]; ok {
			// 若模型被禁用(status!=1)，则直接跳过，不返回给前端
			if meta.Status != 1 {
				continue
			}
			pricing.Description = meta.Description
			pricing.Icon = meta.Icon
			pricing.Tags = meta.Tags
			pricing.VendorID = meta.VendorID
		}
		modelPrice, findPrice := ratio_setting.GetModelPrice(model, false)
		if findPrice {
			pricing.ModelPrice = modelPrice
			pricing.QuotaType = 1
		} else {
			modelRatio, _, _ := ratio_setting.GetModelRatio(model)
			pricing.ModelRatio = modelRatio
			pricing.CompletionRatio = ratio_setting.GetCompletionRatio(model)
			pricing.QuotaType = 0
		}
		if cacheRatio, ok := ratio_setting.GetCacheRatio(model); ok {
			pricing.CacheRatio = &cacheRatio
		}
		if createCacheRatio, ok := ratio_setting.GetCreateCacheRatio(model); ok {
			pricing.CreateCacheRatio = &createCacheRatio
		}
		if imageRatio, ok := ratio_setting.GetImageRatio(model); ok {
			pricing.ImageRatio = &imageRatio
		}
		if ratio_setting.ContainsAudioRatio(model) {
			audioRatio := ratio_setting.GetAudioRatio(model)
			pricing.AudioRatio = &audioRatio
		}
		if ratio_setting.ContainsAudioCompletionRatio(model) {
			audioCompletionRatio := ratio_setting.GetAudioCompletionRatio(model)
			pricing.AudioCompletionRatio = &audioCompletionRatio
		}
		attachVideoPricing(&pricing)
		pricingMap = append(pricingMap, pricing)
	}

	// 防止大更新后数据不通用
	if len(pricingMap) > 0 {
		pricingMap[0].PricingVersion = "5a90f2b86c08bd983a9a2e6d66c255f4eaef9c4bc934386d2b6ae84ef0ff1f1f"
	}

	// 刷新缓存映射，供高并发快速查询
	modelEnableGroupsLock.Lock()
	modelEnableGroups = make(map[string][]string)
	modelQuotaTypeMap = make(map[string]int)
	for _, p := range pricingMap {
		modelEnableGroups[p.ModelName] = p.EnableGroup
		modelQuotaTypeMap[p.ModelName] = p.QuotaType
	}
	modelEnableGroupsLock.Unlock()

	lastGetPricingTime = time.Now()
}

// GetSupportedEndpointMap 返回全局端点到路径的映射
func GetSupportedEndpointMap() map[string]common.EndpointInfo {
	return supportedEndpointMap
}

type bytePlusVideoSaleCalculator struct {
	baseOfficialUSDPerMTokens float64
	baseSaleRMBPerMTokens     float64
}

func newBytePlusVideoSaleCalculator(modelRatio, baseOfficialUSDPerMTokens float64) bytePlusVideoSaleCalculator {
	return bytePlusVideoSaleCalculator{
		baseOfficialUSDPerMTokens: baseOfficialUSDPerMTokens,
		baseSaleRMBPerMTokens:     modelRatio * 2,
	}
}

func (c bytePlusVideoSaleCalculator) saleRMB(officialUSD float64) float64 {
	if c.baseOfficialUSDPerMTokens <= 0 || c.baseSaleRMBPerMTokens <= 0 {
		return 0
	}
	return math.Round(officialUSD*c.baseSaleRMBPerMTokens/c.baseOfficialUSDPerMTokens*10000) / 10000
}

func seedance20VideoRows(calc bytePlusVideoSaleCalculator) []VideoPricingRow {
	return []VideoPricingRow{
		seedance20NoVideoRow(calc, "480p", 7.0, 0.35, 0.07),
		seedance20NoVideoRow(calc, "720p", 7.0, 0.76, 0.15),
		seedance20NoVideoRow(calc, "1080p", 7.7, 1.87, 0.37),
		seedance20NoVideoRow(calc, "4k", 4.0, 3.89, 0.78),
		seedance20VideoInputRow(calc, "480p", 4.3, 0.39, 0.86),
		seedance20VideoInputRow(calc, "720p", 4.3, 0.84, 1.86),
		seedance20VideoInputRow(calc, "1080p", 4.7, 2.06, 4.57),
		seedance20VideoInputRow(calc, "4k", 2.4, 4.20, 9.33),
	}
}

func seedance20FastVideoRows(calc bytePlusVideoSaleCalculator) []VideoPricingRow {
	return []VideoPricingRow{
		seedance20NoVideoRow(calc, "480p", 5.6, 0.28, 0.06),
		seedance20NoVideoRow(calc, "720p", 5.6, 0.60, 0.12),
		seedance20VideoInputRow(calc, "480p", 3.3, 0.30, 0.66),
		seedance20VideoInputRow(calc, "720p", 3.3, 0.64, 1.43),
	}
}

func seedance20NoVideoRow(calc bytePlusVideoSaleCalculator, resolution string, usdPerMTokens, usdPerVideo, usdPerSecond float64) VideoPricingRow {
	return VideoPricingRow{
		Resolution:            resolution,
		Scenario:              "text_or_image_to_video",
		ScenarioLabel:         "文本/图片生成视频",
		OfficialUSDPerMTokens: usdPerMTokens,
		SaleRMBPerMTokens:     calc.saleRMB(usdPerMTokens),
		OfficialUSDPerVideo:   usdPerVideo,
		SaleRMBPerVideo:       calc.saleRMB(usdPerVideo),
		OfficialUSDPerSecond:  usdPerSecond,
		SaleRMBPerSecond:      calc.saleRMB(usdPerSecond),
	}
}

func seedance20VideoInputRow(calc bytePlusVideoSaleCalculator, resolution string, usdPerMTokens, usdPerVideoMin, usdPerVideoMax float64) VideoPricingRow {
	return VideoPricingRow{
		Resolution:             resolution,
		Scenario:               "video_to_video",
		ScenarioLabel:          "视频输入生成视频",
		OfficialUSDPerMTokens:  usdPerMTokens,
		SaleRMBPerMTokens:      calc.saleRMB(usdPerMTokens),
		OfficialUSDPerVideoMin: usdPerVideoMin,
		OfficialUSDPerVideoMax: usdPerVideoMax,
		SaleRMBPerVideoMin:     calc.saleRMB(usdPerVideoMin),
		SaleRMBPerVideoMax:     calc.saleRMB(usdPerVideoMax),
	}
}

func seedance15VideoRows(calc bytePlusVideoSaleCalculator) []VideoPricingRow {
	return []VideoPricingRow{
		seedance15Row(calc, "480p", "audio", "带音频", 2.4, 0.12),
		seedance15Row(calc, "480p", "silent", "无音频", 1.2, 0.06),
		seedance15Row(calc, "480p", "draft_audio", "草稿模式带音频", 2.4, 0.07),
		seedance15Row(calc, "480p", "draft_silent", "草稿模式无音频", 1.2, 0.04),
		seedance15Row(calc, "720p", "audio", "带音频", 2.4, 0.26),
		seedance15Row(calc, "720p", "silent", "无音频", 1.2, 0.13),
		seedance15Row(calc, "1080p", "audio", "带音频", 2.4, 0.58),
		seedance15Row(calc, "1080p", "silent", "无音频", 1.2, 0.29),
	}
}

func seedance15Row(calc bytePlusVideoSaleCalculator, resolution, scenario, scenarioLabel string, usdPerMTokens, usdPerVideo float64) VideoPricingRow {
	const officialDurationSeconds = 5.0
	return VideoPricingRow{
		Resolution:            resolution,
		Scenario:              scenario,
		ScenarioLabel:         scenarioLabel,
		OfficialUSDPerMTokens: usdPerMTokens,
		SaleRMBPerMTokens:     calc.saleRMB(usdPerMTokens),
		OfficialUSDPerVideo:   usdPerVideo,
		SaleRMBPerVideo:       calc.saleRMB(usdPerVideo),
		OfficialUSDPerSecond:  usdPerVideo / officialDurationSeconds,
		SaleRMBPerSecond:      calc.saleRMB(usdPerVideo / officialDurationSeconds),
	}
}

func attachVideoPricing(pricing *Pricing) {
	switch pricing.ModelName {
	case "Seedance-2.0-海外版", "doubao-seedance-2-0-260128":
		calc := newBytePlusVideoSaleCalculator(pricing.ModelRatio, 7.0)
		pricing.VideoPricing = &VideoPricing{
			BillingMode:       "resolution_usage_tokens",
			Currency:          "CNY",
			GroupRatioApplied: false,
			Formula:           "按分辨率、输入类型和实际用量动态计费；¥/条、¥/秒为官方价格折算参考，最终扣费以任务完成后的实际 tokens 为准",
			Rows:              seedance20VideoRows(calc),
		}
	case "Seedance-2.0-fast-海外版", "doubao-seedance-2-0-fast-260128":
		calc := newBytePlusVideoSaleCalculator(pricing.ModelRatio, 5.6)
		pricing.VideoPricing = &VideoPricing{
			BillingMode:       "resolution_usage_tokens",
			Currency:          "CNY",
			GroupRatioApplied: false,
			Formula:           "fast 仅支持 480p/720p；¥/条、¥/秒为官方价格折算参考，最终扣费以任务完成后的实际 tokens 为准",
			Rows:              seedance20FastVideoRows(calc),
		}
	case "Seedance-1.5-pro-海外版", "doubao-seedance-1-5-pro-251215":
		calc := newBytePlusVideoSaleCalculator(pricing.ModelRatio, 2.4)
		pricing.VideoPricing = &VideoPricing{
			BillingMode:       "resolution_audio_usage_tokens",
			Currency:          "CNY",
			GroupRatioApplied: false,
			Formula:           "按分辨率、是否生成音频和实际用量动态计费；¥/条、¥/秒按官方 5 秒视频价格折算，generate_audio=false 时使用无音频售价",
			Rows:              seedance15VideoRows(calc),
		}
	}
}
