package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"research-ability-assessment/internal/llm"
	"research-ability-assessment/internal/models"
	"research-ability-assessment/internal/repository/postgres"
	"research-ability-assessment/pkg/utils"
	"strings"
	"time"

	"gorm.io/datatypes"
)

type InferenceService struct {
	resultRepo      *postgres.ResultRepo
	evidenceService *EvidenceService
	llmClient       *llm.Client
}

func NewInferenceService(resultRepo *postgres.ResultRepo, evidenceService *EvidenceService) *InferenceService {
	return &InferenceService{
		resultRepo:      resultRepo,
		evidenceService: evidenceService,
	}
}

func NewInferenceServiceWithLLM(resultRepo *postgres.ResultRepo, evidenceService *EvidenceService, llmClient *llm.Client) *InferenceService {
	return &InferenceService{
		resultRepo:      resultRepo,
		evidenceService: evidenceService,
		llmClient:       llmClient,
	}
}

func (s *InferenceService) SetLLMClient(llmClient *llm.Client) {
	s.llmClient = llmClient
}

func (s *InferenceService) CreateInferenceResult(ctx context.Context, result *models.InferenceResult) error {
	result.ID = utils.GenerateTaskID()
	return s.resultRepo.CreateInferenceResult(ctx, result)
}

func (s *InferenceService) GetInferenceResultByID(ctx context.Context, id string) (*models.InferenceResult, error) {
	return s.resultRepo.GetInferenceResultByID(ctx, id)
}

func (s *InferenceService) GetInferenceResultsByTaskID(ctx context.Context, taskID string) ([]models.InferenceResult, error) {
	return s.resultRepo.GetInferenceResultsByTaskID(ctx, taskID)
}

func (s *InferenceService) GetInferenceResultByStudentAndTask(ctx context.Context, studentID string, taskID string) (*models.InferenceResult, error) {
	return s.resultRepo.GetInferenceResultByStudentAndTask(ctx, studentID, taskID)
}

func (s *InferenceService) GetEvidencesByStudentAndTask(ctx context.Context, studentID string, taskID string) ([]models.Evidence, error) {
	return s.evidenceService.GetEvidencesByStudentAndTask(ctx, studentID, taskID)
}

func (s *InferenceService) GetAllInferenceResults(ctx context.Context) ([]models.InferenceResult, error) {
	return s.resultRepo.GetAllInferenceResults(ctx)
}

func (s *InferenceService) GetInferenceResultsByStudentID(ctx context.Context, studentID string) ([]models.InferenceResult, error) {
	return s.resultRepo.GetInferenceResultsByStudentID(ctx, studentID)
}

type GenerateInferenceRequest struct {
	StudentTaskID string
	StudentID     string
	TaskID        string
	StudentInfo   *models.User
	TaskInfo      *models.Task
}

