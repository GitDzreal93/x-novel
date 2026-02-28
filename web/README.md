# X-Novel Web

X-Novel AI 小说创作平台前端应用

## 技术栈

- React 18
- TypeScript
- Vite
- Ant Design 5
- React Router v6
- Zustand (状态管理)
- React Query (数据请求)
- Axios

## 项目结构

```
src/
├── api/          # API 请求
├── assets/       # 静态资源
├── components/   # 组件
│   ├── common/   # 通用组件
│   ├── project/  # 项目相关组件
│   ├── chapter/  # 章节相关组件
│   └── device/   # 设备相关组件
├── pages/        # 页面
├── hooks/        # 自定义 Hooks
├── stores/       # 状态管理
├── types/        # TypeScript 类型
├── utils/        # 工具函数
└── styles/       # 样式文件
```

## 快速开始

### 1. 安装依赖

```bash
npm install
```

### 2. 配置环境变量

复制 `.env.example` 到 `.env`：

```bash
cp .env.example .env
```

根据需要修改 `.env` 中的 API 地址。

### 3. 启动开发服务器

```bash
npm run dev
```

访问 `http://localhost:5173` 查看应用。

### 4. 构建生产版本

```bash
npm run build
```

构建产物将输出到 `dist` 目录。

## 功能页面

### 项目列表 (`/projects`)
- 查看所有项目
- 创建新项目
- 删除项目
- 查看项目进度

### 项目详情 (`/projects/:id`)
- **小说架构**: 5 步骤架构生成和编辑
- **章节大纲**: 章节大纲生成和管理
- **章节写作**: 章节内容生成和编辑

## 开发说明

### 设备识别

应用使用设备 ID 进行识别，无需登录注册。设备 ID 会自动生成并存储在 LocalStorage 中。

### 主题切换

支持浅色/深色主题切换，设置会自动保存。

### API 请求

所有 API 请求都会自动添加设备 ID 请求头。

## 许可证

MIT
