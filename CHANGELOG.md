# 更新日志

## 2026-05-23

### 测试基础设施完善

- **新增 DOCX 提取器单元测试** (`pkg/extractor/docx_test.go`)
  - 7 类质量等级单文件测试 (A-G)
  - 81 文件全量批量测试 (100% 通过)
  - ExtractorChain 混合格式测试 (DOCX + PDF)
  - Benchmark 测试 (单文件 + 小文件)
- **新增 Agent 单元测试** (`internal/agent/evidence_agent_test.go`)
  - PreprocessEvidence 6 项 (空格/换行/截断/中文保留等)
  - classifyWithKeywords 7 项 (文献检索/实验设计/数据分析/批判/创新/空文本/无关文本)
  - assessEvidenceWithRules 4 项 (高质量/低质量/中等/可信度区间)
  - ClassifyEvidence/ExtractKBMInfo 降级路径验证
  - KBM 关键词覆盖完整性和无重叠校验
  - buildRationale/calcCredibility 单元测试
  - **共 26 项测试，全部通过**
- **新增 LLM Client Mock 测试** (`internal/llm/client_test.go`)
  - HTTP Mock Server 驱动的 8 项测试
  - 覆盖: 成功响应/HTTP错误/空Choices/无API Key/未配置环境变量/非法JSON/多Choices
- **CI 自动测试工作流** (`.github/workflows/test.yml`)
  - extractor-tests: DOCX + PDF 提取 + Benchmark
  - unit-tests: 全项目短测试 (`go test -short ./...`)
  - llm-batch-test: LLM 采样测试 (main 分支 push 触发)

### PDF 提取验证

- **108/108 PDF 文件 100% 提取成功**
  - 字符数范围: 196-1,721, 平均 446
  - 覆盖 A-H 八个质量等级 + 3 个研究主题
  - 并发提取测试通过 (19 项测试全通过)

### LLM 批量分类验证

- **新增批量分类工具** `scripts/batch_llm_classify.go`
  - 支持并发 Worker、采样模式、可配置路径
  - 189 文件 (81 DOCX + 108 PDF) 全量测试
  - 输出：KBM 分布/等级分组/可信度梯度/异常检测详情
- **全量测试结果**：
  - 189/189 成功 (100%), 耗时 6m23s, 平均 API 响应 2.03s
  - 异常检测率 98% (49/50), 等级±1以内率 100%
  - 高估率 0% (LLM 评级保守), 可信度 A=87%→H=8% 单调递减

### 论文完善

- **第五章全面重写**: 新增文件提取验证、LLM vs 关键词全量对比、异常检测案例
- **摘要更新** (中/英文): 实验规模从 5 学生 35 证据扩展为两层次验证
- **全文审查**: 修复 12 个跨章节不一致 (LLM 优先策略对齐、虚假引用移除、重复文本纠正)
- **新增 3 张图表**: fig5-4 (LLM vs 关键词对比)、fig5-5 (异常检测)、fig5-6 (可信度梯度)
- **图片清理**: 统一 21 张 PNG 导出，移除重复 .drawio.png 文件

### README 更新

- Go 版本更新至 1.24
- 新增功能特性 (异常检测、多格式提取、离线演示)
- 新增测试章节 (运行方式、覆盖矩阵、测试数据说明)
- 项目结构补充 (extractor、scripts、CI)

---

## 2026-05-05

### 核心架构优化：LLM 融入主流评估管线

#### 文件解析能力扩展（P0）
- **新增 `pkg/extractor/` 包**：统一文件内容提取接口
  - `extractor.go`：`ContentExtractor` 接口与 `ExtractorChain` 链式提取器
  - `pdf.go`：基于 `github.com/ledongthuc/pdf` 的 PDF 文本提取
  - `docx.go`：基于 `archive/zip` + `encoding/xml` 的 DOCX 文本提取（零依赖）
