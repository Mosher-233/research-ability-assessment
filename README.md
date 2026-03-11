# Research Ability Assessment

一个用于研究能力评估的全栈应用，使用 Go + Vue + MySQL + Neo4j + Deepseek LLM。

## 技术栈

### 后端
- **语言**: Go 1.20+
- **Web 框架**: Gin
- **ORM**: GORM
- **关系型数据库**: MySQL 8.0
- **图数据库**: Neo4j 5.23
- **LLM**: Deepseek API

### 前端
- **框架**: Vue 3
- **构建工具**: Vite
- **UI 组件**: Element Plus
- **路由**: Vue Router
- **HTTP 客户端**: Axios

### 部署
- **容器化**: Docker + Docker Compose

## 功能特性

- 🔐 用户认证（教师/学生角色）
- 📋 任务管理与分配
- 📝 证据收集与管理
- 🤖 AI 能力推理与评估
- 📊 报告生成与可视化
- 📈 知识图谱存储

## 快速开始

### 前置要求

- Docker Desktop
- Go 1.20+
- Node.js 16+

### 1. 克隆项目

```bash
git clone <repository-url>
cd research-ability-assessment
```

### 2. 配置环境变量

复制 `.env.example` 为 `.env` 并填入你的配置：

```bash
cp .env.example .env
```

编辑 `.env` 文件，填入你的 Deepseek API Key：

```env
DEEPSEEK_API_KEY=your_actual_api_key_here
```

### 3. 启动数据库

```bash
docker-compose up -d
```

这将启动：
- MySQL (端口 3306)
- Neo4j (端口 7474, 7687)

### 4. 安装依赖并启动后端

```bash
go mod tidy
go run cmd/server/main.go
```

后端服务将在 `http://localhost:8080` 启动。

### 5. 安装依赖并启动前端

```bash
cd frontend
npm install
npm run dev
```

前端服务将在 `http://localhost:3000` 启动。

## 项目结构

```
research-ability-assessment/
├── cmd/
│   └── server/
│       └── main.go              # 后端入口文件
├── configs/
│   ├── config.dev.yaml          # 开发环境配置
│   ├── config.mysql.yaml        # MySQL 配置示例
│   └── config.supabase.yaml     # Supabase 配置示例
├── docs/                        # 项目文档
├── frontend/                    # 前端项目
│   ├── src/
│   │   ├── api/                # API 调用
│   │   ├── router/             # 路由配置
│   │   ├── types/              # TypeScript 类型
│   │   ├── utils/              # 工具函数
│   │   └── views/              # Vue 页面
│   └── package.json
├── internal/
│   ├── agent/                  # AI Agent 实现
│   ├── config/                 # 配置管理
│   ├── handler/                # HTTP 处理器
│   ├── llm/                    # LLM 客户端
│   ├── middleware/             # Gin 中间件
│   ├── models/                 # 数据模型
│   ├── repository/             # 数据访问层
│   └── service/                # 业务逻辑层
├── pkg/
│   ├── logger/                 # 日志工具
│   └── utils/                  # 通用工具
├── .env.example                # 环境变量示例
├── .gitignore                  # Git 忽略文件
├── docker-compose.yml          # Docker Compose 配置
├── go.mod                      # Go 依赖管理
├── go.sum                      # Go 依赖校验
└── README.md                   # 项目说明
```

## API 文档

详细的 API 测试指南请参考 [Postman 测试指南](docs/postman_test_guide.md)。

### 主要 API 端点

#### 认证
- `POST /api/v1/auth/register` - 用户注册
- `POST /api/v1/auth/login` - 用户登录
- `GET /api/v1/user/info` - 获取用户信息

#### 任务管理
- `POST /api/v1/tasks` - 创建任务
- `GET /api/v1/tasks` - 获取任务列表
- `GET /api/v1/tasks/:task_id` - 获取任务详情
- `POST /api/v1/tasks/:task_id/assign` - 分配任务
- `GET /api/v1/tasks/:task_id/students` - 获取任务学生
- `GET /api/v1/tasks/students/list` - 获取学生列表

#### 证据管理
- `POST /api/v1/evidences` - 创建证据
- `GET /api/v1/evidences/:evidence_id` - 获取证据详情
- `GET /api/v1/evidences/student-task/:student_task_id` - 获取证据列表

#### 结果管理
- `GET /api/v1/results/:result_id` - 获取推理结果
- `GET /api/v1/results/task/:task_id` - 获取任务结果
- `GET /api/v1/results/report/student` - 生成学生报告
- `GET /api/v1/results/report/task/:task_id` - 生成任务报告

## 配置说明

### LLM 配置

项目支持 Deepseek API，在 `configs/config.dev.yaml` 中配置：

```yaml
llm:
  provider: deepseek
  api_key: ${DEEPSEEK_API_KEY}
  base_url: https://api.deepseek.com
  model: deepseek-chat
  max_tokens: 1000
  temperature: 0.7
```

### 数据库配置

使用 MySQL 作为关系型数据库：

```yaml
database:
  type: mysql
  host: localhost
  port: 3306
  user: mysqluser
  password: mysqlpassword
  dbname: research_assessment
  sslmode: disable
```

### Neo4j 配置

```yaml
neo4j:
  uri: bolt://localhost:7687
  username: neo4j
  password: password123
```

## 文档

- [数据库选择指南](docs/database_options.md)
- [数据库迁移指南](docs/database_migration_guide.md)
- [Postman 测试指南](docs/postman_test_guide.md)
- [开发指南](docs/development_guide.md)
- [部署指南](docs/deployment_guide.md)
- [用户手册](docs/user_manual.md)

## 常见问题

### Q: 如何切换数据库？

A: 项目支持 PostgreSQL、MySQL 和 Supabase。详细说明请参考 [数据库选择指南](docs/database_options.md)。

### Q: 如何重置测试数据？

A: 可以重启数据库容器：

```bash
docker-compose down -v
docker-compose up -d
```

### Q: LLM API Key 如何保密？

A: API Key 通过 `.env` 文件管理，该文件已添加到 `.gitignore`，不会被提交到 Git。

## 贡献

欢迎提交 Issue 和 Pull Request！

## 许可证

MIT License

## 联系方式

如有问题，请提交 Issue。
