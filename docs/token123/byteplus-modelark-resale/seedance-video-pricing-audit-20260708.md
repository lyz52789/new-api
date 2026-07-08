# Seedance Video Pricing Audit - 2026-07-08

## Scope

Models:

- `Seedance-2.0-海外版`
- `Seedance-2.0-fast-海外版`
- `Seedance-1.5-pro-海外版`

Sources:

- BytePlus ModelArk Pricing page, last updated July 2, 2026 03:56:43.
- Local official page snapshot: `/Users/lyz/byteplus-pricing-snapshot.md`.
- Production Token123 `/api/pricing`.
- Production `options.ModelRatio`.

## Billing Formula

Video task settlement is token based:

```text
actual_quota = total_tokens * ModelRatio * GroupRatio * product(OtherRatios)
```

For the `byteplus` group:

```text
GroupRatio = 1
```

The public RMB/M display price is:

```text
sale_rmb_per_m_tokens = ModelRatio * 2 * GroupRatio * product(OtherRatios)
```

For Seedance 2.0 display rows, `OtherRatio` is the official scenario unit price
divided by the model's official base unit price:

```text
OtherRatio = official_usd_per_m_tokens_for_row / official_base_usd_per_m_tokens
```

The target 10% markup uses:

```text
sale_rmb = official_usd * 7.3 * 1.1
```

Equivalently, because `7.3 * 1.1 = 8.03`:

```text
sale_rmb = official_usd * 8.03
```

## Production DB Values

Read-only production DB check:

| Model | ModelRatio | Base sale RMB/M |
| --- | ---: | ---: |
| Seedance-1.5-pro-海外版 | 19.272 | 38.544 |
| Seedance-2.0-海外版 | 28.105 | 56.210 |
| Seedance-2.0-fast-海外版 | 22.484 | 44.968 |

Seedance 2.0 standard base check:

```text
official_cost = 7.0 USD/M * 7.3 = 51.10 RMB/M
sale = 28.105 * 2 = 56.21 RMB/M
sale / cost = 1.10
```

Seedance 2.0 fast base check:

```text
official_cost = 5.6 USD/M * 7.3 = 40.88 RMB/M
sale = 22.484 * 2 = 44.968 RMB/M
sale / cost = 1.10
```

## Seedance 2.0 Standard Unit Price Audit

| Scenario | Official USD/M | Cost RMB/M | Sale RMB/M | Markup |
| --- | ---: | ---: | ---: | ---: |
| 480p text/image to video | 7.0 | 51.1000 | 56.2100 | 1.10x |
| 720p text/image to video | 7.0 | 51.1000 | 56.2100 | 1.10x |
| 1080p text/image to video | 7.7 | 56.2100 | 61.8310 | 1.10x |
| 4k text/image to video | 4.0 | 29.2000 | 32.1200 | 1.10x |
| 480p video input to video | 4.3 | 31.3900 | 34.5290 | 1.10x |
| 720p video input to video | 4.3 | 31.3900 | 34.5290 | 1.10x |
| 1080p video input to video | 4.7 | 34.3100 | 37.7410 | 1.10x |
| 4k video input to video | 2.4 | 17.5200 | 19.2720 | 1.10x |

## Seedance 2.0 Fast Unit Price Audit

| Scenario | Official USD/M | Cost RMB/M | Sale RMB/M | Markup |
| --- | ---: | ---: | ---: | ---: |
| 480p text/image to video | 5.6 | 40.8800 | 44.9680 | 1.10x |
| 720p text/image to video | 5.6 | 40.8800 | 44.9680 | 1.10x |
| 480p video input to video | 3.3 | 24.0900 | 26.4990 | 1.10x |
| 720p video input to video | 3.3 | 24.0900 | 26.4990 | 1.10x |

## Why Some Rows Look Counterintuitive

`RMB/M tokens` is the unit token price, not the final video price. A higher
resolution can have a lower unit price while still producing a higher final
charge because it consumes more tokens.

For official Seedance 2.0 text/image-to-video examples:

| Resolution | Sale RMB/M | Sale RMB/video | Sale RMB/second |
| --- | ---: | ---: | ---: |
| 480p | 56.2100 | 2.8105 | 0.5621 |
| 720p | 56.2100 | 6.1028 | 1.2045 |
| 1080p | 61.8310 | 15.0161 | 2.9711 |
| 4k | 32.1200 | 31.2367 | 6.2634 |

So 4k has a lower `RMB/M tokens` than 1080p, but a higher per-video and
per-second reference price.

For video-input tasks, official unit token prices are lower than no-video-input
unit prices, but input video contributes to usage. The official per-video range
can therefore be higher than no-video-input, especially at longer input
durations.

## Seedance 1.5 Pro Note

Seedance 1.5 Pro official unit token price varies by audio mode:

- Audio: 2.4 USD/M.
- Silent: 1.2 USD/M.

It does not use a different official USD/M token unit price for 480p, 720p, and
1080p in the published table. Resolution affects the final charge through the
returned `usage.total_tokens`, so the marketplace must not treat the repeated
`RMB/M tokens` value as the per-video price.

The pricing API also exposes official 5-second per-video reference prices and
their per-second equivalents to avoid this ambiguity.

## Conclusion

The production DB prices for `Seedance-2.0-海外版` and
`Seedance-2.0-fast-海外版` are exactly official price plus 10% when using
`USD/CNY = 7.3`.

They should not lose money as long as:

- BytePlus official prices remain unchanged.
- The video task returns `usage.total_tokens` and settlement recalculates from
  that usage.
- The configured `byteplus` group ratio remains `1`.

