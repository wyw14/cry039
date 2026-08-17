# Bug 是什么

同一复核人通过大小写和空白变体可被计为两名复核人，从而绕过双人复核执行迁移。

# 如何触发

运行 `go test ./internal/... -run '^TestReviewerQuorumPipeline' -count=20`。

# 错误信息

测试会报告等价复核身份被当作不同人员、应用执行未返回 incomplete review，或 HTTP 返回 200。
