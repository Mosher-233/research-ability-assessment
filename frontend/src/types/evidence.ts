export interface Evidence {
  id: string
  student_task_id: string
  type: string
  content: string
  file_name?: string
  file_path?: string
  file_type?: string
  file_size?: number
  kbm_name: string
  kbm_level: number
  created_at: string
  updated_at: string
  student_id?: string
  student_name?: string
  task_id?: string
  task_name?: string
}

export interface CreateEvidenceRequest {
  student_task_id: string
  type: string
  content?: string
  kbm_name: string
  kbm_level?: number
  file_name?: string
  file_path?: string
  file_type?: string
  file_size?: number
}

export interface Feedback {
  id: string
  evidence_id: string
  content: string
  kbm_level: number
  strengths?: string
  weaknesses?: string
  suggestions?: string
  file_name?: string
  file_path?: string
  created_at: string
  updated_at: string
}