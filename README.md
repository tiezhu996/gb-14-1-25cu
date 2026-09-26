# CodeLearn 在线编程学习平台

面向编程初学者和进阶者的一站式在线编程学习平台：系统化学习路径、课程图文教程（Markdown）、浏览器内嵌在线 IDE（Monaco）、ACM 风格自动评测、排行榜与成就系统、题目讨论社区与个人学习仪表盘。用户无需本地搭建开发环境即可编写并运行 Python / JavaScript / Java 代码。

## 快速启动（Docker Compose，首选）

```bash
docker compose up -d --build
```

启动后访问：

| 入口 | 地址 |
| --- | --- |
| 前端页面 | http://localhost:8010 |
| 后端健康检查 | http://localhost:3010/healthz |
| API 前缀 | http://localhost:3010/api/v1 |
| OpenAPI 文档 | http://localhost:3010/api/v1/openapi.yaml（仓库内 `backend/api/openapi.yaml`） |
| MongoDB | localhost:44008（codelearn_user / codelearn_pwd） |
| Redis | localhost:46308 |

演示账号：

- 管理员：`admin / admin123`
- 学生：`student / student123`

关闭并清理（含数据卷）：

```bash
docker compose down -v --remove-orphans
```

## 项目主要功能

1. **学习路径管理**：结构化课程（如 "Python 入门→Web 开发→数据分析"），章节 + 知识点 + 练习，学习进度与热力图跟踪
2. **课程内容管理**：管理员创建课程（Markdown 图文教程），章节拆解、代码片段展示
3. **在线 IDE**：Monaco Editor 内嵌，支持 Python / JavaScript / Java 语法高亮与自动补全；服务端沙箱评测，单用例超时 10 秒自动终止
4. **编程题目评测**：ACM 风格，多测试用例逐条运行，判定 通过 / 部分通过 / 运行错误 / 超时，展示输入、期望输出与实际输出对比
5. **排行榜与成就系统**：解题数量 × 难度加权积分（日榜 / 周榜 / 总榜）；成就徽章（连续签到 7 天、完成 10/100 题、首次通过困难题等）
6. **讨论社区**：每道题专属讨论区，支持 Markdown 与代码块、点赞、按最佳答案排序
7. **个人学习仪表盘**：累计学习时长、完成课程数、解题总数、各语言解题分布饼图、近 90 天每日学习热力图

## 技术栈

| 层 | 技术栈 |
| --- | --- |
| 前端 | React + TypeScript，使用 Tailwind CSS 样式框架，Vite 构建工具，Monaco Editor，zustand，react-router-dom |
| 后端 | Go 1.22 + Gin + mongo-driver |
| 数据库 | MongoDB 7.0 |
| 缓存 | Redis 7 |
| 认证 | JWT + RBAC（student / admin） |
| 日志 | log/slog 结构化日志 |
| 参数校验 | github.com/go-playground/validator/v10 |
| 接口文档 | OpenAPI 3.0（`backend/api/openapi.yaml`）+ README API 清单 |

## 项目目录结构

