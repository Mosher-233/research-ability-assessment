# 论文写作指引 — Handoff Package for academic-paper-composer

> 由 academic-paper-strategist 基于项目实际代码与论文框架生成  
> 交付对象：academic-paper-composer  
> 目标：北京邮电大学 信息与通信工程学院 电子信息工程专业 本科毕业论文  
> 格式权威：北京邮电大学2026届本科毕业设计（论文）撰写指导手册

---

## 一、约束条件记录

| 项目 | 约束 |
|------|------|
| 学校 | 北京邮电大学 信息与通信工程学院 |
| 专业 | 电子信息工程 |
| 论文类型 | 工程实践类 |
| 题目 | 大学生研究能力评价AI Agent的研究与实现 |
| 正文最低页数 | ≥30页 |
| 参考文献数量 | ≥20篇（近五年为主） |
| 查重率 | 目标<10%，最高<25% |
| AIGC率 | 可控 |
| 框架文件 | `papers/论文完整框架.md` |
| 图表目录 | `papers/figures/`（18张.drawio） |
| 排版格式 | 对齐学校手册 |
| 最终格式 | .docx |

---

## 二、证据源映射表

### 2.1 代码证据源

| 证据类型 | 文件路径 | 用途章节 |
|---------|---------|---------|
| 项目入口 | `cmd/server/main.go` | §4.3.1 |
| Agent实现 | `internal/agent/control_unit.go` | §4.3.4 |
| Agent实现 | `internal/agent/evidence_agent.go` | §4.3.4, §5.4.4 |
| Agent实现 | `internal/agent/feedback_agent.go` | §4.3.3, §4.3.4 |
| Agent实现 | `internal/agent/inference_agent.go` | §4.3.4 |
| Agent实现 | `internal/agent/logic_unit.go` | §4.3.4 |
| Agent实现 | `internal/agent/storage_unit.go` | §4.3.4 |
| LLM客户端 | `internal/llm/client.go` | §4.3.5 |
| 文件提取器 | `pkg/extractor/docx.go` | §4.3.3, §5.2.2 |
| 文件提取器 | `pkg/extractor/pdf.go` | §4.3.3, §5.2.2 |
| 文件提取器 | `pkg/extractor/extractor.go` | §4.3.3, §5.2.2 |
| 认证服务 | `internal/service/auth_service.go` | §4.3.2 |
| 证据服务 | `internal/service/evidence_service.go` | §4.3.3 |
| 推理服务 | `internal/service/inference_service.go` | §4.3.4 |
| 报告服务 | `internal/service/report_service.go` | §4.3.4 |
| 任务服务 | `internal/service/task_service.go` | §4.3.4 |
| 中间件 | `internal/middleware/auth.go` | §4.3.2 |
| 数据模型 | `internal/models/evidence.go` | §4.2.1, §4.2.1 (evidence_citations) |
| 数据模型 | `internal/models/citation.go` | §4.2.1 |
| 数据模型 | `internal/models/dimension.go` | §4.2.1 |
| 配置管理 | `internal/config/config.go` | §4.3.1 |
| Repository | `internal/repository/postgres/` | §4.3.7 |
| Repository | `internal/repository/neo4j/` | §4.3.7 |
| Redis缓存 | `pkg/cache/redis.go` | §4.3.6 |
| 前端路由 | `frontend/src/router/index.ts` | §4.4.2 |
| 前端页面 | `frontend/src/views/*.vue` | §4.4.3 |
| 前端API | `frontend/src/api/*.ts` | §4.4.1 |
| 前端类型 | `frontend/src/types/*.ts` | §4.4.1 |
| 测试数据 | `scripts/init_db.go` | §5.1, §5.5 |
| 并发测试 | `scripts/concurrent_test.go` | §5.3 |
| **批量LLM测试** | **`scripts/batch_llm_classify.go`** | **§5.3.3, §5.4.4** |
| **PDF提取测试** | **`pkg/extractor/pdf_test.go`** | **§5.2.2** |
| **测试数据(DOCX)** | **`testdata/words/*.docx` (81文件)** | **§5.2.2, §5.4.4** |
| **测试数据(PDF)** | **`testdata/pdfs/*.pdf` (108文件)** | **§5.2.2, §5.4.4** |
| **CI工作流** | **`.github/workflows/test.yml`** | **§5.1** |
| 项目文档 | `README.md` | §1.1, §3.2.3 |

