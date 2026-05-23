# 第四章 系统详细设计与实现

本章是论文的核心章节，详细阐述评价指标体系与量规设计、数据库设计、后端LLM主导的多智能体系统实现（含文件内容提取、证据引用溯源）以及前端界面的具体开发过程。所有实现描述均基于项目代码的实际结构和逻辑。

---

## 4.1 评价指标体系与量规设计

### 4.1.1 四维一级指标体系

基于第二章所述的理论基础——工程教育认证标准、CDIO能力框架和Bloom认知层次理论——本文设计了面向电子信息工程专业大学生的四维研究能力评价指标体系。每个能力维度下定义2-3项行为标志物（Knowledge-Based Markers，KBMs），共计10项已实现KBMs（设计阶段规划12项，其中"可视化与呈现"和"反思与改进"两项待后续扩展），如下所示：

**维度一：文献检索与综述能力**

- **KBM-1：文献检索策略。** 能否制定有效的文献检索关键词和筛选标准，使用专业数据库（如IEEE Xplore、CNKI、Web of Science）进行系统检索，记录检索流程。
- **KBM-2：文献综述质量。** 能否系统梳理研究脉络，从多维度归纳已有成果，识别研究空白和前沿方向，形成结构化的文献综述。
- **KBM-3：文献批判性分析。** 能否对已有研究的局限性、方法论不足和结论偏颇进行理性分析，而非简单接受现有结论。

**维度二：研究设计与实验能力**

- **KBM-4：实验方案合理性。** 能否设计科学合理的实验方案（含对照组设计、数据划分策略、评估指标选择），有效控制实验变量。
- **KBM-5：变量控制。** 能否明确识别和控制实验中的关键变量（自变量、因变量、混淆变量），保持实验条件的一致性。
- **KBM-6：实验实施质量。** 能否按照实验方案规范执行，保证数据采集的完整性和可追溯性，妥善记录实验过程。

**维度三：数据分析与解释能力**

- **KBM-7：数据分析方法选择。** 能否选择合适的统计方法或分析框架处理实验数据（如t检验、ANOVA、ROC-AUC等），理解所用方法的适用条件和局限性。
- **KBM-8：结果解释准确性。** 能否正确解释分析结果，识别数据中的规律和异常，讨论结果与假设的一致性，分析偏差来源。

**维度四：批判性思维与创新能力**

- **KBM-9：问题提出新颖性。** 能否在研究过程中提出独立的、有学术价值的见解、改进思路或新的研究问题。
- **KBM-10：解决方案原创性。** 能否设计具有原创性的技术方案或方法框架，在现有研究基础上有明确的改进或融合创新。

> 【图4-1：评价指标体系图 — papers/figures/fig4-1-indicator-system.drawio】

### 4.1.2 KBM-标准条款-理论依据三重映射

为确保评价指标的科学性和权威性，每项KBM都建立了与工程教育认证标准条款、CDIO能力项和教育理论依据的三重映射关系。表4-1展示了完整的映射关系。

**表4-1：KBM-标准条款-理论依据三重映射表**

| KBM | 名称 | 认证标准条款 | CDIO能力项 | 理论依据 |
|-----|------|------------|-----------|---------|
| KBM-1 | 文献检索策略 | 5.1（文献调研） | 2.1 工程推理 | Bloom层次：应用 |
| KBM-2 | 文献综述质量 | 5.1 | 2.2 实验知识发现 | Bloom层次：分析 |
| KBM-3 | 文献批判性分析 | 6.2（社会影响） | 2.3 系统思维 | Bloom层次：评价 |
| KBM-4 | 实验方案合理性 | 4.1（设计方案） | 2.2 | 实验设计理论 |
| KBM-5 | 变量控制 | 4.1 | 2.2 | 科学方法 |
| KBM-6 | 实验实施质量 | 4.1 | 2.2 | 技术素养 |
| KBM-7 | 数据分析方法选择 | 4.3（数据分析） | 2.2 | 统计学基础 |
| KBM-8 | 结果解释准确性 | 4.3 | 2.1 | 科学推理 |
| KBM-9 | 问题提出新颖性 | 3.1（问题分析） | 2.4 个人技能 | 创造性思维理论 |
| KBM-10 | 解决方案原创性 | 4.1 | 2.4 | Bloom层次：创造 |

### 4.1.3 四级量规（Rubrics）设计

每项KBM制定详细的四级量规，量规设计遵循Bloom认知层次理论的递进逻辑。等级4（优秀）对应"评价/创造"层次，等级3（良好）对应"分析"层次，等级2（合格）对应"理解/应用"层次，等级1（不合格）对应"记忆"层次。

以KBM-4"实验方案合理性"为例，表4-2展示了四级量规的具体内容。

**表4-2：KBM-4"实验方案合理性"四级量规**

| 等级 | 行为描述 | 证据示例（电子信息工程场景） |
|------|---------|--------------------------|
| 4（优秀） | 能独立设计多变量析因实验方案，合理设置对照组和实验组，预测并规避潜在混淆因素 | "考虑到温度对传感器精度的非线性影响，采用双因素方差分析，设置三个温度水平（20℃/35℃/50℃），每组重复测量5次" |
| 3（良好） | 能设计单变量对照实验方案，明确自变量和因变量，描述实验步骤 | "采用单一变量法，保持其他条件不变，仅改变输入信号频率（1kHz-100kHz），测量放大器增益的变化" |
| 2（合格） | 能参照模板设计基本实验方案，但变量定义不够清晰，缺少对照组设计 | "参照实验指导书设计实验，测试信号频率与放大倍数的关系" |
| 1（不合格） | 缺乏实验设计意识，方案混乱或不可执行，无法区分变量类型 | "测试一下电路" |

