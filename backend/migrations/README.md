# migrations

CodeLearn 使用 MongoDB，无传统 SQL 迁移脚本。集合结构由 `internal/model/` 中的 BSON 模型驱动，
启动时由 `cmd/server/main.go` 的 `ensureIndexes` 自动创建唯一索引与查询索引。

集合清单：
- users（username/email 唯一索引）
- courses（title 唯一索引）
- problems（title 唯一索引）
- submissions（user_id/problem_id/status 查询索引）
- discussions（problem_id 查询索引）
- achievements（成就定义）
- user_achievements（user_id+code 唯一索引）
- audit_logs（user_id/entity 查询索引）
- user_stats（user_id 唯一索引）

演示数据（幂等种子，见 `internal/seed/seed.go`）：
- 管理员：admin / admin123
- 学生：student / student123
- 2 门课程、4 道题目、5 个成就定义