### 2.2 外部证据源

| 证据类型 | 来源 | 用途章节 |
|---------|------|---------|
| 任务书 | `papers/任务书.md` | §1.3, §4.1 |
| 开题报告 | `papers/开题报告.md` | §1.1-§1.3, §2.1 |
| 问题与解决方案 | `papers/穆方达-问题与解决方案.docx` | §4.1.2, §5.4 |
| 能力清单 | `papers/ai时代学生工程思维.md` | §1.1, §2.1 |
| 格式手册 | `papers/北京邮电大学...docx` | 排版规则 |

---

## 三、逐章节写作指引

---

### 摘要（500-800字中文 + 300-500词英文）

**允许的声明（基于证据）：**
- 项目名称为"大学生研究能力评价AI Agent"，面向电子信息工程专业
- 技术栈：Go+Gin+Vue3+MySQL+Neo4j+DeepSeek LLM
- Agent架构：ControlUnit + EvidenceAgent + FeedbackAgent + InferenceAgent + LogicUnit + StorageUnit
- 评价指标体系：4维度 × 12 KBM，对齐工程教育认证标准(2024版)和CDIO
- 实验数据：小样本AI-教师对比（5学生35证据）+ 大样本LLM批量验证（189文件，81 DOCX + 108 PDF，覆盖A-H八个质量等级）
- AI与3名教师评分皮尔逊相关系数r=0.87（Cohen's Kappa=0.73）
- LLM异常内容检测率98%（49/50），等级评定±1以内率100%，零高估
- 189文件多格式提取（DOCX+PDF）成功率100%

**禁止的声明：**
- 不声称"首次提出"任何理论
- 不声称"大规模部署"
- 不声称"真实邮百工平台数据"（当前为仿真数据）

**写作要点：**
- 严格按"背景-问题-方法-结果-结论"五段式
- 英文摘要注意：ECD, MAS, LLM, KBM, CDIO 首次出现给出全称
- 关键词7个以内

---

### 第一章 绪论（6-8页）

#### §1.1 研究背景与意义

**§1.1.1 EE专业人才培养挑战 — 允许的声明：**
- 引用《工程教育认证标准(2024版)》，可提及12项毕业要求中的"问题分析""研究"等条款
- EE专业核心课程场景：信号处理、通信系统、嵌入式系统、电子电路设计等
- 可适当描述：一名EE学生典型项目过程 → 文献检索(IEEE Xplore) → 方案设计 → 实验实施 → 数据分析 → 报告撰写

**§1.1.2 传统评价局限性 — 允许的声明：**
- 当前依赖静态成果（实验报告、论文）的普遍现象（基于教育评价文献，不需代码证据）
- 教师工作量问题：EI专业实践教学环节1师配N生的常见比例

**§1.1.3 AI带来新范式 — 允许的声明：**
- 大语言模型（LLM）在文本理解推理方面的能力（引用GPT-3/DeepSeek文献）
- MAS的分布式协作优势（引用Wooldridge）
- ECD框架（引用Mislevy）与AI的结合方向

**§1.1.4 研究意义 — 允许的声明：**
- 理论：ECD+MAS+LLM三层融合
- 实践：服务于高校创新创业平台

**禁止的声明：**
- 不能声称"填补空白"或"首次"
- 不能声称AI完全替代教师

**证据链接：** README.md项目描述、开题报告背景论述、参考文献[1][2][9][10][11][12]

---

#### §1.2 国内外研究现状

写作要点：
- 每个小节约400-600字，简洁综述
- 每一段以一个代表性文献为核心展开
- §1.2.5 必须明确指出3个不足（对齐框架中的三个切入点）

**引用分配（必须真实使用，不可虚构页码/卷号）：**
- [1]-[2]：工程教育认证与CDIO → §1.2.1
- [11]-[12]：MAS与ECD → §1.2.3
- [9]-[22]：LLM与CoT → §1.2.4
- [20]-[21]：AI评价系统 → §1.2.2

---

#### §1.3 研究内容与创新点

写作要点：
- 研究内容严格对齐任务书的四个任务
- 创新点表述保守："提出了...框架"而非"首创"
- 技术路线图引用：`【图1-1：技术路线图，见figures/fig1-1-tech-roadmap.drawio】`

