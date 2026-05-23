package main

import (
	"fmt"
	"log"
	"github.com/Mosher-233/research-ability-assessment/internal/config"
	"github.com/Mosher-233/research-ability-assessment/internal/models"
	"github.com/Mosher-233/research-ability-assessment/pkg/utils"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	log.Println("开始初始化数据库...")

	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	db, err := connectDatabase(cfg)
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}
	log.Println("数据库连接成功")

	log.Println("清空数据库表...")
	clearDatabase(db)

	log.Println("创建测试数据...")
	createTestData(db)

	log.Println("数据库初始化完成！")
}

func connectDatabase(cfg *config.Config) (*gorm.DB, error) {
	switch cfg.Database.Type {
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName)
		return gorm.Open(mysql.Open(dsn), &gorm.Config{})
	default:
		return nil, fmt.Errorf("不支持的数据库类型: %s", cfg.Database.Type)
	}
}

func clearDatabase(db *gorm.DB) {
	db.Exec("SET FOREIGN_KEY_CHECKS = 0")
	db.Exec("TRUNCATE TABLE reports")
	db.Exec("TRUNCATE TABLE inference_results")
	db.Exec("TRUNCATE TABLE feedbacks")
	db.Exec("TRUNCATE TABLE evidences")
	db.Exec("TRUNCATE TABLE student_tasks")
	db.Exec("TRUNCATE TABLE tasks")
	db.Exec("TRUNCATE TABLE students")
	db.Exec("TRUNCATE TABLE teachers")
	db.Exec("TRUNCATE TABLE users")
	db.Exec("SET FOREIGN_KEY_CHECKS = 1")
	log.Println("数据库表清空完成")
}

func createTestData(db *gorm.DB) {
	teachers := createTeachers(db)
	students := createStudents(db)
	task := createResearchTask(db, teachers[0])
	studentTaskIDs := assignTaskToStudents(db, task.ID, students)

	createEvidenceForStudents(db, studentTaskIDs)

	log.Println("\n=== 账号信息 ===")
	log.Println("教师账号（密码均为 123456）:")
	log.Println("  1. 张三老师 - 1@tea.com")
	log.Println("\n学生账号（密码均为 123456）:")
	log.Println("  1. 王五 - 1@stu.com (研究能力较强)")
	log.Println("  2. 赵六 - 2@stu.com (中等水平)")
	log.Println("  3. 李明 - 3@stu.com (优秀水平)")
	log.Println("  4. 钱七 - 4@stu.com (研究设计突出)")
	log.Println("  5. 孙八 - 5@stu.com (数据分析突出)")
}

func createTeachers(db *gorm.DB) []string {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	teacher := &models.User{
		ID:       utils.GenerateUserID("teacher"),
		Name:     "张三老师",
		Email:    "1@tea.com",
		Password: string(hashedPassword),
		Role:     "teacher",
	}
	if err := db.Create(teacher).Error; err != nil {
		log.Fatalf("创建教师失败: %v", err)
	}
	log.Printf("创建教师成功: %s (%s)", teacher.Name, teacher.ID)
	return []string{teacher.ID}
}

func createStudents(db *gorm.DB) []string {
	students := []struct {
		name  string
		email string
	}{
		{"王五", "1@stu.com"},
		{"赵六", "2@stu.com"},
		{"李明", "3@stu.com"},
		{"钱七", "4@stu.com"},
		{"孙八", "5@stu.com"},
	}

	var studentIDs []string
	for _, s := range students {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
		student := &models.User{
			ID:       utils.GenerateUserID("student"),
			Name:     s.name,
			Email:    s.email,
			Password: string(hashedPassword),
			Role:     "student",
		}
		if err := db.Create(student).Error; err != nil {
			log.Printf("创建学生 %s 失败: %v", s.name, err)
		} else {
			studentIDs = append(studentIDs, student.ID)
			log.Printf("创建学生成功: %s (%s)", s.name, student.ID)
		}
	}
	return studentIDs
}

