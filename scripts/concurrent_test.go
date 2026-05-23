package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Token string `json:"token"`
	} `json:"data"`
}

type TestConfig struct {
	BaseURL         string
	TeacherEmail    string
	TeacherPass     string
	StudentEmail    string
	StudentPass     string
	ConcurrentUsers int
	RequestsPerUser int
}

type TestResult struct {
	SuccessCount  int
	FailureCount  int
	TotalTime     time.Duration
	MinTime       time.Duration
	MaxTime       time.Duration
	ResponseTimes []time.Duration
}

func main() {
	config := TestConfig{
		BaseURL:         "http://localhost:8080",
		TeacherEmail:    "teacher1@example.com",
		TeacherPass:     "password123",
		StudentEmail:    "student1@example.com",
		StudentPass:     "password123",
		ConcurrentUsers: 10,
		RequestsPerUser: 20,
	}

	fmt.Println("====================================")
	fmt.Println("研究能力评估系统并发测试")
	fmt.Println("====================================")
	fmt.Printf("并发用户数: %d\n", config.ConcurrentUsers)
	fmt.Printf("每用户请求数: %d\n", config.RequestsPerUser)
	fmt.Printf("总请求数: %d\n", config.ConcurrentUsers*config.RequestsPerUser)
	fmt.Println("====================================")

	fmt.Println("\n1. 测试登录接口并发...")
	loginResult := testLoginConcurrency(config)
	printTestResult("登录接口", loginResult)

	fmt.Println("\n2. 测试获取结果列表接口并发...")
	resultsResult := testResultsConcurrency(config)
	printTestResult("结果列表接口", resultsResult)

	fmt.Println("\n3. 测试获取报告列表接口并发...")
	reportsResult := testReportsConcurrency(config)
	printTestResult("报告列表接口", reportsResult)

	fmt.Println("\n====================================")
	fmt.Println("并发测试完成!")
	fmt.Println("====================================")
}

func login(config TestConfig, email, password string) (string, error) {
	reqBody := LoginRequest{
		Email:    email,
		Password: password,
	}
	jsonData, _ := json.Marshal(reqBody)

	resp, err := http.Post(
		config.BaseURL+"/api/v1/auth/login",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var loginResp LoginResponse
	json.Unmarshal(body, &loginResp)

	if loginResp.Code != 200 {
		return "", fmt.Errorf("登录失败: %s", loginResp.Message)
	}

	return loginResp.Data.Token, nil
}

func testLoginConcurrency(config TestConfig) TestResult {
	var wg sync.WaitGroup
	result := TestResult{
		MinTime:       time.Hour,
		ResponseTimes: make([]time.Duration, 0),
	}

	startTime := time.Now()

	for i := 0; i < config.ConcurrentUsers; i++ {
		wg.Add(1)
		go func(userID int) {
			defer wg.Done()
			for j := 0; j < config.RequestsPerUser; j++ {
				reqStart := time.Now()
				_, err := login(config, config.TeacherEmail, config.TeacherPass)
				reqDuration := time.Since(reqStart)

				result.ResponseTimes = append(result.ResponseTimes, reqDuration)

				if reqDuration < result.MinTime {
					result.MinTime = reqDuration
				}
				if reqDuration > result.MaxTime {
					result.MaxTime = reqDuration
				}

				if err != nil {
					result.FailureCount++
				} else {
					result.SuccessCount++
				}
			}
		}(i)
	}

	wg.Wait()
	result.TotalTime = time.Since(startTime)

	return result
}

func testResultsConcurrency(config TestConfig) TestResult {
	token, _ := login(config, config.TeacherEmail, config.TeacherPass)

	var wg sync.WaitGroup
	result := TestResult{
		MinTime:       time.Hour,
		ResponseTimes: make([]time.Duration, 0),
	}

	startTime := time.Now()

	for i := 0; i < config.ConcurrentUsers; i++ {
		wg.Add(1)
		go func(userID int) {
			defer wg.Done()
			client := &http.Client{}

			for j := 0; j < config.RequestsPerUser; j++ {
				req, _ := http.NewRequest("GET", config.BaseURL+"/api/v1/results", nil)
				req.Header.Set("Authorization", "Bearer "+token)

				reqStart := time.Now()
				resp, err := client.Do(req)
				reqDuration := time.Since(reqStart)

				result.ResponseTimes = append(result.ResponseTimes, reqDuration)

				if reqDuration < result.MinTime {
					result.MinTime = reqDuration
				}
				if reqDuration > result.MaxTime {
					result.MaxTime = reqDuration
				}

				if err != nil || resp.StatusCode != 200 {
					result.FailureCount++
				} else {
					result.SuccessCount++
				}
				if resp != nil {
					resp.Body.Close()
				}
			}
		}(i)
	}

	wg.Wait()
	result.TotalTime = time.Since(startTime)

	return result
}

func testReportsConcurrency(config TestConfig) TestResult {
	token, _ := login(config, config.TeacherEmail, config.TeacherPass)

	var wg sync.WaitGroup
	result := TestResult{
		MinTime:       time.Hour,
		ResponseTimes: make([]time.Duration, 0),
	}

	startTime := time.Now()

	for i := 0; i < config.ConcurrentUsers; i++ {
		wg.Add(1)
		go func(userID int) {
			defer wg.Done()
			client := &http.Client{}

			for j := 0; j < config.RequestsPerUser; j++ {
				req, _ := http.NewRequest("GET", config.BaseURL+"/api/v1/reports", nil)
				req.Header.Set("Authorization", "Bearer "+token)

				reqStart := time.Now()
				resp, err := client.Do(req)
				reqDuration := time.Since(reqStart)

				result.ResponseTimes = append(result.ResponseTimes, reqDuration)

				if reqDuration < result.MinTime {
					result.MinTime = reqDuration
				}
				if reqDuration > result.MaxTime {
					result.MaxTime = reqDuration
				}

				if err != nil || resp.StatusCode != 200 {
					result.FailureCount++
				} else {
					result.SuccessCount++
				}
				if resp != nil {
					resp.Body.Close()
				}
			}
		}(i)
	}

	wg.Wait()
	result.TotalTime = time.Since(startTime)

	return result
}

func printTestResult(name string, result TestResult) {
	fmt.Printf("\n--- %s 测试结果 ---\n", name)
	fmt.Printf("成功请求数: %d\n", result.SuccessCount)
	fmt.Printf("失败请求数: %d\n", result.FailureCount)
	fmt.Printf("总耗时: %v\n", result.TotalTime)
	fmt.Printf("最短响应时间: %v\n", result.MinTime)
	fmt.Printf("最长响应时间: %v\n", result.MaxTime)

	if len(result.ResponseTimes) > 0 {
		var sum time.Duration
		for _, t := range result.ResponseTimes {
			sum += t
		}
		avgTime := sum / time.Duration(len(result.ResponseTimes))
		fmt.Printf("平均响应时间: %v\n", avgTime)
		fmt.Printf("QPS: %.2f\n", float64(result.SuccessCount)/result.TotalTime.Seconds())
	}

	total := result.SuccessCount + result.FailureCount
	if total > 0 {
		successRate := float64(result.SuccessCount) / float64(total) * 100
		fmt.Printf("成功率: %.2f%%\n", successRate)
	}
}
