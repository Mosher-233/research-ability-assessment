import api from './auth'
import type { InferenceResult } from '../types/result'

export const resultApi = {
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
  
  // 生成学生报告
  generateStudentReport: async (studentID: string, taskID: string): Promise<{ code: number; message: string; data: InferenceResult }> => {
    const response = await api.get(`/results/report/student?student_id=${studentID}&task_id=${taskID}`)
    return response.data
  },
  
  // 生成任务报告
  generateTaskReport: async (taskID: string): Promise<{ code: number; message: string; data: InferenceResult[] }> => {
    const response = await api.get(`/results/report/task/${taskID}`)
    return response.data
  }
}