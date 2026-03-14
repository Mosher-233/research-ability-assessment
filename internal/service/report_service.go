package service

import (
	"context"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"research-ability-assessment/internal/models"
	"research-ability-assessment/pkg/utils"
	"strings"
	"time"
)

type ReportService struct {
	inferenceService *InferenceService
	resultRepo       interface {
		CreateInferenceResult(ctx context.Context, result *models.InferenceResult) error
		GetInferenceResultByStudentAndTask(ctx context.Context, studentID string, taskID string) (*models.InferenceResult, error)
	}
}

func NewReportService(inferenceService *InferenceService, resultRepo interface {
	CreateInferenceResult(ctx context.Context, result *models.InferenceResult) error
	GetInferenceResultByStudentAndTask(ctx context.Context, studentID string, taskID string) (*models.InferenceResult, error)
}) *ReportService {
	return &ReportService{
		inferenceService: inferenceService,
		resultRepo:       resultRepo,
	}
}

func (s *ReportService) GenerateReport(ctx context.Context, studentTaskID, studentID, taskID string) (*models.Report, error) {
	log.Printf("GenerateReport: 开始生成报告, StudentTaskID=%s, StudentID=%s, TaskID=%s", studentTaskID, studentID, taskID)

	inferenceResult, err := s.inferenceService.GetInferenceResultByStudentAndTask(ctx, studentID, taskID)
	if err != nil {
		log.Printf("GenerateReport: 没有找到现有推理结果，开始生成新的推理")
		inferenceResult, err = s.inferenceService.GenerateInference(ctx, &GenerateInferenceRequest{
			StudentTaskID: studentTaskID,
			StudentID:     studentID,
			TaskID:        taskID,
		})
		if err != nil {
			return nil, fmt.Errorf("生成推理结果失败: %w", err)
		}
	}

	classStats, err := s.inferenceService.GetClassStats(ctx, taskID)
	if err != nil {
		log.Printf("GenerateReport: 获取班级统计失败: %v", err)
		classStats = &ClassStats{
			ClassSize:         0,
			ClassAverage:      0,
			ClassMaxScore:     0,
			ClassMinScore:     0,
			DimensionAverages: make(map[string]float64),
		}
	}

	rank, percentile, err := s.inferenceService.CalculateRankAndPercentile(ctx, inferenceResult.OverallScore, taskID)
	if err != nil {
		log.Printf("GenerateReport: 计算排名失败: %v", err)
		rank = 0
		percentile = 0
	}

	strengths := s.getStrengths(inferenceResult.DimensionScores)
	weaknesses := s.getWeaknesses(inferenceResult.DimensionScores)
	detailedAnalysis := s.generateDetailedAnalysis(inferenceResult, classStats, rank, percentile)
	suggestions := s.generateSuggestions(inferenceResult, weaknesses)
	radarChartData := s.generateRadarChartData(inferenceResult, classStats)

	classComparison := &models.ClassComparisonData{
		ClassSize:         classStats.ClassSize,
		ClassAverage:      math.Round(classStats.ClassAverage*100) / 100,
		ClassMaxScore:     math.Round(classStats.ClassMaxScore*100) / 100,
		ClassMinScore:     math.Round(classStats.ClassMinScore*100) / 100,
		DimensionAverages: classStats.DimensionAverages,
	}

	report := &models.Report{
		ID:               utils.GenerateTaskID(),
		StudentTaskID:    studentTaskID,
		StudentID:        studentID,
		TaskID:           taskID,
		OverallScore:     inferenceResult.OverallScore,
		OverallLevel:     inferenceResult.OverallLevel,
		DimensionScores:  inferenceResult.DimensionScores,
		ClassComparison:  classComparison,
		Rank:             rank,
		Percentile:       percentile,
		Strengths:        strengths,
		Weaknesses:       weaknesses,
		DetailedAnalysis: detailedAnalysis,
		Suggestions:      suggestions,
		RadarChartData:   radarChartData,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	reportPath, err := s.saveReportToFile(report)
	if err != nil {
		log.Printf("GenerateReport: 保存报告文件失败: %v", err)
	} else {
		report.ReportPath = reportPath
	}

	log.Printf("GenerateReport: 报告生成成功, ReportID=%s", report.ID)
	return report, nil
}

func (s *ReportService) getStrengths(dimensionScores map[string]models.DimensionScore) []string {
	var strengths []string
	for _, score := range dimensionScores {
		if score.Score >= 80 {
			strengths = append(strengths, fmt.Sprintf("%s：表现优秀，得分%.1f", score.Name, score.Score))
		}
	}
	if len(strengths) == 0 {
		strengths = append(strengths, "整体表现良好，继续保持")
	}
	return strengths
}

func (s *ReportService) getWeaknesses(dimensionScores map[string]models.DimensionScore) []string {
	var weaknesses []string
	for _, score := range dimensionScores {
		if score.Score < 70 {
			weaknesses = append(weaknesses, fmt.Sprintf("%s：需要提升，得分%.1f", score.Name, score.Score))
		}
	}
	if len(weaknesses) == 0 {
		weaknesses = append(weaknesses, "无明显短板")
	}
	return weaknesses
}

func (s *ReportService) generateDetailedAnalysis(inferenceResult *models.InferenceResult, classStats *ClassStats, rank int, percentile float64) map[string]string {
	analysis := make(map[string]string)

	analysis["overview"] = fmt.Sprintf("学生在本次任务中的综合表现为%s，总得分为%.1f分。", inferenceResult.OverallLevel, inferenceResult.OverallScore)

	if classStats.ClassSize > 0 {
		analysis["class_comparison"] = fmt.Sprintf("班级平均分为%.1f分，该学生在班级中排名第%d位，超过了%.1f%%的同学。", classStats.ClassAverage, rank, percentile)
	} else {
		analysis["class_comparison"] = "暂无班级对比数据"
	}

	for dimID, score := range inferenceResult.DimensionScores {
		classAvg := classStats.DimensionAverages[dimID]
		if classAvg > 0 {
			if score.Score >= classAvg {
				analysis[dimID] = fmt.Sprintf("%s：得分%.1f，高于班级平均水平%.1f分", score.Name, score.Score, score.Score-classAvg)
			} else {
				analysis[dimID] = fmt.Sprintf("%s：得分%.1f，低于班级平均水平%.1f分", score.Name, score.Score, classAvg-score.Score)
			}
		} else {
			analysis[dimID] = fmt.Sprintf("%s：得分%.1f", score.Name, score.Score)
		}
	}

	return analysis
}

func (s *ReportService) generateSuggestions(inferenceResult *models.InferenceResult, weaknesses []string) []models.ImprovementSuggestion {
	var suggestions []models.ImprovementSuggestion

	resourceMap := map[string][]models.LearningResource{
		"literature_review": {
			{Type: "book", Title: "如何做好文献综述", Author: "陈向明", ISBN: "9787301302872"},
			{Type: "online", Title: "文献检索与管理", URL: "https://www.coursera.org/courses?query=literature%20review"},
		},
		"research_design": {
			{Type: "book", Title: "研究设计与方法", Author: "温忠麟", ISBN: "9787040557497"},
			{Type: "course", Title: "社会科学研究方法", Duration: "40小时"},
		},
		"data_analysis": {
			{Type: "book", Title: "数据分析与应用", Author: "薛薇", ISBN: "9787300290816"},
			{Type: "online", Title: "Python数据分析", URL: "https://www.kaggle.com/learn"},
		},
		"critical_thinking": {
			{Type: "book", Title: "批判性思维工具", Author: "理查德·保罗", ISBN: "9787111422358"},
			{Type: "course", Title: "批判性思维导论", Duration: "20小时"},
		},
	}

	priority := 1
	for dimID, score := range inferenceResult.DimensionScores {
		if score.Score < 75 {
			targetScore := math.Min(score.Score+15, 100)
			suggestion := models.ImprovementSuggestion{
				ID:            dimID,
				Dimension:     dimID,
				DimensionName: score.Name,
				CurrentScore:  score.Score,
				TargetScore:   targetScore,
				Suggestion:    s.getDimensionSuggestion(dimID, score.Score),
				ActionItems:   s.getActionItems(dimID),
				Resources:     resourceMap[dimID],
				Priority:      priority,
			}
			suggestions = append(suggestions, suggestion)
			priority++
		}
	}

	if len(suggestions) == 0 {
		suggestions = append(suggestions, models.ImprovementSuggestion{
			ID:            "general",
			Dimension:     "general",
			DimensionName: "综合提升",
			CurrentScore:  inferenceResult.OverallScore,
			TargetScore:   math.Min(inferenceResult.OverallScore+5, 100),
			Suggestion:    "继续保持现有优势，尝试在各个维度寻求更高突破",
			ActionItems: []string{
				"参与更高难度的研究项目",
				"尝试发表学术论文",
				"指导低年级同学",
			},
			Resources: []models.LearningResource{
				{Type: "book", Title: "研究是一门艺术", Author: "韦恩·C·布斯", ISBN: "9787513328832"},
			},
			Priority: 1,
		})
	}

	return suggestions
}

func (s *ReportService) getDimensionSuggestion(dimID string, score float64) string {
	suggestionMap := map[string]string{
		"literature_review": "建议加强文献检索技巧，学习如何系统梳理和分析相关研究，注重批判性地评价文献。",
		"research_design":   "建议深入学习研究设计方法，提高实验或调查方案的科学性和可行性，加强变量控制。",
		"data_analysis":     "建议学习更多数据分析方法，提高数据处理和可视化能力，注重结果的合理解释。",
		"critical_thinking": "建议多进行批判性思考练习，学会从多角度分析问题，培养创新思维能力。",
	}
	return suggestionMap[dimID]
}

func (s *ReportService) getActionItems(dimID string) []string {
	actionMap := map[string][]string{
		"literature_review": {
			"每周阅读3-5篇核心期刊论文",
			"练习撰写文献综述段落",
			"学习使用文献管理软件",
		},
		"research_design": {
			"分析5个优秀研究案例的设计",
			"练习设计小型研究方案",
			"请教老师或同学的反馈",
		},
		"data_analysis": {
			"完成5个数据分析实战项目",
			"学习1-2种数据分析工具",
			"练习撰写数据分析报告",
		},
		"critical_thinking": {
			"参加学术辩论或讨论活动",
			"每周写一篇批判性思考笔记",
			"阅读哲学和逻辑学相关书籍",
		},
	}
	return actionMap[dimID]
}

func (s *ReportService) generateRadarChartData(inferenceResult *models.InferenceResult, classStats *ClassStats) *models.RadarChartData {
	var labels []string
	var values []float64
	var classAverages []float64

	for dimID, score := range inferenceResult.DimensionScores {
		labels = append(labels, score.Name)
		values = append(values, score.Score)
		classAverages = append(classAverages, classStats.DimensionAverages[dimID])
	}

	return &models.RadarChartData{
		Labels:        labels,
		Values:        values,
		ClassAverages: classAverages,
	}
}

func (s *ReportService) saveReportToFile(report *models.Report) (string, error) {
	uploadsDir := "uploads/reports"
	if err := os.MkdirAll(uploadsDir, 0755); err != nil {
		return "", err
	}

	filename := fmt.Sprintf("report_%s.txt", report.ID)
	filepath := filepath.Join(uploadsDir, filename)

	var content strings.Builder
	content.WriteString("=")
	content.WriteString(strings.Repeat("=", 58))
	content.WriteString("=\n")
	content.WriteString(fmt.Sprintf("                          研究能力评价报告\n"))
	content.WriteString("=")
	content.WriteString(strings.Repeat("=", 58))
	content.WriteString("=\n\n")
	content.WriteString(fmt.Sprintf("报告编号: %s\n", report.ID))
	content.WriteString(fmt.Sprintf("学生ID: %s\n", report.StudentID))
	content.WriteString(fmt.Sprintf("任务ID: %s\n", report.TaskID))
	content.WriteString(fmt.Sprintf("生成时间: %s\n\n", report.CreatedAt.Format("2006-01-02 15:04:05")))

	content.WriteString("一、综合评价\n")
	content.WriteString(strings.Repeat("-", 60))
	content.WriteString("\n")
	content.WriteString(fmt.Sprintf("总体得分: %.1f分\n", report.OverallScore))
	content.WriteString(fmt.Sprintf("评价等级: %s\n", report.OverallLevel))
	if report.Rank > 0 {
		content.WriteString(fmt.Sprintf("班级排名: 第%d名\n", report.Rank))
		content.WriteString(fmt.Sprintf("超越比例: %.1f%%\n", report.Percentile))
	}
	content.WriteString("\n")

	content.WriteString("二、维度得分\n")
	content.WriteString(strings.Repeat("-", 60))
	content.WriteString("\n")
	for _, score := range report.DimensionScores {
		content.WriteString(fmt.Sprintf("%s: %.1f分 (%s)\n", score.Name, score.Score, score.Level))
	}
	content.WriteString("\n")

	content.WriteString("三、优势分析\n")
	content.WriteString(strings.Repeat("-", 60))
	content.WriteString("\n")
	for i, strength := range report.Strengths {
		content.WriteString(fmt.Sprintf("%d. %s\n", i+1, strength))
	}
	content.WriteString("\n")

	content.WriteString("四、待提升方向\n")
	content.WriteString(strings.Repeat("-", 60))
	content.WriteString("\n")
	for i, weakness := range report.Weaknesses {
		content.WriteString(fmt.Sprintf("%d. %s\n", i+1, weakness))
	}
	content.WriteString("\n")

	content.WriteString("五、改进建议\n")
	content.WriteString(strings.Repeat("-", 60))
	content.WriteString("\n")
	for _, suggestion := range report.Suggestions {
		content.WriteString(fmt.Sprintf("\n【%s】(优先级: %d)\n", suggestion.DimensionName, suggestion.Priority))
		content.WriteString(fmt.Sprintf("当前得分: %.1f → 目标得分: %.1f\n", suggestion.CurrentScore, suggestion.TargetScore))
		content.WriteString(fmt.Sprintf("建议: %s\n", suggestion.Suggestion))
		content.WriteString("行动项:\n")
		for _, item := range suggestion.ActionItems {
			content.WriteString(fmt.Sprintf("  - %s\n", item))
		}
	}
	content.WriteString("\n")

	content.WriteString("=")
	content.WriteString(strings.Repeat("=", 58))
	content.WriteString("=\n")
	content.WriteString("                        报告生成完毕\n")
	content.WriteString("=")
	content.WriteString(strings.Repeat("=", 58))
	content.WriteString("=\n")

	if err := os.WriteFile(filepath, []byte(content.String()), 0644); err != nil {
		return "", err
	}

	return filepath, nil
}

func (s *ReportService) GenerateStudentReport(ctx context.Context, studentID string, taskID string) (*models.InferenceResult, error) {
	return s.inferenceService.GetInferenceResultByStudentAndTask(ctx, studentID, taskID)
}

func (s *ReportService) GenerateTaskReport(ctx context.Context, taskID string) ([]models.InferenceResult, error) {
	return s.inferenceService.GetInferenceResultsByTaskID(ctx, taskID)
}