func (s *InferenceService) GenerateInference(ctx context.Context, req *GenerateInferenceRequest) (*models.InferenceResult, error) {
	log.Printf("GenerateInference: 开始生成推理结果, StudentTaskID=%s", req.StudentTaskID)

	evidences, err := s.evidenceService.GetEvidencesByStudentTaskID(ctx, req.StudentTaskID)
	if err != nil {
		log.Printf("GenerateInference: 获取证据失败: %v", err)
		return nil, fmt.Errorf("获取证据失败: %w", err)
	}

	log.Printf("GenerateInference: 找到 %d 个证据", len(evidences))

	if len(evidences) == 0 {
		log.Printf("GenerateInference: 没有找到证据，返回错误")
		return nil, fmt.Errorf("没有找到证据，无法生成评估结果")
	}

	dimensions := s.getDefaultDimensions()
	dimensionScores := make(map[string]models.DimensionScore)
	var totalWeightedScore float64
	var totalWeight float64
	var totalConfidence float64

	for _, dim := range dimensions {
		score, conf := s.calculateDimensionScore(evidences, dim)
		level := s.getLevelFromScore(score)

		dimensionScore := models.DimensionScore{
			Name:        dim.Name,
			Score:       score,
			Level:       level,
			Details:     "",
			EvidenceIDs: s.getEvidenceIDsForDimension(evidences, dim.ID),
		}

		dimensionScores[dim.ID] = dimensionScore
		totalWeightedScore += score * dim.Weight
		totalWeight += dim.Weight
		totalConfidence += conf
	}

	overallScore := totalWeightedScore / totalWeight
	overallLevel := s.getLevelFromScore(overallScore)

	dimensionScoresJSON, err := json.Marshal(dimensionScores)
	if err != nil {
		log.Printf("GenerateInference: 序列化维度得分失败: %v", err)
		return nil, fmt.Errorf("序列化维度得分失败: %w", err)
	}

	result := &models.InferenceResult{
		ID:              utils.GenerateTaskID(),
		StudentID:       req.StudentID,
		TaskID:          req.TaskID,
		OverallScore:    math.Round(overallScore*100) / 100,
		OverallLevel:    overallLevel,
		DimensionScores: datatypes.JSON(dimensionScoresJSON),
		Reasoning:       s.generateReasoning(overallScore, overallLevel, dimensionScores),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	log.Printf("GenerateInference: 准备保存结果, ID=%s, OverallScore=%.2f, OverallLevel=%s",
		result.ID, result.OverallScore, result.OverallLevel)

	if err := s.resultRepo.CreateInferenceResult(ctx, result); err != nil {
		log.Printf("GenerateInference: 保存结果失败: %v", err)
		return nil, fmt.Errorf("保存结果失败: %w", err)
	}

	log.Printf("GenerateInference: 推理结果生成成功")
	return result, nil
}

func (s *InferenceService) getDefaultDimensions() []models.Dimension {
	return []models.Dimension{
		{
			ID:          "literature_review",
			Name:        "文献综述",
			Description: "文献检索、综述撰写能力",
			Weight:      0.25,
		},
		{
			ID:          "research_design",
			Name:        "研究设计",
			Description: "研究方案设计、实验规划能力",
			Weight:      0.25,
		},
		{
			ID:          "data_analysis",
			Name:        "数据分析",
			Description: "数据处理、统计分析能力",
			Weight:      0.25,
		},
		{
			ID:          "critical_thinking",
			Name:        "批判性思维",
			Description: "批判性思考、创新思维能力",
			Weight:      0.25,
		},
	}
}

func (s *InferenceService) calculateDimensionScore(evidences []models.Evidence, dim models.Dimension) (float64, float64) {
	var totalScore float64
	var count int

	for _, evidence := range evidences {
		if evidence.KBMName == dim.ID || strings.Contains(evidence.Type, dim.ID) {
			if evidence.KBMLevel > 0 {
				score := float64(evidence.KBMLevel) * 20
				totalScore += score
				count++
			} else {
				totalScore += 60
				count++
			}
		}
	}

	if count == 0 {
		return 50, 0.5
	}

	avgScore := totalScore / float64(count)
	confidence := 0.7 + (float64(count) * 0.05)
	if confidence > 0.95 {
		confidence = 0.95
	}

	return math.Min(avgScore, 100), confidence
}

func (s *InferenceService) getLevelFromScore(score float64) string {
	switch {
	case score >= 90:
		return "优秀"
	case score >= 80:
		return "良好"
	case score >= 70:
		return "中等"
	case score >= 60:
		return "及格"
	default:
		return "不及格"
	}
}

func (s *InferenceService) getEvidenceIDsForDimension(evidences []models.Evidence, dimID string) []string {
	var ids []string
	for _, evidence := range evidences {
		if evidence.KBMName == dimID || strings.Contains(evidence.Type, dimID) {
			ids = append(ids, evidence.ID)
		}
	}
	return ids
}

func (s *InferenceService) generateReasoning(overallScore float64, overallLevel string, dimensionScores map[string]models.DimensionScore) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("基于收集到的 %d 个证据，对学生的研究能力进行了综合评估。",
		len(dimensionScores)))
	sb.WriteString(fmt.Sprintf("总体得分为 %.2f，等级为%s。", overallScore, overallLevel))

	var strengths []string
	var weaknesses []string

	for _, score := range dimensionScores {
		if score.Score >= 80 {
			strengths = append(strengths, fmt.Sprintf("%s(%.1f分)", score.Name, score.Score))
		} else if score.Score < 70 {
			weaknesses = append(weaknesses, fmt.Sprintf("%s(%.1f分)", score.Name, score.Score))
		}
	}

	if len(strengths) > 0 {
		sb.WriteString(fmt.Sprintf(" 优势维度：%s。", strings.Join(strengths, "、")))
	}

	if len(weaknesses) > 0 {
		sb.WriteString(fmt.Sprintf(" 待提升维度：%s。", strings.Join(weaknesses, "、")))
	}

	sb.WriteString(" 建议学生在保持优势的同时，针对待提升维度进行重点改进。")

	return sb.String()
}

type ClassStats struct {
	ClassSize         int
	ClassAverage      float64
	ClassMaxScore     float64
	ClassMinScore     float64
	DimensionAverages map[string]float64
}

