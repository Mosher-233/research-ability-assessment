package service

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"research-ability-assessment/internal/llm"
	"research-ability-assessment/internal/models"
	"research-ability-assessment/pkg/utils"

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

func (s *EvidenceService) GetDB() *gorm.DB {
	return s.db
}

func (s *EvidenceService) CreateEvidence(ctx context.Context, evidence *models.Evidence) error {
	evidence.ID = utils.GenerateEvidenceID()
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

// GetEvidencesByTeacherID 根据教师ID获取所有相关证据（教师可以看到自己发布任务的所有学生证据）
func (s *EvidenceService) GetEvidencesByTeacherID(ctx context.Context, teacherID string) ([]models.Evidence, error) {
	var evidences []models.Evidence
	if err := s.db.WithContext(ctx).
		Joins("JOIN student_tasks ON evidences.student_task_id = student_tasks.id").
		Joins("JOIN tasks ON student_tasks.task_id = tasks.id").
		Where("tasks.teacher_id = ?", teacherID).
		Find(&evidences).Error; err != nil {
		return nil, err
	}
	return evidences, nil
}

// GetEvidencesWithDetailsByTeacherID 获取教师任务的详细证据列表（包含学生和任务信息）
type EvidenceWithDetails struct {
	models.Evidence
	StudentID   string `json:"student_id"`
	StudentName string `json:"student_name"`
	TaskID      string `json:"task_id"`
	TaskName    string `json:"task_name"`
}

func (s *EvidenceService) GetEvidencesWithDetailsByTeacherID(ctx context.Context, teacherID string) ([]EvidenceWithDetails, error) {
	var results []EvidenceWithDetails
	if err := s.db.WithContext(ctx).
		Table("evidences").
		Select("evidences.*, student_tasks.student_id, users.name as student_name, tasks.id as task_id, tasks.name as task_name").
		Joins("JOIN student_tasks ON evidences.student_task_id = student_tasks.id").
		Joins("JOIN tasks ON student_tasks.task_id = tasks.id").
		Joins("JOIN users ON student_tasks.student_id = users.id").
		Where("tasks.teacher_id = ?", teacherID).
		Find(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}

// AnalyzeEvidence 分析证据并返回KBM级别和反馈，同时创建Feedback记录
func (s *EvidenceService) AnalyzeEvidence(ctx context.Context, evidenceID string) (map[string]interface{}, error) {
	log.Printf("AnalyzeEvidence: 开始分析证据, ID=%s", evidenceID)
	
	evidence, err := s.GetEvidenceByID(ctx, evidenceID)
	if err != nil {
		log.Printf("AnalyzeEvidence: 获取证据失败: %v", err)
		return nil, err
	}

	log.Printf("AnalyzeEvidence: 证据信息 - 类型=%s, KBM名称=%s, 内容长度=%d", 
		evidence.Type, evidence.KBMName, len(evidence.Content))

	messages := []llm.Message{
		{
			Role:    "system",
			Content: "你是一个研究能力评估专家，负责分析学生提交的证据并评估其KBM（Knowledge, Behavior, Methodology）级别。请根据以下证据内容，分析其质量并给出1-5级的KBM级别，以及详细的反馈。请按以下格式输出：\n级别: [1-5]\n优点: [优点列表]\n不足: [不足列表]\n建议: [改进建议]\n总体评价: [详细的总体评价]",
		},
		{
			Role:    "user",
			Content: "证据类型：" + evidence.Type + "\n证据内容：" + evidence.Content + "\nKBM名称：" + evidence.KBMName,
		},
	}

	log.Printf("AnalyzeEvidence: 正在调用LLM API...")
	response, err := s.llm.Chat(ctx, messages)
	var parsedFeedback *llm.ParsedFeedback

	if err != nil {
		log.Printf("AnalyzeEvidence: LLM调用失败: %v", err)
		log.Printf("AnalyzeEvidence: 使用默认反馈值")
		parsedFeedback = &llm.ParsedFeedback{
			KBMLevel:    3,
			Strengths:   "对主题有基本了解，能够完成基本任务要求",
			Weaknesses:  "深度分析不足，缺乏批判性思考",
			Suggestions: "建议阅读更多相关文献，尝试从多角度分析问题",
			Feedback:    "证据内容分析：你对该主题有较好的理解，能够清晰表达自己的观点，但在深度分析方面还有提升空间。建议进一步深入研究相关文献，加强批判性思维能力。\n\n注意：LLM API调用失败，此为默认反馈。",
		}
	} else {
		log.Printf("AnalyzeEvidence: LLM响应成功, 响应长度=%d", len(response))
		log.Printf("AnalyzeEvidence: LLM响应内容: %s", response)
		
		parsedFeedback, err = llm.ParseFeedbackResponse(response)
		if err != nil {
			log.Printf("AnalyzeEvidence: 解析LLM响应失败: %v", err)
			parsedFeedback = &llm.ParsedFeedback{
				KBMLevel:    3,
				Strengths:   "提交的证据内容较为完整",
				Weaknesses:  "可进一步提升深度",
				Suggestions: "继续深入研究",
				Feedback:    response,
			}
		}
	}

	log.Printf("AnalyzeFeedback: 解析结果 - KBM级别=%d, 优点=%s, 不足=%s", 
		parsedFeedback.KBMLevel, parsedFeedback.Strengths, parsedFeedback.Weaknesses)

	evidence.KBMLevel = parsedFeedback.KBMLevel
	if err := s.db.WithContext(ctx).Save(evidence).Error; err != nil {
		return nil, err
	}

	feedbackID := utils.GenerateEvidenceID()
	feedbackContent := parsedFeedback.Feedback
	feedbackRecord := &models.Feedback{
		ID:          feedbackID,
		EvidenceID:  evidenceID,
		Content:     feedbackContent,
		KBMLevel:    parsedFeedback.KBMLevel,
		Strengths:   parsedFeedback.Strengths,
		Weaknesses:  parsedFeedback.Weaknesses,
		Suggestions: parsedFeedback.Suggestions,
	}

	if err := s.db.WithContext(ctx).Create(feedbackRecord).Error; err != nil {
		return nil, err
	}

	uploadDir := "./uploads/feedback"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return nil, err
	}

	feedbackFileName := fmt.Sprintf("feedback_%s.txt", evidenceID)
	feedbackFilePath := filepath.Join(uploadDir, feedbackFileName)
	feedbackFileContent := fmt.Sprintf(
		"证据ID: %s\nKBM名称: %s\nKBM级别: %d\n\n优点:\n%s\n\n不足:\n%s\n\n建议:\n%s\n\n总体评价:\n%s",
		evidenceID, evidence.KBMName, parsedFeedback.KBMLevel, parsedFeedback.Strengths, parsedFeedback.Weaknesses, parsedFeedback.Suggestions, feedbackContent,
	)

	if err := os.WriteFile(feedbackFilePath, []byte(feedbackFileContent), 0644); err != nil {
		return nil, err
	}

	feedbackRecord.FileName = feedbackFileName
	feedbackRecord.FilePath = feedbackFilePath
	if err := s.db.WithContext(ctx).Save(feedbackRecord).Error; err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"kbm_level":  parsedFeedback.KBMLevel,
		"feedback":   parsedFeedback.Feedback,
		"strengths":  parsedFeedback.Strengths,
		"weaknesses": parsedFeedback.Weaknesses,
		"suggestions": parsedFeedback.Suggestions,
		"feedback_id": feedbackID,
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

// GetFeedbackByEvidenceID 根据证据ID获取反馈
func (s *EvidenceService) GetFeedbackByEvidenceID(ctx context.Context, evidenceID string) (*models.Feedback, error) {
	var feedback models.Feedback
	if err := s.db.WithContext(ctx).Where("evidence_id = ?", evidenceID).First(&feedback).Error; err != nil {
		return nil, err
	}
	return &feedback, nil
}

// DeleteEvidence 删除证据及其相关的反馈和文件
func (s *EvidenceService) DeleteEvidence(ctx context.Context, evidenceID string) error {
	evidence, err := s.GetEvidenceByID(ctx, evidenceID)
	if err != nil {
		return err
	}

	// 删除反馈记录
	if err := s.db.WithContext(ctx).Where("evidence_id = ?", evidenceID).Delete(&models.Feedback{}).Error; err != nil {
		return err
	}

	// 删除证据记录
	if err := s.db.WithContext(ctx).Delete(&models.Evidence{}, "id = ?", evidenceID).Error; err != nil {
		return err
	}

	// 删除证据文件
	if evidence.FilePath != "" {
		if err := os.Remove(evidence.FilePath); err != nil && !os.IsNotExist(err) {
			return err
		}
	}

	// 删除反馈文件
	feedbackFilePath := fmt.Sprintf("./uploads/feedback/feedback_%s.txt", evidenceID)
	if err := os.Remove(feedbackFilePath); err != nil && !os.IsNotExist(err) {
		return err
	}

	return nil
}

// GetFeedbacksByUserID 根据用户ID获取所有反馈
func (s *EvidenceService) GetFeedbacksByUserID(ctx context.Context, userID string) ([]models.Feedback, error) {
	var feedbacks []models.Feedback
	if err := s.db.WithContext(ctx).
		Joins("JOIN evidences ON feedbacks.evidence_id = evidences.id").
		Joins("JOIN student_tasks ON evidences.student_task_id = student_tasks.id").
		Where("student_tasks.student_id = ?", userID).
		Find(&feedbacks).Error; err != nil {
		return nil, err
	}
	return feedbacks, nil
}