完整的10项已实现KBM四级量规定义表见附录B（设计阶段规划的KBM-9"可视化与呈现"和KBM-12"反思与改进"见§4.1.1正文，待后续扩展纳入实现）。

---

## 4.2 数据库设计

### 4.2.1 关系型数据库设计（MySQL）

系统采用MySQL作为主关系型数据库，使用GORM（Go Object-Relational Mapping）框架进行数据操作。数据库设计围绕"教师-任务-学生-证据-评估结果"的核心业务链路展开。

> 【图4-2：ER图 — papers/figures/fig4-2-er-diagram.drawio】

**表4-3：核心数据表说明**

| 表名 | 主要字段 | 说明 |
|------|---------|------|
| users | id, name, email, password, role, created_at, updated_at | 用户基础表，role字段区分teacher/student |
| tasks | id, name, description, course_id, teacher_id(fk), start_date, end_date, status | 评价任务表 |
| student_tasks | id, task_id(fk), student_id(fk), status, progress | 任务-学生关联表 |
| evidences | id, student_task_id(fk), type, content, file_name, file_path, kbm_name, kbm_level | 证据表，存储文本证据内容或文件路径 |
| feedbacks | id, evidence_id(fk), content, kbm_level, strengths, weaknesses, suggestions | 单证据AI反馈表 |
| inference_results | id, student_id(fk), task_id(fk), student_task_id(fk), overall_score, overall_level, dimension_scores(JSON), reasoning | 推理结果表 |
| reports | id, student_id(fk), task_id(fk), overall_score, overall_level, dimension_scores(JSON), class_comparison(JSON), rank, percentile, strengths(JSON), weaknesses(JSON), suggestions(JSON), radar_chart_data(JSON) | 综合报告表 |

数据库实体关系图如图4-2所示。一条评价任务可以被分配给多名学生（通过student_tasks关联表），每名学生针对被分配的任务可以提交多条证据（evidences表），每条证据经过AI分析后生成对应的反馈记录（feedbacks表）。当教师或学生触发评估时，系统对该学生在某任务下的所有证据进行综合推理，生成推理结果（inference_results表）和综合报告（reports表）。

所有JSON类型的字段（如dimension_scores、class_comparison等）利用MySQL的JSON列类型存储，支持灵活的结构化查询。

### 4.2.2 图数据库设计（Neo4j）

系统引入Neo4j图数据库来存储和查询学生-维度-得分的图谱关系，为知识图谱可视化提供数据支持。图谱采用简单的节点-关系模型：

- **节点类型1：Student。** 属性包含studentID和name，代表一名学生。
- **节点类型2：Dimension。** 属性包含dimID和dimName（文献检索、研究设计、数据分析、批判思维），代表一个评价维度。
- **关系：SCORES。** 从Student节点指向Dimension节点，属性score存储该学生在该维度的得分。

> 【图4-3：Neo4j图谱模型 — papers/figures/fig4-3-neo4j-model.drawio】

典型的Cypher查询语句如下：

```cypher
MATCH (s:Student {id: $sid})-[r:SCORES]->(d:Dimension)
RETURN d.name, r.score
```

该查询返回指定学生在四个评价维度上的得分，前端可以利用返回数据绘制能力雷达图。

---

## 4.3 后端系统实现

### 4.3.1 项目结构设计与Go模块划分

后端系统采用Go语言开发，遵循Go社区公认的项目布局规范。工程入口位于`cmd/server/main.go`，核心业务代码位于`internal/`目录，公共工具位于`pkg/`目录。

```
research-ability-assessment/
├── cmd/server/main.go              # 应用入口：配置加载、数据库连接、服务/Agent/路由初始化
├── internal/
│   ├── agent/                      # 多智能体模块（核心）
│   │   ├── control_unit.go         # ControlUnit — 编排引擎
│   │   ├── evidence_agent.go       # EvidenceAgent — 证据预处理
│   │   ├── feedback_agent.go       # FeedbackAgent — LLM单证据评估
│   │   ├── inference_agent.go      # InferenceAgent — LLM多维推理
│   │   ├── logic_unit.go           # LogicUnit — 得分计算
│   │   └── storage_unit.go         # StorageUnit — 数据持久化
│   ├── config/config.go            # Viper+YAML配置管理
│   ├── handler/                    # HTTP处理器
│   ├── llm/                        # LLM客户端
│   │   ├── client.go               # OpenAI兼容API封装
│   │   └── parser.go               # 响应解析器
│   ├── middleware/                  # Gin中间件
│   ├── models/                     # GORM数据模型
│   ├── repository/                 # 数据访问层
│   └── service/                    # 业务逻辑层
├── pkg/                            # 公共包
│   ├── cache/redis.go              # Redis缓存
│   ├── extractor/                   # 多格式文件内容提取
│   │   ├── extractor.go             # ExtractorChain提取链
│   │   ├── docx.go                  # DOCX提取器
│   │   └── pdf.go                   # PDF提取器
│   └── utils/id_generator.go       # ID生成器
├── configs/config.dev.yaml         # 开发环境配置
├── scripts/                        # 辅助脚本
│   ├── init_db.go                  # 数据库种子数据
│   ├── concurrent_test.go          # 并发性能测试
│   └── batch_llm_classify.go       # LLM批量分类回归测试
├── testdata/                       # 测试数据（189文件）
│   ├── words/*.docx                # 81个DOCX测试文件
│   └── pdfs/*.pdf                  # 108个PDF测试文件
└── .github/workflows/test.yml      # CI自动测试工作流
```

