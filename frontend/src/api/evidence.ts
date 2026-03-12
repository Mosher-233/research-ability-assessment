import api from './auth'
import type { Evidence, CreateEvidenceRequest } from '../types/evidence'

export const evidenceApi = {
  // 创建证据
  createEvidence: async (data: CreateEvidenceRequest): Promise<{ code: number; message: string; data: Evidence }> => {
    const response = await api.post('/evidences', data)
    return response.data
  },
  
  // 分析证据（AI）
  analyzeEvidence: async (evidenceId: string): Promise<{ code: number; message: string; data: { kbm_level: number; feedback: string } }> => {
    const response = await api.post(`/evidences/${evidenceId}/analyze`)
    return response.data
  },
  
  // 获取证据列表
  getEvidences: async (): Promise<{ code: number; message: string; data: Evidence[] }> => {
    const response = await api.get('/evidences')
    return response.data
  },
  
  // 获取证据详情
  getEvidenceById: async (evidenceId: string): Promise<{ code: number; message: string; data: Evidence }> => {
    const response = await api.get(`/evidences/${evidenceId}`)
    return response.data
  },
  
  // 删除证据
  deleteEvidence: async (evidenceId: string): Promise<{ code: number; message: string; data: any }> => {
    const response = await api.delete(`/evidences/${evidenceId}`)
    return response.data
  },
  
  // 获取学生任务的证据列表
  getEvidencesByStudentTaskId: async (studentTaskId: string): Promise<{ code: number; message: string; data: Evidence[] }> => {
    const response = await api.get(`/evidences/student-task/${studentTaskId}`)
    return response.data
  },
  
  // 根据学生ID和任务ID获取证据列表
  getEvidencesByStudentAndTask: async (studentId: string, taskId: string): Promise<{ code: number; message: string; data: Evidence[] }> => {
    const response = await api.get(`/evidences/student-task?student_id=${studentId}&task_id=${taskId}`)
    return response.data
  }
}