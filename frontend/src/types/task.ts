export interface Task {
  id: string
  name: string
  description: string
  course_id: string
  teacher_id: string
  start_date: string
  end_date: string
  status: string
  created_at: string
  updated_at: string
  teacher?: {
    id: string
    name: string
    email: string
  }
  student_count?: number
  completed_count?: number
}

export interface StudentTask {
  id: string
  task_id: string
  student_id: string
  status: string
  progress: number
  created_at: string
  updated_at: string
  task?: Task
  student?: {
    id: string
    name: string
    student_id: string
  }
}

export interface CreateTaskRequest {
  name: string
  description: string
  course_id: string
  start_date: string
  end_date: string
}

export interface AssignTaskRequest {
  student_ids: string[]
}

export interface TaskStatistics {
  total_students: number
  completed: number
  processing: number
  pending: number
}