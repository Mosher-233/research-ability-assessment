# 研究能力评价系统开发文档

## 1. 系统架构

### 1.1 技术栈

- **后端**：
  - 语言：Go 1.20+
  - 框架：Gin 1.9.0+
  - 数据库：PostgreSQL 14+, Neo4j 5.0+
  - 认证：JWT
  - LLM集成：OpenAI API

- **前端**：
  - 框架：Vue 3
  - 语言：TypeScript
  - UI库：Element Plus
  - 路由：Vue Router
  - HTTP客户端：Axios

### 1.2 系统架构图

```
┌─────────────────┐       ┌─────────────────┐       ┌─────────────────┐
│    前端应用     │◄──────┤    后端API      │◄──────┤    数据库       │
│  Vue 3 + TS    │       │  Go + Gin       │       │ PostgreSQL +   │
│  Element Plus  │       │  JWT认证        │       │ Neo4j          │
└─────────────────┘       └─────────────────┘       └─────────────────┘
                               │
                               ▼
                      ┌─────────────────┐
                      │    AI Agent     │
                      │  证据收集       │
                      │  能力推理       │
                      │  反馈生成       │
                      └─────────────────┘
                               │
                               ▼
                      ┌─────────────────┐
                      │    LLM服务      │
                      │  OpenAI API     │
                      └─────────────────┘
```

## 2. 目录结构

### 2.1 后端目录结构

```
research-ability-assessment/
├── cmd/                  # 命令行工具
│   └── server/           # 服务器启动
│       └── main.go       # 主程序入口
├── configs/              # 配置文件
│   └── config.dev.yaml   # 开发环境配置
├── internal/             # 内部包
│   ├── agent/            # AI Agent实现
│   │   ├── control_unit.go    # 控制单元
│   │   ├── evidence_agent.go  # 证据收集Agent
│   │   ├── feedback_agent.go  # 反馈生成Agent
│   │   ├── inference_agent.go # 能力推理Agent
│   │   ├── io_unit.go         # IO单元
│   │   ├── logic_unit.go      # 逻辑单元
│   │   └── storage_unit.go    # 存储单元
│   ├── config/           # 配置管理
│   │   └── config.go     # 配置加载和管理
│   ├── handler/          # API处理器
│   │   ├── auth_handler.go     # 认证相关API
│   │   ├── evidence_handler.go # 证据相关API
│   │   ├── result_handler.go   # 结果相关API
│   │   └── task_handler.go     # 任务相关API
│   ├── llm/              # 大语言模型集成
│   │   ├── client.go     # LLM客户端
│   │   ├── parser.go     # 结果解析器
│   │   └── prompts.go    # 提示词模板
│   ├── middleware/       # 中间件
│   │   ├── auth.go       # 认证中间件
│   │   ├── cors.go       # CORS中间件
│   │   └── logging.go    # 日志中间件
│   ├── models/           # 数据模型
│   │   ├── evidence.go   # 证据模型
│   │   ├── result.go     # 结果模型
│   │   ├── task.go       # 任务模型
│   │   └── user.go       # 用户模型
│   ├── repository/       # 数据访问层
│   │   ├── neo4j/        # Neo4j仓库
│   │   │   └── graph_repo.go  # 图数据库操作
│   │   └── postgres/     # PostgreSQL仓库
│   │       ├── result_repo.go # 结果数据操作
│   │       ├── task_repo.go   # 任务数据操作
│   │       └── user_repo.go   # 用户数据操作
│   └── service/          # 业务逻辑层
│       ├── auth_service.go     # 认证服务
│       ├── evidence_service.go # 证据服务
│       ├── inference_service.go # 推理服务
│       ├── report_service.go   # 报告服务
│       └── task_service.go     # 任务服务
├── pkg/                  # 公共包
│   ├── logger/           # 日志工具
│   │   └── logger.go     # 日志配置
│   └── utils/            # 工具函数
│       ├── jwt.go        # JWT工具
│       └── validator.go  # 数据验证
├── go.mod                # Go模块文件
└── go.sum                # 依赖版本锁定
```