func (s *InferenceService) GetClassStats(ctx context.Context, taskID string) (*ClassStats, error) {
	results, err := s.resultRepo.GetInferenceResultsByTaskID(ctx, taskID)
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return &ClassStats{
			ClassSize:         0,
			ClassAverage:      0,
			ClassMaxScore:     0,
			ClassMinScore:     0,
			DimensionAverages: make(map[string]float64),
		}, nil
	}

	var totalScore float64
	maxScore := results[0].OverallScore
	minScore := results[0].OverallScore
	dimensionTotals := make(map[string]float64)
	dimensionCounts := make(map[string]int)

	for _, result := range results {
		totalScore += result.OverallScore
		if result.OverallScore > maxScore {
			maxScore = result.OverallScore
		}
		if result.OverallScore < minScore {
			minScore = result.OverallScore
		}

		var dimensionScores map[string]models.DimensionScore
		if len(result.DimensionScores) > 0 {
			_ = json.Unmarshal(result.DimensionScores, &dimensionScores)
			for id, score := range dimensionScores {
				dimensionTotals[id] += score.Score
				dimensionCounts[id]++
			}
		}
	}

	dimensionAverages := make(map[string]float64)
	for id, total := range dimensionTotals {
		if count := dimensionCounts[id]; count > 0 {
			dimensionAverages[id] = total / float64(count)
		}
	}

	return &ClassStats{
		ClassSize:         len(results),
		ClassAverage:      totalScore / float64(len(results)),
		ClassMaxScore:     maxScore,
		ClassMinScore:     minScore,
		DimensionAverages: dimensionAverages,
	}, nil
}

func (s *InferenceService) CalculateRankAndPercentile(ctx context.Context, studentScore float64, taskID string) (int, float64, error) {
	results, err := s.resultRepo.GetInferenceResultsByTaskID(ctx, taskID)
	if err != nil {
		return 0, 0, err
	}

	if len(results) == 0 {
		return 0, 0, nil
	}

	rank := 1
	belowOrEqual := 0

	for _, result := range results {
		if result.OverallScore > studentScore {
			rank++
		}
		if result.OverallScore <= studentScore {
			belowOrEqual++
		}
	}

	percentile := float64(belowOrEqual) / float64(len(results)) * 100

	return rank, math.Round(percentile*100) / 100, nil
}

func (s *InferenceService) GenerateInferenceWithLLM(ctx context.Context, req *GenerateInferenceRequest) (*models.InferenceResult, error) {
	if s.llmClient == nil {
		log.Printf("GenerateInferenceWithLLM: LLM客户端未设置，使用简化方法")
		return s.GenerateInference(ctx, req)
	}

	log.Printf("GenerateInferenceWithLLM: 使用LLM生成推理结果")

	evidences, err := s.evidenceService.GetEvidencesByStudentTaskID(ctx, req.StudentTaskID)
	if err != nil {
		return nil, fmt.Errorf("获取证据失败: %w", err)
	}

	if len(evidences) == 0 {
		return nil, fmt.Errorf("没有找到证据，无法生成评估结果")
	}

	evidenceContext := s.buildEvidenceContext(evidences)

	messages := []llm.Message{
		{
			Role: "system",
			Content: `你是一个专业的大学生研究能力评价专家，负责基于证据对学生的研究能力进行综合评估。

评分标准：
- 优秀 (90-100分)：在该维度表现出色，展现出高水平的研究能力
- 良好 (80-89分)：在该维度表现较好，具备较强的研究能力
- 中等 (70-79分)：在该维度表现一般，具备基本的研究能力
- 及格 (60-69分)：在该维度表现较弱，需要进一步提升
- 不及格 (0-59分)：在该维度表现不足，存在明显缺陷

评估维度：
1. 文献综述 (权重0.25)：文献检索、综述撰写能力
2. 研究设计 (权重0.25)：研究方案设计、实验规划能力
3. 数据分析 (权重0.25)：数据处理、统计分析能力
4. 批判性思维 (权重0.25)：批判性思考、创新思维能力

请以JSON格式返回评估结果，格式如下：
{
  "dimension_scores": {
    "literature_review": {"score": 85.5, "level": "良好", "reasoning": "评分理由"},
    "research_design": {"score": 80.0, "level": "良好", "reasoning": "评分理由"},
    "data_analysis": {"score": 75.5, "level": "中等", "reasoning": "评分理由"},
    "critical_thinking": {"score": 70.0, "level": "中等", "reasoning": "评分理由"}
  },
  "overall_reasoning": "总体评价理由"
}`,
		},
		{
			Role:    "user",
			Content: evidenceContext,
		},
	}

	response, err := s.llmClient.Chat(ctx, messages)
	if err != nil {
		log.Printf("GenerateInferenceWithLLM: LLM调用失败，使用简化方法: %v", err)
		return s.GenerateInference(ctx, req)
	}

	result, err := s.parseLLMResponse(response, req)
	if err != nil {
		log.Printf("GenerateInferenceWithLLM: 解析LLM响应失败，使用简化方法: %v", err)
		return s.GenerateInference(ctx, req)
	}

	log.Printf("GenerateInferenceWithLLM: 准备保存结果, ID=%s, OverallScore=%.2f, OverallLevel=%s",
		result.ID, result.OverallScore, result.OverallLevel)

	if err := s.resultRepo.CreateInferenceResult(ctx, result); err != nil {
		log.Printf("GenerateInferenceWithLLM: 保存结果失败: %v", err)
		return nil, fmt.Errorf("保存结果失败: %w", err)
	}

	log.Printf("GenerateInferenceWithLLM: 推理结果生成成功")
	return result, nil
}