---

### 第二章 相关理论与技术基础（6-8页）

#### §2.1 工程教育中的研究能力评价理论

**§2.1.1 — 允许的声明：**
- 引用工程教育认证标准具体条款编号
- CDIO Syllabus 2.0的能力维度编号（2.1-2.4）

**§2.1.2 ECD三层模型 — 引用图：** `fig2-1-ecd-model.drawio`

**§2.1.3 — 允许的声明：**
- 四个维度的具体定义
- 与EE专业的场景关联（IEEE Xplore、MATLAB、电路仿真等）

**§2.1.4 KBM概念 — 允许的声明：**
- KBM完整的英文全称和定义
- 四级量规的Bloom理论依据

**禁止的声明：**
- 不要在本节展开12个KBM全部定义（这是§4.1的内容）
- 不要在此处引入任何代码实现细节

---

#### §2.2 多智能体系统理论

- 以Wooldridge[11]的Agent定义为核心
- 解释主控-从属架构模式（为本系统的ControlUnit编排做铺垫）
- 对比黑板架构和联邦架构（一两段即可）

**证据链接：** `internal/agent/` 的包结构体现了MAS的设计思想

---

#### §2.3 大语言模型与提示词工程

**§2.3.1 Transformer — 允许的声明：**
- 注意力机制的概念描述（水平相当即可，不需要数学公式）
- 引用 Vaswani 等[26]

**§2.3.2 — 允许的声明：**
- Few-shot/Zero-shot的概念
- Chain-of-Thought的原理（引用Wei等[22]）
- 本文中FeedbackAgent采用CoT风格提示词

**证据链接：** `internal/agent/feedback_agent.go` 中的提示词模板结构

---

### 第三章 系统需求分析与总体设计（6-8页）

#### §3.1 系统需求分析

- 功能性需求采用用例描述方式（教师做什么、学生做什么）
- 非功能性需求精确量化：
  - 单生评估<30s ✓（基于DeepSeek API 2.8s/次调用）
  - 支持≥10教师+50学生并发 ✓（基于`concurrent_test.go`）
  - API响应<2s（不含LLM）✓（基于性能测试数据）

**证据链接：** `README.md` 功能特性描述、`cmd/server/main.go` API路由

---

#### §3.2 系统总体架构设计

**引用图：** `fig3-1-system-architecture.drawio`（五层架构图）

**技术选型表必须声明：**

| 选项 | 选择 | 理由（必须基于实际，不可编造） |
|------|------|------|
| 后端语言 | Go | 并发模型(goroutine)适合Agent并行任务；Gin框架性能好 |
| 数据库 | MySQL + Neo4j | 关系型存结构化数据，图数据库存图谱关系 |
| LLM | DeepSeek | API兼容OpenAI格式，中文能力强，性价比高 |
| 前端 | Vue3 | Element Plus组件生态完善 |
| 缓存 | Redis | 降低重复LLM调用成本 |

**禁止的声明：**
- 不编造与其他方案的量化对比数据
- 不声称"最佳方案"

---

#### §3.3 多智能体系统架构设计

**引用图：** `fig3-2-agent-architecture.drawio`（Agent组件图）
**引用图：** `fig3-3-agent-sequence.drawio`（序列图）
**引用图：** `fig3-4-data-flow.drawio`（数据流图）

**Agent职责表：**
每个Agent的职责描述必须与 `internal/agent/` 下实际代码一致。

---

### 第四章 系统详细设计与实现（10-14页）

#### §4.1 评价指标体系与量规设计

**§4.1.1 — 完整列出4维度12KBM**

**§4.1.2 — KBM-标准条款-理论依据三重映射表**
- 12行×6列的完整表格
- 此处应直接引用任务书中的毕业要求指标点编号（3.1, 3.3, 4.1, 4.3, 5.1, 6.2, 10.1, 10.2, 10.3, 12.1）

**§4.1.3 — 以KBM-5为例展示四级量规**
- 等级4/3/2/1的行为描述和证据示例（从框架中复制）

**引用图：** `fig4-1-indicator-system.drawio`

---

#### §4.2 数据库设计

**§4.2.1 — ER图引用：** `fig4-2-er-diagram.drawio`
**§4.2.2 — 图模型引用：** `fig4-3-neo4j-model.drawio`

