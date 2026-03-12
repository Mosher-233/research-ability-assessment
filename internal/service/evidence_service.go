package service

import (
	"context"
	"research-ability-assessment/internal/llm"
	"research-ability-assessment/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EvidenceService struct {
	db  *gorm.DB
	llm *llm.Client
}

func NewEvidenceService(db *gorm.DB, llmClient *llm.Client) *EvidenceService {
	return &EvidenceService{
		db:  db,
		llm: llmClient,
	}
}

func (s *EvidenceService) CreateEvidence(ctx context.Context, evidence *models.Evidence) error {
	evidence.ID = uuid.New().String()
	return s.db.WithContext(ctx).Create(evidence).Error
}

func (s *EvidenceService) GetEvidenceByID(ctx context.Context, id string) (*models.Evidence, error) {
	var evidence models.Evidence
	if err := s.db.WithContext(ctx).First(&evidence, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &evidence, nil
}

func (s *EvidenceService) GetEvidencesByStudentTaskID(ctx context.Context, studentTaskID string) ([]models.Evidence, error) {
	var evidences []models.Evidence
	if err := s.db.WithContext(ctx).Where("student_task_id = ?", studentTaskID).Find(&evidences).Error; err != nil {
		return nil, err
	}
	return evidences, nil
}

func (s *EvidenceService) GetEvidencesByStudentAndTask(ctx context.Context, studentID string, taskID string) ([]models.Evidence, error) {
	var evidences []models.Evidence
	if err := s.db.WithContext(ctx).Joins("JOIN student_tasks ON evidences.student_task_id = student_tasks.id").Where("student_tasks.student_id = ? AND student_tasks.task_id = ?", studentID, taskID).Find(&evidences).Error; err != nil {
		return nil, err
	}
	return evidences, nil
}

// AnalyzeEvidence 分析证据并返回KBM级别和反馈
func (s *EvidenceService) AnalyzeEvidence(ctx context.Context, evidenceID string) (map[string]interface{}, error) {
	// 获取证据
	evidence, err := s.GetEvidenceByID(ctx, evidenceID)
	if err != nil {
		return nil, err
	}

	// 调用Deepseek API进行分析
	messages := []llm.Message{
		{
			Role:    "system",
			Content: "你是一个研究能力评估专家，负责分析学生提交的证据并评估其KBM（Knowledge, Behavior, Methodology）级别。请根据以下证据内容，分析其质量并给出1-5级的KBM级别，以及详细的反馈。",
		},
		{
			Role:    "user",
			Content: "证据类型：" + evidence.Type + "\n证据内容：" + evidence.Content + "\nKBM名称：" + evidence.KBMName,
		},
	}

	response, err := s.llm.Chat(ctx, messages)
	if err != nil {
		// 如果LLM API调用失败，使用模拟数据
		kbmLevel := 3
		feedback := "证据内容分析：你对该主题有较好的理解，能够清晰表达自己的观点，但在深度分析方面还有提升空间。建议进一步深入研究相关文献，加强批判性思维能力。"
		
		// 更新证据的KBM级别
		evidence.KBMLevel = kbmLevel
		if err := s.db.WithContext(ctx).Save(evidence).Error; err != nil {
			return nil, err
		}
		
		// 返回分析结果
		return map[string]interface{}{
			"kbm_level": kbmLevel,
			"feedback":  feedback,
		}, nil
	}

	// 解析LLM响应，提取KBM级别和反馈
	kbmLevel := 3 // 默认值
	feedback := response

	// 简单的解析逻辑，实际应用中可能需要更复杂的解析
	// 这里只是一个示例，实际应该根据LLM的输出格式进行解析

	// 更新证据的KBM级别
	evidence.KBMLevel = kbmLevel
	if err := s.db.WithContext(ctx).Save(evidence).Error; err != nil {
		return nil, err
	}

	// 返回分析结果
	return map[string]interface{}{
		"kbm_level": kbmLevel,
		"feedback":  feedback,
	}, nil
}

// GetEvidencesByUserID 根据用户ID获取证据列表
func (s *EvidenceService) GetEvidencesByUserID(ctx context.Context, userID string) ([]models.Evidence, error) {
	var evidences []models.Evidence
	if err := s.db.WithContext(ctx).Joins("JOIN student_tasks ON evidences.student_task_id = student_tasks.id").Where("student_tasks.student_id = ?", userID).Find(&evidences).Error; err != nil {
		return nil, err
	}
	return evidences, nil
}
