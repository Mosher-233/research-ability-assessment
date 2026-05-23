# 数据库选择指南

本项目现在支持多种数据库选项！以下是详细的配置说明。

## 一、数据库选项说明

### 关系型数据库（选择其一）

1. **MySQL**（默认，推荐）
2. **PostgreSQL**
3. **Supabase**（基于 PostgreSQL 的云服务）

### 图数据库（必须）

- **Neo4j**（用于知识图谱存储）

---

## 二、MySQL（默认配置）

### 使用配置文件：`configs/config.dev.yaml`

```yaml
database:
  type: mysql
  host: localhost
  port: 3306
  user: mysqluser
  password: mysqlpassword
  dbname: research_assessment
  sslmode: disable
```

### Docker 启动

使用项目根目录的 `docker-compose.yml`：

```bash
docker-compose up -d
```

---

## 三、PostgreSQL 配置

### 步骤：

1. **创建 PostgreSQL 配置**

```yaml
database:
  type: postgres
  host: localhost
  port: 5432
  user: postgres
  password: postgres
  dbname: research_assessment
  sslmode: disable
```

2. **创建数据库**

```sql
CREATE DATABASE research_assessment;
```

3. **设置环境变量或修改配置文件**

```bash
# 方式 1：设置环境变量
$env:APP_ENV="dev"

# 方式 2：直接修改 configs/config.dev.yaml
```

4. **启动后端**

```bash
go run cmd/server/main.go
```

---

## 四、Supabase 配置

### 步骤：

1. **注册 Supabase 账号**

访问 [supabase.com](https://supabase.com) 并创建新项目

2. **获取连接信息**

在 Supabase 项目设置中找到：
- Project ID
- Database Password
- Connection Pooler（如果使用）

3. **修改配置文件 `configs/config.dev.yaml`**

```yaml
database:
  type: postgres
  host: your-project-id.supabase.co
  port: 5432
  user: postgres
  password: your_supabase_password
  dbname: postgres
  sslmode: require
```

4. **启动后端**

```bash
go run cmd/server/main.go
```

---

## 五、Neo4j 配置（必须）

无论选择哪种关系型数据库，都需要配置 Neo4j 图数据库。

### Docker 启动（推荐）

使用项目中的 `docker-compose.yml` 已经包含了 Neo4j：

```bash
docker-compose up -d
```

Neo4j Web 界面：http://localhost:7474
默认用户名：neo4j
默认密码：password123

### 配置文件

```yaml
neo4j:
  uri: bolt://localhost:7687
  username: neo4j
  password: password123
```

---

## 六、快速开始

### 方案 A：全部使用 Docker（最简单）

```bash
# 1. 启动所有数据库
docker-compose up -d

# 2. 配置环境变量
copy .env.example .env
# 编辑 .env，填入 DeepSeek API Key

# 3. 安装 Go 依赖
go mod tidy

# 4. 启动后端
go run cmd/server/main.go

# 5. （可选）初始化测试数据
go run scripts/init_db.go

# 6. 启动前端（新终端）
cd frontend
npm install
npm run dev
```

### 方案 B：使用 PostgreSQL + Docker

```bash
# 1. 手动安装 PostgreSQL 或修改 docker-compose.yml
# 2. 修改 configs/config.dev.yaml，将 database.type 改为 postgres
# 3. 启动 Neo4j
docker-compose up -d neo4j

# 4. 启动后端
go run cmd/server/main.go
```

### 方案 C：使用 Supabase

```bash
# 1. 注册 Supabase 并创建项目

# 2. 修改 configs/config.dev.yaml，填入 Supabase 连接信息

# 3. 启动 Neo4j（本地或使用 Docker）
docker-compose up -d neo4j

# 4. 启动后端
go run cmd/server/main.go
```

---

## 七、环境变量配置

你也可以通过环境变量来选择配置文件：

```bash
# Windows PowerShell
$env:APP_ENV="dev"      # 默认（加载 configs/config.dev.yaml）

# 然后启动
go run cmd/server/main.go
```

---

## 八、数据库迁移说明

### PostgreSQL vs MySQL 兼容性

由于项目使用 GORM，大多数功能在两种数据库间是兼容的，但有些差异：

- **自动迁移**：GORM 会自动处理
- **数据类型**：GORM 会自动转换
- **JSON 字段**：都支持
- **时间字段**：都支持

### Neo4j 的作用

Neo4j 用于存储：
- 学生节点
- 能力维度节点
- 学生-能力评分关系
- 知识图谱可视化数据

这部分不能替换为其他数据库。

---

## 九、常见问题

### Q: 可以不用 Neo4j 吗？

A: 目前不行，项目的知识图谱功能依赖 Neo4j。

### Q: Supabase 免费版够用吗？

A: 是的，Supabase 免费版对于开发和小型项目完全够用。

### Q: MySQL 和 PostgreSQL 性能差异大吗？

A: 对于这个项目来说差异不大，选择你熟悉的即可。

### Q: 如何切换数据库？

A: 修改配置文件中的 `database.type` 字段，或者使用不同的配置文件。