### 2.2 前端目录结构

```
research-ability-assessment/frontend/
├── src/                  # 源代码
│   ├── api/              # API服务
│   │   ├── auth.ts       # 认证相关API
│   │   ├── evidence.ts   # 证据相关API
│   │   ├── result.ts     # 结果相关API
│   │   └── task.ts       # 任务相关API
│   ├── router/           # 路由配置
│   │   └── index.ts      # 路由定义
│   ├── types/            # 类型定义
│   │   ├── evidence.ts   # 证据类型
│   │   ├── result.ts     # 结果类型
│   │   ├── task.ts       # 任务类型
│   │   └── user.ts       # 用户类型
│   ├── utils/            # 工具函数
│   │   ├── chart.ts      # 图表工具
│   │   ├── format.ts     # 格式化工具
│   │   └── validator.ts  # 验证工具
│   ├── views/            # 页面组件
│   │   ├── Dashboard.vue # 仪表盘
│   │   ├── EvidenceManagement.vue # 证据管理
│   │   ├── Login.vue     # 登录页面
│   │   ├── Register.vue  # 注册页面
│   │   ├── ResultManagement.vue # 结果管理
│   │   ├── TaskManagement.vue # 任务管理
│   │   └── ReportManagement.vue # 报告管理
│   ├── App.vue           # 根组件
│   └── main.ts           # 入口文件
├── index.html            # HTML模板
├── package.json          # 项目依赖
├── tsconfig.json         # TypeScript配置
├── tsconfig.node.json    # Node.js TypeScript配置
└── vite.config.ts        # Vite配置
```

## 3. 核心功能模块

### 3.1 后端核心模块

#### 3.1.1 认证模块

认证模块负责用户的注册、登录和身份验证。使用JWT进行身份认证，确保API访问的安全性。

**主要功能**：
- 用户注册
- 用户登录
- 获取用户信息
- JWT令牌验证

#### 3.1.2 任务模块

任务模块负责研究任务的创建、管理和分配。教师可以创建任务并分配给学生，学生可以查看分配给自己的任务。

**主要功能**：
- 创建任务
- 获取任务列表
- 获取任务详情
- 分配任务给学生
- 获取学生任务列表

#### 3.1.3 证据模块

证据模块负责研究证据的提交和管理。学生可以提交研究证据，教师可以查看和管理证据。

**主要功能**：
- 提交证据
- 获取证据详情
- 获取学生任务的证据列表
- 根据学生ID和任务ID获取证据列表

#### 3.1.4 推理模块

推理模块负责使用AI Agent对学生的研究能力进行评估。AI Agent会分析学生提交的证据，评估学生的研究能力维度。

**主要功能**：
- 证据收集和分析
- 能力维度评估
- 推理结果生成

#### 3.1.5 报告模块

报告模块负责生成研究能力评价报告。根据评估结果，生成详细的评价报告，包括能力分析、改进建议等内容。

**主要功能**：
- 生成学生报告
- 生成任务报告
- 下载报告

### 3.2 前端核心模块

#### 3.2.1 登录注册模块

登录注册模块负责用户的登录和注册功能。用户可以通过邮箱和密码登录系统，也可以注册新账号。

#### 3.2.2 任务管理模块

任务管理模块负责任务的创建、管理和分配。教师可以创建任务并分配给学生，学生可以查看分配给自己的任务。

#### 3.2.3 证据管理模块

证据管理模块负责证据的提交和管理。学生可以提交研究证据，教师可以查看和管理证据。

#### 3.2.4 结果管理模块

结果管理模块负责查看和管理评估结果。教师可以查看所有学生的评估结果，学生只能查看自己的评估结果。

#### 3.2.5 报告管理模块

报告管理模块负责生成和管理评价报告。教师可以生成详细的研究能力评价报告，并查看和下载报告。

## 4. API接口设计

### 4.1 认证接口

| 路径 | 方法 | 功能 | 权限 |
|------|------|------|------|
| `/api/v1/auth/register` | POST | 用户注册 | 公共 |
| `/api/v1/auth/login` | POST | 用户登录 | 公共 |
| `/api/v1/user/info` | GET | 获取用户信息 | 受保护 |