- **证据上传支持所有常见学术格式**：PDF、DOCX、TXT、MD
- **证据模型扩展**：新增 `SourceType`、`ExtractionMetadata`、`ReviewedBy`、`ReviewedAt` 字段

#### LLM 融入主要 Agent 管线（P0/P1）
- **InferenceAgent 重写**：LLM 作为主推理路径，规则引擎为降级方案
  - 新增 `NewInferenceAgentWithLLM()` 构造函数
  - 引入完整 Rubrics（评分量规）体系，每个维度五级标准
  - LLM 响应包含 `evidence_quotes`（证据原文引用）
  - 自动从推理结果提取并生成 `EvidenceCitation` 记录
- **EvidenceAgent 重写**：LLM 优先分类与评估
  - 新增 `classifyWithLLM()` 和 `extractKBMWithLLM()` 方法
  - 保留关键词匹配 `classifyWithKeywords()` 作为降级
- **ControlUnit**：评估后自动保存引用记录到数据库

#### LLM 上下文丰富与解析鲁棒性（P1）
- **`inference_service.go`**：
  - `buildEvidenceContext()` 按维度分组展示证据，包含完整元数据
  - `parseLLMResponse()` 支持 markdown 代码块提取，添加重试机制（最多2次）
  - 新增 `extractJSON()` 通用 JSON 提取函数，按括号深度匹配
- **`feedback_agent.go`**：
  - LLM 反馈 prompt 包含评估推理的详细上下文
  - `parseLLMFeedbackResponse()` 使用统一 JSON 提取逻辑
  - 规则反馈按分数区间细分，在每个等级内显示当前得分

#### 证据引用与可溯源性（P1/P2）
- **新增 `internal/models/citation.go`**：`EvidenceCitation` 模型
  - 追踪结果ID → 维度ID → 证据ID → 引用原文段落的完整链路
- **新增 `ResultRepo` 引用方法**：`CreateCitation`、`CreateCitations`、`GetCitationsByResultID`、`GetCitationsByResultIDs`
- **`StorageUnit` 新增 `StoreCitations()` 方法**
- **结果查询 API 自动附带 citations 数据**

#### LLM 客户端增强
- **`llm/client.go`**：新增 `ResponseFormat` 支持，可约束 LLM 输出 JSON 格式

---

## 2026-04-26

### 文档体系维护

#### README.md 完善
- **项目结构修正**：移除不存在的 `config.mysql.yaml`、`config.supabase.yaml` 引用，补充实际存在的 `scripts/`、`CHANGELOG.md`、`pkg/cache/` 目录
- **技术栈补全**：添加 Redis（可选）、JWT 认证、TypeScript、Vite 4、Vue Router 4、ECharts 6
- **快速开始优化**：将 `cp` 改为 `copy`（适配 Windows），新增第5步 `go run scripts/init_db.go` 测试数据初始化
- **LLM 配置修正**：`max_tokens: 1000` → `2048`（与实际代码一致）
- **文档索引补全**：添加 `demo_script.md`、`result_and_report_management.md`、`concurrent_testing_guide.md` 链接
- **FAQ 扩展**：修正数据库支持说明，新增测试账号表格

#### docs/demo_script.md 修正
- **数据库修正**：PostgreSQL + Neo4j → MySQL + Neo4j
- **登录凭证修正**：教师 `admin@example.com/password123` → `1@tea.com/123456`，学生 `student@example.com/password123` → `1@stu.com/123456`
- **评估维度修正**：移除不存在的"学术写作能力"，保持4维度与实际代码一致
- **LLM 引用修正**：OpenAI API / GPT-3.5-turbo → DeepSeek API / deepseek-chat

#### docs/development_guide.md 修正
- **技术栈修正**：PostgreSQL 14+ → MySQL 8.0+，Neo4j 5.0+ → 5.23+，OpenAI API → DeepSeek API
- **架构图修正**：PostgreSQL → MySQL，OpenAI API → DeepSeek API
- **目录结构修正**：移除不存在的 `io_unit.go`、`prompts.go`、`pkg/logger/`、`utils/jwt.go`、`utils/validator.go`；补充 `agent_handler.go`、`dimension.go`、`feedback.go`、`report.go`、`pkg/cache/redis.go`、`pkg/utils/id_generator.go`