> 【图4-4：main.go初始化流程图 — papers/figures/fig4-4-main-init.drawio】

main.go的启动流程按照严格有序的初始化管线执行：首先通过Viper加载YAML配置文件并展开环境变量；其次连接MySQL/PostgreSQL关系型数据库和Neo4j图数据库；然后通过GORM的AutoMigrate自动迁移数据表结构；接着依次初始化Repository层、LLM客户端、Service层、Agent层和Handler层；最后配置路由与中间件并启动HTTP服务。

这种有序的初始化管线设计确保了各模块之间的依赖关系清晰——下层模块（Repository）优先于上层模块（Service），上层模块在初始化时通过构造函数注入下层的依赖，形成单向依赖关系，有利于系统的测试和维护。

### 4.3.2 认证模块实现

系统采用JWT实现无状态的身份认证。认证流程分为注册、登录和请求验证三个阶段。

> 【图4-5：JWT认证流程图 — papers/figures/fig4-5-jwt-flow.drawio】

**注册。** 用户提交姓名、邮箱、密码和角色（teacher或student）。系统使用`golang.org/x/crypto/bcrypt`对明文密码进行哈希处理，哈希cost设置为10。哈希后的密码密文存储到MySQL的users表，系统中从不存储明文密码。注册成功后，系统直接生成并返回JWT令牌，用户无需二次登录。

**登录。** 用户提交邮箱和明文密码。系统从users表查询对应记录，使用`bcrypt.CompareHashAndPassword`将用户提交的密码与数据库中存储的bcrypt哈希值进行比对。验证通过后，系统生成JWT令牌，令牌payload中包含user_id和role两个claims，使用HMAC-SHA256算法签名，有效期为24小时。

**请求验证。** `AuthMiddleware`中间件从HTTP请求头`Authorization: Bearer <token>`中提取JWT令牌，调用`authService.ValidateToken`进行签名验证和过期检查。验证通过后，将user_id和role注入Gin框架的Context对象，供后续Handler使用。对于需要角色权限控制的接口（如只有教师可以创建任务），Handler从Context中读取role并判断。

### 4.3.3 证据管理与自动分类模块实现

证据管理是学生端与AI评估之间的桥梁，也是本文"有理有据"证据提取理念的第一个落地环节。证据创建时，若学生未手动指定KBM名称，系统自动调用EvidenceAgent的`ExtractKBMInfo`方法完成三项操作：文本预处理、LLM语义级KBM分类与等级评定（LLM不可用时降级为关键词匹配）、证据可信度评估。

**证据创建时的自动分类（Handler层）。** `CreateEvidence`接口的`kbm_name`字段为可选——前端可以留空。当证据创建请求中`kbm_name`为空且`content`不为空时，Handler调用`evidenceAgent.ExtractKBMInfo(content)`，自动完成KBM分类和等级预估：

```go
if evidence.KBMName == "" && evidence.Content != "" {
    info := h.evidenceAgent.ExtractKBMInfo(evidence.Content)
    evidence.KBMName = info.KBMName
    if evidence.KBMLevel == 0 {
        evidence.KBMLevel = info.Level
    }
}
```

**证据AI分析（Service层）。** 当学生或教师触发AI分析时——通过`POST /api/v1/evidences/:id/analyze`——系统调用`EvidenceService.AnalyzeEvidence`方法，将证据内容发送至DeepSeek LLM API进行语义层面的深度评估。LLM以"研究能力评估专家"的身份，基于System Prompt中嵌入的KBM四级量规标准，对证据文本进行阅读和分析，输出KBM级别和结构化反馈（优点、不足、建议、总体评价）。LLM响应经`llm.ParseFeedbackResponse`解析后，更新证据的kbm_level字段并创建feedback记录持久化。

FeedbackAgent构建的System Prompt定义了评价专家的角色和输出格式约束，User Prompt嵌入证据原文和量规标准。LLM返回的文本通过基于行的规则解析器提取KBM级别和结构化反馈字段，解析失败时使用默认值（KBMLevel=3）保证流程不中断。

> 【图4-6：FeedbackAgent工作流程图 — papers/figures/fig4-6-feedback-agent-flow.drawio】

#### EvidenceAgent：证据预处理与自动分类（核心实现）

EvidenceAgent是本文"有理有据"证据提取理念的核心实现载体。它采用LLM优先、规则兜底的双模策略，将原始学生文本转化为附带分类依据和可信度标签的结构化证据标注。

**PreprocessEvidence——文本标准化预处理。** 该方法对原始证据文本执行三步清洗操作：（1）去除首尾空白字符并合并非ASCII空格为单个空格；（2）统一换行符——丢弃`\r`，保留`\n`，压缩3个以上连续空行为2个；（3）按rune截断至2000字符，防止过长文本溢出LLM的上下文窗口。

**ClassifyEvidence——LLM优先KBM分类（含规则降级）。** 该方法是双模策略的直接体现。主路径`classifyWithLLM()`：将证据文本与10个KBM类别定义发送至LLM，由LLM进行语义级分类并输出包含`kbm_name`、`rationale`和`matched_keywords`的结构化JSON。LLM能够理解语义等价关系——例如识别出"使用PICO框架制定检索策略"和"在IEEE Xplore中用复合关键词检索"描述的是同一类KBM行为，即使两者用词完全不同。降级路径`classifyWithKeywords()`：当LLM不可用时，自动切换至关键词精确匹配——系统维护10组共120个领域关键词，统计每个KBM的命中数，取最高者作为分类结果。该降级路径确保在LLM服务不可用时，系统仍能基于确定性规则完成KBM分类。

