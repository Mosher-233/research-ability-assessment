# 并发测试指南

本指南介绍如何对研究能力评估系统进行并发测试。

## 目录
1. [使用内置Go测试脚本](#使用内置go测试脚本)
2. [使用Apache Bench](#使用apache-bench)
3. [使用k6](#使用k6)
4. [使用wrk](#使用wrk)
5. [Go标准库测试](#go标准库测试)

## 使用内置Go测试脚本

项目已包含一个专门的并发测试脚本：

```bash
cd scripts
go run concurrent_test.go
```

该脚本会测试：
- 登录接口并发
- 获取结果列表接口并发
- 获取报告列表接口并发

测试结果会显示：
- 成功/失败请求数
- 响应时间统计（最小、最大、平均）
- QPS（每秒查询数）
- 成功率

## 使用Apache Bench

Apache Bench 是最简单的并发测试工具之一：

### 安装
- Windows: 从 Apache 官网下载或使用 Chocolatey
- macOS: `brew install httpd`
- Linux: `apt-get install apache2-utils`

### 测试登录接口
```bash
# 创建登录请求JSON文件
echo '{"email":"teacher1@example.com","password":"password123"}' > login.json

# 并发测试：100个请求，10个并发
ab -n 100 -c 10 -p login.json -T application/json http://localhost:8080/api/v1/auth/login
```

### 测试需要认证的接口
```bash
# 首先获取token
TOKEN="your_token_here"

# 使用header进行认证
ab -n 100 -c 10 -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/v1/results
```

## 使用k6

k6 是一个现代化的负载测试工具，支持JavaScript脚本：

### 安装
- Windows: `choco install k6`
- macOS: `brew install k6`
- Linux: 从 https://k6.io/docs/getting-started/installation/ 下载

### 创建测试脚本 (k6_test.js)
```javascript
import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
  vus: 50,           // 虚拟用户数
  duration: '30s',   // 持续时间
};

const BASE_URL = 'http://localhost:8080';

export default function () {
  // 登录
  let loginRes = http.post(`${BASE_URL}/api/v1/auth/login`, JSON.stringify({
    email: 'teacher1@example.com',
    password: 'password123'
  }), { headers: { 'Content-Type': 'application/json' } });
  
  check(loginRes, { 'login status 200': (r) => r.status === 200 });
  
  let token = loginRes.json('data.token');
  
  // 使用token访问受保护的接口
  let resultsRes = http.get(`${BASE_URL}/api/v1/results`, {
    headers: { 'Authorization': `Bearer ${token}` }
  });
  
  check(resultsRes, { 'get results status 200': (r) => r.status === 200 });
  
  sleep(1);
}
```

### 运行测试
```bash
k6 run k6_test.js
```

## 使用wrk

wrk 是一个高性能的HTTP压测工具：

### 安装
```bash
# macOS
brew install wrk

# Linux
git clone https://github.com/wg/wrk.git
cd wrk
make
sudo cp wrk /usr/local/bin
```

### 运行测试
```bash
# 基本测试：12线程，400连接，持续30秒
wrk -t12 -c400 -d30s http://localhost:8080/health

# 使用Lua脚本进行复杂测试
wrk -t12 -c400 -d30s -s script.lua http://localhost:8080/api/v1/results
```

## Go标准库测试

也可以使用Go的标准库编写并发测试：

### 创建测试文件 (main_test.go)
```go
package main

import (
	"net/http"
	"sync"
	"testing"
)

func TestConcurrentRequests(t *testing.T) {
	const numRequests = 100
	const concurrency = 10

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, concurrency)

	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		semaphore <- struct{}{}

		go func(requestID int) {
			defer wg.Done()
			defer func() { <-semaphore }()

			resp, err := http.Get("http://localhost:8080/health")
			if err != nil {
				t.Errorf("Request %d failed: %v", requestID, err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				t.Errorf("Request %d: expected status 200, got %d", requestID, resp.StatusCode)
			}
		}(i)
	}

	wg.Wait()
}
```

### 运行测试
```bash
go test -v -run TestConcurrentRequests
```

## 性能指标说明

在进行并发测试时，关注以下指标：

1. **QPS (Queries Per Second)** - 每秒处理的请求数
2. **响应时间** - 最小、平均、最大、P95、P99响应时间
3. **错误率** - 失败请求的百分比
4. **吞吐量** - 系统在单位时间内处理的总数据量
5. **并发用户数** - 系统可以同时处理的用户数量

## 测试建议

1. **从低并发开始** - 先用小并发数测试，逐步增加
2. **监控系统资源** - 测试时监控CPU、内存、数据库连接数
3. **测试关键路径** - 重点测试登录、创建任务、生成报告等核心功能
4. **多次测试取平均** - 运行多次测试，取平均值作为结果
5. **记录环境配置** - 记录测试时的硬件配置、数据库配置等

## 常见问题

### 端口被占用
```bash
# Windows
netstat -ano | findstr :8080
taskkill /F /PID <进程ID>

# Linux/macOS
lsof -ti :8080 | xargs kill -9
```

### 数据库连接数不足
检查并调整MySQL的最大连接数：
```sql
SHOW VARIABLES LIKE 'max_connections';
SET GLOBAL max_connections = 200;
```