func createResearchTask(db *gorm.DB, teacherID string) *models.Task {
	task := &models.Task{
		ID:          utils.GenerateTaskID(),
		Name:        "大学生研究能力综合评价",
		Description: "基于学生提交的创新创业项目材料（结题报告、项目计划书），对学生进行多维度研究能力评价。评价维度包括：文献综述能力、研究设计能力、数据分析能力和批判性思维能力。",
		CourseID:    "INNOVATE-2024",
		TeacherID:   teacherID,
		StartDate:   time.Now(),
		EndDate:     time.Now().AddDate(0, 1, 0),
		Status:      "active",
	}
	if err := db.Create(task).Error; err != nil {
		log.Fatalf("创建任务失败: %v", err)
	}
	log.Printf("创建任务成功: %s (%s)", task.Name, task.ID)
	return task
}

func assignTaskToStudents(db *gorm.DB, taskID string, studentIDs []string) []string {
	var studentTaskIDs []string
	for _, studentID := range studentIDs {
		studentTask := &models.StudentTask{
			ID:        utils.GenerateStudentTaskID(),
			TaskID:    taskID,
			StudentID: studentID,
			Status:    "pending",
			Progress:  0,
		}
		if err := db.Create(studentTask).Error; err != nil {
			log.Printf("分配任务给学生 %s 失败: %v", studentID, err)
		} else {
			studentTaskIDs = append(studentTaskIDs, studentTask.ID)
		}
	}
	log.Printf("已将任务分配给 %d 名学生", len(studentTaskIDs))
	return studentTaskIDs
}

type evidenceInput struct {
	Type     string
	KBMName  string
	KBMLevel int
	Content  string
}

