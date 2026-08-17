# Bug 是什么

迁移后的源区和目标区统计发生串区，未解决与严重阈值计算错误，两份统计证明还会写到同一文件。

# 如何触发

运行 `go test ./internal/... -run '^TestStatsRecalculationPipeline' -count=20`。

# 错误信息

测试会报告源区与目标区统计相同、仓储使用错误范围，或两个统计附件路径均为 `stats-a.csv`。