**ExtractKBMInfo——LLM优先等级评定（含规则降级）。** 该方法对已分类的证据进行等级评定和可信度判断。主路径`extractKBMWithLLM()`：将证据文本与四级Rubrics量规发送至LLM，LLM对照量规给出等级（1-4）、可信度分数和详细评估理由（rationale），以JSON格式返回。降级路径`assessEvidenceWithRules()`：当LLM不可用时，采用三维评分因子——文本长度（1-4级）、关键词丰富度（0-10）和结构规范性（1-3）——综合计算等级（≥8→L4, ≥5→L3, ≥3→L2, <3→L1）和可信度分数（`Richness/6×0.5 + Structure/3×0.2 + Total/11×0.3`）。

**Rationale——结构化分类依据。** LLM主路径返回的rationale包含详细的语义评估理由，引用证据中的具体内容作为判断依据（如"证据内容展示了PICO框架的完整应用，检索流程清晰...因此评定为文献检索策略L4"）。规则降级路径生成格式如"分类为「文献检索策略」，等级评定为「良好」(L3)[规则引擎]。匹配关键词(4个): PICO, 检索式, 核心文献, 筛选标准。评分因子: 长度3 + 关键词丰富度4 + 结构2 = 综合9。可信度: 59%。"——rationale末尾的`[规则引擎]`标记使使用者能够区分该条评估是LLM还是规则生成。

**FeedbackAgent的LLM调用流程。** 当学生触发AI分析时，FeedbackAgent构建System Prompt（定义评价专家角色和输出格式约束）和User Prompt（嵌入证据原文和量规标准），调用`llmClient.Chat(ctx, messages)`将请求发送至DeepSeek API。LLM返回的结构化反馈经`llm.Parser`解析后，提取KBM等级、优点、不足、建议和总体评价，存储到MySQL的feedbacks表中。

### 4.3.4 多智能体系统核心实现

多智能体系统是本文的核心技术实现，由ControlUnit作为编排引擎，EvidenceAgent、InferenceAgent、FeedbackAgent、LogicUnit和StorageUnit五个专业Agent/单元协作完成从证据收集到评估报告输出的完整流水线。各Agent以依赖注入方式在ControlUnit的构造函数中组装，遵循单向依赖原则。

#### ControlUnit编排引擎

ControlUnit是系统的顶层编排器，聚合了任务管理服务、推理服务及三个Agent和一个StorageUnit的引用。其实际结构体定义如下：

```go
type ControlUnit struct {
    taskService      *service.TaskService
    inferenceService *service.InferenceService
    inferenceAgent   *InferenceAgent
    feedbackAgent    *FeedbackAgent
    storage          *StorageUnit
}
```

ControlUnit的`ExecuteEvaluation(ctx, taskID, studentID)`方法是完整评估流程的入口，该方法将任务状态管理、推理计算、结果持久化和反馈生成编排为一条5阶段流水线：

```go
func (c *ControlUnit) ExecuteEvaluation(ctx context.Context, taskID string, studentID string) (*EvaluationResult, error) {
    // Phase 1: 初始化任务状态
    task := &EvaluationTask{TaskID: taskID, StudentID: studentID, Progress: 0, Status: "processing"}
    c.taskService.UpdateStudentTaskStatus(ctx, task.TaskID, task.StudentID, "processing", task.Progress)

    // Phase 2: InferenceAgent执行维度推理与评分计算（进度 0→50→50）
    task.Progress = 50
    c.taskService.UpdateStudentTaskStatus(ctx, task.TaskID, task.StudentID, "processing", task.Progress)
    result, err := c.inferenceAgent.InferAbility(ctx, task.StudentID, task.TaskID)
    if err != nil { return nil, fmt.Errorf("能力推理失败: %w", err) }

    // Phase 3: StorageUnit持久化结果到MySQL并更新Neo4j图谱（进度 50→70）
    inferenceResult := &models.InferenceResult{...}
    c.storage.StoreInferenceResult(ctx, inferenceResult)
    for dimension, score := range result.DimensionScores {
        c.storage.UpdateKnowledgeGraph(ctx, task.StudentID, dimension, score.Score)
    }

    // Phase 4: FeedbackAgent生成诊断性反馈（进度 70→70）
    task.Progress = 70
    c.taskService.UpdateStudentTaskStatus(ctx, task.TaskID, task.StudentID, "processing", task.Progress)
    feedbackResult, _ := c.generateFeedback(ctx, inferenceResult)

    // Phase 5: 标记评估完成（进度 70→100）
    task.Progress = 100; task.Status = "completed"
    c.taskService.UpdateStudentTaskStatus(ctx, task.TaskID, task.StudentID, "completed", task.Progress)

    return &EvaluationResult{...}, nil
}
```

> 【图4-7：ControlUnit.ExecuteEvaluation流程图 — papers/figures/fig4-7-evaluate-task-flow.drawio】

#### EvidenceAgent实现

EvidenceAgent负责证据的收集、预处理和LLM优先的自动分类。其`CollectEvidence(ctx, studentID, taskID)`方法通过`EvidenceService`查询指定学生在指定任务下的所有已提交证据记录（evidence表），返回按KBM名称标记的结构化证据列表。对于新提交的证据，其`ExtractKBMInfo(content)`方法自动完成LLM优先的KBM分类与等级评定（LLM不可用时降级为关键词匹配+规则评分）——详见§4.3.3 EvidenceAgent实现。`AutoClassifyAndUpdate(ctx, evidence)`方法对数据库中尚未标注KBM的证据进行批量自动分类并回写更新。这种设计将证据的初始KBM标注自动化，减少了学生手动标注的负担，同时保留了手动覆盖的灵活性。

#### InferenceAgent实现

