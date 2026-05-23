package main

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ============================================================
// 模拟仿真测试数据生成脚本
// 覆盖场景：高质量(L4/L5)、良好(L3)、合格(L2)、不合格(L1)、答非所问、异常输入
// 共6名学生 × 6-8条证据 = 40条证据
// ============================================================

var dsn = "root:rootpassword@tcp(127.0.0.1:3306)/research_assessment?charset=utf8mb4&parseTime=True&loc=Local"

// ----- Models (对齐项目实际结构) -----

type User struct {
	ID        string `gorm:"primaryKey;type:varchar(36)"`
	Email     string `gorm:"uniqueIndex;type:varchar(255);not null"`
	Password  string `gorm:"type:varchar(255);not null"`
	Name      string `gorm:"type:varchar(100);not null"`
	Role      string `gorm:"type:varchar(20);not null;default:'student'"` // teacher / student
	CreatedAt time.Time
}

type Task struct {
	ID          string    `gorm:"primaryKey;type:varchar(36)"`
	Name        string    `gorm:"type:varchar(255);not null"`
	Description string    `gorm:"type:text"`
	CourseID    string    `gorm:"type:varchar(36);not null"`
	TeacherID   string    `gorm:"type:varchar(36);index"`
	StartDate   time.Time `gorm:"not null"`
	EndDate     time.Time `gorm:"not null"`
	Status      string    `gorm:"type:varchar(20);not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type StudentTask struct {
	ID        string `gorm:"primaryKey;type:varchar(36)"`
	TaskID    string `gorm:"type:varchar(36);index"`
	StudentID string `gorm:"type:varchar(36);index"`
	Status    string `gorm:"type:varchar(20);default:'pending'"` // pending / processing / completed
	Progress  int    `gorm:"default:0"`
	CreatedAt time.Time
}

type Evidence struct {
	ID            string `gorm:"primaryKey;type:varchar(36)"`
	StudentTaskID string `gorm:"type:varchar(36);index;not null"`
	Type          string `gorm:"type:varchar(50);not null;default:'text'"`
	Content       string `gorm:"type:text;not null"`
	KBMName       string `gorm:"type:varchar(50);not null"`
	KBMLevel      int    `gorm:"not null"` // 1-4
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// ----- Data -----

type StudentInfo struct {
	ID       string
	Name     string
	Email    string
	Password string
}

type EvidenceInfo struct {
	KBMName   string
	Dimension string
	KBMLevel  int
	Content   string
}

var teacher = StudentInfo{ID: uuid.New().String(), Name: "詹志强", Email: "1@tea.com", Password: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy"}

var students = []StudentInfo{
	{ID: uuid.New().String(), Name: "周博远", Email: "zhouby@stu.com", Password: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy"},
	{ID: uuid.New().String(), Name: "陈思涵", Email: "chensh@stu.com", Password: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy"},
	{ID: uuid.New().String(), Name: "林浩然", Email: "linhr@stu.com", Password: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy"},
	{ID: uuid.New().String(), Name: "张晓婷", Email: "zhangxt@stu.com", Password: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy"},
	{ID: uuid.New().String(), Name: "王铭宇", Email: "wangmy@stu.com", Password: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy"},
	{ID: uuid.New().String(), Name: "刘雨桐", Email: "liuyt@stu.com", Password: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy"},
}

var task = struct {
	ID          string
	Name        string
	Description string
	CourseID    string
}{
	ID:          uuid.New().String(),
	Name:        "基于深度学习的语音情感识别系统设计与实现",
	Description: "设计并实现一个基于深度学习的语音情感识别系统。要求：调研语音情感识别领域研究现状，完成文献综述；设计合理实验方案；实现系统原型并进行充分实验验证；对实验结果进行深入分析。",
	CourseID:    "DS2024001",
}

// 每个学生的证据列表
var studentEvidences = map[string][]EvidenceInfo{
	"zhouby@stu.com": {
		// 周博远 - 优秀 (8条L4/L5)
		{"文献检索策略", "文献检索与综述", 4, "采用PICO框架制定检索策略，从IEEE Xplore、ACM DL、Web of Science、CNKI四个数据库进行系统检索。检索式：(speech emotion recognition OR SER) AND (deep learning OR neural network) AND (transformer OR self-attention)。首次检索获得287篇文献，通过双人独立筛选（Cohen's Kappa=0.82）和PRISMA流程图记录，最终纳入46篇核心文献进行分析。检索过程完整记录了每个数据库的命中数、筛选标准和排除理由。"},
		{"文献综述质量", "文献检索与综述", 4, "完成系统综述，从三个维度对46篇文献进行梳理：（1）特征提取方法演进：从手工特征（MFCC、韵律特征）到深度特征（CNN、Transformer）；（2）模型架构对比：CRNN、Transformer、Conformer等架构在IEMOCAP数据集上的WA对比；（3）跨语种/跨域泛化挑战。制作了详细的对比表格，列出了15篇代表性论文的方法、数据集和性能指标。明确指出了三个研究空白：多模态融合不足、跨语种泛化差、实时性有待提升。"},
		{"实验方案合理性", "研究设计与实验", 4, "设计多阶段实验方案。数据集：IEMOCAP（四分类）+ CASIA（中文两分类）。数据划分：7:2:1（训练/验证/测试），5折交叉验证。基线模型：Logistic Regression、Random Forest、CNN-LSTM、Transformer。消融实验：分别移除注意力机制、数据增强、多任务学习模块，评估各组件贡献。统计检验：使用McNemar检验评估模型间差异的显著性（p<0.05）。变量控制：统一随机种子（42）、统一数据预处理流程、统一训练轮次（100 epochs）。"},
		{"变量控制", "研究设计与实验", 4, "实验环境统一配置：GPU为NVIDIA RTX 4090，CUDA 12.1，PyTorch 2.1。所有模型使用相同的数据预处理pipeline（包括相同的归一化参数）。随机种子固定为42，每个实验重复运行10次取平均值±标准差。学习率统一使用CosineAnnealing调度器，初始值0.001。记录了GPU温度、显存占用和训练时长等环境参数，确保实验可复现。"},
		{"数据分析方法选择", "数据分析与解释", 4, "采用多种评估指标：加权准确率（WA）、非加权准确率（UA）、F1-score（macro）、ROC-AUC。使用SHAP值分析特征重要性，发现MFCC的第13维（高频能量）对愤怒情绪识别贡献最大。绘制注意力热图可视化模型对语音不同时间段的关注权重。使用配对样本t检验比较本模型与基线模型的性能差异（p=0.003<0.05，差异显著）。使用McNemar检验验证消融实验中各模块的统计显著性。"},
		{"结果解释准确性", "数据分析与解释", 4, "实验结果显示本模型在IEMOCAP上WA达到76.2%，优于基线模型。深入分析发现：（1）Transformer在长时依赖建模上优于CNN-LSTM，但对短时爆发性情感（如惊讶）识别效果较差；（2）SHAP分析表明韵律特征（基频、能量）对情感分类的贡献大于频谱特征；（3）注意力热图显示模型在情感转折点（如从平静到愤怒）的关注度显著提高。讨论了局限：仅使用英文数据集，跨语种泛化能力未验证；实时推理延迟（45ms/帧）尚未满足实时性要求。"},
		{"问题提出新颖性", "批判性思维与创新", 4, "提出研究问题：能否通过引入语音情感变化的时序动态特征来提升连续语音情感识别的准确率？该问题的创新性在于：（1）现有方法主要关注帧级特征，忽略了情感在时间维度上的渐变过程；（2）参考文献[Chen et al., 2023]发现情感转换点附近的声学特征存在显著变化，但未被充分利用。本方案设计了一个基于Transformer编码器的时序情感变化检测模块，首次将情感转换点作为辅助任务进行联合训练。"},
		{"解决方案原创性", "批判性思维与创新", 4, "提出双流融合架构：流1为常规语音特征编码器（MelSpec→Conv→Transformer），流2为时序情感变化检测器（检测情感转换点并生成变化向量）。两流通过门控融合机制整合。与现有方法相比的创新点：（1）首次将情感转换点检测作为辅助任务；（2）门控融合机制可自适应调整两流的权重；（3）在IEMOCAP上WA提升2.1个百分点。绘制了完整的系统架构图和数据流图。"},
	},
	"chensh@stu.com": {
		// 陈思涵 - 良好 (6条L3 + 2条L2)
		{"文献检索策略", "文献检索与综述", 3, "从Google Scholar和CNKI两个数据库检索语音情感识别相关文献。使用关键词'语音情感识别 深度学习'和'speech emotion recognition deep learning'，共检索到85篇文献。通过阅读标题和摘要筛选出32篇相关文献。检索过程记录了每个数据库的命中数和筛选标准。"},
		{"文献综述质量", "文献检索与综述", 3, "对32篇文献进行了多维度梳理，从特征提取方法和模型架构两个角度进行了归纳。制作了对比表格，列出了8篇代表性论文的方法和性能。指出了当前研究的两个不足：跨语种泛化差和实时性有待提升。"},
		{"实验方案合理性", "研究设计与实验", 3, "使用IEMOCAP数据集，按8:1:1划分训练/验证/测试集。选择CNN-LSTM和Transformer两个模型进行对比实验。使用加权准确率和F1-score作为评估指标。设置了学习率为0.001，batch size为32，训练100个epoch。"},
		{"变量控制", "研究设计与实验", 3, "实验在GPU RTX 3060上进行，使用PyTorch 2.0框架。固定随机种子为42。两个模型使用相同的数据预处理流程和训练配置。记录了训练过程中的loss曲线和验证集准确率变化。"},
		{"数据分析方法选择", "数据分析与解释", 3, "使用加权准确率（WA）、非加权准确率（UA）和F1-score（macro）评估模型性能。绘制了混淆矩阵分析各类别的识别准确率。使用t检验比较两个模型的性能差异（p=0.02<0.05）。绘制了训练过程的loss曲线和准确率曲线。"},
		{"结果解释准确性", "数据分析与解释", 3, "Transformer模型在IEMOCAP上WA达到73.5%，CNN-LSTM为70.1%。分析发现Transformer在愤怒和快乐两类上表现较好，但在悲伤和恐惧上识别率较低。混淆矩阵显示悲伤容易被误分类为中性。讨论了局限：仅使用单一数据集，未进行跨数据集验证。"},
		{"问题提出新颖性", "批判性思维与创新", 3, "研究问题：如何改进CNN-LSTM模型在语音情感识别中的性能？该问题有一定针对性，参考文献指出CNN-LSTM在长时依赖建模上存在不足。计划通过引入注意力机制来改进模型。"},
		{"解决方案原创性", "批判性思维与创新", 2, "在CNN-LSTM模型中加入了自注意力层，使模型能够关注语音序列中的重要时间步。参考了[Zhao et al., 2022]的注意力机制设计。模型结构为：MelSpec输入→CNN特征提取→BiLSTM序列建模→Self-Attention→全连接分类。"},
	},
	"linhr@stu.com": {
		// 林浩然 - 合格 (4条L2 + 2条L1)
		{"文献检索策略", "文献检索与综述", 2, "在百度学术上搜索了'语音情感识别'关键词，找到了大约20篇相关论文。下载了其中5篇仔细阅读。"},
		{"文献综述质量", "文献检索与综述", 2, "阅读了5篇论文后进行了归纳，主要介绍了语音情感识别的基本概念、常用的数据集和几种常见的深度学习模型。"},
		{"实验方案合理性", "研究设计与实验", 2, "参考课程实验指导书，使用IEMOCAP数据集，搭建了一个CNN模型进行语音情感分类。训练集和测试集按9:1划分。使用交叉熵损失函数和Adam优化器。"},
		{"数据分析方法选择", "数据分析与解释", 2, "使用准确率作为评估指标，模型在测试集上达到62.3%。绘制了训练过程的loss曲线，loss从2.3下降到0.8。"},
		{"文献批判性分析", "文献检索与综述", 1, "参考文献都很有价值，提出的深度学习方法都很先进，我打算按照其中一篇论文的方法来实现。"},
		{"问题提出新颖性", "批判性思维与创新", 1, "我想做一个语音情感识别系统，因为这是课程设计的要求。"},
	},
	"zhangxt@stu.com": {
		// 张晓婷 - 不合格 (2条L1 + 4条答非所问)
		{"文献检索策略", "文献检索与综述", 1, "我在网上搜了一些关于语音识别的资料，找到了几篇博客文章和百度百科的介绍。"},
		{"实验方案合理性", "研究设计与实验", 1, "今天天气不错，我打算去图书馆借几本关于Python编程的书，然后开始写代码。我觉得Python是一门很好的编程语言，适合初学者学习。希望能在下周之前完成作业。"},
		{"数据分析方法选择", "数据分析与解释", 1, "我下载了IEMOCAP数据集但是打不开，文件格式好像不对。然后我试了用Audacity打开音频文件，可以播放但是不知道怎么提取特征。请问老师这个数据集应该怎么用？"},
		{"文献综述质量", "文献检索与综述", 1, "语音情感识别就是让计算机听懂人的情绪。我觉得这个方向很有意义，因为以后可以用来做智能客服。我查了一下，好像已经有很多人在做这个了。"},
		{"解决方案原创性", "批判性思维与创新", 1, "我打算用一个现成的开源项目来做，在GitHub上找到了一个star比较多的项目，直接跑一下看看效果。"},
		{"变量控制", "研究设计与实验", 1, "我觉得这个课程设计太难了，能不能换一个简单的题目？或者直接用现成的代码跑一下就行了？"},
	},
	"wangmy@stu.com": {
		// 王铭宇 - 混合 (2条L4 + 2条L3 + 3条答非所问)
		{"文献检索策略", "文献检索与综述", 4, "使用PICO框架从IEEE Xplore、ACM DL和Web of Science三个数据库检索。检索式：(speech emotion recognition) AND (transformer OR attention) AND (multimodal)。检索到156篇，经筛选纳入28篇。使用PRISMA流程图记录筛选过程，双人独立筛选Cohen's Kappa=0.78。"},
		{"实验方案合理性", "研究设计与实验", 4, "设计三组对照实验：（1）仅语音特征（MelSpec）；（2）仅文本特征（BERT嵌入）；（3）多模态融合（注意力机制）。数据集IEMOCAP，5折交叉验证。基线模型包括CNN-LSTM和纯Transformer。消融实验：移除模态注意力、移除跨模态对齐模块。使用McNemar检验评估统计显著性。"},
		{"数据分析方法选择", "数据分析与解释", 1, "我昨天跑代码的时候电脑突然蓝屏了，然后数据就丢了。我现在正在重新跑实验，但是GPU好像被别人占用了，排队要等很久。有没有办法不用GPU也能跑这个模型？"},
		{"结果解释准确性", "数据分析与解释", 3, "多模态融合模型WA达到74.8%，优于单模态语音（69.2%）和单模态文本（71.5%）。分析发现文本模态对中性情绪识别贡献较大，语音模态对高唤醒情绪（愤怒、快乐）识别效果更好。消融实验表明跨模态注意力模块贡献了约2.3个百分点的提升。"},
		{"问题提出新颖性", "批判性思维与创新", 1, "我觉得可以做一个APP，让用户对着手机说话然后告诉他们是什么情绪。这样就可以用在心理咨询或者智能客服里面。我觉得这个想法很有商业价值，可以创业。"},
		{"解决方案原创性", "批判性思维与创新", 3, "设计了基于跨模态注意力机制的融合方案。语音编码器使用Conformer架构，文本编码器使用BERT-base。两模态通过交叉注意力层进行对齐，再通过门控融合生成最终预测。参考了[Tsao et al., 2022]的跨模态注意力设计并进行了改进。"},
	},
	"liuyt@stu.com": {
		// 刘雨桐 - 异常输入 (2条L2 + 4条异常)
		{"文献检索策略", "文献检索与综述", 2, "在知网搜索了相关文献，找到了15篇左右的论文，主要看了其中3篇。"},
		{"实验方案合理性", "研究设计与实验", 1, ""},
		{"数据分析方法选择", "数据分析与解释", 1, "哈哈哈哈哈哈哈哈哈哈哈哈哈哈哈哈哈哈哈哈哈哈哈哈哈哈哈哈哈哈哈哈哈哈哈哈哈哈哈哈哈哈哈哈哈哈哈哈哈哈"},
		{"文献综述质量", "文献检索与综述", 1, "<script>alert('xss')</script>SELECT * FROM users WHERE 1=1; DROP TABLE evidences;"},
		{"结果解释准确性", "数据分析与解释", 2, "模型准确率是65%，比随机猜测好一些。loss曲线下降得还算正常。我觉得结果一般般，可能是数据集的问题。"},
		{"解决方案原创性", "批判性思维与创新", 1, "我参考了https://github.com/someone/speech-emotion这个项目的代码，直接clone下来跑了一下，效果还不错。具体实现细节请看那个项目的README。"},
	},
}

func main() {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	// 自动迁移
	db.AutoMigrate(&User{}, &Task{}, &StudentTask{}, &Evidence{})

	// 清空旧数据
	db.Exec("DELETE FROM evidences")
	db.Exec("DELETE FROM student_tasks")
	db.Exec("DELETE FROM tasks")
	db.Exec("DELETE FROM users WHERE email LIKE '%@stu.com' OR email LIKE '%@tea.com'")

	// 插入教师
	teacherUser := User{ID: teacher.ID, Email: teacher.Email, Password: teacher.Password, Name: teacher.Name, Role: "teacher", CreatedAt: time.Now()}
	if err := db.Create(&teacherUser).Error; err != nil {
		log.Fatalf("创建教师失败: %v", err)
	}
	fmt.Printf("✅ 教师已创建: %s (%s)\n", teacher.Name, teacher.Email)

	// 插入任务
	taskRecord := Task{
		ID:          task.ID,
		Name:        task.Name,
		Description: task.Description,
		CourseID:    task.CourseID,
		TeacherID:   teacher.ID,
		StartDate:   time.Now(),
		EndDate:     time.Now().AddDate(0, 3, 0),
		Status:      "active",
		CreatedAt:   time.Now(),
	}
	if err := db.Create(&taskRecord).Error; err != nil {
		log.Fatalf("创建任务失败: %v", err)
	}
	fmt.Printf("✅ 任务已创建: %s\n", task.Name)

	// 插入学生和证据
	totalEvidences := 0
	for _, s := range students {
		studentUser := User{ID: s.ID, Email: s.Email, Password: s.Password, Name: s.Name, Role: "student", CreatedAt: time.Now()}
		if err := db.Create(&studentUser).Error; err != nil {
			log.Fatalf("创建学生失败: %v", err)
		}

		// 创建学生-任务关联
		st := StudentTask{ID: uuid.New().String(), TaskID: task.ID, StudentID: s.ID, Status: "pending", Progress: 0, CreatedAt: time.Now()}
		if err := db.Create(&st).Error; err != nil {
			log.Fatalf("创建学生任务失败: %v", err)
		}

		// 插入证据
		evidences := studentEvidences[s.Email]
		count := 0
		for _, e := range evidences {
			ev := Evidence{
				ID:            uuid.New().String(),
				StudentTaskID: st.ID,
				Type:          "text",
				KBMName:       e.KBMName,
				KBMLevel:      e.KBMLevel,
				Content:       e.Content,
				CreatedAt:     time.Now(),
			}
			if err := db.Create(&ev).Error; err != nil {
				log.Printf("⚠️ 创建证据失败 [%s]: %v", s.Name, err)
				continue
			}
			count++
		}
		totalEvidences += count
		fmt.Printf("✅ 学生 %s (%s): %d 条证据已创建\n", s.Name, s.Email, count)
	}

	fmt.Printf("\n🎉 数据生成完成！共 %d 名学生, %d 条证据\n", len(students), totalEvidences)
	fmt.Println("\n📊 场景覆盖：")
	fmt.Println("   周博远: 8条高质量(L4) — 测试系统对优秀研究的识别")
	fmt.Println("   陈思涵: 6条良好(L3)+2条合格(L2) — 测试系统对中等水平的区分")
	fmt.Println("   林浩然: 4条合格(L2)+2条不合格(L1) — 测试系统对低质量的识别")
	fmt.Println("   张晓婷: 2条不合格(L1)+4条答非所问 — 测试系统对答非所问的识别")
	fmt.Println("   王铭宇: 2条L4+2条L3+3条答非所问 — 测试系统混合场景区分能力")
	fmt.Println("   刘雨桐: 2条L2+4条异常输入 — 测试系统容错和安全过滤能力")
}
