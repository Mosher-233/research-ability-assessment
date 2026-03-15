// 格式化日期
export const formatDate = (dateStr: string): string => {
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-CN')
}

// 格式化时间
export const formatDateTime = (dateStr: string): string => {
  const date = new Date(dateStr)
  return date.toLocaleString('zh-CN')
}

// 格式化得分
export const formatScore = (score: number): string => {
  const normalizedScore = score <= 1 ? score * 100 : score
  return normalizedScore.toFixed(0) + '分'
}

// 获取状态类型
export const getStatusType = (status: string): string => {
  const typeMap: Record<string, string> = {
    active: 'success',
    completed: 'info',
    archived: 'warning',
    pending: 'info',
    processing: 'primary'
  }
  return typeMap[status] || 'info'
}

// 获取状态文本
export const getStatusText = (status: string): string => {
  const textMap: Record<string, string> = {
    active: '进行中',
    completed: '已完成',
    archived: '已归档',
    pending: '待处理',
    processing: '处理中'
  }
  return textMap[status] || status
}

// 获取进度百分比
export const getProgressPercentage = (task: { student_count?: number; completed_count?: number }): number => {
  if (!task.student_count || task.student_count === 0) return 0
  return Math.round((task.completed_count || 0) / task.student_count * 100)
}

// 获取进度颜色
export const getProgressColor = (percentage: number): string => {
  if (percentage === 100) return '#67c23a'
  if (percentage >= 50) return '#409eff'
  return '#e6a23c'
}