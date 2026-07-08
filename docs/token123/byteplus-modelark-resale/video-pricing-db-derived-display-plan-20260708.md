# Seedance 视频价格展示改为读取数据库倍率

日期：2026-07-08

## 背景

`/api/pricing` 的 `video_pricing.rows` 之前用代码常量生成展示价格：

```text
官方 USD 价格 * 7.3 * 2.2
```

这导致每次 Token123 只调整数据库里的视频模型 `ModelRatio` 后，实际扣费已经变化，但模型广场的视频分辨率价格表仍显示旧价。2026-07-08 将 Seedance 2.0 从 2.2x 调整到 1.1x 后，这个问题已经出现：

- 实际扣费：生产库 `ModelRatio` 已变为 1.1x。
- 展示价格：代码常量仍按 2.2x 显示。

## 决策

不新增数据库字段。

原因：

- 视频任务实际扣费唯一价格输入已经是 `options.ModelRatio`。
- `updatePricing()` 生成 `Pricing` 时已经把当前 DB `ModelRatio` 写入 `Pricing.ModelRatio`。
- `attachVideoPricing()` 在同一个对象上补 `video_pricing.rows`，可以直接使用该 `ModelRatio`。
- 再新增一个“视频展示倍率”字段会制造第二套价格源，反而容易和实际扣费不一致。

## 新公式

实际视频任务扣费公式：

```text
actualQuota = totalTokens * ModelRatio * GroupRatio * OtherRatio
```

生产展示金额换算：

```text
QuotaPerUnit = 500000
RMB/M = 1,000,000 * ModelRatio * GroupRatio * OtherRatio / 500000
      = 2 * ModelRatio * GroupRatio * OtherRatio
```

`/api/pricing` 的视频价格表不应用 hardcoded markup，而应从 DB 当前倍率推导：

```text
基准售价 RMB/M = 2 * ModelRatio
场景售价 RMB/M = 基准售价 RMB/M * 场景官方 USD/M / 基准官方 USD/M
```

对于 Seedance 2.0 Fast：

```text
ModelRatio = 22.484
基准官方价 = 5.6 USD/M
基准售价 = 22.484 * 2 = 44.968 RMB/M

视频输入售价 = 44.968 * 3.3 / 5.6 = 26.499 RMB/M
```

对于 Seedance 2.0 标准版：

```text
ModelRatio = 28.105
基准官方价 = 7.0 USD/M
基准售价 = 28.105 * 2 = 56.210 RMB/M

1080p 无视频输入售价 = 56.210 * 7.7 / 7.0 = 61.831 RMB/M
```

同一套相对倍率也用于 `sale_rmb_per_video` 和 `sale_rmb_per_second`：

```text
场景售价 RMB = 官方场景 USD * 基准售价 RMB/M / 基准官方 USD/M
```

这等价于保留当前 DB 配置隐含的 RMB/USD 加价口径，但不再把汇率或加价倍率写死在代码里。

## 代码变更

目标文件：

```text
model/pricing.go
model/pricing_test.go
```

变更点：

- 删除 `bytePlusVideoMarkup = 2.2`。
- 删除 `bytePlusVideoUsdCnyRate = 7.3` 的展示依赖。
- 删除 `VideoPricing` 内部隐藏字段 `UsdCnyRate`、`Markup`。
- 新增 `bytePlusVideoSaleCalculator`，从当前 `Pricing.ModelRatio` 和官方基准 USD/M 计算展示价。
- `Seedance-2.0-海外版`、`Seedance-2.0-fast-海外版`、`Seedance-1.5-pro-海外版` 全部复用同一推导方式。
- 测试覆盖：
  - Fast 当前 1.1x DB 倍率显示 `44.968 RMB/M`。
  - Fast 如果 DB 仍是旧倍率会显示 `89.936 RMB/M`，证明展示跟随 DB 倍率。
  - 标准版 1080p 按官方 `7.7 / 7.0` 相对价显示 `61.831 RMB/M`。
  - 公共 JSON 不暴露官方成本、汇率或 markup 字段。

## 后续调整流程

以后调整 Seedance 视频模型价格时，只需要走原生产 DB 工作流：

1. 按目标售价计算新的 `ModelRatio`。
2. 更新生产库 `options.ModelRatio`。
3. 重启 `new-api` 或刷新配置缓存。
4. 验证 `/api/pricing`：
   - `model_ratio` 为新值。
   - `video_pricing.rows` 自动按新值变化。
5. 不需要再改代码里的展示倍率。

## 边界

- 代码里仍保留 BytePlus 官方场景 USD/M 单价表，因为这些是官方计费规则和 `OtherRatio` 的依据；只有 BytePlus 官方价格变更时才需要更新。
- 当前推导不额外乘 `GroupRatio`，保持既有 `group_ratio_applied=false` 语义。生产 Seedance 使用 `byteplus` 分组且 `GroupRatio=1`，与实际扣费一致。
- 如果未来同一公开视频模型同时卖给多个不同 `GroupRatio` 的分组，需要把价格接口扩展为按分组返回展示价；这不是本次需求。