InferenceAgent是系统的核心推理引擎，负责将过程性证据转化为四维度的量化评分。采用LLM优先、规则兜底的双模策略。其结构体聚合了EvidenceAgent、LogicUnit和LLM Client三个底层组件：

```go
type InferenceAgent struct {
    evidenceAgent *EvidenceAgent
    logicUnit     *LogicUnit
    llmClient     *llm.Client
}
```

InferenceAgent的`InferAbility`方法首先尝试LLM路径，LLM不可用时自动降级至规则计算。具体推理流程如下：

**第一步——证据收集与维度分组。** 调用`evidenceAgent.CollectEvidence`获取该学生在任务下的全部证据。然后通过`groupEvidencesByDimension`方法，依据每条证据的KBM名称，将所有证据分配到对应的四个评价维度组中。

**第二步——LLM主导多维推理（主路径）。** 调用`inferWithLLM()`方法，构建包含完整Rubrics量规的System Prompt和按维度分组的证据User Prompt，将请求发送至LLM。LLM在每个维度对照量规进行语义级评分（0-100分），并引用证据原文作为判据（evidence_quotes）。LLM返回结构化JSON，包含各维度的score、level、reasoning和evidence_quotes字段。系统自动解析LLM响应，从evidence_quotes提取证据引用生成EvidenceCitation记录（详见§4.3.6双模架构设计决策），构建"评估结论→维度→证据→原文段落"的完整溯源链。

**第三步——LogicUnit规则降级（备用路径）。** LLM调用失败（网络超时、API错误或JSON解析失败）时，自动回退至LogicUnit的确定性计算：KBMLevel × 20 → 维度内算术平均 → 四维度等权重（各0.25）加权 → GetLevelFromScore等级判定。该降级路径确保在任何条件下系统都能完成评分，且计算过程完全透明可复现。

**第四步——等级映射与推理依据生成。** 无论走LLM路径还是规则路径，最终由`GetLevelFromScore`函数将综合得分映射为等级标签（≥90优秀、≥75良好、≥60合格、<60不合格）。`generateReasoning`方法基于推理结果生成结构化推理依据文本。

**设计决策说明。** InferenceAgent选择LLM优先策略的核心考量在于语义理解深度：LLM能够理解"使用PICO框架构建检索策略"与"使用百度搜索关键词"之间的本质差异——前者属于L4级文献检索，后者仅为L1-L2级——这种基于研究规范知识的语义判断是规则映射（KBMLevel×20）无法做到的。同时，系统通过规则降级路径（LogicUnit）确保LLM不可用时的可用性——降级路径的公式计算完全透明、可复现、无随机性，在确保基础服务质量的同时为LLM的语义理解优势提供发挥空间。

#### FeedbackAgent实现

FeedbackAgent负责生成面向学生和教师的诊断性反馈报告。它在ControlUnit的评估流程中位于推理评分之后——即先完成定量评分，再基于评分结果生成定性的改进建议。其结构体持有LLM客户端的引用：

```go
type FeedbackAgent struct {
    llmClient *llm.Client
}
```

FeedbackAgent的`GenerateFeedback`方法采用"LLM优先+规则兜底"的双轨策略：

**LLM生成路径。** 构建包含学生各维度得分、等级和推理依据的提示词，调用`llmClient.Chat`向DeepSeek API请求生成针对性的反馈文本。LLM需要输出每个维度的改进建议（ImprovementSuggestion），包含当前得分、目标得分、具体建议和可操作的行动项（ActionItems）。提示词中明确要求LLM基于"得分差距分析"的逻辑生成反馈——差距越大，建议越具体和可操作。

**规则兜底路径。** 当LLM调用失败（网络超时、API错误或返回格式异常）时，`generateRuleFeedback`方法根据各维度的得分区间（优秀/良好/合格/不合格），从预定义的规则模板库中选取对应的反馈文本。例如，对低于60分的维度自动生成"建议从基础文献检索和实验设计方法论学习入手"等通用建议。该设计确保了系统在任何条件下都能返回有意义的反馈输出，而非因LLM调用失败而中断评估流程。

**设计决策说明。** FeedbackAgent的生成结果不作为推理评分环节的输入——即InferenceAgent不从FeedbackAgent获取中间分析结果。这一设计的核心考量在于角色分离：推理评分负责"从证据到分数的客观映射"，反馈生成负责"从分数到建议的诊断分析"。两者在信息流上形成串行而非循环依赖，避免了LLM生成文本的语义噪声反哺评分环节，保持了评分逻辑的独立性和稳定性。

#### LogicUnit实现

LogicUnit是纯计算组件，不依赖任何外部服务或LLM，在LLM不可用时充当降级评分路径，提供三个核心计算方法：

**证据评分（EvaluateEvidence）。** 将单条证据的KBM级别（1-4）按`Level × 20`映射为百分制分数，上限为100分。无KBM级别时默认返回60分。该方法仅在LLM分类失败后的降级路径中使用。

**维度得分（CalculateDimensionScore）。** 计算一个维度下所有证据分数的算术平均值。如果某维度下没有有效证据，返回默认值50分。

**综合得分（CalculateOverallScore）。** 读取每个维度的权重（当前四维度均配置为0.25），计算加权平均值（round保留两位小数）。

**等级判定（GetLevelFromScore）。** 在`models`包中定义，基于以下阈值将数值分数转换为等级标签：≥90分为"优秀"、≥75分为"良好"、≥60分为"合格"、<60分为"不合格"。

#### StorageUnit实现

StorageUnit封装了所有持久化操作，聚合了MySQL的ResultRepo（注：Go包名沿用"postgres"，实际后端使用MySQL）和Neo4j的GraphRepo：

