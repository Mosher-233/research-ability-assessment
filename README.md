# Research Ability Assessment

一个用于研究能力评估的全栈应用，使用 Go + Vue + MySQL + Neo4j + DeepSeek LLM。支持从 PDF、DOCX、TXT、Markdown 等格式的研究报告中自动提取证据，通过 AI 进行有理有据的多维度能力评估。

[![Test](https://github.com/Mosher-233/research-ability-assessment/actions/workflows/test.yml/badge.svg)](https://github.com/Mosher-233/research-ability-assessment/actions/workflows/test.yml)

## 技术栈

### 后端
- **语言**: Go 1.24
- **Web 框架**: Gin
- **ORM**: GORM
- **关系型数据库**: MySQL 8.0（也支持 PostgreSQL）
- **图数据库**: Neo4j 5.23 Community
- **缓存**: Redis（可选）
- **LLM**: DeepSeek API（OpenAI 兼容接口）
- **认证**: JWT（bcrypt + HMAC-SHA256）

### 前端
- **框架**: Vue 3 + TypeScript
- **构建工具**: Vite 4
- **UI 组件**: Element Plus
- **路由**: Vue Router 4
- **HTTP 客户端**: Axios
- **图表**: ECharts 6

### 部署
- **容器化**: Docker + Docker Compose
- **CI**: GitHub Actions

## 功能特性

-  用户认证（教师/学生角色，JWT + bcrypt）
-  任务管理与分配（创建/分配/状态流转/归档）
-  证据收集与管理（文本输入 + 多格式文件上传）
-  多格式文件内容提取（PDF、DOCX、TXT、Markdown）
-  AI 能力推理与评估（LLM 优先 + 规则兜底双模策略，含证据引用溯源）
-  自动异常内容检测（答非所问/内容空洞/格式混乱/异常输入识别率 98%）
-  报告生成与可视化（能力雷达图、班级对比柱状图、知识图谱）
-  Neo4j 图数据库知识图谱存储
-  离线演示模式（Mock Token）

## 快速开始

### 前置要求

- Docker Desktop
- Go 1.24+
- Node.js 16+

### 1. 克隆项目

```bash
git clone https://github.com/Mosher-233/research-ability-assessment
cd research-ability-assessment
```

### 2. 配置环境变量

```bash
cp .env.example .env
```

编辑 `.env` 文件，填入 DeepSeek API Key：

```env
DEEPSEEK_API_KEY=your_actual_api_key_here
```

### 3. 启动数据库

```bash
docker-compose up -d
```

启动 MySQL (3306)、Neo4j (7474/7687)、Redis (6379)。

### 4. 启动后端

```bash
go mod tidy
go run cmd/server/main.go
```

后端服务: `http://localhost:8080`

### 5. (可选) 初始化测试数据

```bash
go run scripts/init_db.go
```

创建预设教师和学生账号（1 教师 + 5 学生 + 35 条证据）。

### 6. 启动前端

```bash
cd frontend
npm install
npm run dev
```

前端服务: `http://localhost:3000`

## 测试

### 运行所有测试

```bash
# 快速测试（跳过批量文件测试）
go test -short ./...

# 完整测试（含 189 文件批量提取验证）
go test -v -count=1 ./pkg/extractor/ ./internal/agent/ ./internal/llm/

# 批量 LLM 回归测试（需配置 API Key）
go run scripts/batch_llm_classify.go --sample 20 --workers 3

# 并发性能测试（需启动后端服务）
go run scripts/concurrent_test.go
```

### 测试覆盖

| 模块 | 测试文件 | 说明 |
|------|---------|------|
| `pkg/extractor/` | `pdf_test.go`, `docx_test.go` | PDF(108) + DOCX(81) 提取测试，批量/并发/边界/benchmark |
| `internal/agent/` | `evidence_agent_test.go` | PreprocessEvidence, classifyWithKeywords, assessEvidenceWithRules 等 26 项 |
| `internal/llm/` | `client_test.go` | HTTP Mock 测试（成功/错误/空响应/API Key 验证等 8 项）|

## 项目结构

```
research-ability-assessment/
├── cmd/server/main.go              # 应用入口
├── configs/config.dev.yaml         # 开发环境配置
├── docs/                           # 项目文档
├── frontend/                       # Vue3 前端
├── internal/
│   ├── agent/                      # 多智能体模块（核心）
│   │   ├── control_unit.go         # ControlUnit 编排引擎
│   │   ├── evidence_agent.go       # EvidenceAgent (LLM优先分类+规则兜底)
│   │   ├── feedback_agent.go       # FeedbackAgent (LLM诊断反馈)
│   │   ├── inference_agent.go      # InferenceAgent (LLM主导多维推理)
│   │   ├── logic_unit.go           # LogicUnit (规则兜底计算)
│   │   └── storage_unit.go         # StorageUnit (持久化+图谱+引用)
│   ├── config/config.go            # Viper+YAML 配置管理
│   ├── handler/                    # HTTP 处理器
│   ├── llm/                        # LLM 客户端 (OpenAI兼容)
│   ├── middleware/                  # Gin 中间件 (Auth/CORS/Logging)
│   ├── models/                     # GORM 数据模型 + Citation + Dimension
│   ├── repository/                 # 数据访问层 (MySQL + Neo4j)
│   └── service/                    # 业务逻辑层
├── pkg/
│   ├── cache/redis.go              # Redis 缓存 (优雅降级)
│   ├── extractor/                  # 多格式文件提取
│   │   ├── extractor.go            # ExtractorChain
│   │   ├── docx.go                 # DOCX 提取器
│   │   └── pdf.go                  # PDF 提取器
│   └── utils/
├── scripts/
│   ├── init_db.go                  # 数据库种子
│   ├── concurrent_test.go          # 并发性能测试
│   └── batch_llm_classify.go       # LLM 批量分类回归测试
├── .github/workflows/test.yml      # CI 自动测试
├── docker-compose.yml
└── README.md
```

## 测试账号

运行 `go run scripts/init_db.go` 后可用：

| 角色 | 邮箱 | 密码 |
|------|------|------|
| 教师 | `1@tea.com` | `123456` |
| 学生 | `1@stu.com` ~ `5@stu.com` | `123456` |

## 文档

- [用户手册](docs/user_manual.md)
- [开发指南](docs/development_guide.md)
- [部署指南](docs/deployment_guide.md)
- [演示脚本](docs/demo_script.md)
- [数据库选择](docs/database_options.md)
- [并发测试指南](docs/concurrent_testing_guide.md)
- [结果与报告管理](docs/result_and_report_management.md)
- [CHANGELOG](CHANGELOG.md)
