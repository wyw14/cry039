# Bug 是什么

反馈迁移会重写原始提交时间并截断既有处理时间线，查询结果还可能反向污染仓储中的审计记录。

# 如何触发

运行 `go test ./internal/... -run '^TestMigrationHistoryPipeline' -count=20`。

# 错误信息

测试会报告 `submission audit was rewritten`、`application lost feedback history` 或 `repository truncated history`。
