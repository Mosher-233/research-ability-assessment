export interface Evidence {
  id: string
  student_task_id: string
  type: string
  content: string
  kbm_name: string
  kbm_level: number
  created_at: string
  updated_at: string
  student_task?: {
    id: string
    task_id: string
    student_id: string
  }
}

export interface CreateEvidenceRequest {
  student_task_id: string
  type: string
  content: string
  kbm_name: string
  kbm_level: number
}