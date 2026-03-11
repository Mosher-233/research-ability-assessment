import api from './auth'
import type { Evidence, CreateEvidenceRequest } from '../types/evidence'

export const evidenceApi = {
  // 创建证据
  createEvidence: async (data: CreateEvidenceRequest): Promise<{ code: number; message: string; data: Evidence }> => {
    const response = await api.post('/evidences', data)
    return response.data
  },
  
  // 获取证据详情
  getEvidenceByID: async (evidenceID: string): Promise<{ code: number; message: string; data: Evidence }> => {
    const response = await api.get(`/evidences/${evidenceID}`)
    return response.data
  },
  
  // 获取学生任务的证据列表
  getEvidencesByStudentTaskID: async (studentTaskID: string): Promise<{ code: number; message: string; data: Evidence[] }> => {
    const response = await api.get(`/evidences/student-task/${studentTaskID}`)
    return response.data
  },
  
  // 根据学生ID和任务ID获取证据列表
  getEvidencesByStudentAndTask: async (studentID: string, taskID: string): Promise<{ code: number; message: string; data: Evidence[] }> => {
    const response = await api.get(`/evidences/student-task?student_id=${studentID}&task_id=${taskID}`)
    return response.data
  }
}