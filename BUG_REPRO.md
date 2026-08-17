# Bug 是什么

撤销窗口按时区小时字段错误计算且被应用层缩短，撤销范围和证明还包含目标区原有的无关反馈。

# 如何触发

运行 `go test ./internal/... -run '^TestUndoWindowPipeline' -count=20`。

# 错误信息

20 分钟内撤销会报告 `migration undo window expired`，或测试报告撤销证明包含无关反馈。
