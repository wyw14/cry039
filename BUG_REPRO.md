# Bug 是什么

相同幂等键和摘要的迁移重试会再次执行、追加重复反馈并错误增加接口迁移数量。

# 如何触发

运行 `go test ./internal/... -run '^TestIdempotentMigrationPipeline' -count=20`。

# 错误信息

测试会报告迁移版本或时间线再次增长、仓储追加重复行，或接口返回 `migrated: 2`。
