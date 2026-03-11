export interface User {
  id: string
  name: string
  email: string
  role: string
  created_at: string
  updated_at: string
}

export interface Teacher extends User {
  department: string
  title: string
}

export interface Student extends User {
  student_id: string
  major: string
  grade: string
}

export interface LoginRequest {
  email: string
  password: string
}

export interface RegisterRequest {
  name: string
  email: string
  password: string
  role: string
}

export interface AuthResponse {
  code: number
  message: string
  data: {
    token: string
    user: User
  }
}