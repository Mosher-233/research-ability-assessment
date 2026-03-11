// 验证邮箱格式
export const validateEmail = (email: string): boolean => {
  const pattern = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/
  return pattern.test(email)
}

// 验证密码强度
export const validatePassword = (password: string): boolean => {
  // 密码长度至少6位
  return password.length >= 6
}

// 验证学生ID格式
export const validateStudentID = (studentID: string): boolean => {
  // 学生ID通常是数字，长度在8-12位之间
  const pattern = /^\d{8,12}$/
  return pattern.test(studentID)
}

// 验证角色
export const validateRole = (role: string): boolean => {
  return role === 'teacher' || role === 'student'
}

// 验证日期范围
export const validateDateRange = (startDate: string, endDate: string): boolean => {
  const start = new Date(startDate)
  const end = new Date(endDate)
  return start <= end
}