### 4.2 任务接口

| 路径 | 方法 | 功能 | 权限 |
|------|------|------|------|
| `/api/v1/tasks` | POST | 创建任务 | 教师 |
| `/api/v1/tasks` | GET | 获取任务列表 | 教师 |
| `/api/v1/tasks/:task_id` | GET | 获取任务详情 | 教师/学生 |
| `/api/v1/tasks/:task_id/assign` | POST | 分配任务 | 教师 |
| `/api/v1/tasks/:task_id/students` | GET | 获取学生任务列表 | 教师 |
| `/api/v1/tasks/students/list` | GET | 获取学生列表 | 教师 |

### 4.3 证据接口

| 路径 | 方法 | 功能 | 权限 |
|------|------|------|------|
| `/api/v1/evidences` | POST | 提交证据 | 学生 |
| `/api/v1/evidences/:evidence_id` | GET | 获取证据详情 | 教师/学生 |
| `/api/v1/evidences/student-task/:student_task_id` | GET | 获取学生任务的证据列表 | 教师/学生 |
| `/api/v1/evidences/student-task` | GET | 根据学生ID和任务ID获取证据列表 | 教师/学生 |

### 4.4 结果接口

| 路径 | 方法 | 功能 | 权限 |
|------|------|------|------|
| `/api/v1/results/:result_id` | GET | 获取推理结果详情 | 教师/学生 |
| `/api/v1/results/task/:task_id` | GET | 获取任务的推理结果列表 | 教师 |
| `/api/v1/results/student-task` | GET | 根据学生ID和任务ID获取推理结果 | 教师/学生 |
| `/api/v1/results/report/student` | GET | 生成学生报告 | 教师 |
| `/api/v1/results/report/task/:task_id` | GET | 生成任务报告 | 教师 |

## 5. 数据模型设计

### 5.1 后端数据模型

#### 5.1.1 用户模型

```go
type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"size:255;not null"`
	Email     string    `json:"email" gorm:"size:255;uniqueIndex;not null"`
	Password  string    `json:"-" gorm:"size:255;not null"`
	Role      string    `json:"role" gorm:"size:50;not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Teacher struct {
	User
	Department string `json:"department" gorm:"size:255;not null"`
	Title      string `json:"title" gorm:"size:255;not null"`
}

type Student struct {
	User
	StudentID string `json:"student_id" gorm:"size:50;uniqueIndex;not null"`
	Major     string `json:"major" gorm:"size:255;not null"`
	Grade     string `json:"grade" gorm:"size:50;not null"`
}
```

#### 5.1.2 任务模型

```go
type Task struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"size:255;not null"`
	Description string    `json:"description" gorm:"type:text;not null"`
	CourseID    string    `json:"course_id" gorm:"size:100;not null"`
	TeacherID   uint      `json:"teacher_id" gorm:"not null"`
	StartDate   time.Time `json:"start_date" gorm:"not null"`
	EndDate     time.Time `json:"end_date" gorm:"not null"`
	Status      string    `json:"status" gorm:"size:50;not null;default:'pending'"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Teacher     User      `json:"teacher,omitempty" gorm:"foreignKey:TeacherID"`
}

type StudentTask struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	TaskID       uint      `json:"task_id" gorm:"not null"`
	StudentID    uint      `json:"student_id" gorm:"not null"`
	Status       string    `json:"status" gorm:"size:50;not null;default:'pending'"`
	Progress     float64   `json:"progress" gorm:"not null;default:0"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Task         Task      `json:"task,omitempty" gorm:"foreignKey:TaskID"`
	Student      User      `json:"student,omitempty" gorm:"foreignKey:StudentID"`
}
```

#### 5.1.3 证据模型

```go
type Evidence struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	StudentTaskID uint      `json:"student_task_id" gorm:"not null"`
	Type          string    `json:"type" gorm:"size:100;not null"`
	Content       string    `json:"content" gorm:"type:text;not null"`
	KBMName       string    `json:"kbm_name" gorm:"size:255;not null"`
	KBMLevel      int       `json:"kbm_level" gorm:"not null"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	StudentTask   StudentTask `json:"student_task,omitempty" gorm:"foreignKey:StudentTaskID"`
}
```