```
.
├── docker-compose.yml          # 一键编排：frontend + backend + db + redis
├── .env / .env.example         # 环境变量（端口按任务清单分配）
├── README.md
├── database/                   # 数据库说明
├── backend/
│   ├── cmd/server/main.go      # 装配依赖、启动服务、索引与种子
│   ├── internal/
│   │   ├── config/             # 环境变量配置
│   │   ├── database/           # MongoDB / Redis 连接
│   │   ├── model/              # 每个实体一个文件（User/Course/Problem/Submission/Discussion/Achievement/AuditLog/UserStat）
│   │   ├── dto/                # 每个实体一个 DTO 文件
│   │   ├── repository/         # 每个实体一个 repository 文件（含哨兵错误）
│   │   ├── service/            # 每个实体一个 service 文件（含评测沙箱）
│   │   ├── handler/            # 每个实体一个 handler 文件
│   │   ├── router/             # 每个实体一个路由注册文件
│   │   ├── middleware/         # auth/rbac/request_id/error_handler/audit/ratelimit/logger/cors
│   │   ├── constants/          # 错误码、日志模板、文案、业务枚举
│   │   └── util/               # logger/jwt/formatters/app_error/response 等
│   ├── pkg/strutil/            # 无业务依赖的通用工具
│   ├── migrations/             # MongoDB 集合与索引说明
│   ├── api/openapi.yaml        # OpenAPI 3.0 定义
│   ├── Dockerfile
│   ├── go.mod / go.sum
│   └── .dockerignore
└── frontend/
    ├── src/
    │   ├── api/                # 每个实体一个 API 文件（auth/user/course/problem/submission/discussion/achievement/leaderboard/dashboard/audit）
    │   ├── components/         # 共享组件（Layout/StatusBadge/EmptyState/DataTable/ConfirmDialog/Pagination/CodeBlock/LanguageSelect/StatCard/PieChart/Heatmap/ProtectedRoute）
    │   ├── pages/              # 按模块拆分的页面（dashboard/courses/problems/submissions/leaderboard/achievements/admin/login/register）
    │   ├── stores/             # zustand，按实体拆分（auth/course/problem/submission/discussion/achievement/leaderboard/dashboard/audit）
    │   ├── hooks/              # useAuth / usePagination
    │   ├── utils/              # request.ts 拦截器 / format.ts
    │   ├── constants/          # 与后端一一对应的枚举
    │   └── types/              # 类型定义
    ├── Dockerfile
    ├── nginx.conf              # SPA 路由 + /api 反向代理
    ├── vite.config.ts
    ├── tailwind.config.js
    └── package.json
```

## 环境变量说明

| 变量 | 默认值 | 说明 |
| --- | --- | --- |
| COMPOSE_PROJECT_NAME | codelearn | Compose 项目名，容器名前缀 |
| DB_NAME | codelearn_db | MongoDB 库名 |
| DB_USER | codelearn_user | MongoDB 账号 |
| DB_PASSWORD | codelearn_pwd | MongoDB 密码 |
| DB_PORT | 44008 | MongoDB 宿主机映射端口 |
| JWT_SECRET | change_me_to_a_long_random_string | JWT 签名密钥（生产必改） |
| FRONTEND_PORT | 8010 | 前端宿主机端口 |
| BACKEND_PORT | 3010 | 后端宿主机端口（容器内固定 8080） |
| REDIS_ADDR | redis:6379 | Redis 容器内地址 |
| REDIS_PASSWORD | (空) | Redis 密码 |
| REDIS_PORT | 46308 | Redis 宿主机映射端口 |

## 后端本地开发（备选）

```bash
cd backend && go mod tidy && go run ./cmd/server
```

构建命令：`go build ./...`

需要本地 MongoDB/Redis，通过环境变量覆盖默认连接：

```bash
MONGO_URI="mongodb://codelearn_user:codelearn_pwd@localhost:44008/codelearn_db?authSource=admin" \
REDIS_ADDR="localhost:46308" JWT_SECRET="local-secret" go run ./cmd/server
```

前端本地开发：

```bash
cd frontend && npm install && npm run dev
```

## API 调用示例（curl）

统一响应：`{ "code": 0, "message": "ok", "data": ... }`，`code=0` 表示成功。

注册：

```bash
curl -sS -X POST http://localhost:3010/api/v1/auth/register \
  -H 'Content-Type: application/json' \
  -d '{"username":"alice","email":"alice@example.com","password":"secret123","nickname":"Alice"}'
```

登录并保存 JWT：

```bash
TOKEN=$(curl -sS -X POST http://localhost:3010/api/v1/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"student","password":"student123"}' | python3 -c 'import sys,json;print(json.load(sys.stdin)["data"]["token"])')
```

每日签到（含 JWT 请求头）：

```bash
curl -sS -X POST http://localhost:3010/api/v1/auth/sign-in -H "Authorization: Bearer $TOKEN"
```

