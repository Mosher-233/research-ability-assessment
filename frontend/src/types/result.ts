export interface DimensionScore {
  name: string
  score: number
  level: string
  details: string
  evidence_ids: string[]
}

export interface InferenceResult {
  id: string
  student_id: string
  task_id: string
  overall_score: number
  overall_level: string
  dimension_scores: Record<string, DimensionScore>
  reasoning: string
  created_at: string
  updated_at: string
  student?: {
    id: string
    name: string
    student_id: string
  }
  task?: {
    id: string
    name: string
  }
}