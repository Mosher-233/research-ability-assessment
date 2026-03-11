package agent

import (
	"context"
	"research-ability-assessment/internal/models"
	"research-ability-assessment/internal/service"
)

type EvidenceAgent struct {
	evidenceService *service.EvidenceService
	ioUnit         *IOUnit
}

func NewEvidenceAgent(evidenceService *service.EvidenceService, ioUnit *IOUnit) *EvidenceAgent {
	return &EvidenceAgent{
		evidenceService: evidenceService,
		ioUnit:         ioUnit,
	}
}

func (a *EvidenceAgent) CollectEvidence(ctx context.Context, studentID string, taskID string) ([]*models.Evidence, error) {
	// 从IO单元获取证据
	evidences, err := a.ioUnit.GetEvidences(ctx, studentID, taskID)
	if err != nil {
		return nil, err
	}
	
	// 这里可以添加证据处理的逻辑
	// 例如对证据进行分类、筛选或增强
	
	return evidences, nil
}

func (a *EvidenceAgent) ProcessEvidence(ctx context.Context, evidence *models.Evidence) error {
	// 这里可以添加证据处理的逻辑
	// 例如提取关键信息、进行标准化处理等
	
	return nil
}