#### 5.1.4 结果模型

```go
type InferenceResult struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	StudentID       uint      `json:"student_id" gorm:"not null"`
	TaskID          uint      `json:"task_id" gorm:"not null"`
	OverallScore    float64   `json:"overall_score" gorm:"not null"`
	OverallLevel    string    `json:"overall_level" gorm:"size:50;not null"`
	DimensionScores string    `json:"dimension_scores" gorm:"type:text;not null"` // JSON格式存储
	Reasoning       string    `json:"reasoning" gorm:"type:text;not null"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	Student         User      `json:"student,omitempty" gorm:"foreignKey:StudentID"`
	Task            Task      `json:"task,omitempty" gorm:"foreignKey:TaskID"`
}
```

### 5.2 前端类型定义

#### 5.2.1 用户类型

```typescript
export interface User {
  id: string
  name: string
  email: string
  role: string
  created_at: string
  updated_at: string
}

export interface Teacher extends User {
  department: string
  title: string
}

export interface Student extends User {
  student_id: string
  major: string
  grade: string
}

export interface LoginRequest {
  email: string
  password: string
}

export interface RegisterRequest {
  name: string
  email: string
  password: string
  role: string
}

export interface AuthResponse {
  code: number
  message: string
  data: {
    token: string
    user: User
  }
}
```

#### 5.2.2 任务类型

```typescript
export interface Task {
  id: string
  name: string
  description: string
  course_id: string
  teacher_id: string
  start_date: string
  end_date: string
  status: string
  created_at: string
  updated_at: string
  teacher?: {
    id: string
    name: string
    email: string
  }
  student_count?: number
  completed_count?: number
}

export interface StudentTask {
  id: string
  task_id: string
  student_id: string
  status: string
  progress: number
  created_at: string
  updated_at: string
  task?: Task
  student?: {
    id: string
    name: string
    student_id: string
  }
}

export interface CreateTaskRequest {
  name: string
  description: string
  course_id: string
  start_date: string
  end_date: string
}

export interface AssignTaskRequest {
  student_ids: string[]
}

export interface TaskStatistics {
  total_students: number
  completed: number
  processing: number
  pending: number
}
```

#### 5.2.3 证据类型

```typescript
export interface Evidence {
  id: string
  student_task_id: string
  type: string
  content: string
  kbm_name: string
  kbm_level: number
  created_at: string
  updated_at: string
  student_task?: {
    id: string
    task_id: string
    student_id: string
  }
}

export interface CreateEvidenceRequest {
  student_task_id: string
  type: string
  content: string
  kbm_name: string
  kbm_level: number
}
```

#### 5.2.4 结果类型

```typescript
export interface DimensionScore {
  name: string
  score: number
  level: string
  details: string
  evidence_ids: string[]
}

