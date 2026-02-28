#!/bin/bash

echo "🚀 启动 X-Novel 开发环境..."

# 检查 PostgreSQL 是否运行
if ! command -v psql &> /dev/null; then
    echo "❌ 未检测到 PostgreSQL，请先安装 PostgreSQL"
    echo "   macOS: brew install postgresql"
    echo "   Ubuntu: sudo apt-get install postgresql"
    exit 1
fi

# 启动后端
echo "📦 启动后端服务..."
cd server
if [ ! -f .env ]; then
    echo "创建 .env 文件..."
    cp .env.example .env
    echo "⚠️  请编辑 server/.env 配置数据库连接和 API Keys"
fi
go run cmd/server/main.go &
SERVER_PID=$!
cd ..

# 等待后端启动
echo "⏳ 等待后端服务启动..."
sleep 3

# 启动前端
echo "🎨 启动前端应用..."
cd web
if [ ! -f .env ]; then
    cp .env.example .env
fi
npm run dev &
WEB_PID=$!
cd ..

echo ""
echo "✅ 开发环境启动成功！"
echo "   后端: http://localhost:8080"
echo "   前端: http://localhost:5173"
echo ""
echo "按 Ctrl+C 停止所有服务"

# 捕获退出信号
trap "echo '正在停止服务...'; kill $SERVER_PID $WEB_PID; exit" INT TERM

# 等待进程
wait
