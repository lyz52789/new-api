# Seedance 素材库网关设计

## 背景与结论

当前 Seedance 视频适配器可把公网图片 URL、视频 URL 和 `asset://<id>` 字符串透传给 BytePlus，但当前分支没有素材库 API。因此，调用方无法通过 new-api 把自己上传到对象存储的图片或视频登记到 Seedance 素材库，也无法安全地列出和管理自己的素材。

BytePlus 素材库不接收 multipart 文件。`CreateAsset` 接收一个公网可下载的 URL，由 BytePlus 异步拉取并处理素材。调用方必须先把本地文件上传到现有对象存储或 CDN，再调用素材库网关。

## 目标

1. 为普通 new-api Token 提供素材创建、列表、查询、更新和删除接口。
2. 图片与视频均可登记，分别使用 `AssetType: "Image"` 和 `AssetType: "Video"`。
3. 每个 new-api 用户拥有独立的 BytePlus `AIGC` 素材组；客户端不能指定或越过该分组。
4. `CreateAsset` 返回素材 ID；客户端轮询 `GetAsset` 至 `Status == "Active"`，随后在 Seedance 视频请求中使用 `asset://<id>`。
5. BytePlus AK/SK 只保存在系统设置中，不返回给普通用户或管理端浏览器。
6. 数据库迁移同时兼容 SQLite、MySQL 5.7.8+ 和 PostgreSQL 9.6+。

## 非目标

- new-api 不保存上传文件，也不新增对象存储或 multipart 上传能力。
- 不实现多出口、模板化协议转换、素材操作计费或独立限流；这些能力不属于当前用户流程。
- 不公开素材组 CRUD。用户组由网关自动创建和维护。
- 不处理 `LivenessFace` 真人授权组；本需求使用 Seedance `AIGC` 图片/视频参考素材。

## API

所有接口使用 `middleware.TokenAuth()`，路径保持与 BytePlus Action 名称一致：

| 方法与路径 | 用途 |
| --- | --- |
| `POST /doubao/open/CreateAsset` | 用公网 URL 登记图片或视频 |
| `POST /doubao/open/ListAssets` | 列出当前用户素材 |
| `POST /doubao/open/GetAsset` | 查询当前用户素材及处理状态 |
| `POST /doubao/open/UpdateAsset` | 更新当前用户素材元数据 |
| `POST /doubao/open/DeleteAsset` | 删除当前用户素材 |

客户端不能控制 `GroupId` 或 `ProjectName`。服务端始终用当前 Token 的用户 ID 查询绑定，并覆盖请求中的相关字段。

`CreateAsset` 至少校验：

- `URL` 必须为绝对 `http` 或 `https` URL；
- `AssetType` 只允许 `Image` 或 `Video`；
- 空 URL、私有协议和未知类型返回 HTTP 400。

## 组件边界

### 系统配置

`setting/system_setting/volc_asset.go` 定义单一 BytePlus 出口：Access Key、Secret Key、Region、Project Name、Group Type。Secret Key 采用只写更新语义：读取系统选项时返回空字符串，管理员提交空值时保留旧值。

`web/src/components/settings/VolcAssetSetting.jsx` 在当前系统设置页提供上述字段。页面只显示 Access Key、Region、Project Name、Group Type 以及“Secret Key 已配置”状态；Secret Key 输入框永不回填真实值。

### 用户分组绑定

`model/volc_asset_group.go` 保存 `user_id -> group_id`：

```text
volc_asset_user_groups
  id          primary key
  user_id     unique index
  group_id    text
  created_at  unix timestamp
  updated_at  unix timestamp
```

仅使用 GORM `Where`、`First`、`Create` 和事务，避免数据库方言相关 SQL。

首次素材调用时，服务层执行：

1. 查询当前用户绑定；
2. 未找到则调用 BytePlus `CreateAssetGroup`，名称为 `newapi-user-<userID>`；
3. 保存绑定；
4. 若并发创建触发唯一索引冲突，重新读取已存在绑定。

### BytePlus 客户端

`relay/channel/task/doubao/asset_client.go` 负责：

- 使用项目的 `common.Marshal` / `common.Unmarshal`；
- 构造 `Action`、`Version=2024-01-01` 请求；
- 生成 Volcengine HMAC-SHA256 签名；
- 限制上游响应体大小；
- 解析 `ResponseMetadata.Error` 与 `Result`；
- 返回结构化错误，不直接写 Gin 响应。

客户端通过小接口注入 handler，测试时使用假客户端，不发真实网络请求。

### HTTP 服务与归属校验

`relay/channel/task/doubao/asset_handler.go` 负责请求解析、当前用户识别、分组注入及响应转换。

- `ListAssets` 强制 `Filter.GroupIds = [用户分组]`。
- `CreateAsset` 强制 `GroupId = 用户分组`。
- `GetAsset`、`UpdateAsset`、`DeleteAsset` 先调用 `GetAsset`，验证返回的 `GroupId` 等于用户分组；不相等时返回 404。
- 上游业务错误保留 code/message；网络错误转换为 502；配置缺失返回 503。

### 视频生成引用

现有 `requestPayload.Content` 的 `image_url.url` 和 `video_url.url` 不限制 URL scheme，因此可直接透传 `asset://<id>`。新增适配器测试证明图片和视频素材引用不会被修改或丢弃，并更新调用文档，明确必须等待素材 `Active`。

## 数据流

```text
本地文件
  -> 现有对象存储/CDN
  -> POST CreateAsset {URL, AssetType}
  -> 自动创建/复用用户 AIGC 分组
  -> BytePlus 返回 asset id
  -> POST GetAsset 轮询至 Active
  -> POST /v1/video/generations，content 中使用 asset://<id>
  -> Seedance 异步生成视频
```

## 安全与错误处理

- 素材接口只接受 Token 鉴权，用户 ID 取自鉴权上下文，不接受客户端声明。
- 任何素材详情、修改或删除前都做归属校验。
- AK/SK 不写日志；响应错误不包含签名头或凭证。
- 上游响应最大 10 MiB，避免无界读取。
- URL 校验只允许公网可表达的 HTTP(S) 形式；真正的网络可达性和素材格式由 BytePlus 异步处理结果决定。

## 测试与验收

1. 配置测试：序列化时 Secret Key 脱敏，空 Secret Key 更新保留旧值。
2. 模型测试：保存/读取用户绑定，并发唯一冲突后可读取已存在绑定。
3. 签名测试：固定时间下 canonical request 和 Authorization 稳定。
4. Handler 测试：用户分组被强制注入；非法 URL/类型被拒绝；跨组 Get/Update/Delete 返回 404。
5. 适配器测试：`asset://` 图片和视频引用完整进入上游 `content`。
6. 路由测试或路由表审计：五个端点均受 `TokenAuth` 保护。
7. 运行相关 Go 测试、全量 `go test ./...`，并对前端配置改动运行 Bun 构建。
