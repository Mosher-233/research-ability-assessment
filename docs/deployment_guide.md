# 研究能力评价系统部署指南

## 1. 环境要求

### 1.1 后端环境

- **Go**：1.20 或更高版本
- **PostgreSQL**：14.0 或更高版本
- **Neo4j**：5.0 或更高版本
- **Redis**：6.0 或更高版本（可选，用于缓存）

### 1.2 前端环境

- **Node.js**：16.0 或更高版本
- **npm**：8.0 或更高版本

## 2. 开发环境部署

### 2.1 后端部署

#### 2.1.1 安装依赖

1. 克隆项目代码：
   ```bash
   git clone <项目仓库地址>
   cd research-ability-assessment
   ```

2. 安装 Go 依赖：
   ```bash
   go mod tidy
   ```

#### 2.1.2 配置数据库

1. 启动 PostgreSQL 服务
2. 创建数据库：
   ```sql
   CREATE DATABASE research_assessment;
   ```
3. 创建用户并授权：
   ```sql
   CREATE USER postgres WITH PASSWORD 'postgres';
   GRANT ALL PRIVILEGES ON DATABASE research_assessment TO postgres;
   ```

#### 2.1.3 配置 Neo4j

1. 启动 Neo4j 服务
2. 访问 Neo4j 浏览器（默认地址：http://localhost:7474）
3. 登录并修改默认密码（默认用户名：neo4j，默认密码：neo4j）

#### 2.1.4 配置系统

1. 复制配置文件：
   ```bash
   cp configs/config.dev.yaml configs/config.yaml
   ```

2. 编辑配置文件 `configs/config.yaml`，设置数据库连接信息：
   ```yaml
   database:
     host: localhost
     port: 5432
     user: postgres
     password: postgres
     dbname: research_assessment
     sslmode: disable

   neo4j:
     uri: bolt://localhost:7687
     username: neo4j
     password: <your-neo4j-password>
   ```

3. 配置 LLM API 密钥（可选）：
   ```yaml
   llm:
     provider: openai
     api_key: <your-openai-api-key>
     base_url: https://api.openai.com/v1
     model: gpt-3.5-turbo
     max_tokens: 1000
     temperature: 0.7
   ```

#### 2.1.5 启动后端服务

```bash
go run cmd/server/main.go
```

后端服务将在 `http://localhost:8080` 启动。

### 2.2 前端部署

#### 2.2.1 安装依赖

1. 进入前端目录：
   ```bash
   cd frontend
   ```

2. 安装 npm 依赖：
   ```bash
   npm install
   ```

#### 2.2.2 配置前端

编辑 `vite.config.ts` 文件，设置后端 API 代理：

```typescript
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  server: {
    port: 3000,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api/, '')
      }
    }
  }
})
```

#### 2.2.3 启动前端开发服务器

```bash
npm run dev
```

前端服务将在 `http://localhost:3000` 启动。

## 3. 生产环境部署

### 3.1 后端部署

#### 3.1.1 构建可执行文件

```bash
go build -o server cmd/server/main.go
```

#### 3.1.2 配置生产环境

1. 创建生产环境配置文件：
   ```bash
   cp configs/config.dev.yaml configs/config.prod.yaml
   ```

2. 编辑配置文件 `configs/config.prod.yaml`，设置生产环境的数据库连接信息和其他配置。

#### 3.1.3 启动后端服务

```bash
./server --config=configs/config.prod.yaml
```

### 3.2 前端部署

#### 3.2.1 构建生产版本

```bash
npm run build
```

构建产物将生成在 `dist` 目录中。

#### 3.2.2 部署静态文件

将 `dist` 目录中的文件部署到静态文件服务器，如 Nginx、Apache 等。

#### 3.2.3 配置反向代理

在 Nginx 中配置反向代理，将 API 请求转发到后端服务：

```nginx
server {
    listen 80;
    server_name example.com;

    location / {
        root /path/to/dist;
        index index.html;
        try_files $uri $uri/ /index.html;
    }

    location /api/ {
        proxy_pass http://localhost:8080/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

## 4. Docker 部署

### 4.1 构建 Docker 镜像

#### 4.1.1 后端 Dockerfile

```Dockerfile
FROM golang:1.20-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o server cmd/server/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/server /app/server
COPY configs/ /app/configs/

EXPOSE 8080

CMD ["./server"]
```

#### 4.1.2 前端 Dockerfile

```Dockerfile
FROM node:16-alpine AS builder

WORKDIR /app

COPY package*.json ./
RUN npm install

COPY . .

RUN npm run build

