# X-Novel Server

X-Novel AI 小说创作平台后端服务

## 技术栈

- Go 1.21+
- Gin Web Framework
- GORM
- PostgreSQL
- Zap Logger
- Viper Configuration

## 项目结构

```
.
├── cmd/server/       # 入口文件
├── internal/
│   ├── api/         # API 层
│   │   ├── handler/ # 处理器
│   │   ├── middleware/ # 中间件
│   │   └── router/ # 路由
│   ├── service/     # 业务逻辑层
│   ├── repository/  # 数据访问层
│   ├── model/       # 数据模型
│   ├── dto/         # 数据传输对象
│   ├── llm/         # LLM 适配器
│   └── config/      # 配置
├── pkg/             # 工具包
├── migrations/      # 数据库迁移
└── prompts/        # AI 提示词模板
```

## 快速开始

### 1. 安装依赖

```bash
go mod tidy
```

### 2. 配置数据库

创建 PostgreSQL 数据库：

```sql
CREATE DATABASE x_novel;
```

### 3. 配置环境变量

复制 `.env.example` 到 `.env` 并修改配置：

```bash
cp .env.example .env
```

### 4. 运行服务

```bash
go run cmd/server/main.go
```

服务将在 `http://localhost:8080` 启动

## API 文档

### 设备相关

- `GET /api/v1/device/info` - 获取设备信息
- `GET /api/v1/device/settings` - 获取设备设置
- `PUT /api/v1/device/settings` - 更新设备设置

### 项目相关

- `GET /api/v1/projects` - 获取项目列表
- `POST /api/v1/projects` - 创建项目
- `GET /api/v1/projects/:id` - 获取项目详情
- `PUT /api/v1/projects/:id` - 更新项目
- `DELETE /api/v1/projects/:id` - 删除项目
- `POST /api/v1/projects/:id/architecture/generate` - 生成小说架构
- `POST /api/v1/projects/:id/blueprint/generate` - 生成章节大纲
- `GET /api/v1/projects/:id/export/:format` - 导出项目

### 章节相关

- `GET /api/v1/projects/:id/chapters` - 获取章节列表
- `POST /api/v1/projects/:id/chapters` - 创建章节
- `GET /api/v1/projects/:id/chapters/:number` - 获取章节详情
- `PUT /api/v1/projects/:id/chapters/:number` - 更新章节
- `POST /api/v1/projects/:id/chapters/:number/generate` - 生成章节内容
- `POST /api/v1/projects/:id/chapters/:number/finalize` - 定稿章节
- `POST /api/v1/projects/:id/chapters/:number/enrich` - 扩写章节

## 开发

### 数据库迁移

启动服务时会自动执行数据库迁移。

### 日志

日志使用 Zap，支持 JSON 和 Console 两种格式。

### 配置

配置文件为 `config.yaml`，支持环境变量覆盖。

## 许可证

MIT
