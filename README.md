# 办公环境反馈纠错与归档平台

该离线平台纠正区域合并、拆分或录入错误造成的反馈归属问题。迁移申请支持候选预览、区域兼容校验、两位不同人员复核、批次幂等执行、区域统计重算、限时撤销及不可覆盖的处理时间线；反馈原始提交时间始终保留。

从 `.env.example` 配置监听地址、PostgreSQL 和撤销窗口，执行 `migrations/001_feedback.sql` 与 `scripts/seed.sql` 后运行 `go run ./cmd/server`。`docker compose up --build` 提供本地容器方案。演示区域为 north-2 与 north-3。执行 API 是 `POST /api/v1/migrations/execute`，所有业务错误返回 code、message、request_id。

领域层定义迁移、复核、撤销和统计语义；应用层把幂等键与原子替换绑定；仓储层含并发安全的离线实现及 pgx 连接；敏感备注在服务层脱敏，附件路径由平台适配器限制。Vue 页面覆盖区域、反馈、迁移、审批、统计差异和历史。

已实际验证 `go test ./...`、`go test -race ./...` 与 `go vet ./...`。前端执行 `cd web && npm ci && npm test && npm run build`。健康探针为 `/healthz` 和 `/readyz`。