func (s *InferenceService) buildEvidenceContext(evidences []models.Evidence) string {
	var sb strings.Builder

	sb.WriteString("学生证据信息：\n\n")

	for i, evidence := range evidences {
		sb.WriteString(fmt.Sprintf("证据 %d:\n", i+1))
		sb.WriteString(fmt.Sprintf("  ID: %s\n", evidence.ID))
		sb.WriteString(fmt.Sprintf("  类型: %s\n", evidence.Type))
		sb.WriteString(fmt.Sprintf("  KBM名称: %s\n", evidence.KBMName))
		if evidence.KBMLevel > 0 {
			sb.WriteString(fmt.Sprintf("  KBM级别: %d\n", evidence.KBMLevel))
		}
		if evidence.Content != "" {
			sb.WriteString(fmt.Sprintf("  内容: %s\n", evidence.Content))
		}
		if evidence.FileName != "" {
			sb.WriteString(fmt.Sprintf("  文件: %s\n", evidence.FileName))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

type LLMResponse struct {
	DimensionScores map[string]struct {
		Score     float64 `json:"score"`
		Level     string  `json:"level"`
		Reasoning string  `json:"reasoning"`
	} `json:"dimension_scores"`
	OverallReasoning string `json:"overall_reasoning"`
}

func (s *InferenceService) parseLLMResponse(response string, req *GenerateInferenceRequest) (*models.InferenceResult, error) {
	jsonStart := strings.Index(response, "{")
	jsonEnd := strings.LastIndex(response, "}")

	if jsonStart == -1 || jsonEnd == -1 {
		return nil, fmt.Errorf("未找到JSON响应")
	}

	jsonStr := response[jsonStart : jsonEnd+1]

	var llmResp LLMResponse
	if err := json.Unmarshal([]byte(jsonStr), &llmResp); err != nil {
		return nil, fmt.Errorf("解析JSON失败: %w", err)
	}

	dimensionScores := make(map[string]models.DimensionScore)
	var totalWeightedScore float64
	dimensions := s.getDefaultDimensions()
	weightMap := make(map[string]float64)
	nameMap := make(map[string]string)
	for _, dim := range dimensions {
		weightMap[dim.ID] = dim.Weight
		nameMap[dim.ID] = dim.Name
	}

	for dimID, score := range llmResp.DimensionScores {
		weight := weightMap[dimID]
		if weight == 0 {
			weight = 0.25
		}

		dimName := nameMap[dimID]
		if dimName == "" {
			dimName = dimID
		}

		dimensionScores[dimID] = models.DimensionScore{
			Name:    dimName,
			Score:   score.Score,
			Level:   score.Level,
			Details: "",
		}

		totalWeightedScore += score.Score * weight
	}

	overallScore := totalWeightedScore
	overallLevel := s.getLevelFromScore(overallScore)

	dimensionScoresJSON, err := json.Marshal(dimensionScores)
	if err != nil {
		log.Printf("parseLLMResponse: 序列化维度得分失败: %v", err)
		return nil, fmt.Errorf("序列化维度得分失败: %w", err)
	}

	result := &models.InferenceResult{
		ID:              utils.GenerateTaskID(),
		StudentID:       req.StudentID,
		TaskID:          req.TaskID,
		OverallScore:    math.Round(overallScore*100) / 100,
		OverallLevel:    overallLevel,
		DimensionScores: datatypes.JSON(dimensionScoresJSON),
		Reasoning:       llmResp.OverallReasoning,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	return result, nil
}