export interface InferenceResult {
  id: string
  student_id: string
  task_id: string
  overall_score: number
  overall_level: string
  dimension_scores: Record<string, DimensionScore>
  reasoning: string
  created_at: string
  updated_at: string
  student?: {
    id: string
    name: string
    student_id: string
  }
  task?: {
    id: string
    name: string
  }
}
```

## 6. AI Agent设计

### 6.1 Agent架构

AI Agent采用模块化设计，由以下组件组成：

- **控制单元**：协调各个Agent的工作流程
- **证据收集Agent**：收集和分析学生提交的证据
- **能力推理Agent**：根据证据评估学生的研究能力
- **反馈生成Agent**：生成详细的评估反馈和改进建议
- **IO单元**：处理与外部系统的交互
- **逻辑单元**：执行推理逻辑
- **存储单元**：管理数据存储

### 6.2 工作流程

1. **证据收集**：收集学生提交的研究证据
2. **证据分析**：分析证据内容，提取关键信息
3. **能力评估**：根据证据评估学生的研究能力维度
4. **结果生成**：生成详细的评估结果和推理过程
5. **反馈生成**：生成改进建议和个性化反馈
6. **报告生成**：生成完整的研究能力评价报告

### 6.3 LLM集成

系统集成了OpenAI API，使用GPT-3.5-turbo模型进行能力评估和反馈生成。通过精心设计的提示词模板，引导LLM生成准确、客观的评估结果。

## 7. 系统部署

### 7.1 开发环境部署

#### 7.1.1 后端部署

1. 克隆项目代码
2. 安装依赖：`go mod tidy`
3. 配置数据库连接（修改 `configs/config.dev.yaml`）
4. 启动后端服务：`go run cmd/server/main.go`

#### 7.1.2 前端部署

1. 进入前端目录：`cd frontend`
2. 安装依赖：`npm install`
3. 启动开发服务器：`npm run dev`

### 7.2 生产环境部署

#### 7.2.1 后端部署

1. 构建可执行文件：`go build -o server cmd/server/main.go`
2. 配置生产环境配置文件
3. 启动服务器：`./server`

#### 7.2.2 前端部署

1. 构建生产版本：`npm run build`
2. 将构建产物部署到静态文件服务器
3. 配置反向代理，将API请求转发到后端服务

## 8. 开发规范

### 8.1 代码规范

- **Go代码**：遵循Go语言标准规范，使用`gofmt`和`golint`进行代码格式化和检查
- **TypeScript代码**：遵循ESLint和Prettier规范，确保代码风格一致

### 8.2 命名规范

- **文件命名**：使用小写字母和下划线，如`auth_handler.go`
- **函数命名**：使用驼峰命名法，如`CreateTask`
- **变量命名**：使用驼峰命名法，如`userID`
- **常量命名**：使用全大写字母和下划线，如`MAX_TOKEN_LENGTH`

### 8.3 注释规范

- **函数注释**：每个函数都应该有详细的注释，说明函数的功能、参数和返回值
- **代码注释**：对于复杂的代码逻辑，应该添加注释说明
- **文档注释**：重要的模块和组件应该有文档注释

### 8.4 测试规范

- **单元测试**：每个模块都应该有相应的单元测试
- **集成测试**：测试模块之间的集成
- **端到端测试**：测试完整的业务流程

## 9. 故障排除

### 9.1 常见问题

#### 9.1.1 数据库连接问题

**症状**：后端服务启动失败，提示数据库连接错误
**解决方案**：
- 检查数据库服务是否运行
- 检查数据库连接配置是否正确
- 检查数据库用户权限是否正确

#### 9.1.2 API调用失败

**症状**：前端API调用失败，返回404或500错误
**解决方案**：
- 检查后端服务是否运行
- 检查API路径是否正确
- 检查请求参数是否正确
- 检查后端日志，查看具体错误信息

#### 9.1.3 认证失败

**症状**：API调用返回401错误，提示认证失败
**解决方案**：
- 检查JWT令牌是否有效
- 检查令牌是否过期
- 检查令牌是否正确传递

#### 9.1.4 LLM调用失败

**症状**：评估过程失败，提示LLM调用错误
**解决方案**：
- 检查OpenAI API密钥是否正确
- 检查网络连接是否正常
- 检查LLM请求参数是否正确

## 10. 未来规划

### 10.1 功能扩展

- **支持更多证据类型**：如视频、音频、代码等
- **增强AI评估能力**：使用更先进的LLM模型，提高评估准确性
- **添加协作功能**：支持教师之间的协作评估
- **添加数据可视化**：提供更丰富的数据可视化功能

### 10.2 技术优化

- **性能优化**：优化数据库查询和API响应时间
- **安全性增强**：加强系统安全性，防止常见的安全漏洞
- **可扩展性**：提高系统的可扩展性，支持更多用户和数据
- **可维护性**：改善代码结构，提高可维护性

### 10.3 集成与生态

- **与学习管理系统集成**：与常见的学习管理系统集成
- **与学术数据库集成**：集成学术数据库，提供更丰富的参考资料
- **开放API**：提供开放API，支持第三方应用集成