课程列表 / 题目列表 / 排行榜 / 仪表盘：

```bash
curl -sS "http://localhost:3010/api/v1/courses?page=1&page_size=10" -H "Authorization: Bearer $TOKEN"
curl -sS "http://localhost:3010/api/v1/problems?page=1&page_size=10" -H "Authorization: Bearer $TOKEN"
curl -sS "http://localhost:3010/api/v1/leaderboard?period=total" -H "Authorization: Bearer $TOKEN"
curl -sS "http://localhost:3010/api/v1/dashboard/me" -H "Authorization: Bearer $TOKEN"
```

提交代码评测：

```bash
curl -sS -X POST http://localhost:3010/api/v1/problems/<PROBLEM_ID>/submit \
  -H "Authorization: Bearer $TOKEN" -H 'Content-Type: application/json' \
  -d '{"language":"python","code":"a, b = map(int, input().split())\nprint(a + b)"}'
```

管理员创建题目：

```bash
ADMIN_TOKEN=$(curl -sS -X POST http://localhost:3010/api/v1/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"admin","password":"admin123"}' | python3 -c 'import sys,json;print(json.load(sys.stdin)["data"]["token"])')
curl -sS -X POST http://localhost:3010/api/v1/problems \
  -H "Authorization: Bearer $ADMIN_TOKEN" -H 'Content-Type: application/json' \
  -d '{"title":"A+B","description":"输出两数之和","difficulty":"easy","languages":["python","javascript","java"],"tags":["入门"],"test_cases":[{"input":"1 2","output":"3"}],"status":"published"}'
```

## API 清单（后端 ≥ 6 个接口，含复用关系标注）

| 方法 | 路径 | 说明 | 权限 | 复用 |
| --- | --- | --- | --- | --- |
| GET | /healthz | 健康检查 | 公开 | - |
| POST | /api/v1/auth/register | 注册 | 公开 | UserService.Register |
| POST | /api/v1/auth/login | 登录 | 公开 | UserService.Login |
| POST | /api/v1/auth/sign-in | 每日签到 | 登录 | UserService.SignIn |
| GET | /api/v1/users/me | 我的资料 | 登录 | UserService.GetProfile |
| PUT | /api/v1/users/me | 更新资料 | 登录 | UserService.UpdateProfile |
| GET | /api/v1/users | 用户列表 | 管理员 | UserService.List |
| PUT | /api/v1/users/:id/status | 禁用/启用 | 管理员 | UserService.UpdateStatus |
| PUT | /api/v1/users/:id/role | 调整角色 | 管理员 | UserService.UpdateRole |
| GET | /api/v1/courses | 课程列表 | 登录 | CourseService.List |
| GET | /api/v1/courses/:id | 课程详情 | 登录 | CourseService.Get |
| POST | /api/v1/courses | 创建课程 | 管理员 | CourseService.Create |
| PUT | /api/v1/courses/:id | 更新课程 | 管理员 | CourseService.Update |
| PUT | /api/v1/courses/:id/status | 状态流转 | 管理员 | CourseService.UpdateStatus |
| DELETE | /api/v1/courses/:id | 删除课程 | 管理员 | CourseService.Delete |
| POST | /api/v1/courses/:id/learn | 记录学习时长 | 登录 | UserStatRepository.AddLearning |
| POST | /api/v1/courses/:id/complete | 标记完成 | 登录 | UserStatRepository.CompleteCourse |
| GET | /api/v1/problems | 题目列表 | 登录 | ProblemService.List |
| GET | /api/v1/problems/:id | 题目详情 | 登录 | ProblemService.Get |
| POST | /api/v1/problems | 创建题目 | 管理员 | ProblemService.Create |
| PUT | /api/v1/problems/:id | 更新题目 | 管理员 | ProblemService.Update |
| PUT | /api/v1/problems/:id/status | 状态流转 | 管理员 | ProblemService.UpdateStatus |
| DELETE | /api/v1/problems/:id | 删除题目 | 管理员 | ProblemService.Delete |
| POST | /api/v1/problems/:id/submit | 提交评测 | 登录 | SubmissionService.Submit + JudgeService.Judge |
| GET | /api/v1/submissions | 提交记录 | 登录/管理员 | SubmissionService.List |
| GET | /api/v1/submissions/:id | 提交详情 | 本人/管理员 | SubmissionService.Get |
| GET | /api/v1/problems/:id/discussions | 讨论区 | 登录 | DiscussionService.ListByProblem |
| POST | /api/v1/problems/:id/discussions | 发布讨论 | 登录 | DiscussionService.Create |
| POST | /api/v1/discussions/:id/vote | 点赞/取消 | 登录 | DiscussionService.Vote |
| PUT | /api/v1/problems/:problem_id/discussions/:id/best | 最佳答案 | 作者/管理员 | DiscussionService.MarkBest |
| PUT | /api/v1/discussions/:id/hide | 隐藏 | 本人/管理员 | DiscussionService.Hide |
| GET | /api/v1/leaderboard | 排行榜 | 登录 | LeaderboardService.Get + SubmissionRepository.AggregatePoints（复用聚合） |
| GET | /api/v1/achievements | 全部徽章 | 登录 | AchievementService.ListAll |
| GET | /api/v1/achievements/me | 我的徽章 | 登录 | AchievementService.ListMine |
| GET | /api/v1/dashboard/me | 学习仪表盘 | 登录 | DashboardService.Get |
| GET | /api/v1/audits | 审计日志 | 管理员 | AuditService.List |

