# Seedance 2.0 海外版 10% 加价价格说明

日期：2026-07-08

## 目标

把 Token123 线上两个 BytePlus Seedance 2.0 视频模型从当前 `2.2x` 成本售价调整为 `1.1x` 成本售价：

- `Seedance-2.0-海外版`
- `Seedance-2.0-fast-海外版`

本次只改生产库 `options.ModelRatio` 中这两个模型的基准倍率，不改渠道、分组、模型广场 metadata、`USDExchangeRate`、`Price`、`CompletionRatio`、`CacheRatio`、`ModelPrice`。

## 线上现状

2026-07-08 只读生产查询结果：

| 项 | 当前值 |
| --- | ---: |
| `GroupRatio.byteplus` | `1` |
| `USDExchangeRate` | `1` |
| `Price` | `1` |
| `ModelRatio.Seedance-2.0-海外版` | `56.210` |
| `ModelRatio.Seedance-2.0-fast-海外版` | `44.968` |

换算验证：

```text
Seedance-2.0-海外版 当前售价 = 56.210 * 2 * 1 = 112.420 RMB/M
112.420 = 7.0 * 7.3 * 2.2

Seedance-2.0-fast-海外版 当前售价 = 44.968 * 2 * 1 = 89.936 RMB/M
89.936 = 5.6 * 7.3 * 2.2
```

## 价格来源和口径

BytePlus 官方 ModelArk Pricing 页：

- URL: `https://docs.byteplus.com/en/docs/ModelArk/1544106`
- 2026-07-08 打开页面显示 `Last updated: July 2, 2026 03:56:43`
- 本地官方快照：`/Users/lyz/byteplus-pricing-snapshot.md`

官方单位价：

| 模型 | 场景 | 官方 USD/M tokens |
| --- | --- | ---: |
| Seedance 2.0 | 480p/720p，无视频输入 | `7.0` |
| Seedance 2.0 | 480p/720p，有视频输入 | `4.3` |
| Seedance 2.0 | 1080p，无视频输入 | `7.7` |
| Seedance 2.0 | 1080p，有视频输入 | `4.7` |
| Seedance 2.0 | 4K，无视频输入 | `4.0` |
| Seedance 2.0 | 4K，有视频输入 | `2.4` |
| Seedance 2.0 Fast | 480p/720p，无视频输入 | `5.6` |
| Seedance 2.0 Fast | 480p/720p，有视频输入 | `3.3` |

Token123 继续沿用人民币成本汇率口径：

```text
成本 RMB/M = 官方 USD/M * 7.3
售价 RMB/M = 官方 USD/M * 7.3 * 1.1
ModelRatio = 售价 RMB/M / 2 / GroupRatio
GroupRatio.byteplus = 1
```

## 调整后价格

| 模型 | 基准官方 USD/M | 成本 RMB/M | 调整后售价 RMB/M | 调整后 ModelRatio | 成本加价 |
| --- | ---: | ---: | ---: | ---: | ---: |
| `Seedance-2.0-海外版` | `7.0` | `51.100` | `56.210` | `28.105` | `10%` |
| `Seedance-2.0-fast-海外版` | `5.6` | `40.880` | `44.968` | `22.484` | `10%` |

Seedance 2.0 标准版仍由代码按请求分辨率/视频输入写入 `OtherRatios.seedance_unit_price`，因此各场景最终单位价为：

| 场景 | 官方 USD/M | 成本 RMB/M | 调整后售价 RMB/M | 毛利 RMB/M |
| --- | ---: | ---: | ---: | ---: |
| 480p/720p，无视频输入 | `7.0` | `51.100` | `56.210` | `5.110` |
| 480p/720p，有视频输入 | `4.3` | `31.390` | `34.529` | `3.139` |
| 1080p，无视频输入 | `7.7` | `56.210` | `61.831` | `5.621` |
| 1080p，有视频输入 | `4.7` | `34.310` | `37.741` | `3.431` |
| 4K，无视频输入 | `4.0` | `29.200` | `32.120` | `2.920` |
| 4K，有视频输入 | `2.4` | `17.520` | `19.272` | `1.752` |