#### docs/deployment_guide.md 修正
- **全篇数据库修正**：PostgreSQL 全面替换为 MySQL（版本号、端口、用户名密码、Docker Compose 示例）
- **Neo4j 修正**：版本 5.0 → 5.23，默认密码 `neo4j` → `password123`
- **LLM 修正**：OpenAI/GPT-3.5-turbo → DeepSeek/deepseek-chat
- **配置文件路径修正**：`configs/config.yaml` → `configs/config.dev.yaml`
- **Vite proxy 修正**：`rewrite` 逻辑与实际 `vite.config.ts` 保持一致
- **简化日志配置**：移除不存在的 `logger` 配置项

#### docs/user_manual.md 修正
- **环境要求修正**：PostgreSQL 14+, Neo4j 5.0+ → MySQL 8.0+, Neo4j 5.23+
- **系统架构修正**：PostgreSQL → MySQL
- **评估维度修正**：移除"学术写作能力"

#### docs/database_options.md 修正
- **默认数据库**：PostgreSQL（默认）→ MySQL（默认，与 docker-compose.yml 一致）
- **章节顺序重排**：MySQL → PostgreSQL → Supabase
- **配置方式统一**：移除不存在的 `config.mysql.yaml`、`config.supabase.yaml` 引用，统一使用 `config.dev.yaml`
- **Neo4j 密码修正**：`neo4jpassword` → `password123`
- **快速开始完善**：添加 `.env` 配置和 `init_db.go` 步骤

#### .gitignore 维护
- **IDE 目录**：新增 `.trae/` 忽略规则（Trae IDE 技能/配置目录，类似 `.vscode/`、`.idea/`）
- **测试覆盖率**：新增 `coverage/` 和 `*.coverprofile` 防御性规则
- **Python 产物**：新增 `__pycache__/` 和 `*.py[cod]`（项目中存在 Python 脚本）
- **Dump 文件**：新增 `*.dump` 防御性规则

---

## 2026-03-15

### 核心问题修复与功能完善

#### LLM集成问题修复
- **修复LLM API硬编码响应问题**：移除`llm/client.go`中的默认硬编码文本，确保AI真正被调用
  - 增加API超时时间从30秒到120秒
  - 增加max_tokens到2048
  - 添加详细日志记录
- **修复正则表达式兼容性问题**：完全重写`llm/parser.go`的解析逻辑，改用逐行解析
  - 解决Go regexp不支持正向先行断言的问题
  - 添加`foundLevel`标志，只解析第一个出现的KBM级别
- **修复JSON数据类型处理问题**：安装`gorm.io/datatypes`包
  - 更新所有相关模型字段为`datatypes.JSON`类型
  - 修复所有JSON序列化和反序列化错误
- **修复推理结果保存问题**：在`GenerateInferenceWithLLM`函数中添加保存逻辑
  - 确保LLM响应解析后正确保存到数据库

#### 删除证据功能完善
- **添加删除证据API**：实现完整的删除功能
  - `internal/service/evidence_service.go`：添加`DeleteEvidence`方法
  - `internal/handler/evidence_handler.go`：添加删除证据处理函数
  - `cmd/server/main.go`：添加删除证据路由

#### 报告生成与管理功能完善
- **完整的报告API实现**：
  - 添加`/reports`路由获取所有报告
  - 添加`/reports/student`路由获取学生自己的报告
  - 添加`enrichResultsWithDetails`和`enrichReportsWithDetails`函数补全信息
- **报告详情对话框**：完善报告管理页面的查看功能
  - 信息栏显示完整信息（报告编号、学生ID、学生姓名、任务ID、任务名称、生成时间）
  - 综合评价：总体得分、等级、班级排名、超越比例
  - 能力分析：各维度得分进度条
  - 能力雷达图：使用ECharts绘制
  - 优势分析：自动识别优势
  - 待提升方向：识别不足
  - 改进建议：个性化建议、行动项、学习资源

