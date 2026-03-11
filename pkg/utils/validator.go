package utils

import (
	"regexp"
)

// ValidateEmail 验证邮箱格式
func ValidateEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	regex := regexp.MustCompile(pattern)
	return regex.MatchString(email)
}

// ValidatePassword 验证密码强度
func ValidatePassword(password string) bool {
	// 密码长度至少6位
	return len(password) >= 6
}

// ValidateStudentID 验证学生ID格式
func ValidateStudentID(studentID string) bool {
	// 学生ID通常是数字，长度在8-12位之间
	pattern := `^\d{8,12}$`
	regex := regexp.MustCompile(pattern)
	return regex.MatchString(studentID)
}

// ValidateRole 验证角色
func ValidateRole(role string) bool {
	return role == "teacher" || role == "student"
}