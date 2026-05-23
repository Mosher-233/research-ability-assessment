package models

const (
	DimLiteratureReview  = "literature_review"
	DimResearchDesign    = "research_design"
	DimDataAnalysis      = "data_analysis"
	DimCriticalThinking  = "critical_thinking"

	DimNameLiteratureReview = "文献综述"
	DimNameResearchDesign   = "研究设计"
	DimNameDataAnalysis     = "数据分析"
	DimNameCriticalThinking = "批判性思维"
)

var DefaultDimensions = []Dimension{
	{
		ID:          DimLiteratureReview,
		Name:        DimNameLiteratureReview,
		Description: "文献检索与综述撰写能力，涵盖文献检索策略制定、文献质量评估与批判性分析、文献综述的组织与撰写",
		Weight:      0.25,
	},
	{
		ID:          DimResearchDesign,
		Name:        DimNameResearchDesign,
		Description: "研究方案设计与实验规划能力，涵盖实验方案的科学性、变量的有效控制、实验实施的规范性",
		Weight:      0.25,
	},
	{
		ID:          DimDataAnalysis,
		Name:        DimNameDataAnalysis,
		Description: "数据处理与统计分析能力，涵盖数据分析方法的合理选择、分析过程的严谨性和结果解释的准确性",
		Weight:      0.25,
	},
	{
		ID:          DimCriticalThinking,
		Name:        DimNameCriticalThinking,
		Description: "批判性思考与创新思维能力，涵盖问题的独立提出、多角度分析和解决方案的原创性",
		Weight:      0.25,
	},
}

var KBMMapping = map[string]string{
	"文献检索策略":   DimLiteratureReview,
	"文献综述质量":   DimLiteratureReview,
	"文献批判性分析": DimLiteratureReview,

	"实验方案合理性": DimResearchDesign,
	"变量控制":     DimResearchDesign,
	"实验实施质量":   DimResearchDesign,

	"数据分析方法选择": DimDataAnalysis,
	"结果解释准确性":   DimDataAnalysis,

	"问题提出新颖性": DimCriticalThinking,
	"解决方案原创性":   DimCriticalThinking,
}

func GetDimensionByKBM(kbmName string) string {
	if dim, ok := KBMMapping[kbmName]; ok {
		return dim
	}
	return DimCriticalThinking
}

func GetDimensionName(dimID string) string {
	names := map[string]string{
		DimLiteratureReview: DimNameLiteratureReview,
		DimResearchDesign:   DimNameResearchDesign,
		DimDataAnalysis:     DimNameDataAnalysis,
		DimCriticalThinking: DimNameCriticalThinking,
	}
	if name, ok := names[dimID]; ok {
		return name
	}
	return dimID
}

func GetDimensionWeight(dimID string) float64 {
	weights := map[string]float64{
		DimLiteratureReview: 0.25,
		DimResearchDesign:   0.25,
		DimDataAnalysis:     0.25,
		DimCriticalThinking: 0.25,
	}
	if weight, ok := weights[dimID]; ok {
		return weight
	}
	return 0.25
}

func GetLevelFromScore(score float64) string {
	switch {
	case score >= 90:
		return "优秀"
	case score >= 75:
		return "良好"
	case score >= 60:
		return "合格"
	default:
		return "不合格"
	}
}

var KBMLabels = []string{
	"文献检索策略",
	"文献综述质量",
	"文献批判性分析",
	"实验方案合理性",
	"变量控制",
	"实验实施质量",
	"数据分析方法选择",
	"结果解释准确性",
	"问题提出新颖性",
	"解决方案原创性",
}

var KBMToLabel map[string]string

func init() {
	KBMToLabel = make(map[string]string)
	for _, label := range KBMLabels {
		KBMToLabel[label] = label
	}
}
