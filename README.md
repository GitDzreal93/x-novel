# X-Novel AI 小说创作平台

基于 AI 的智能小说创作辅助工具，帮助作者从灵感涌现到完成出版的全流程创作。

## 项目特点

- **无需登录注册**：基于设备 ID 识别，打开即用
- **AI 辅助创作**：集成多种大模型，支持架构生成、大纲生成、章节写作
- **完整创作流程**：项目管理 → 架构设计 → 大纲规划 → 章节写作 → 导出发布

## 技术栈

### 后端
- Go 1.21+
- Gin Web Framework
- GORM + PostgreSQL
- 多 LLM 适配器（OpenAI、Anthropic、DeepSeek 等）

### 前端
- React 18 + TypeScript
- Vite
- Ant Design 5
- Zustand + React Query
- React Router v6

## 项目结构

```
x-novel/
├── server/          # 后端服务
│   ├── cmd/         # 入口文件
│   ├── internal/    # 源代码
│   ├── pkg/         # 工具包
│   └── migrations/  # 数据库迁移
├── web/             # 前端应用
│   ├── src/         # 源代码
│   └── public/      # 静态资源
└── 竞品项目/        # 竞品分析
    └── docs/        # 文档
```

## 快速开始

### 环境要求

- Go 1.21+
- Node.js 18+
- PostgreSQL 15+

### 1. 数据库准备

创建 PostgreSQL 数据库：

```sql
CREATE DATABASE x_novel;
```

### 2. 启动后端

```bash
cd server
cp .env.example .env
# 编辑 .env 配置数据库连接和 API Keys
go mod tidy
go run cmd/server/main.go
```

后端服务将在 `http://localhost:8080` 启动

### 3. 启动前端

```bash
cd web
npm install
npm run dev
```

前端应用将在 `http://localhost:5173` 启动

## 核心功能

### MVP（第一阶段）
- ✅ 项目管理：创建、编辑、删除项目
- ✅ 小说架构：基于雪花写作法的 5 步骤架构生成
- ✅ 章节大纲：支持长篇分块生成的章节大纲
- ✅ 章节写作：AI 辅助生成章节内容
- ✅ 数据导出：支持 TXT 和 Markdown 格式

### 后续规划
- ⏳ 关系图谱：角色关系可视化（灵感罗盘）
- ⏳ 错误检测：错别字、逻辑一致性检测
- ⏳ 智能对话：ChatGPT 风格的灵感激发对话
- ⏳ 写作助手：润色、扩写、灵感建议
- ⏳ AI 审阅：质量评分、优缺点分析
- ⏳ 多平台发布：一键发布到主流小说平台

## 开发进度

详见 [功能清单-追踪版.md](./竞品项目/docs/功能清单-追踪版.md)

## 文档

- [竞品分析报告](./竞品项目/docs/竞品分析报告.md)
- [技术方案](./竞品项目/docs/技术方案-前后端分离版.md)
- [开发计划](./竞品项目/docs/开发计划.md)
- [AI 提示词模板分析](./竞品项目/docs/AI提示词模板分析.md)

## 贡献

欢迎提交 Issue 和 Pull Request！

## 许可证

MIT