```go
type StorageUnit struct {
    resultRepo *postgres.ResultRepo
    graphRepo  *neo4j.GraphRepo
}
```

`StoreInferenceResult`方法将推理结果（学生ID、任务ID、综合得分、等级、维度得分JSON、推理依据）写入PostgreSQL的inference_results表。`UpdateKnowledgeGraph`方法采用MERGE操作确保Neo4j中的Student节点和Dimension节点在首次写入时自动创建，随后通过`HAS_SCORE`关系记录该维度的最新得分。这种"关系型数据库做主存储+图数据库做可视化加速"的双存储架构设计，使评估结果的批量查询和分析报表生成利用MySQL的索引能力，而前端雷达图等可视化需求的即时查询则利用Neo4j的图遍历效率。

### 4.3.5 LLM集成实现

LLM客户端`llm.Client`的设计目标是实现与OpenAI Chat Completions API标准的兼容，以支持灵活的Provider切换。

**接口兼容性。** 客户端实现了OpenAI API的标准请求和响应格式。请求体包含model、messages（role+content数组）、max_tokens和temperature等标准字段。这种设计允许系统在DeepSeek、OpenAI GPT-4、通义千问等任何兼容OpenAI格式的LLM服务之间无缝切换——只需修改配置文件中的base_url和api_key即可。

**超时控制。** HTTP客户端设置了120秒的连接和读取超时。这一超时值的设定基于对LLM API典型响应时间的分析：批量分类测试（189文件）中DeepSeek API的平均响应时间约为2.03秒/次，单证据分析（含HTTP往返和提示词构建开销）约12.6秒。120秒的超时窗口为网络波动和LLM服务高负载等最坏情况留出了充足的缓冲。

**响应解析。** LLM的输出具有不确定性——即使明确要求JSON格式输出，实际响应中仍可能夹杂Markdown代码块标记、解释性前缀或后缀文本。系统采用两层解析策略：（1）`extractJSON()`函数（位于`inference_agent.go`）通过括号深度匹配在响应文本中定位JSON对象的起止位置（`{`到`}`），提取后进行JSON反序列化——该方法用于InferenceAgent和EvidenceAgent的LLM JSON响应解析；（2）`ParseFeedbackResponse()`方法（位于`llm/parser.go`）采用基于行的关键词匹配策略——逐行扫描响应文本，通过识别中文关键词（"优点""不足""建议"等）将响应内容分配到对应字段，用于FeedbackAgent的反馈文本解析。如果解析失败，解析器返回明确的错误信息或默认值，由上层Agent决定是重试还是返回降级评估结果。

### 4.3.6 关键设计决策讨论

本节对系统实现中的三项关键设计决策进行专题讨论，阐明选择依据和工程考量。

**Neo4j图数据库的引入合理性。** 系统引入Neo4j的图谱模型是简单的`(:Student)-[:HAS_SCORE]->(:Dimension)`二类节点结构，从数据存储功能而言，该结构确实可以通过PostgreSQL的JSON字段或键值对表实现等效存储。选择引入Neo4j的主要考量包括以下三点。其一，**查询语义的直观性**：评分关系天然构成一个以学生为中心的星型图谱，使用Cypher的`MATCH (s:Student)-[r:HAS_SCORE]->(d:Dimension)`查询模式直接对应"获取某学生各维度得分"的业务语义，比关系型数据库的多表JOIN或JSON展开更直观。其二，**可视化集成的便捷性**：Neo4j Browser和前端ECharts雷达图的对接可以直接利用图谱查询返回的节点-关系数据，避免了"关系型查询→中间转换→图谱结构"的额外序列化开销。其三，**未来图谱扩展的预留空间**：当前简单的图谱模型为后续扩展留下了明确的路径——例如引入表示"课程→任务→证据"的层级节点以构建完整的评估知识图谱，或引入学生之间的协作关系边以支持团队评价场景。这种前瞻性的架构预留避免了后期"先全用MySQL、再迁移到图数据库"的高成本重构。需要指出的是，在当前系统规模下，Neo4j是可选而非必需的组件——如果仅需雷达图展示功能，可以完全通过PostgreSQL的维度得分JSON字段实现。

**四维度等权重设计的考量。** 系统将文献检索与综述、研究设计与实验、数据分析与解释、批判性思维与创新四个维度配置为各0.25的等权重计算综合得分。这一选择并非默认假定四个维度在研究能力中具有同等重要性，而是基于以下工程考量。其一，**基线公平性**：作为本科研究能力评价的初步探索，等权重提供了最中立的基线——在缺乏大规模实证数据支持差异化权重之前，不偏向某一特定维度是最保守、最不易引入系统性偏差的选择。其二，**可配置性**：系统的`GetDimensionWeight(dimID)`函数并非硬编码等权重，而是通过查询`models.DefaultDimensions`中每个维度的Weight字段获取。这意味着未来在积累足够评价数据后，可以基于因子分析或层次分析法等方法学实证调整各维度权重，系统代码无需修改，仅需更新模型配置。其三，**教育导向的可解释性**：等权重设计向学生传递了清晰的教育导向——四个维度的能力同等重要，不应偏废。在本科阶段这一能力发展的奠基期，避免引导学生只关注"高分维度"而忽视其他方面。