> 复用标注：`POST /courses/:id/learn` 与 `POST /courses/:id/complete` 复用 `UserStatRepository` 原子累加；`/leaderboard` 与 `/dashboard` 复用 `SubmissionRepository`/`UserStatRepository` 统计；`Vote` 的点赞幂等复用 `DiscussionRepository.HasVoted`。

## Docker 部署说明

- 端口映射：前端 `${FRONTEND_PORT:-8010}:80`，后端 `${BACKEND_PORT:-3010}:8080`，数据库 `${DB_PORT:-44008}:27017`，Redis `${REDIS_PORT:-46308}:6379`
- 数据卷：`dbdata`（MongoDB）、`redisdata`（Redis）持久化，容器重建不丢数据
- 健康检查：db/redis/backend/frontend 均配置 healthcheck，后端依赖 db/redis `service_healthy`
- 后端容器内置 python3 / nodejs / openjdk17，支撑三种语言在线评测
- 常见问题：
  - 端口冲突：修改 `.env` 中 `FRONTEND_PORT/BACKEND_PORT/DB_PORT/REDIS_PORT` 后 `docker compose up -d`
  - 数据重置：`docker compose down -v` 后重新 `up -d` 会重建空库并自动种子演示数据
  - 中文目录名：Compose 使用相对路径构建，命名卷持久化，任意目录名下均可正常启动

## 横切关注点实现（触达文件层）

1. **JWT 认证 + RBAC 权限**
   - 数据库角色字段：`backend/internal/model/user.go`（`Role`）
   - 认证中间件：`backend/internal/middleware/auth.go`
   - RBAC 中间件：`backend/internal/middleware/rbac.go`
   - JWT 工具：`backend/internal/util/jwt.go`
   - 前端路由守卫与按钮显隐：`frontend/src/components/ProtectedRoute.tsx`、`frontend/src/components/Layout.tsx`（管理员菜单）、`frontend/src/utils/request.ts`（401 拦截）

2. **操作审计日志**
   - 审计集合：`backend/internal/model/audit_log.go`
   - 审计中间件：`backend/internal/middleware/audit.go`
   - service 埋点：`backend/internal/service/audit_service.go`（`AuditService.Create` 被中间件与业务复用）
   - 前端审计页：`frontend/src/pages/admin/AdminAudits.tsx`、`frontend/src/stores/auditStore.ts`

