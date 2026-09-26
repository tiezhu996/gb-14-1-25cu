# database

本目录说明数据库脚本与初始化方式。

CodeLearn 使用 MongoDB（mongo-driver），无 SQL 迁移脚本：

- 集合结构与索引：由后端 `backend/internal/model/` 的 BSON 模型与 `cmd/server/main.go` 的 `ensureIndexes` 自动创建。
- 演示数据种子：`backend/internal/seed/seed.go`（幂等，`SEED_ADMIN=true` 时启动自动执行）。

演示账号：
- 管理员：`admin / admin123`
- 学生：`student / student123`

端口：MongoDB 宿主机 44008（任务清单分配），容器内 27017。