**LLM主导评分与规则引擎降级的双模架构。** 系统采用"LLM优先、规则兜底"的双模策略处理评分与反馈。InferenceAgent的`InferAbility()`方法首先尝试LLM路径——将完整证据上下文和四级Rubrics量规发送至LLM，要求LLM逐维度评分并引用证据原文作为判据。LLM返回的结构化JSON中包含各维度的score、level、reasoning和evidence_quotes字段。系统自动解析LLM响应，提取evidence_quotes生成EvidenceCitation记录存入数据库。当LLM调用失败（网络超时、API错误、JSON解析失败）时，系统自动降级为LogicUnit的确定性公式（KBMLevel×20→维度平均→等权重加权）完成规则化评分。这一设计的核心考量如下。其一，**语义理解的深度**：LLM能够理解"使用了PICO框架构建检索策略"与"使用了百度搜索关键词"之间的本质差异——前者属于L4级文献检索，后者仅为L1-L2级——这种基于研究规范知识的语义判断是关键词匹配无法做到的。其二，**可追溯性的兼顾**：LLM路径要求在每个维度的reasoning中引用证据原文的具体段落（evidence_quotes），使评分有据可查；规则降级路径同样提供了透明的计算链条（KBMLevel×20→平均→加权）。其三，**可用性的保障**：双模设计确保系统在任何情况下都能完成评估——LLM可用时发挥其语义理解优势，LLM不可用时规则引擎保证基础服务质量。

### 4.3.7 RESTful API接口设计

系统采用RESTful设计风格，所有API响应统一使用以下JSON格式：

```json
{"code": 200, "message": "操作成功", "data": {...}}
```

**表4-4：核心API端点汇总（完整43条路由见附录C）**

| 方法 | 路径 | 功能 | 认证要求 |
|------|------|------|---------|
| POST | /api/v1/auth/register | 用户注册 | 否 |
| POST | /api/v1/auth/login | 用户登录 | 否 |
| GET | /api/v1/user/info | 获取当前用户信息 | 是 |
| POST | /api/v1/tasks | 创建任务 | 是（教师） |
| GET | /api/v1/tasks | 获取任务列表 | 是（教师） |
| GET | /api/v1/tasks/:task_id | 获取任务详情 | 是 |
| POST | /api/v1/tasks/:task_id/assign | 分配任务给学生 | 是（教师） |
| POST | /api/v1/evidences | 创建文本证据 | 是（学生） |
| POST | /api/v1/evidences/upload | 上传文件证据 | 是（学生） |
| POST | /api/v1/evidences/:id/analyze | AI分析单条证据 | 是 |
| GET | /api/v1/evidences/:id/feedback | 查看证据AI反馈 | 是 |
| POST | /api/v1/results/generate/student | 生成推理结果 | 是 |
| GET | /api/v1/results/report/student | 生成学生综合报告 | 是 |
| GET | /api/v1/results/report/task/:task_id | 生成任务汇总报告 | 是（教师） |
| POST | /api/v1/agent/evaluate | 直接调用Agent评估 | 是 |
| GET | /api/v1/graph/student/:student_id | 获取知识图谱数据 | 是 |

API接口的路由配置在`cmd/server/main.go`的`setupRouter`函数中集中管理。所有受保护的接口（标注"是"的接口）均需要在请求头中携带有效的JWT令牌。标注"是（教师）"的接口额外校验当前用户角色为teacher，学生角色的请求将被拒绝。

---

## 4.4 前端系统实现

### 4.4.1 前端项目结构与技术选型

前端采用Vue 3 + TypeScript + Vite 4构建，技术选型服务于"数据展示+管理操作"的核心交互模式：

- **Element Plus。** 作为UI组件库，提供成熟的表格、表单、对话框、导航、标签页等后台管理组件，保证了界面风格的一致性和交互的标准化。
- **ECharts 6。** 作为数据可视化引擎，用于绘制能力雷达图和班级对比柱状图等核心可视化图表。
- **Vue Router 4。** 管理单页面应用（SPA）的路由跳转和页面切换。
- **Axios。** 封装HTTP请求，在`api/auth.ts`中创建单一axios实例并配置请求拦截器，统一注入Bearer Token：

```typescript
const api = axios.create({ baseURL: '/api/v1' })
api.interceptors.request.use(config => {
  const token = localStorage.getItem('token')
  if (token) config.headers.Authorization = `Bearer ${token}`
  return config
})
```

所有业务API模块（task.ts、evidence.ts、result.ts）均通过`import api from './auth'`复用该实例，保证了Token注入逻辑的一致性，避免了重复配置。

### 4.4.2 路由设计与权限控制

> 【图4-8：前端路由结构图 — papers/figures/fig4-8-frontend-routes.drawio】

前端路由采用嵌套结构，根路径`/`重定向到`/dashboard/tasks`（任务管理页面）。整个应用的路由层级如下：

- `/login` → Login.vue — 登录页面
- `/register` → Register.vue — 注册页面
- `/dashboard` → Dashboard.vue — 包含侧边导航和顶部信息栏的布局容器
  - `/dashboard/tasks` → TaskManagement.vue — 任务创建、分配和管理
  - `/dashboard/evidences` → EvidenceManagement.vue — 证据提交和AI分析
  - `/dashboard/results` → ResultManagement.vue — 评估结果查看
  - `/dashboard/reports` → ReportManagement.vue — 综合报告查看

权限控制在前端通过Vue Router的导航守卫实现。在进入每个受保护路由前，守卫检查localStorage中是否存在有效的JWT令牌。在Dashboard组件的侧边导航中，根据从JWT令牌中解析出的role字段动态渲染不同的菜单项——教师端显示"任务管理""结果管理""报告管理"，学生端显示"证据管理""结果管理""报告管理"。

### 4.4.3 核心页面实现

