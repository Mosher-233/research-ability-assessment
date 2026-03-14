package utils

import (
	"fmt"
	"sync"
	"time"
)

var (
	userCounter    int64 = 1000
	taskCounter    int64 = 1000
	evidenceCounter int64 = 1000
	teacherIDPrefix  string = "T"
	studentIDPrefix  string = "S"
	taskIDPrefix     string = "TK"
	evidenceIDPrefix   string = "EV"
	userMu           sync.Mutex
	taskMu           sync.Mutex
	evidenceMu       sync.Mutex
)

func GenerateUserID(role string) string {
	userMu.Lock()
	defer userMu.Unlock()
	
	userCounter++
	timestamp := time.Now().Format("20060102")
	
	if role == "teacher" {
		return fmt.Sprintf("%s%s%05d", teacherIDPrefix, timestamp, userCounter)
	}
	return fmt.Sprintf("%s%s%05d", studentIDPrefix, timestamp, userCounter)
}

func GenerateTaskID() string {
	taskMu.Lock()
	defer taskMu.Unlock()
	
	taskCounter++
	timestamp := time.Now().Format("20060102")
	
	return fmt.Sprintf("%s%s%05d", taskIDPrefix, timestamp, taskCounter)
}

func GenerateStudentTaskID() string {
	return GenerateTaskID()
}

func GenerateEvidenceID() string {
	evidenceMu.Lock()
	defer evidenceMu.Unlock()
	
	evidenceCounter++
	timestamp := time.Now().Format("20060102")
	
	return fmt.Sprintf("%s%s%05d", evidenceIDPrefix, timestamp, evidenceCounter)
}
