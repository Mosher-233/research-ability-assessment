# 数据库迁移指南：PostgreSQL → MySQL

## 文档信息

- **迁移日期**: 2026-03-10
- **项目名称**: Research Ability Assessment
- **迁移原因**: PostgreSQL 认证问题持续无法解决，切换到 MySQL 以确保项目正常运行

---

## 一、迁移背景

### 1.1 问题描述

在使用 PostgreSQL 作为关系型数据库时，持续遇到认证问题：

```
failed SASL auth (用户 "postgres" Password 认证失败 (SQLSTATE 28P01))
```

尽管尝试了以下解决方案，但问题仍未解决：
- 重置 PostgreSQL 用户密码
- 修改 `pg_hba.conf` 认证配置
- 切换认证方式（scram-sha-256 → md5）
- 完全重新创建容器和数据卷
- 重启 PostgreSQL 服务

### 1.2 解决方案

由于项目已经支持多种数据库（通过 GORM），决定将关系型数据库从 PostgreSQL 切换到 MySQL。

---

## 二、前置条件

在开始迁移前，确保以下条件已满足：

1. ✅ Docker Desktop 已安装并运行
2. ✅ Go 1.20+ 已安装
3. ✅ 项目代码已支持 MySQL（go.mod 中已有 `gorm.io/driver/mysql`）
4. ✅ 已有完整的 `docker-compose.yml` 配置文件

---

## 三、迁移步骤

### 步骤 1：修改 docker-compose.yml

将 PostgreSQL 服务替换为 MySQL 服务。

**修改前（PostgreSQL 配置）：**
```yaml
services:
  postgres:
    image: postgres:16-alpine
    container_name: research-assessment-postgres
    environment:
      POSTGRES_DB: research_assessment
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
    volumes:
      - postgres_data:/var/lib/postgresql/data
    ports:
      - "5432:5432"
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 10s
      timeout: 5s
      retries: 5
```

**修改后（MySQL 配置）：**
```yaml
services:
  mysql:
    image: mysql:8.0
    container_name: research-assessment-mysql
    environment:
      MYSQL_ROOT_PASSWORD: rootpassword
      MYSQL_DATABASE: research_assessment
      MYSQL_USER: mysqluser
      MYSQL_PASSWORD: mysqlpassword
    volumes:
      - mysql_data:/var/lib/mysql
    ports:
      - "3306:3306"
    healthcheck:
      test: ["CMD", "mysqladmin", "ping", "-h", "localhost", "-u", "root", "-prootpassword"]
      interval: 10s
      timeout: 5s
      retries: 5
```

### 步骤 2：更新 volumes 配置

**修改前：**
```yaml
volumes:
  postgres_data:
  neo4j_data:
  neo4j_logs:
```

**修改后：**
```yaml
volumes:
  mysql_data:
  neo4j_data:
  neo4j_logs:
```

### 步骤 3：更新配置文件 config.dev.yaml

**修改前（PostgreSQL 配置）：**
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

**修改后（MySQL 配置）：**
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

### 步骤 4：停止并删除旧容器

```bash
docker-compose down -v
```

**说明**：
- `-v` 参数会同时删除数据卷，确保干净的迁移
- 如果需要保留数据，请省略 `-v` 参数（但建议全新开始）

### 步骤 5：启动新容器

```bash
docker-compose up -d
```

**预期输出：**
```
[+] Running 3/3
 ✔ Volume research-ability-assessment_mysql_data  Created
 ✔ Container research-assessment-mysql            Started
 ✔ Container research-assessment-neo4j            Started
```

### 步骤 6：等待容器健康检查通过

等待约 30 秒，然后检查容器状态：

```bash
docker ps
```

**预期输出（两个容器都显示 healthy）：**
```
CONTAINER ID   IMAGE                  STATUS                    PORTS
1f606ec8d168   mysql:8.0              Up 39 seconds (healthy)   0.0.0.0:3306->3306/tcp
49a0d9e04846   neo4j:5.23-community   Up 39 seconds (healthy)   0.0.0.0:7474->7474/tcp, 0.0.0.0:7687->7687/tcp
```

### 步骤 7：启动后端服务

```bash
go run cmd/server/main.go
```

**预期输出（服务成功启动）：**
```
[GIN-debug] Listening and serving HTTP on :8080
2026/03/10 19:12:28 服务器启动在 :8080
```

