package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"github.com/Mosher-233/research-ability-assessment/internal/llm"
	"github.com/Mosher-233/research-ability-assessment/internal/models"
	"github.com/Mosher-233/research-ability-assessment/internal/repository/postgres"
	"github.com/Mosher-233/research-ability-assessment/pkg/cache"
	"github.com/Mosher-233/research-ability-assessment/pkg/utils"
	"strings"
	"time"

	"gorm.io/datatypes"
)

type InferenceService struct {
	resultRepo      *postgres.ResultRepo
	evidenceService *EvidenceService
	llmClient       *llm.Client
	cache           *cache.RedisCache
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

func (s *InferenceService) SetCache(c *cache.RedisCache) {
	s.cache = c
}

func (s *InferenceService) CreateInferenceResult(ctx context.Context, result *models.InferenceResult) error {
	result.ID = utils.GenerateEvidenceID()
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

func (s *InferenceService) GetInferenceResultsByTeacherID(ctx context.Context, teacherID string) ([]models.InferenceResult, error) {
	return s.resultRepo.GetInferenceResultsByTeacherID(ctx, teacherID)
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
	return models.DefaultDimensions
}

func (s *InferenceService) calculateDimensionScore(evidences []models.Evidence, dim models.Dimension) (float64, float64) {
	var totalScore float64
	var count int

	for _, evidence := range evidences {
		mappedDim := models.GetDimensionByKBM(evidence.KBMName)
		if mappedDim == dim.ID {
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
	return models.GetLevelFromScore(score)
}

func (s *InferenceService) getEvidenceIDsForDimension(evidences []models.Evidence, dimID string) []string {
	var ids []string
	for _, evidence := range evidences {
		if models.GetDimensionByKBM(evidence.KBMName) == dimID {
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
	cacheKey := "class_stats:" + taskID

	if s.cache != nil && s.cache.IsAvailable() {
		var cached ClassStats
		if found, _ := s.cache.Get(ctx, cacheKey, &cached); found {
			return &cached, nil
		}
	}

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

	stats := &ClassStats{
		ClassSize:         len(results),
		ClassAverage:      totalScore / float64(len(results)),
		ClassMaxScore:     maxScore,
		ClassMinScore:     minScore,
		DimensionAverages: dimensionAverages,
	}

	if s.cache != nil && s.cache.IsAvailable() {
		s.cache.Set(ctx, cacheKey, stats, 5*time.Minute)
	}

	return stats, nil
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
		log.Printf("GenerateInferenceWithLLM: 解析LLM响应失败(第1次)，重试中: %v", err)
		// Retry once
		response, err = s.llmClient.Chat(ctx, messages)
		if err != nil {
			log.Printf("GenerateInferenceWithLLM: 重试LLM调用失败，使用简化方法: %v", err)
			return s.GenerateInference(ctx, req)
		}
		result, err = s.parseLLMResponse(response, req)
		if err != nil {
			log.Printf("GenerateInferenceWithLLM: 重试解析仍失败，使用简化方法: %v", err)
			return s.GenerateInference(ctx, req)
		}
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

	sb.WriteString("## 学生证据列表\n\n")
	sb.WriteString(fmt.Sprintf("共提交 %d 条证据。\n\n", len(evidences)))

	// Group evidences by KBM dimension
	dimGroups := make(map[string][]models.Evidence)
	for _, ev := range evidences {
		dim := models.GetDimensionByKBM(ev.KBMName)
		dimGroups[dim] = append(dimGroups[dim], ev)
	}

	for dimID, dimEvidences := range dimGroups {
		dimName := models.GetDimensionName(dimID)
		sb.WriteString(fmt.Sprintf("### %s维度 (%d条证据)\n\n", dimName, len(dimEvidences)))
		for i, ev := range dimEvidences {
			sb.WriteString(fmt.Sprintf("**证据 %d** | KBM: %s | 类型: %s | 来源: %s", i+1, ev.KBMName, ev.Type, ev.SourceType))
			if ev.FileName != "" {
				sb.WriteString(fmt.Sprintf(" | 文件: %s (%s)", ev.FileName, ev.FileType))
			}
			sb.WriteString("\n\n")
			if ev.Content != "" {
				content := ev.Content
				if len([]rune(content)) > 3000 {
					content = string([]rune(content)[:3000]) + "\n...(内容截断)"
				}
				sb.WriteString(content)
				sb.WriteString("\n\n")
			}
			sb.WriteString("---\n\n")
		}
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
	jsonStr := extractJSON(response)
	if jsonStr == "" {
		return nil, fmt.Errorf("未在LLM响应中找到有效JSON")
	}

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
			Details: score.Reasoning,
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

// extractJSON extracts a valid JSON object from text that may contain markdown wrappers.
func extractJSON(s string) string {
	// Try to find ```json code blocks first
	start := strings.Index(s, "```json")
	if start >= 0 {
		rest := s[start+7:]
		if end := strings.Index(rest, "```"); end >= 0 {
			s = rest[:end]
		}
	} else if start = strings.Index(s, "```"); start >= 0 {
		rest := s[start+3:]
		if end := strings.Index(rest, "```"); end >= 0 {
			s = rest[:end]
		}
	}

	// Find the outermost JSON object
	start = strings.Index(s, "{")
	if start < 0 {
		return ""
	}
	depth := 0
	for i := start; i < len(s); i++ {
		switch s[i] {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return s[start : i+1]
			}
		}
	}
	return ""
}
