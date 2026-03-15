import api from './auth'
import type { InferenceResult } from '../types/result'

export interface Report {
  id: string
  student_task_id: string
  student_id: string
  task_id: string
  overall_score: number
  overall_level: string
  dimension_scores: any
  class_comparison: any
  rank: number
  percentile: number
  strengths: any
  weaknesses: any
  detailed_analysis: any
  suggestions: any
  radar_chart_data: any
  report_path: string
  created_at: string
  updated_at: string
  student_name?: string
  task_name?: string
}

export const resultApi = {
  // 获取所有推理结果
  getResults: async (): Promise<{ code: number; message: string; data: InferenceResult[] }> => {
    const response = await api.get('/results')
    return response.data
  },
  
  // 获取学生的推理结果
  getStudentResults: async (): Promise<{ code: number; message: string; data: InferenceResult[] }> => {
    const response = await api.get('/results/student')
    return response.data
  },
  
  // 获取推理结果详情
  getInferenceResultByID: async (resultID: string): Promise<{ code: number; message: string; data: InferenceResult }> => {
    const response = await api.get(`/results/${resultID}`)
    return response.data
  },
  
  // 获取任务的推理结果列表
  getInferenceResultsByTaskID: async (taskID: string): Promise<{ code: number; message: string; data: InferenceResult[] }> => {
    const response = await api.get(`/results/task/${taskID}`)
    return response.data
  },
  
  // 根据学生ID和任务ID获取推理结果
  getInferenceResultByStudentAndTask: async (studentID: string, taskID: string): Promise<{ code: number; message: string; data: InferenceResult }> => {
    const response = await api.get(`/results/student-task?student_id=${studentID}&task_id=${taskID}`)
    return response.data
  },
  
  // 生成学生推理结果
  generateStudentInference: async (taskID: string): Promise<{ code: number; message: string; data: InferenceResult }> => {
    const response = await api.post(`/results/generate/student?task_id=${taskID}`)
    return response.data
  },
  
  // 生成学生报告
  generateStudentReport: async (studentID: string, taskID: string): Promise<{ code: number; message: string; data: Report }> => {
    const response = await api.get(`/results/report/student?student_id=${studentID}&task_id=${taskID}`)
    return response.data
  },
  
  // 生成任务报告
  generateTaskReport: async (taskID: string): Promise<{ code: number; message: string; data: InferenceResult[] }> => {
    const response = await api.get(`/results/report/task/${taskID}`)
    return response.data
  },

  // 获取所有报告
  getReports: async (): Promise<{ code: number; message: string; data: Report[] }> => {
    const response = await api.get('/reports')
    return response.data
  },

  // 获取学生的报告
  getStudentReports: async (): Promise<{ code: number; message: string; data: Report[] }> => {
    const response = await api.get('/reports/student')
    return response.data
  }
}