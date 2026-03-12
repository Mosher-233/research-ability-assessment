import api from './auth'
import type { Task, StudentTask, CreateTaskRequest, AssignTaskRequest } from '../types/task'
import type { Student } from '../types/user'

export const taskApi = {
  // 创建任务
  createTask: async (data: CreateTaskRequest): Promise<{ code: number; message: string; data: Task }> => {
    const response = await api.post('/tasks', data)
    return response.data
  },
  
  // 获取任务列表
  getTasks: async (): Promise<{ code: number; message: string; data: Task[] }> => {
    const response = await api.get('/tasks')
    return response.data
  },
  
  // 获取学生分配的任务
  getAssignedTasks: async (): Promise<{ code: number; message: string; data: Task[] }> => {
    const response = await api.get('/tasks/students/assigned')
    return response.data
  },
  
  // 获取学生任务列表
  getStudentTasks: async (taskID: string): Promise<{ code: number; message: string; data: StudentTask[] }> => {
    const response = await api.get(`/tasks/${taskID}/students`)
    return response.data
  },
  
  // 获取任务详情
  getTaskByID: async (taskID: string): Promise<{ code: number; message: string; data: Task }> => {
    const response = await api.get(`/tasks/${taskID}`)
    return response.data
  },
  
  // 分配任务给学生
  assignTask: async (taskID: string, data: AssignTaskRequest): Promise<{ code: number; message: string; data: any }> => {
    const response = await api.post(`/tasks/${taskID}/assign`, data)
    return response.data
  },
  
  // 获取任务的学生列表
  getTaskStudents: async (taskID: string): Promise<{ code: number; message: string; data: StudentTask[] }> => {
    const response = await api.get(`/tasks/${taskID}/students`)
    return response.data
  },
  
  // 获取学生列表
  getStudents: async (): Promise<{ code: number; message: string; data: Student[] }> => {
    const response = await api.get('/tasks/students/list')
    return response.data
  },
  
  // 更新任务状态
  updateTaskStatus: async (taskID: string, status: string): Promise<{ code: number; message: string; data: any }> => {
    const response = await api.put(`/tasks/${taskID}/status`, { status })
    return response.data
  }
}