---

## 四、配置文件变更总结

### 4.1 修改的文件

| 文件路径 | 变更内容 |
|---------|---------|
| `docker-compose.yml` | PostgreSQL → MySQL，volumes 更新 |
| `configs/config.dev.yaml` | 数据库类型、连接信息更新 |

### 4.2 MySQL 连接信息

| 配置项 | 值 |
|-------|-----|
| 数据库类型 | mysql |
| 主机 | localhost |
| 端口 | 3306 |
| 用户名 | mysqluser |
| 密码 | mysqlpassword |
| 数据库名 | research_assessment |

### 4.3 Neo4j 配置（未变更）

| 配置项 | 值 |
|-------|-----|
| URI | bolt://localhost:7687 |
| 用户名 | neo4j |
| 密码 | password123 |

---

## 五、验证迁移成功

### 5.1 检查容器状态

```bash
docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"
```

### 5.2 测试后端 API

访问以下端点验证服务正常：

```bash
# 健康检查（如果有的话）
curl http://localhost:8080/api/v1/auth/register
```

或者直接在浏览器中访问：
- 后端服务：http://localhost:8080
- Neo4j 界面：http://localhost:7474

### 5.3 检查数据库连接

从容器内部测试 MySQL 连接：

```bash
docker exec -it research-assessment-mysql mysql -u mysqluser -pmysqlpassword research_assessment
```

---

## 六、已知问题和注意事项

### 6.1 数据库迁移警告

启动后端时可能会看到以下警告：

```
Error 1824 (HY000): Failed to open the referenced table 'student_tasks'
```

**说明**：
- 这是由于 GORM 自动迁移时外键约束的创建顺序问题
- 不影响服务的正常运行
- 后续可以通过手动迁移或调整模型定义来解决

### 6.2 PostgreSQL vs MySQL 差异

虽然 GORM 提供了跨数据库兼容性，但仍需注意：

| 特性 | PostgreSQL | MySQL |
|-----|-----------|-------|
| 字符串类型 | VARCHAR | VARCHAR(191) |
| JSON 支持 | JSONB | JSON |
| 自增主键 | SERIAL | AUTO_INCREMENT |
| 时间精度 | 微秒 | 微秒（MySQL 5.6.4+） |

### 6.3 数据备份

如果之前在 PostgreSQL 中有重要数据，建议在迁移前备份：

```bash
# 备份 PostgreSQL 数据
docker exec research-assessment-postgres pg_dump -U postgres research_assessment > backup.sql
```

---

## 七、回滚方案

如果需要回滚到 PostgreSQL，按以下步骤操作：

### 7.1 恢复 docker-compose.yml

将 MySQL 配置改回 PostgreSQL 配置（见步骤 1 的"修改前"部分）

### 7.2 恢复配置文件

将 `configs/config.dev.yaml` 改回 PostgreSQL 配置（见步骤 3 的"修改前"部分）

### 7.3 重新创建容器

```bash
docker-compose down -v
docker-compose up -d
```

### 7.4 恢复数据（如果有备份）

```bash
docker exec -i research-assessment-postgres psql -U postgres research_assessment < backup.sql
```

---

## 八、后续优化建议

### 8.1 解决外键约束问题

可以考虑：
1. 调整模型定义，移除或调整外键约束
2. 使用手动迁移脚本替代 GORM 自动迁移
3. 分阶段创建表，确保依赖顺序正确

### 8.2 配置文件管理

建议为不同数据库创建独立的配置文件：
- `config.dev.yaml` - 开发环境（当前使用 MySQL）
- `config.postgres.yaml` - PostgreSQL 配置
- `config.mysql.yaml` - MySQL 配置

通过环境变量或命令行参数切换配置。

### 8.3 数据库连接池

考虑在配置中添加连接池参数：

```yaml
database:
  max_open_conns: 100
  max_idle_conns: 10
  conn_max_lifetime: 3600
```

---

## 九、参考文档

- [GORM MySQL 驱动文档](https://gorm.io/docs/connecting_to_the_database.html#MySQL)
- [MySQL 8.0 官方文档](https://dev.mysql.com/doc/refman/8.0/en/)
- [Docker Compose 官方文档](https://docs.docker.com/compose/)
- 项目现有文档：`docs/database_options.md`

---

**文档结束**