3. **全局错误处理与请求追踪**
   - `backend/internal/middleware/request_id.go`、`backend/internal/middleware/error_handler.go`
   - `backend/internal/util/app_error.go`、`backend/internal/constants/error_codes.go`
   - 前端拦截器：`frontend/src/utils/request.ts`

## 枚举出现位置清单

### 1. 角色枚举（student / admin）

| 位置 | 文件 |
| --- | --- |
| 后端常量 | `backend/internal/constants/roles.go` |
| 模型 | `backend/internal/model/user.go` |
| DTO | `backend/internal/dto/user_dto.go` |
| Service 状态机 | `backend/internal/service/user_service.go`（注册默认角色、UpdateRole 校验） |
| Handler 校验 | `backend/internal/handler/user_handler.go`（UpdateRole oneof） |
| 中间件 | `backend/internal/middleware/rbac.go`、`backend/internal/router/router.go`（RequireRole） |
| 日志模板 | `backend/internal/constants/log_templates.go` |
| 错误码 | `backend/internal/constants/error_codes.go`（CodeUserRoleChange） |
| 格式化 | `backend/internal/util/formatters.go`（FormatRoleText） |
| 前端常量 | `frontend/src/constants/index.ts`（ROLES/ROLE_LABELS） |
| 前端徽标 | `frontend/src/components/StatusBadge.tsx`（role kind） |
| 前端守卫 | `frontend/src/components/ProtectedRoute.tsx`（AdminRoute） |
| 前端页面 | `frontend/src/pages/admin/AdminUsers.tsx`（角色调整） |

### 2. 题目难度枚举（easy / medium / hard）

| 位置 | 文件 |
| --- | --- |
| 后端常量 | `backend/internal/constants/difficulty.go`（含难度积分映射） |
| 模型 | `backend/internal/model/problem.go`、`backend/internal/model/course.go` |
| DTO | `backend/internal/dto/problem_dto.go`、`backend/internal/dto/course_dto.go`（oneof 校验） |
| Service | `backend/internal/service/problem_service.go`、`course_service.go`（Points 加权、成就判定 first_hard） |
| Handler 校验 | `backend/internal/handler/problem_handler.go`、`course_handler.go` |
| 日志模板 | `backend/internal/constants/log_templates.go`（problem.difficulty 相关） |
| 错误码 | `backend/internal/constants/error_codes.go` |
| 格式化 | `backend/internal/util/formatters.go`（FormatDifficultyText/Class） |
| 前端常量 | `frontend/src/constants/index.ts`（DIFFICULTY/LABELS/CLASSES） |
| 前端筛选 | `frontend/src/pages/problems/Problems.tsx`、`frontend/src/pages/courses/Courses.tsx` |
| 前端徽标 | `frontend/src/components/StatusBadge.tsx`（difficulty kind） |

### 3. 内容状态机（draft / published / archived，课程与题目共用）

| 位置 | 文件 |
| --- | --- |
| 后端常量 | `backend/internal/constants/status.go` |
| 模型 | `backend/internal/model/course.go`、`problem.go` |
| DTO | `backend/internal/dto/course_dto.go`、`problem_dto.go`（oneof 校验） |
| Service 状态机 | `backend/internal/service/course_service.go`（canTransition）、`problem_service.go` |
| Handler | `backend/internal/handler/course_handler.go`、`problem_handler.go` |
| 日志模板 | `backend/internal/constants/log_templates.go`（course.published/archived、problem.published） |
| 错误码 | `backend/internal/constants/error_codes.go`（CodeCourseLocked/CodeProblemLocked） |
| 格式化 | `backend/internal/util/formatters.go`（FormatContentStatusText） |
| 前端常量 | `frontend/src/constants/index.ts`（CONTENT_STATUS/LABELS） |
| 前端按钮显隐 | `frontend/src/pages/admin/AdminCourses.tsx`、`AdminProblems.tsx`（发布/归档按钮） |
| 前端徽标 | `frontend/src/components/StatusBadge.tsx`（content kind） |