#### 雷达图与评分优化
- **100分制转换**：
  - 更新`frontend/src/utils/chart.ts`：雷达图使用100分制
  - 更新`frontend/src/utils/format.ts`：格式化显示100分制
  - 归一化处理：自动检测1.0系数制，转换为100分制
- **维度名称映射**：添加完整的维度ID到中文名称映射
  - literature_review → 文献综述
  - research_design → 研究设计
  - data_analysis → 数据分析
  - critical_thinking → 批判性思维
  - 支持多种维度ID格式

#### 前端UI/UX改进
- **登录页面**：
  - 确保登录和注册按钮宽度完全一致
  - 使用flexbox垂直布局和深度选择器
- **报告管理页面**：
  - 学生选择：同时显示学生姓名和学生ID
  - 任务选择：同时显示任务名称和任务ID
  - 报告列表：显示完整的学生姓名和任务名称
- **结果管理页面**：
  - 学生任务选择区域：学生可以选择任务并生成推理结果
  - 结果信息栏：显示完整的学生信息和任务信息
  - 维度得分表：正确显示维度名称、得分、等级、详情

#### 后端服务完善
- **学生推理结果生成**：
  - 添加`POST /api/v1/results/generate/student`路由
  - 学生可以为自己的任务生成推理结果
- **报告服务扩展**：
  - 实现`GetAllReports`方法
  - 实现`GetReportsByStudentID`方法
- **推理服务优化**：
  - 修复LLM响应解析时维度名称错误赋值为等级的问题（`inference_service.go:502-512`）
  - 添加`nameMap`确保维度名称正确映射

#### 技术债务修复
- **编译错误修复**：解决所有类型错误、未使用变量、导入缺失问题
- **数据库模型**：
  - 更新`result.go`和`report.go`模型，使用`datatypes.JSON`
  - 确保所有JSON字段正确序列化
- **路由配置**：`cmd/server/main.go`中添加所有新路由

#### 问题修复总结
- ✅ 不同证据AI反馈完全相同的问题（真正调用LLM）
- ✅ 删除证据失败的问题
- ✅ 结果管理没有生成结果的问题
- ✅ 报告生成失败的问题
- ✅ 结果管理信息栏没有对应学生信息的问题
- ✅ 任务选择只显示ID的问题（同时显示名称和ID）
- ✅ 学生选择只显示ID的问题（同时显示姓名和ID）
- ✅ 登录和注册按钮宽度不一致的问题
- ✅ 报告管理能力分析显示等级而非维度名称的问题
- ✅ 雷达图1.00系数制转换为100分制的问题
- ✅ 报告显示"无相关数据"的问题

---

## 2026-03-14

### 核心功能：结果管理和报告管理完善

- **完整的推理服务（InferenceService）实现**
  - 实现`GenerateInference()`方法，基于证据生成综合评估结果
  - 实现`GetClassStats()`方法，获取班级统计数据（平均分、最高分、最低分、维度平均分）
  - 实现`CalculateRankAndPercentile()`方法，计算学生排名和百分位
  - 支持4个默认评估维度：文献综述、研究设计、数据分析、批判性思维（各权重0.25）
  - 支持简化推理和LLM智能推理两种模式
- **完整的报告服务（ReportService）实现**
  - 实现`GenerateReport()`方法，生成详尽的研究能力评价报告
  - 班级对比分析：班级人数、平均分、最高分、最低分、各维度平均分
  - 排名分析：班级排名、超越比例（百分位）
  - 优势劣势分析：自动识别≥80分的优势维度和<70分的待提升维度
  - 个性化改进建议：针对每个待提升维度提供具体建议、可执行行动项、推荐学习资源
  - 能力雷达图数据生成
  - 自动生成并保存格式化的TXT报告文件到uploads/reports目录

### 数据模型扩展

