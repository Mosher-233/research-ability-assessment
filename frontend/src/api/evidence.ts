import api from './auth'
import type { Evidence, CreateEvidenceRequest, Feedback } from '../types/evidence'

export const evidenceApi = {
  createEvidence: async (data: CreateEvidenceRequest): Promise<{ code: number; message: string; data: Evidence }> => {
    const response = await api.post('/evidences', data)
    return response.data
  },

  uploadEvidenceFile: async (formData: FormData): Promise<{ code: number; message: string; data: Evidence }> => {
    const response = await api.post('/evidences/upload', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
    return response.data
  },

  analyzeEvidence: async (evidenceId: string): Promise<{ 
    code: number; 
    message: string; 
    data: { 
      kbm_level: number; 
      feedback: string;
      strengths?: string;
      weaknesses?: string;
      suggestions?: string;
      feedback_id?: string;
    } 
  }> => {
    const response = await api.post(`/evidences/${evidenceId}/analyze`)
    return response.data
  },

  getEvidences: async (): Promise<{ code: number; message: string; data: Evidence[] }> => {
    const response = await api.get('/evidences')
    return response.data
  },

  getEvidenceById: async (evidenceId: string): Promise<{ code: number; message: string; data: Evidence }> => {
    const response = await api.get(`/evidences/${evidenceId}`)
    return response.data
  },

  downloadEvidenceFile: async (evidenceId: string): Promise<void> => {
    const response = await api.get(`/evidences/${evidenceId}/download`, {
      responseType: 'blob'
    })
    const url = window.URL.createObjectURL(new Blob([response.data]))
    const link = document.createElement('a')
    link.href = url
    link.setAttribute('download', `evidence_${evidenceId}`)
    document.body.appendChild(link)
    link.click()
    link.remove()
  },

  deleteEvidence: async (evidenceId: string): Promise<{ code: number; message: string; data: any }> => {
    const response = await api.delete(`/evidences/${evidenceId}`)
    return response.data
  },

  getEvidencesByStudentTaskId: async (studentTaskId: string): Promise<{ code: number; message: string; data: Evidence[] }> => {
    const response = await api.get(`/evidences/student-task/${studentTaskId}`)
    return response.data
  },

  getEvidencesByStudentAndTask: async (studentId: string, taskId: string): Promise<{ code: number; message: string; data: Evidence[] }> => {
    const response = await api.get(`/evidences/student-task?student_id=${studentId}&task_id=${taskId}`)
    return response.data
  },

  getFeedbackByEvidenceId: async (evidenceId: string): Promise<{ code: number; message: string; data: Feedback }> => {
    const response = await api.get(`/evidences/${evidenceId}/feedback`)
    return response.data
  },

  getFeedbacks: async (): Promise<{ code: number; message: string; data: Feedback[] }> => {
    const response = await api.get('/evidences/feedbacks/list')
    return response.data
  }
}