FROM nginx:alpine

COPY --from=builder /app/dist /usr/share/nginx/html

COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
```

#### 4.1.3 构建镜像

```bash
# 构建后端镜像
docker build -t research-ability-assessment-backend .

# 构建前端镜像
docker build -t research-ability-assessment-frontend ./frontend
```

### 4.2 使用 Docker Compose

创建 `docker-compose.yml` 文件：

```yaml
version: '3.8'

services:
  db:
    image: postgres:14-alpine
    environment:
      POSTGRES_DB: research_assessment
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
    volumes:
      - postgres_data:/var/lib/postgresql/data
    ports:
      - "5432:5432"

  neo4j:
    image: neo4j:5.0
    environment:
      NEO4J_AUTH: neo4j/password
    volumes:
      - neo4j_data:/data
    ports:
      - "7474:7474"
      - "7687:7687"

  backend:
    build: .
    environment:
      - CONFIG_PATH=configs/config.prod.yaml
    depends_on:
      - db
      - neo4j
    ports:
      - "8080:8080"

  frontend:
    build: ./frontend
    depends_on:
      - backend
    ports:
      - "80:80"

volumes:
  postgres_data:
  neo4j_data:
```

启动服务：

```bash
docker-compose up -d
```

## 5. 系统初始化

### 5.1 数据库迁移

后端服务启动时会自动执行数据库迁移，创建所需的表结构。

### 5.2 创建初始用户

1. 访问前端登录页面（默认地址：http://localhost:3000）
2. 点击「注册」按钮
3. 填写注册信息，创建管理员账号
4. 使用管理员账号登录系统

### 5.3 配置系统

1. 登录系统后，进入「系统设置」页面
2. 配置系统参数，如 LLM API 密钥、评估维度等
3. 保存配置

## 6. 系统监控

### 6.1 日志管理

后端服务的日志默认输出到控制台，可通过配置文件设置日志文件路径：

```yaml
logger:
  level: info
  file: logs/app.log
```

### 6.2 健康检查

系统提供健康检查接口，可用于监控系统状态：

```bash
curl http://localhost:8080/api/v1/health
```

### 6.3 性能监控

建议使用 Prometheus 和 Grafana 监控系统性能，配置示例：

```yaml
# prometheus.yml
scrape_configs:
  - job_name: 'research-ability-assessment'
    static_configs:
      - targets: ['backend:8080']
```

## 7. 故障排除

### 7.1 常见问题

#### 7.1.1 数据库连接失败

**症状**：后端服务启动失败，提示数据库连接错误
**解决方案**：
- 检查数据库服务是否运行
- 检查数据库连接配置是否正确
- 检查数据库用户权限是否正确

#### 7.1.2 Neo4j 连接失败

**症状**：后端服务启动失败，提示 Neo4j 连接错误
**解决方案**：
- 检查 Neo4j 服务是否运行
- 检查 Neo4j 连接配置是否正确
- 检查 Neo4j 用户密码是否正确

#### 7.1.3 前端 API 调用失败

**症状**：前端页面加载失败，控制台提示 API 调用错误
**解决方案**：
- 检查后端服务是否运行
- 检查前端代理配置是否正确
- 检查 API 路径是否正确

#### 7.1.4 LLM 调用失败

**症状**：评估过程失败，提示 LLM 调用错误
**解决方案**：
- 检查 OpenAI API 密钥是否正确
- 检查网络连接是否正常
- 检查 LLM 请求参数是否正确

## 8. 系统更新

### 8.1 代码更新

1. 拉取最新代码：
   ```bash
   git pull
   ```

2. 重新构建和部署：
   ```bash
   # 后端
   go build -o server cmd/server/main.go
   ./server

   # 前端
   cd frontend
   npm install
   npm run build
   # 部署构建产物
   ```

### 8.2 数据库更新

如果数据库结构发生变化，系统会自动执行数据库迁移。

## 9. 安全注意事项

### 9.1 环境变量

敏感信息如数据库密码、API 密钥等应通过环境变量或配置文件管理，避免硬编码在代码中。

### 9.2 HTTPS

在生产环境中，应配置 HTTPS 以确保数据传输安全。

### 9.3 访问控制

确保系统有适当的访问控制机制，防止未授权访问。

### 9.4 输入验证

对所有用户输入进行验证，防止 SQL 注入、XSS 等攻击。

### 9.5 定期备份

定期备份数据库和配置文件，确保数据安全。

## 10. 联系方式

如有任何部署问题或建议，请联系系统管理员：
- 邮箱：admin@example.com
- 电话：12345678900