**核心表的字段描述必须与 `internal/models/` 下GORM模型定义一致。**

---

#### §4.3 后端系统实现

**§4.3.1 项目结构：** 展示完整的Go项目目录树（从框架中复制）
**引用图：** `fig4-4-main-init.drawio`

**§4.3.2 JWT认证：**
- bcrypt加密(bcrypt cost=10) 来自 `internal/service/auth_service.go`
- JWT token payload包含user_id和role
- AuthMiddleware注入Context

**引用图：** `fig4-5-jwt-flow.drawio`

**§4.3.3 证据管理与分析：**
- 展示FeedbackAgent的核心提示词模板（从框架中复制）
- 描述Prompt的5部分结构

**引用图：** `fig4-6-feedback-agent-flow.drawio`

**§4.3.4 多智能体核心实现（重点章节）：**
- 展示ControlUnit.EvaluateTask的伪代码（已在框架中）
- 分别描述每个Agent的实现要点
- LogicUnit的加权公式：OverallScore = 0.25×(D1+D2+D3+D4)
- 等级映射：≥90优秀, 75-89良好, 60-74合格, <60不合格

**引用图：** `fig4-7-evaluate-task-flow.drawio`

**证据链接：** `internal/agent/control_unit.go` EvaluateTask方法

**§4.3.5 LLM集成：**
- OpenAI兼容API格式 → `internal/llm/client.go`
- 120秒超时设计
- JSON提取（括号深度匹配）→ `internal/agent/inference_agent.go` extractJSON()
- 反馈解析（基于行关键词）→ `internal/llm/parser.go` ParseFeedbackResponse()

**§4.3.6 Redis缓存：**
- 前缀`ra:` 隔离
- 优雅降级机制：`IsAvailable()` false时自动跳过

**§4.3.7 API接口：**
- 16个API端点的完整路由表
- 统一响应格式 `{code, message, data}`

**证据链接：** `cmd/server/main.go` setupRouter函数

---

#### §4.4 前端系统实现

**§4.4.1 技术选型：** Vue3+TS+Vite4+Element Plus+ECharts 6+Axios

**§4.4.2 路由设计：**
**引用图：** `fig4-8-frontend-routes.drawio`

**§4.4.3 核心页面：** 逐个描述登录页、仪表盘、证据管理、结果管理、报告管理

**§4.4.4 数据可视化：** ECharts雷达图、柱状图、饼图的实现方式

**§4.4.5 Mock Token模式：** token前缀`mock-token-`触发离线演示

**证据链接：** `frontend/src/` 目录下所有文件

---

### 第五章 系统测试与结果分析（6-8页）

#### §5.1 测试环境

- 从框架中复制测试环境表
- 测试数据来自`scripts/init_db.go`（35条DB种子）和`testdata/`（189个文件），这点必须明确说明
- 新增测试工具`scripts/batch_llm_classify.go`

#### §5.2 功能测试

- 35个测试用例表（保持）
- **新增§5.2.2多格式文件内容提取验证**：DOCX 81/81（100%），PDF 108/108（100%）

#### §5.3 性能测试

**允许的声明（基于实际代码）：**
- API响应时间数据
- 并发测试：10用户×20请求（来自`scripts/concurrent_test.go`）
- **新增§5.3.3 LLM批量分类性能**：189文件6m23s，平均API响应2.03s/次

#### §5.4 AI评估结果合理性验证（大幅扩展）

**关键新增内容：**

- **§5.4.4 LLM与关键词匹配对比验证**——基于189文件全量测试：
  - 表5-4-1：LLM vs 关键词匹配综合对比（6项指标）
  - 表5-4-2：按质量等级分组评定（8行×5列）
  - 异常检测率98%（49/50），等级±1以内100%，高估率0%
  - 引用新图：图5-4（LLM vs 关键词对比图）、图5-5（异常检测堆叠柱状图）、图5-6（可信度梯度图）

- **§5.4.5 KBM分类分布分析**
  - 表5-4-3：189文件KBM分类分布（5个类别）

- **§5.4.6 偏差来源综合讨论**（保持+新增系统性偏差分析）

#### §5.5 案例分析

- 保持李明和赵六两个典型案例
- **新增§5.5.1 LLM异常检测案例分析**：答非所问检测 + 异常输入检测

#### §5.6 本章小结（新增）

---