Seedance 2.0 Fast 最终单位价：

| 场景 | 官方 USD/M | 成本 RMB/M | 调整后售价 RMB/M | 毛利 RMB/M |
| --- | ---: | ---: | ---: | ---: |
| 480p/720p，无视频输入 | `5.6` | `40.880` | `44.968` | `4.088` |
| 480p/720p，有视频输入 | `3.3` | `24.090` | `26.499` | `2.409` |

## 不亏钱判断

按 API 调用成本口径，本次价格不会亏钱：

```text
售价 = 官方成本 * 7.3 * 1.1
毛利 = 官方成本 * 7.3 * 0.1
毛利率(按售价) = 1 - 1 / 1.1 = 9.09%
```

也就是说，只要 BytePlus 仍按上述官方 USD/M 计费，且内部结算仍按 `1 USD = 7.3 RMB` 口径核算，每个视频 token 场景都有 10% 成本加价。

风险边界：

- 这里证明的是模型调用成本不亏，不包含支付手续费、坏账、优惠券、税费、汇兑波动、人工运营成本。
- 若实际综合成本超过售价的 `9.09%` 毛利空间，或者实际换汇成本高于 `8.03 RMB/USD`，则全口径利润可能变薄或转负。
- 如果 BytePlus 官方价格变更，应重新按官方新价计算后再同步 Token123。

## 实际扣费公式

视频任务完成后，Token123/new-api 用 BytePlus 返回的 `usage.total_tokens` 重新计算实际扣费额度：

```text
actualQuota = totalTokens * ModelRatio * GroupRatio * OtherRatio
```

当前生产：

```text
QuotaPerUnit = 500000
quota_display_type = CNY
USDExchangeRate = 1
```

所以折算成每 100 万 tokens 展示价时：

```text
RMB/M = 1,000,000 * ModelRatio * GroupRatio * OtherRatio / 500,000
      = 2 * ModelRatio * GroupRatio * OtherRatio
```

这就是数据库 `ModelRatio` 要写目标 RMB/M 一半的原因。不是业务价格再额外乘 2，而是额度点数到金额的系统换算。

## 生产执行结果

执行时间：2026-07-08

已执行：

```text
docs/token123/byteplus-modelark-resale/sql/20260708-seedance20-110-markup-rollout.sql
docker restart new-api
docs/token123/byteplus-modelark-resale/sql/20260708-seedance20-110-markup-verify.sql
```

执行结果：

| 模型 | 官方 USD/M | 成本 RMB/M | 实际售价 RMB/M | 成本倍数 |
| --- | ---: | ---: | ---: | ---: |
| `Seedance-2.0-海外版` | `7.0` | `51.100000` | `56.210000` | `1.100000` |
| `Seedance-2.0-fast-海外版` | `5.6` | `40.880000` | `44.968000` | `1.100000` |

生产验证：

```text
ModelRatio.Seedance-2.0-海外版 = 28.105
ModelRatio.Seedance-2.0-fast-海外版 = 22.484
channel = byteplus-modelark-ap-sea-video, type=54, status=1
abilities = 两个模型均 byteplus enabled=true
backup = token123_change_log.token123_seedance20_110_markup_20260708
new-api = healthy
/api/status = success=true
```

注意：本次按“原线上模型价格 DB 工作流”已经更新实际扣费。变更执行后，生产 `/api/pricing` 的 `model_ratio` 已返回新值；当 [video-pricing-db-derived-display-plan-20260708.md](./video-pricing-db-derived-display-plan-20260708.md) 对应代码发布后，`video_pricing.rows` 也会从数据库 `ModelRatio` 自动推导，不再需要代码常量同步展示倍率。

## 执行文件

- Rollout: `docs/token123/byteplus-modelark-resale/sql/20260708-seedance20-110-markup-rollout.sql`
- Verify: `docs/token123/byteplus-modelark-resale/sql/20260708-seedance20-110-markup-verify.sql`
- Rollback: `docs/token123/byteplus-modelark-resale/sql/20260708-seedance20-110-markup-rollback.sql`
- Backup table: `token123_change_log.token123_seedance20_110_markup_20260708`