- 新增`internal/models/report.go`：完整的报告相关数据模型
  - `Report`模型：包含综合评价、班级对比、排名、建议等完整信息
  - `ClassComparisonData`：班级对比数据
  - `ImprovementSuggestion`：改进建议模型
  - `LearningResource`：学习资源模型
  - `RadarChartData`：雷达图数据模型
  - `Dimension`：评估维度模型
- 新增`internal/models/feedback.go`：反馈模型
- 更新`internal/models/evidence.go`：添加文件支持字段（FileName、FilePath、FileType、FileSize）

### 基础设施完善

- **有序ID生成器（pkg/utils/id\_generator.go）**
  - 新增`GenerateEvidenceID()`函数：证据ID格式为EV+日期+序号
  - 保持教师ID（T+日期+序号）、学生ID（S+日期+序号）、任务ID（TK+日期+序号）的有序性
- **仓库层扩展（internal/repository/postgres/result\_repo.go）**
  - 新增`CreateReport()`方法
  - 新增`GetReportByID()`方法
  - 新增`GetReportByStudentAndTask()`方法
  - 新增`GetReportsByTaskID()`方法
  - 新增`GetReportsByStudentID()`方法
  - 新增`GetAllReports()`方法
- **API层扩展**
  - 新增`POST /api/v1/results/generate`：生成推理结果
  - 新增`POST /api/v1/reports/generate`：生成完整报告
  - 更新`internal/handler/result_handler.go`：添加完整的结果和报告处理器
- **数据库迁移**
  - 在`cmd/server/main.go`中添加Report模型的自动迁移
  - 在`migrateDatabase()`函数中注册\&models.Report{}

### 前端功能完善

- 更新`frontend/src/views/EvidenceManagement.vue`：
  - 支持任务去重
  - 支持文件上传
  - 支持AI反馈查看
  - 教师界面实时查看学生证据
- 更新`frontend/src/views/TaskManagement.vue`：
  - 修复学生ID重复显示问题
  - 添加学生任务列表去重逻辑

### 问题修复

- **学生ID重复显示问题**：前端添加去重逻辑，后端分配任务时检查是否已分配
- **创建任务和分配任务的500错误**：将日期字段从time.Time改为string，手动解析日期
- **证据上传问题**：完善文件上传和下载API
- **后端编译错误**：修复所有类型错误和未使用变量问题

### 文档完善

- 新增`docs/result_and_report_management.md`：
  - 从学生角度详细说明结果管理和报告管理的使用流程和关注点
  - 从教师角度详细说明教学应用场景和关注点
  - 从产品角度说明设计原则和未来扩展方向
  - 包含系统实现细节、API接口说明、核心服务介绍
  - 提供真实使用场景示例

### 初始化脚本

- 新增`scripts/init_db.go`：数据库初始化脚本，支持创建测试用户和清空表

## 2026-3-09

### 后端更新

- 修复任务仓库中的Preload调用，确保正确加载关联数据
- 修复学生列表获取功能，确保能够正确获取学生信息
- 添加任务归档API，支持任务状态的更新和管理
- 集成AI API，为学生提供智能建议和反馈
- 优化数据库查询性能，减少响应时间

### 前端更新

- 实现任务管理页面的Tab切换功能，支持查看不同状态的任务
- 添加任务归档功能，允许用户归档已完成的任务
- 修复雷达图渲染问题，确保图表正确显示
- 优化界面响应速度，提升用户体验
- 完善表单验证逻辑，减少错误输入

### 技术修复

- 解决数据库迁移问题，确保数据结构的一致性
- 修复模型关联关系错误，确保数据完整性
- 优化代码结构，提高代码可读性和可维护性
- 修复部分API接口的错误处理逻辑
- 完善日志记录系统，便于问题排查

### 其他改进

- 增加单元测试覆盖率，提高代码质量
- 优化项目配置，简化部署流程
- 完善文档说明，便于团队协作
- 增加错误处理机制，提高系统稳定性