func createEvidenceForStudents(db *gorm.DB, studentTaskIDs []string) {
	evidences := map[int][]evidenceInput{
		0: {
			{
				Type: "结题报告", KBMName: "文献检索策略", KBMLevel: 4,
				Content: "在项目启动初期,我们通过CNKI和Web of Science数据库检索了2019-2024年间关于'智能垃圾分类'的文献共87篇,使用了'深度学习+垃圾分类'、'图像识别+废物管理'等复合关键词组合。在此基础上,进一步筛选出与本研究直接相关的核心文献32篇。检索策略涵盖了中文和英文数据库,确保了文献来源的全面性。",
			},
			{
				Type: "项目计划书", KBMName: "文献综述质量", KBMLevel: 3,
				Content: "文献综述部分从垃圾分类的国内外现状、现有技术路线、评价指标三个维度进行了系统梳理。国外方面重点分析了美国、日本等发达国家在智能垃圾分类系统方面的研究进展;国内方面对比了不同高校和企业的研发成果。综述指出,当前研究主要集中在单模态识别上,多模态融合的研究相对较少,这为本项目的创新点提供了理论支撑。",
			},
			{
				Type: "结题报告", KBMName: "实验方案合理性", KBMLevel: 4,
				Content: "实验方案设计采用对照实验法,设置了三个实验组(仅图像识别组、仅传感器组、多模态融合组),每组30次实验,共90次实验。同时设置了对照组(传统人工分类)作为基准。实验变量包括:垃圾种类(塑料瓶、玻璃瓶、纸箱、厨余垃圾)、光线条件(强光/弱光)、垃圾摆放角度(正面/侧面/倾斜)。方案经指导教师审核后执行。",
			},
			{
				Type: "项目计划书", KBMName: "变量控制", KBMLevel: 4,
				Content: "为确保实验结果的可靠性,我们对以下变量进行了严格控制:1)实验环境:在固定实验室进行,温度控制在22±1°C,湿度控制在50%±5%;2)设备一致性:所有实验使用同一套硬件设备(树莓派4B+摄像头模块);3)垃圾样本标准化:从学校食堂和宿舍楼统一收集,按照统一标准进行分类和预处理。每个实验条件重复5次取平均值。",
			},
			{
				Type: "结题报告", KBMName: "数据分析方法选择", KBMLevel: 3,
				Content: "数据分析采用Python的scikit-learn和pandas库完成。主要分析方法包括:1)描述性统计:计算各实验组的识别准确率均值、标准差;2)单因素方差分析(ANOVA):检验不同实验组之间是否存在显著差异;3)混淆矩阵分析:评估不同垃圾种类的分类效果。所有统计检验的显著性水平设为α=0.05。",
			},
			{
				Type: "结题报告", KBMName: "问题提出新颖性", KBMLevel: 4,
				Content: "本研究提出的研究问题具有创新性:不同于现有研究仅关注垃圾分类准确率的单一指标,我们提出了'分类效率-能耗比'这一新的综合评价指标。这一指标的引入源于在宿舍区实地观察中发现,高精度分类往往伴随高能耗,对于实际部署场景而言、两者需要权衡。该研究问题的提出切中了智慧校园建设的实际痛点。",
			},
			{
				Type: "项目计划书", KBMName: "解决方案原创性", KBMLevel: 4,
				Content: "我们提出的解决方案融合了视觉识别和物理传感器两个信号源,通过自注意力机制进行特征融合。相比于市面上已有的纯视觉方案,我们的方案在弱光条件下仍有较高准确率;相比于纯传感器方案,我们的方案可识别的垃圾种类更多。技术方案中使用的联邦学习框架确保各宿舍楼的数据隐私,这一设计得到了学校信息化办公室的认可。",
			},
			{
				Type: "结题报告", KBMName: "结果解释准确性", KBMLevel: 3,
				Content: "实验结果表明,多模态融合组的平均识别准确率达到92.3%,显著高于仅图像识别组(85.1%)和仅传感器组(78.6%)(F=15.47,p<0.001)。在弱光条件下,多模态融合组的准确率仅下降3.2个百分点,而纯图像组下降了12.5个百分点,说明传感器信息对光照有较强补偿作用。但需要指出的是,对于黑色塑料袋的识别准确率在三个实验组中均偏低,可能与光学传感器反射特性有关。",
			},
		},
		1: {
			{
				Type: "结题报告", KBMName: "文献检索策略", KBMLevel: 3,
				Content: "通过万方数据库和Google Scholar检索了关于'校园二手交易平台'的相关文献,检索关键词包括'校园二手'、'C2C交易平台'、'大学生消费行为'等。共检索到文献45篇,经阅读摘要筛选后保留20篇作为参考文献。检索范围以中文文献为主,英文文献较少涉及。",
			},
			{
				Type: "项目计划书", KBMName: "文献综述质量", KBMLevel: 3,
				Content: "文献综述从校园二手交易的需求分析、现有平台功能对比、用户行为研究三个方面进行了梳理。分析了闲鱼、转转等主流平台的功能特点,并将其与校园场景的特殊需求进行对比,指出了校园二手交易的三个特殊性:交易双方信任度高、物品类别集中、物流成本低。文献综述逻辑清晰,但对国外高校二手交易平台的研究涉及较少。",
			},
			{
				Type: "项目计划书", KBMName: "实验方案合理性", KBMLevel: 3,
				Content: "系统测试方案包括功能测试和用户测试两个部分。功能测试覆盖了系统的所有核心功能模块(用户注册、商品发布、交易匹配、评价反馈);用户测试邀请了30名大学生进行为期一周的体验,通过问卷调查收集使用反馈。方案设计基本合理,但在性能和安全性测试方面考虑不足。",
			},
			{
				Type: "结题报告", KBMName: "数据分析方法选择", KBMLevel: 2,
				Content: "对用户反馈数据进行了简单的统计分析,主要计算了各功能模块的满意度得分和用户使用频率分布。分析方法较为基础,未进行深入的统计检验或关联分析。问卷数据的信度和效度分析也未进行。",
			},
			{
				Type: "结题报告", KBMName: "结果解释准确性", KBMLevel: 2,
				Content: "用户反馈显示满意度为4.2/5,但未解释该得分的统计学意义。交易匹配功能的满意度为3.8/5,低于预期,可能是因为系统中早期的商品数据量不足导致匹配精度下降,但具体原因缺乏数据支持。对'用户活跃度下降'现象的解释仅基于主观判断,缺少客观数据分析。",
			},
			{
				Type: "项目计划书", KBMName: "问题提出新颖性", KBMLevel: 2,
				Content: "提出了'基于校园关系链的信誉评价机制'这一想法,相比现有平台仅依赖交易记录的评价体系,添加了同学关系维度。但该想法在现有文献中已有类似研究,新颖性有限。",
			},
			{
				Type: "结题报告", KBMName: "变量控制", KBMLevel: 2,
				Content: "用户测试中参与了30名学生,但未控制参与者的年级、专业、性别等人口统计学变量的影响。不同年级的学生可能有不同的二手交易需求,这一潜在变量未被纳入分析。",
			},
		},
		2: {
			{
				Type: "结题报告", KBMName: "文献检索策略", KBMLevel: 5,
				Content: "本研究采用了系统性的文献检索策略。首先在Web of Science、IEEE Xplore、ACM Digital Library、CNKI四个数据库中检索了2018-2024年间关于'教育数据挖掘与学习分析'的文献。检索策略采用PICO框架构建:P(大学生群体)、I(机器学习预测模型)、C(传统统计方法)、O(学业预警准确率),共检索到文献126篇。经过双人独立筛选和交叉比对,最终纳入46篇核心文献。使用PRISMA流程图记录了文献筛选全过程。",
			},
			{
				Type: "项目计划书", KBMName: "文献综述质量", KBMLevel: 5,
				Content: "文献综述采用系统综述方法(Systematic Review),从以下维度进行了结构化分析:1)预测模型的时间演变(传统统计模型→机器学习模型→深度学习模型);2)输入特征的分类(学业成绩类、行为数据类、心理测评类);3)评价指标的统一性(F1-score、AUC、召回率的比较);4)模型的公平性和可解释性。综述最终通过对比分析表展现了15个代表性研究的方法、样本量和预测效果,清晰指出当前研究的三个不足:样本量普遍偏小、缺乏多模态数据融合、模型可解释性不足。",
			},
			{
				Type: "项目计划书", KBMName: "问题提出新颖性", KBMLevel: 5,
				Content: "本研究的创新之处在于:提出了'行为时序特征+学业静态特征'的双流融合预测框架。现有研究大多仅使用单一时间点的学业数据或行为数据进行预测,忽略了学生行为的时序动态特性。我们通过分析学生在学习管理系统(LMS)中的周度行为序列数据,发现'在线时长下降速率'和'作业提交延迟天数'这两个时间衍生特征对成绩预警的预测能力远强于静态特征。这一发现为学业预警研究提供了新的视角。",
			},
			{
				Type: "项目计划书", KBMName: "解决方案原创性", KBMLevel: 5,
				Content: "我们设计了基于Transformer的时序行为编码器(Behavioral Sequence Encoder),将学生每周的LMS行为数据转化为固定长度的向量表示,再与学业成绩特征进行跨模态注意力融合。相比于现有方法,该方案的核心优势在于:1)不需要手动构造时序特征,模型自动学习行为模式;2)注意力权重可视化提供了模型可解释性;3)支持增量学习,可适应学期中数据逐步增加的场景。该架构设计参考了Transformer在NLP中的成功经验,首次将其应用于教育时序行为建模。",
			},
			{
				Type: "结题报告", KBMName: "实验方案合理性", KBMLevel: 5,
				Content: "实验采用严格的机器学习实验设计:1)数据集划分:7:2:1的三层划分(训练/验证/测试),确保了模型评估的独立性;2)交叉验证:5折交叉验证评估模型的泛化能力;3)基线对比:与逻辑回归、随机森林、XGBoost、LSTM四个基线模型进行对比;4)消融实验:分别移除时序特征、学业特征和注意力机制,评估各组件的贡献;5)统计检验:使用McNemar检验评估模型间预测差异的显著性。实验方案的严谨性确保了结论的可信度。",
			},
			{
				Type: "结题报告", KBMName: "变量控制", KBMLevel: 5,
				Content: "实验中的变量控制非常严格:1)所有模型使用完全相同的训练/测试数据划分;2)超参数调优使用相同的搜索空间(网格搜索)和评价指标(F1-score);3)重复实验10次取平均值和标准差,以减小随机性影响;4)控制了学期阶段变量的影响,分别在学期第4、8、12、16周进行预测,观察预测效果的时间动态;5)对样本不平衡问题进行了处理,使用SMOTE过采样方法确保正负类样本均衡。",
			},
			{
				Type: "结题报告", KBMName: "数据分析方法选择", KBMLevel: 5,
				Content: "数据分析方法的选取严谨且多样:1)采用ROC-AUC作为主要评价指标(适合不平衡数据集),同时报告F1-score和准确率;2)使用Shapley值分析(SHAP)对特征重要性进行解释;3)通过注意力热图可视化模型关注的行为时段;4)使用配对样本t检验比较不同时间点的预测效果;5)采用混淆矩阵和分类报告详细分析不同类型学生的预测效果。所有统计分析均使用Python的scikit-learn、statsmodels和自定义代码实现,代码已托管至GitHub仓库以供复现。",
			},
			{
				Type: "结题报告", KBMName: "结果解释准确性", KBMLevel: 5,
				Content: "实验结果表明,双流融合模型的F1-score达到了0.84,显著优于最佳基线模型XGBoost(0.78)(McNemar检验,p=0.003)。SHAP分析显示,'在线时长下降速率'是第四周对最终挂科预测贡献最大的特征,这与教育心理学中的'学习投入下降先于成绩下降'理论一致。注意力可视化进一步揭示,模型主要通过学期第5-8周的行为变化进行预测,这个阶段恰好是课程中期,学生的行为模式变化最为明显。这些结果的解释既有统计学依据,也有理论支撑。",
			},
		},
		3: {
			{
				Type: "项目计划书", KBMName: "实验方案合理性", KBMLevel: 5,
				Content: "针对'智慧教室环境监测系统'项目,我们设计了多层次的测试方案:1)单元测试:对温度传感器、湿度传感器、CO2传感器进行独立校准和精度验证,与标准仪器进行对比;2)集成测试:在三个不同类型的教室(大阶梯教室、小班教室、实验室)部署系统,进行连续7天×24小时的数据采集;3)对比测试:与商业环境监测仪(型号TES-1370)进行同步数据对比,计算偏差范围;4)压力测试:模拟100个节点同时上传数据的场景,测试网关的并发处理能力。该方案从精度、稳定性、可扩展性三个维度全面评估了系统性能。",
			},
			{
				Type: "结题报告", KBMName: "变量控制", KBMLevel: 5,
				Content: "实验中的变量控制设计精细:1)为消除传感器位置偏差,每个教室部署3个传感器节点(前、中、后排),数据取均值;2)为控制时间变量的影响,连续采集7天数据,同时在每个整点时刻记录数据;3)为控制教室使用状态的影响,分别记录了上课/下课时段的数据;4)传感器的预热时间统一设置为30分钟。所有数据采集严格按照预定的采样方案执行,减少了人为主观因素对数据质量的影响。",
			},
			{
				Type: "结题报告", KBMName: "实验实施质量", KBMLevel: 5,
				Content: "整个实验实施过程严格按照实验方案执行,得到了实验室管理员和任课教师的配合。系统连续运行期间(168小时),数据采集成功率达到98.7%。仅有的1.3%丢包主要发生在第5天网络波动期间,已通过本地缓存+断点续传机制成功补传。实验过程的所有原始数据、日志文件和操作记录均已存档,保证了实验的可追溯性和可复现性。",
			},
			{
				Type: "结题报告", KBMName: "数据分析方法选择", KBMLevel: 4,
				Content: "数据分析采用Python生态工具链:1)使用pandas进行数据清洗和缺失值处理(线性插值法);2)使用matplotlib和seaborn绘制环境参数的时间变化曲线;3)计算各环境参数的日均值、峰谷差和变异系数(CV);4)使用Pearson相关系数分析温度、湿度、CO2浓度之间的关联关系;5)使用K-means聚类分析不同教室类型的环境特征模式。分析方法的选择充分考虑了数据的时间序列特性和传感器的测量误差。",
			},
			{
				Type: "项目计划书", KBMName: "问题提出新颖性", KBMLevel: 4,
				Content: "本研究提出将环境参数与学生学习状态进行关联分析的新视角,即通过分析教室环境参数的变化趋势,间接反映课堂教学活动的节奏和模式。例如,CO2浓度的累计速率可能反映课堂人数的密度和通风情况,温度和湿度的变化可能影响学生的舒适度和注意力。这一研究视角不同于传统环境监测项目仅关注环境参数本身,而是尝试建立环境-学习效果之间的桥梁。",
			},
			{
				Type: "项目计划书", KBMName: "文献批判性分析", KBMLevel: 3,
				Content: "在文献综述中,我们发现多数智慧教室研究聚焦于硬件设备的部署和技术架构的设计,较少有研究关注环境参数与学习效果之间的因果关系,也更少有人从传感器数据的质量(精度、稳定性、可靠性)角度评估系统。另外,现有文献中实验样本的规模普遍偏小,多数仅为1-2间教室的短期测试,缺乏大规模、长时间的验证性研究。这些发现为本项目的实验设计提供了重要的参考和改进方向。",
			},
		},
		4: {
			{
				Type: "项目计划书", KBMName: "数据分析方法选择", KBMLevel: 5,
				Content: "在'基于多源数据的大学生心理健康预警系统'项目中,我们构建了完整的心理数据分析流水线:1)数据预处理:使用Min-Max标准化处理不同量纲的问卷得分,使用One-Hot编码处理分类变量;2)特征工程:从原始问卷得分中构建了复合因子得分(抑郁因子、焦虑因子、压力因子),从校园卡流水数据中提取了行为特征(就餐规律性、社交广度、作息规律);3)建模:比较了逻辑回归、SVM、随机森林、XGBoost和LightGBM五种分类器,使用网格搜索和交叉验证进行超参数调优;4)模型评估:采用分层K折交叉验证评估模型稳定性,使用AUC-ROC作为主要评价指标;5)模型解释:使用LIME和SHAP分别进行局部和全局的可解释性分析。",
			},
			{
				Type: "结题报告", KBMName: "结果解释准确性", KBMLevel: 5,
				Content: "LightGBM模型在测试集上取得了AUC=0.89的最佳效果。SHAP分析揭示,'就餐规律性的标准差'和'深夜门禁次数'是预测心理健康风险最重要的两个行为特征。这与心理学文献中'行为节律紊乱是心理问题的早期信号'的理论一致。LIME分析表明,模型对高风险样本的预测置信度较高,对中风险样本的预测存在一定的模糊性,这提示了后续可以引入不确定性量化方法(如蒙特卡洛Dropout)来改进中风险区间的预测。结果的解释既有数据驱动,也关联了理论依据。",
			},
			{
				Type: "项目计划书", KBMName: "文献检索策略", KBMLevel: 4,
				Content: "文献检索采用多阶段策略:第一阶段在PubMed、PsycINFO和CNKI中检索关于'大学生心理健康'、'预警模型'、'机器学习'等关键词的文献;第二阶段根据第一阶段检索到的文献的引用关系,通过引文追溯法扩展检索范围,重点关注高引论文和最新发表的顶会论文;第三阶段检索了一些灰色文献,包括教育部关于学生心理健康的相关政策文件和技术报告。最终纳入分析的文献共63篇,时间跨度为2015-2024年。",
			},
			{
				Type: "项目计划书", KBMName: "文献批判性分析", KBMLevel: 4,
				Content: "在文献综述中,我们从技术方法和研究伦理两个角度进行了批判性分析。技术上,发现现有预警模型存在三个主要问题:1)过度依赖事后问卷调查数据,缺乏实时的行为信号;2)模型在不同高校间的泛化能力未经验证;3)多数研究仅关注预测准确率,忽略了模型的公平性(如不同性别、专业背景的学生是否被公平对待)。伦理上,注意到数据隐私和数据使用知情同意是关键的合规要求,但多数技术文献对此讨论不足。这些批判性发现为本项目的设计提供了重要指导。",
			},
			{
				Type: "结题报告", KBMName: "问题提出新颖性", KBMLevel: 4,
				Content: "本研究提出了一个新颖的研究视角:将校园行为大数据(校园卡/门禁/WiFi连接)作为心理健康状态的观测窗口。传统心理学依赖于定期问卷测评,但问卷存在社会期望偏差、填写疲劳和延时性等问题。我们假设:日常行为的微小变化(如就餐时间漂移、社交范围收缩、作息紊乱)可能是心理状态变化的早期行为标志。这一视角将心理学量表这一'主观测度'与校园行为数据这一'客观测度'进行了整合,提供了一种低成本、持续的监测方式。",
			},
		},
	}

	for i, studentTaskID := range studentTaskIDs {
		if i >= len(evidences) {
			break
		}
		createStudentEvidences(db, studentTaskID, evidences[i])
	}
}

func createStudentEvidences(db *gorm.DB, studentTaskID string, evs []evidenceInput) {
	for _, ev := range evs {
		evidence := &models.Evidence{
			ID:            utils.GenerateEvidenceID(),
			StudentTaskID: studentTaskID,
			Type:          ev.Type,
			Content:       ev.Content,
			KBMName:       ev.KBMName,
			KBMLevel:      ev.KBMLevel,
		}
		if err := db.Create(evidence).Error; err != nil {
			log.Printf("创建证据 [%s] 失败: %v", ev.KBMName, err)
		}
	}
	log.Printf("为 StudentTask %s 创建了 %d 条证据", studentTaskID, len(evs))
}
