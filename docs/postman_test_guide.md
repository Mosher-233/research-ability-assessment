# Postman API 测试指南

## 文档信息

- **创建日期**: 2026-03-11
- **项目名称**: Research Ability Assessment
- **后端地址**: http://localhost:8080

---

## 目录

1. [前置准备](#前置准备)
2. [环境配置](#环境配置)
3. [认证 API 测试](#认证-api-测试)
4. [任务管理 API 测试](#任务管理-api-测试)
5. [证据管理 API 测试](#证据管理-api-测试)
6. [结果管理 API 测试](#结果管理-api-测试)
7. [常见问题](#常见问题)

---

## 前置准备

### 1. 确保服务运行

在开始测试前，确保以下服务正在运行：

| 服务 | 状态 | 地址 |
|------|------|------|
| MySQL 数据库 | ✅ 运行中 | localhost:3306 |
| Neo4j 数据库 | ✅ 运行中 | http://localhost:7474 |
| 后端服务 | ✅ 运行中 | http://localhost:8080 |
| 前端服务 | ✅ 运行中 | http://localhost:3000 |

### 2. 安装 Postman

下载并安装 [Postman](https://www.postman.com/downloads/)

---

## 环境配置

### 创建 Postman 环境

1. 打开 Postman
2. 点击右上角的 "Environments" → "Create Environment"
3. 环境名称：`Research Ability Assessment - Dev`
4. 添加以下变量：

| 变量名 | 初始值 | 当前值 |
|--------|--------|--------|
| `base_url` | `http://localhost:8080` | `http://localhost:8080` |
| `api_version` | `api/v1` | `api/v1` |
| `token` | (留空) | (登录后自动填充) |
| `teacher_id` | (留空) | (注册后自动填充) |
| `student_id` | (留空) | (注册后自动填充) |
| `task_id` | (留空) | (创建任务后自动填充) |

5. 点击 "Save" 保存环境
6. 在右上角选择该环境

---

## 认证 API 测试

### 1. 注册教师账号

**请求信息：**
- Method: `POST`
- URL: `{{base_url}}/{{api_version}}/auth/register`
- Headers:
  ```
  Content-Type: application/json
  ```
- Body (raw JSON):
  ```json
  {
    "name": "张老师",
    "email": "teacher@example.com",
    "password": "123456",
    "role": "teacher"
  }
  ```

**预期响应 (200 OK):**
```json
{
  "code": 200,
  "message": "注册成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": "uuid-here",
      "name": "张老师",
      "email": "teacher@example.com",
      "role": "teacher",
      "created_at": "2026-03-11T...",
      "updated_at": "2026-03-11T..."
    }
  }
}
```

**操作：**
1. 将响应中的 `token` 复制到环境变量 `token`
2. 将响应中的 `user.id` 复制到环境变量 `teacher_id`

---

### 2. 注册学生账号

**请求信息：**
- Method: `POST`
- URL: `{{base_url}}/{{api_version}}/auth/register`
- Headers:
  ```
  Content-Type: application/json
  ```
- Body (raw JSON):
  ```json
  {
    "name": "李同学",
    "email": "student@example.com",
    "password": "123456",
    "role": "student"
  }
  ```

**预期响应 (200 OK):**
```json
{
  "code": 200,
  "message": "注册成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": "uuid-here",
      "name": "李同学",
      "email": "student@example.com",
      "role": "student",
      "created_at": "2026-03-11T...",
      "updated_at": "2026-03-11T..."
    }
  }
}
```

**操作：**
- 将响应中的 `user.id` 复制到环境变量 `student_id`

---

### 3. 教师登录

**请求信息：**
- Method: `POST`
- URL: `{{base_url}}/{{api_version}}/auth/login`
- Headers:
  ```
  Content-Type: application/json
  ```
- Body (raw JSON):
  ```json
  {
    "email": "teacher@example.com",
    "password": "123456"
  }
  ```

**预期响应 (200 OK):**
```json
{
  "code": 200,
  "message": "登录成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": "uuid-here",
      "name": "张老师",
      "email": "teacher@example.com",
      "role": "teacher",
      "created_at": "2026-03-11T...",
      "updated_at": "2026-03-11T..."
    }
  }
}
```

**操作：**
- 确保环境变量 `token` 已更新为新的 token

---

### 4. 获取教师信息 (需要认证)

**请求信息：**
- Method: `GET`
- URL: `{{base_url}}/{{api_version}}/user/info`
- Headers:
  ```
  Authorization: Bearer {{token}}
  ```

**预期响应 (200 OK):**
```json
{
  "code": 200,
  "message": "获取用户信息成功",
  "data": {
    "id": "uuid-here",
    "name": "张老师",
    "email": "teacher@example.com",
    "role": "teacher",
    "created_at": "2026-03-11T...",
    "updated_at": "2026-03-11T..."
  }
}
```

---

## 任务管理 API 测试

### 5. 创建任务 (需要认证)

**请求信息：**
- Method: `POST`
- URL: `{{base_url}}/{{api_version}}/tasks`
- Headers:
  ```
  Content-Type: application/json
  Authorization: Bearer {{token}}
  ```
- Body (raw JSON):
  ```json
  {
    "title": "数学能力评估",
    "description": "评估学生的数学逻辑思维能力",
    "due_date": "2026-04-01T00:00:00Z"
  }
  ```

**预期响应 (200/201 OK):**
```json
{
  "code": 200,
  "message": "创建成功",
  "data": {
    "id": "task-uuid-here",
    "title": "数学能力评估",
    "description": "评估学生的数学逻辑思维能力",
    "due_date": "2026-04-01T00:00:00Z",
    "teacher_id": "{{teacher_id}}",
    "created_at": "2026-03-11T...",
    "updated_at": "2026-03-11T..."
  }
}
```

**操作：**
- 将响应中的 `data.id` 复制到环境变量 `task_id`

---

### 6. 获取教师的任务列表 (需要认证)

**请求信息：**
- Method: `GET`
- URL: `{{base_url}}/{{api_version}}/tasks`
- Headers:
  ```
  Authorization: Bearer {{token}}
  ```

**预期响应 (200 OK):**
```json
{
  "code": 200,
  "message": "获取成功",
  "data": [
    {
      "id": "{{task_id}}",
      "title": "数学能力评估",
      "description": "评估学生的数学逻辑思维能力",
      "due_date": "2026-04-01T00:00:00Z",
      "teacher_id": "{{teacher_id}}",
      "created_at": "2026-03-11T...",
      "updated_at": "2026-03-11T..."
    }
  ]
}
```

---

### 7. 获取任务详情 (需要认证)

**请求信息：**
- Method: `GET`
- URL: `{{base_url}}/{{api_version}}/tasks/{{task_id}}`
- Headers:
  ```
  Authorization: Bearer {{token}}
  ```

**预期响应 (200 OK):**
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "id": "{{task_id}}",
    "title": "数学能力评估",
    "description": "评估学生的数学逻辑思维能力",
    "due_date": "2026-04-01T00:00:00Z",
    "teacher_id": "{{teacher_id}}",
    "created_at": "2026-03-11T...",
    "updated_at": "2026-03-11T..."
  }
}
```

---

### 8. 分配任务给学生 (需要认证)

**请求信息：**
- Method: `POST`
- URL: `{{base_url}}/{{api_version}}/tasks/{{task_id}}/assign`
- Headers:
  ```
  Content-Type: application/json
  Authorization: Bearer {{token}}
  ```
- Body (raw JSON):
  ```json
  {
    "student_id": "{{student_id}}"
  }
  ```

**预期响应 (200 OK):**
```json
{
  "code": 200,
  "message": "分配成功",
  "data": {
    "id": "student-task-uuid",
    "task_id": "{{task_id}}",
    "student_id": "{{student_id}}",
    "status": "assigned",
    "created_at": "2026-03-11T...",
    "updated_at": "2026-03-11T..."
  }
}
```

---

### 9. 获取任务的学生列表 (需要认证)

**请求信息：**
- Method: `GET`
- URL: `{{base_url}}/{{api_version}}/tasks/{{task_id}}/students`
- Headers:
  ```
  Authorization: Bearer {{token}}
  ```

**预期响应 (200 OK):**
```json
{
  "code": 200,
  "message": "获取成功",
  "data": [
    {
      "id": "student-task-uuid",
      "task_id": "{{task_id}}",
      "student_id": "{{student_id}}",
      "student": {
        "id": "{{student_id}}",
        "name": "李同学",
        "email": "student@example.com"
      },
      "status": "assigned",
      "created_at": "2026-03-11T...",
      "updated_at": "2026-03-11T..."
    }
  ]
}
```

---

### 10. 获取学生列表 (需要认证)

**请求信息：**
- Method: `GET`
- URL: `{{base_url}}/{{api_version}}/tasks/students/list`
- Headers:
  ```
  Authorization: Bearer {{token}}
  ```

**预期响应 (200 OK):**
```json
{
  "code": 200,
  "message": "获取成功",
  "data": [
    {
      "id": "{{student_id}}",
      "name": "李同学",
      "email": "student@example.com",
      "role": "student"
    }
  ]
}
```

---

## 证据管理 API 测试

### 11. 创建证据 (需要认证)

**请求信息：**
- Method: `POST`
- URL: `{{base_url}}/{{api_version}}/evidences`
- Headers:
  ```
  Content-Type: application/json
  Authorization: Bearer {{token}}
  ```
- Body (raw JSON):
  ```json
  {
    "student_task_id": "student-task-uuid",
    "title": "数学作业完成情况",
    "description": "学生完成了所有数学作业",
    "content": "学生在作业中表现出色，解题思路清晰...",
    "evidence_type": "homework",
    "rating": 90
  }
  ```

**预期响应 (200 OK):**
```json
{
  "code": 200,
  "message": "创建成功",
  "data": {
    "id": "evidence-uuid",
    "student_task_id": "student-task-uuid",
    "title": "数学作业完成情况",
    "description": "学生完成了所有数学作业",
    "content": "学生在作业中表现出色，解题思路清晰...",
    "evidence_type": "homework",
    "rating": 90,
    "created_at": "2026-03-11T...",
    "updated_at": "2026-03-11T..."
  }
}
```

---

### 12. 获取证据详情 (需要认证)

**请求信息：**
- Method: `GET`
- URL: `{{base_url}}/{{api_version}}/evidences/{{evidence_id}}`
- Headers:
  ```
  Authorization: Bearer {{token}}
  ```

**预期响应 (200 OK):**
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "id": "evidence-uuid",
    "student_task_id": "student-task-uuid",
    "title": "数学作业完成情况",
    "description": "学生完成了所有数学作业",
    "content": "学生在作业中表现出色，解题思路清晰...",
    "evidence_type": "homework",
    "rating": 90,
    "created_at": "2026-03-11T...",
    "updated_at": "2026-03-11T..."
  }
}
```

---

### 13. 获取学生任务的证据列表 (需要认证)

**请求信息：**
- Method: `GET`
- URL: `{{base_url}}/{{api_version}}/evidences/student-task/{{student_task_id}}`
- Headers:
  ```
  Authorization: Bearer {{token}}
  ```

**预期响应 (200 OK):**
```json
{
  "code": 200,
  "message": "获取成功",
  "data": [
    {
      "id": "evidence-uuid",
      "student_task_id": "student-task-uuid",
      "title": "数学作业完成情况",
      "description": "学生完成了所有数学作业",
      "content": "学生在作业中表现出色，解题思路清晰...",
      "evidence_type": "homework",
      "rating": 90,
      "created_at": "2026-03-11T...",
      "updated_at": "2026-03-11T..."
    }
  ]
}
```

---

## 结果管理 API 测试

### 14. 获取推理结果 (需要认证)

**请求信息：**
- Method: `GET`
- URL: `{{base_url}}/{{api_version}}/results/{{result_id}}`
- Headers:
  ```
  Authorization: Bearer {{token}}
  ```

**预期响应 (200 OK):**
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "id": "result-uuid",
    "student_task_id": "student-task-uuid",
    "result_data": {
      "ability_scores": {
        "logical_thinking": 85,
        "problem_solving": 80,
        "creativity": 75
      },
      "overall_score": 80,
      "recommendations": [
        "继续加强逻辑思维训练",
        "多参与创造性活动"
      ]
    },
    "created_at": "2026-03-11T...",
    "updated_at": "2026-03-11T..."
  }
}
```

---

### 15. 获取任务的所有结果 (需要认证)

**请求信息：**
- Method: `GET`
- URL: `{{base_url}}/{{api_version}}/results/task/{{task_id}}`
- Headers:
  ```
  Authorization: Bearer {{token}}
  ```

**预期响应 (200 OK):**
```json
{
  "code": 200,
  "message": "获取成功",
  "data": [
    {
      "id": "result-uuid",
      "student_task_id": "student-task-uuid",
      "student": {
        "id": "{{student_id}}",
        "name": "李同学"
      },
      "result_data": {...},
      "created_at": "2026-03-11T...",
      "updated_at": "2026-03-11T..."
    }
  ]
}
```

---

### 16. 生成学生报告 (需要认证)

**请求信息：**
- Method: `GET`
- URL: `{{base_url}}/{{api_version}}/results/report/student?student_id={{student_id}}`
- Headers:
  ```
  Authorization: Bearer {{token}}
  ```

**预期响应 (200 OK):**
```json
{
  "code": 200,
  "message": "生成成功",
  "data": {
    "student_id": "{{student_id}}",
    "student_name": "李同学",
    "reports": [
      {
        "task_id": "{{task_id}}",
        "task_title": "数学能力评估",
        "overall_score": 80,
        "ability_dimensions": {...},
        "generated_at": "2026-03-11T..."
      }
    ],
    "summary": {
      "average_score": 80,
      "strengths": ["逻辑思维"],
      "improvements": ["创造力"]
    }
  }
}
```

---

### 17. 生成任务报告 (需要认证)

**请求信息：**
- Method: `GET`
- URL: `{{base_url}}/{{api_version}}/results/report/task/{{task_id}}`
- Headers:
  ```
  Authorization: Bearer {{token}}
  ```

**预期响应 (200 OK):**
```json
{
  "code": 200,
  "message": "生成成功",
  "data": {
    "task_id": "{{task_id}}",
    "task_title": "数学能力评估",
    "student_count": 1,
    "average_score": 80,
    "score_distribution": {
      "90-100": 0,
      "80-89": 1,
      "70-79": 0,
      "60-69": 0,
      "below_60": 0
    },
    "student_results": [
      {
        "student_id": "{{student_id}}",
        "student_name": "李同学",
        "score": 80,
        "rank": 1
      }
    ],
    "generated_at": "2026-03-11T..."
  }
}
```

---

## Postman 集合导入

为了更方便测试，您可以创建一个 Postman Collection。以下是集合的基本结构：

### Collection 结构

```
Research Ability Assessment
├── 认证
│   ├── 注册教师
│   ├── 注册学生
│   ├── 教师登录
│   └── 获取用户信息
├── 任务管理
│   ├── 创建任务
│   ├── 获取任务列表
│   ├── 获取任务详情
│   ├── 分配任务
│   ├── 获取任务学生
│   └── 获取学生列表
├── 证据管理
│   ├── 创建证据
│   ├── 获取证据详情
│   └── 获取证据列表
└── 结果管理
    ├── 获取结果
    ├── 获取任务结果
    ├── 生成学生报告
    └── 生成任务报告
```

---

## 常见问题

### Q1: 收到 "未授权" 错误？

**A:** 确保：
1. 已正确设置 `Authorization` header，格式为 `Bearer {{token}}`
2. Token 没有过期
3. 环境变量 `token` 已正确设置

### Q2: 收到 "请求参数错误"？

**A:** 检查：
1. 请求体格式是否正确（JSON）
2. 必填字段是否都已提供
3. 字段类型是否正确

### Q3: 如何重置测试数据？

**A:** 可以：
1. 重启 MySQL 容器：`docker-compose restart mysql`
2. 或者完全重新创建：`docker-compose down -v && docker-compose up -d`

### Q4: 前端和后端如何配合测试？

**A:** 
1. 先使用 Postman 测试后端 API
2. 确保 API 正常工作后，在浏览器中打开前端 `http://localhost:3000`
3. 使用相同的账号在前端登录

---

## 附录

### 测试账号参考

| 角色 | 邮箱 | 密码 | 姓名 |
|------|------|------|------|
| 教师 | teacher@example.com | 123456 | 张老师 |
| 学生 | student@example.com | 123456 | 李同学 |

### 快速开始脚本

您也可以使用 curl 命令快速测试：

```bash
# 注册教师
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"张老师","email":"teacher@example.com","password":"123456","role":"teacher"}'

# 教师登录
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"teacher@example.com","password":"123456"}'
```

---

**文档结束**