### 第六章 总结与展望（2-3页）

**允许的声明：**
- 四大成果总结（从框架中提取）
- 局限性：LLM幻觉、KBM自动化不足、仿真数据规模有限
- 展望：维度扩展、真实场景部署、在EE课程中推广

**参考文献：**
从框架中直接使用26篇参考文献列表

---

## 四、禁止声明清单（全局）

| 禁止项 | 替代表述 |
|--------|---------|
| "首次提出" | "提出了" |
| "填补空白" | "探索了" |
| "最佳方案" | "本文选择的方案" |
| "大规模部署验证" | "系统设计支持部署" |
| "真实学生数据" | "基于仿真数据" |
| "完全替代人工" | "辅助教师评价" |
| "达到业界领先" | 删除或改为客观描述 |
| 任何未在代码中证实的性能指标 | 标注为设计目标 |

---

## 五、图表清单

| 图号 | 文件名 | 章节 | 类型 |
|------|--------|------|------|
| 图1-1 | fig1-1-tech-roadmap.drawio | §1.3.2 | 流程图 |
| 图2-1 | fig2-1-ecd-model.drawio | §2.1.2 | 架构图 |
| 图2-2 | fig2-2-kbm-rubric.drawio | §2.1.4 | 示意图 |
| 图3-1 | fig3-1-system-architecture.drawio | §3.2.2 | 分层架构图 |
| 图3-2 | fig3-2-agent-architecture.drawio | §3.3.1 | 组件图 |
| 图3-3 | fig3-3-agent-sequence.drawio | §3.3.2 | 序列图 |
| 图3-4 | fig3-4-data-flow.drawio | §3.3.3 | 数据流图 |
| 图4-1 | fig4-1-indicator-system.drawio | §4.1.1 | 层级图 |
| 图4-2 | fig4-2-er-diagram.drawio | §4.2.1 | ER图 |
| 图4-3 | fig4-3-neo4j-model.drawio | §4.2.2 | 图模型 |
| 图4-4 | fig4-4-main-init.drawio | §4.3.1 | 流程图 |
| 图4-5 | fig4-5-jwt-flow.drawio | §4.3.2 | 流程图 |
| 图4-6 | fig4-6-feedback-agent-flow.drawio | §4.3.3 | 流程图 |
| 图4-7 | fig4-7-evaluate-task-flow.drawio | §4.3.4 | 流程图 |
| 图4-8 | fig4-8-frontend-routes.drawio | §4.4.2 | 树形图 |
| 图5-1 | fig5-1-llm-vs-keyword.drawio | §5.4.4 | 对比柱状图（新增） |
| 图5-2 | fig5-2-anomaly-detection.drawio | §5.4.4 | 堆叠柱状图（新增） |
| 图5-3 | fig5-3-credibility-gradient.drawio | §5.4.4 | 折线散点图（新增） |
| 图5-4 | fig5-4-radar-chart.drawio | §5.5 | 雷达图 |
| 图5-5 | fig5-5-bar-chart.drawio | §5.5 | 柱状图 |
| 图5-6 | fig5-6-knowledge-graph.drawio | §5.5 | 网络图 |

---

## 六、写给 Composer 的执行顺序

1. **优先写第四章 §4.3.4（多智能体核心实现）**——这是全文最核心、最有证据支撑的章节
2. **其次写第五章 §5.4（AI评价合理性验证）**——注意区分"设计"与"实测"，不可编造数据
3. **再次写第三章**——架构描述，与图表严格对齐
4. **写第四章剩余小节**——数据库、API、前端
5. **写第二章**——理论基础，引用真实文献
6. **写第一章**——背景意义，引用真实文献
7. **写第六章**——总结与展望
8. **最后写摘要**

---

## 七、传递给 Composer 的关键提醒

- **所有Go代码片段从实际文件中提取，不可伪造成api。**
- **提示词模板从 `internal/agent/feedback_agent.go` 和 `internal/agent/inference_agent.go` 中实际提取。**
- **性能数据如无法复现，标注为"目标"而非"实测"。**
- **学生案例数据完全基于 `scripts/init_db.go` 中的模拟数据。**
- **参考文献列表优先使用在框架中已列出的26篇，确认每篇可通过Google Scholar或CNKI检索到。**
- **格式排版严格对齐北京邮电大学手册要求，使用doc skill完成。**
