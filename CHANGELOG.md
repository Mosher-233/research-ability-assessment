# 更新日志

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