**证据管理页面（EvidenceManagement.vue）。** 学生端使用频率最高的页面。页面主体包含两部分：上半部分为证据提交区——包含文本输入框（Element Plus的el-input textarea组件）和文件上传按钮（el-upload组件）；下半部分为已提交证据列表——使用el-table展示证据的摘要信息，每条证据提供"AI分析"操作按钮。点击"AI分析"后，前端发送`POST /api/v1/evidences/:id/analyze`请求，后端触发FeedbackAgent的评估流程。分析完成后，证据状态标签从"待分析"变为"已分析"，点击状态标签弹出对话框，展示FeedbackAgent生成的详细反馈内容——包括KBM评估级别、优点、不足、改进建议和总体评价。

**结果管理页面（ResultManagement.vue）。** 展示学生的评估推理结果列表。每条结果以el-card卡片组件呈现，卡片顶部显示学生姓名和任务名称，中部显示总体得分（大号数字+颜色编码——绿色≥90、蓝色75-89、橙色60-74、红色<60）和对应的等级徽章，下部以水平柱状条的形式展示四个维度的得分。点击卡片进入详情弹窗，展示完整的推理链文本和维度得分表格。

**报告管理页面（ReportManagement.vue）。** 提供综合报告的生成和查看功能。报告生成按钮触发`GET /api/v1/results/report/student`请求。报告以el-dialog弹窗呈现，内部分为四个标签页：（1）综合评价——总分、等级、班级排名和百分位；（2）能力雷达图——ECharts渲染的四维度雷达图；（3）详细分析——优势维度和待提升方向；（4）改进建议——个性化发展建议列表。

**结果管理页面中雷达图初始化的关键实现**（摘自ResultManagement.vue，`<script setup lang="ts">`）：

```typescript
import * as echarts from 'echarts'
import { generateRadarData } from '../utils/chart'

let radarChart: echarts.ECharts | null = null

// 对话框打开时延迟初始化（等待DOM渲染），关闭时销毁
watch(showResultDialog, (newVal) => {
  if (newVal) {
    setTimeout(() => {
      radarChart = echarts.init(radarChartRef.value!)
      const chartData = generateRadarData(resultDetails.value.dimension_scores)
      radarChart.setOption({
        radar: chartData.radar,
        series: [{
          type: 'radar', data: chartData.series[0].data,
          areaStyle: { opacity: 0.2 }, lineStyle: { width: 2 },
          itemStyle: { color: '#409EFF' }
        }]
      })
    }, 100)
  } else {
    radarChart?.dispose()
  }
})
```

ECharts实例的生命周期与对话框的visible状态绑定——仅在对话框打开时初始化、关闭时dispose销毁，避免图表实例的内存泄漏和未渲染DOM上的初始化错误。

**登录页的演示模式降级逻辑**（摘自Login.vue）：

```typescript
const handleLogin = async () => {
  try {
    const res = await login(loginForm.email, loginForm.password)
    localStorage.setItem('token', res.data.token)
    router.push('/dashboard')
  } catch (err) {
    // 后端不可达时自动降级到演示模式
    ElMessage.warning('后端服务未启动，切换到演示模式')
    localStorage.setItem('token', 'mock-token-' + Date.now())
    localStorage.setItem('user', JSON.stringify({name:'演示用户', role:'teacher'}))
    router.push('/dashboard')
  }
}
```

当后端API调用因网络不可达或服务未启动而失败时，前端自动切换到Mock Token模式——在localStorage中写入模拟令牌和用户信息，后续所有页面基于内嵌的模拟数据独立渲染，系统界面和交互流程的展示不依赖后端服务的可用性。

### 4.4.4 数据可视化与图表工具

系统使用ECharts 6实现能力雷达图、班级对比柱状图和等级分布饼图。图表数据生成逻辑封装在`utils/chart.ts`中：

```typescript
// 雷达图数据生成
export function generateRadarData(scores: Record<string, number>) {
  const dimensions = ['文献检索', '研究设计', '数据分析', '批判思维']
  const maxVal = Math.max(...Object.values(scores), 1)
  // 自动归一化：若所有分数≤1则视为比例，直接使用；否则归一化到[0,1]
  const normalized = maxVal <= 1 ? scores :
    Object.fromEntries(Object.entries(scores).map(([k, v]) => [k, v / 100]))
  return {
    radar: { indicator: dimensions.map(d => ({ name: d, max: 1 })) },
    series: [{ data: [{ value: dimensions.map(d => normalized[d] || 0) }] }]
  }
}
```

雷达图的四个坐标轴对应四个评价维度，自动检测分数区间进行归一化——100分制数据自动缩放到[0,1]区间供ECharts雷达坐标系使用。图表在`dialog`打开时才初始化、关闭时dispose销毁，严格管理内存。

### 4.4.5 路由守卫与前端类型系统

**路由守卫**（摘自`router/index.ts`）在每次导航前检查localStorage中的JWT令牌：

```typescript
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  if (to.path !== '/login' && to.path !== '/register' && !token) {
    next('/login')  // 未登录重定向
  } else {
    next()
  }
})
```

前端项目未使用Vuex/Pinia等状态管理库，而是通过localStorage存储token和user信息、组件内部的`ref()`管理页面状态，整体状态管理简洁轻量。

**TypeScript类型定义**为所有API交互提供编译时类型安全。以推理结果类型为例（摘自`types/result.ts`）：

```typescript
interface DimensionScore {
  name: string; score: number; level: string
  details: string; evidence_ids: string[]
}
interface InferenceResult {
  id: string; student_name: string; task_name: string
  overall_score: number; overall_level: string
  dimension_scores: Record<string, DimensionScore>
  reasoning: string; created_at: string
}
```

全量类型定义覆盖了User、Task、Evidence、Feedback、InferenceResult和Report六个核心领域，总计约120行TypeScript接口定义，确保前后端数据契约的严格一致。