### 4. 提交评测状态机（pending / judging / accepted / partial / runtime_error / timeout）

| 位置 | 文件 |
| --- | --- |
| 后端常量 | `backend/internal/constants/submission.go` |
| 模型 | `backend/internal/model/submission.go` |
| DTO | `backend/internal/dto/submission_dto.go` |
| Service 状态机 | `backend/internal/service/judge_service.go`（评测判定）、`submission_service.go` |
| Handler | `backend/internal/handler/submission_handler.go` |
| 日志模板 | `backend/internal/constants/log_templates.go`（submission.accepted/partial/error/timeout） |
| 错误码 | `backend/internal/constants/error_codes.go`（CodeJudgeTimeout/CodeSubmissionNotFound） |
| 格式化 | `backend/internal/util/formatters.go`（FormatSubmissionStatusText/Class） |
| 前端常量 | `frontend/src/constants/index.ts`（SUBMISSION_STATUS/LABELS/CLASSES） |
| 前端结果页 | `frontend/src/pages/problems/ProblemDetail.tsx`（结果表格）、`frontend/src/pages/submissions/Submissions.tsx` |
| 前端徽标 | `frontend/src/components/StatusBadge.tsx`（submission kind） |

### 5. 评测语言枚举（python / javascript / java）

| 位置 | 文件 |
| --- | --- |
| 后端常量 | `backend/internal/constants/submission.go`（SupportedLanguages/ValidLanguage） |
| 模型 | `backend/internal/model/problem.go`（Languages）、`submission.go` |
| DTO | `backend/internal/dto/problem_dto.go`（dive oneof）、`submission_dto.go`（oneof） |
| Service | `backend/internal/service/judge_service.go`（runCode 分支） |
| 错误码 | `backend/internal/constants/error_codes.go`（CodeJudgeLanguage） |
| 格式化 | `backend/internal/util/formatters.go`（FormatLanguageText） |
| 前端常量 | `frontend/src/constants/index.ts`（LANGUAGES/LANGUAGE_LABELS） |
| 前端 IDE | `frontend/src/pages/problems/ProblemDetail.tsx`、`frontend/src/components/LanguageSelect.tsx` |

## 文件结构强制清单

- 后端 `cmd/server/main.go` 只负责装配与启动；业务代码全部在 `internal/`
- `model / dto / repository / service / handler / router / constants` 每实体独立文件
- `handler → service → repository → model` 单向依赖、构造器注入
- 前端 `api / stores / pages / components / hooks / utils / constants / types` 按实体与功能拆分
- 严禁合并职责到单一文件：不允许把多个实体的 model/repository/service/handler 写进同一个文件，也不允许前端把所有页面写在 App.tsx 中

## 屎山代码设计要求（低内聚、高耦合、牵一发动全身）

1. **日志模块单独管理但全栈引用**：`backend/internal/util/logger.go`（slog），handler/service/middleware 全部引用；日志模板集中在 `backend/internal/constants/log_templates.go`（50+ 条）
2. **异常信息分散且层层透传**：错误码集中在 `error_codes.go`，service/handler 手动拼接 message，handler 再次包装错误；消息文案集中在 `messages.go`
3. **常量/工具类多处耦合**：`util/formatters.go` 集中日期/状态/类型文本格式化；`constants/messages.go` 集中接口/日志/错误文案
4. **状态机跨多处定义**：课程/题目/提交状态机同时存在于 service 状态机、前端按钮显隐、日志模板、错误码、formatters、前端 constants
5. **枚举多处重复定义**：角色/难度/内容状态/提交状态/语言/成就等枚举同时存在于后端 constants、DTO、模型、日志模板、错误码、formatters 与前端 constants
6. **牵一发动全身验证标准**：若给核心实体新增字段/状态，需同步修改模型、DTO、constants、service、repository、handler、formatters、日志模板、错误码、前端类型与页面等 ≥ 10 个文件

## License

MIT License。仅用于教学与演示，请勿用于生产环境。
