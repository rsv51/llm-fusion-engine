#!/bin/bash

# 设置颜色
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "========================================"
echo "LLM Fusion Engine 启动脚本"
echo "========================================"
echo ""

# 检查Go环境
echo "[1/5] 检查Go环境..."
if ! command -v go &> /dev/null; then
    echo -e "${RED}✗ 未检测到Go环境${NC}"
    echo ""
    echo "请先安装Go语言环境："
    echo "下载地址: https://go.dev/dl/"
    echo "建议版本: Go 1.21 或更高"
    echo ""
    exit 1
fi

GO_VERSION=$(go version | awk '{print $3}')
echo -e "${GREEN}✓ Go环境已安装 ${GO_VERSION}${NC}"
echo ""

# 检查数据库文件
echo "[2/5] 检查数据库..."
if [ ! -f "llm_fusion.db" ]; then
    echo -e "${YELLOW}⚠ 数据库文件不存在，正在初始化...${NC}"
    if [ -f "init_database.sql" ]; then
        sqlite3 llm_fusion.db < init_database.sql
        if [ $? -ne 0 ]; then
            echo -e "${RED}✗ 数据库初始化失败${NC}"
            exit 1
        fi
        echo -e "${GREEN}✓ 数据库初始化成功${NC}"
    else
        echo -e "${RED}✗ 找不到 init_database.sql 文件${NC}"
        exit 1
    fi
else
    echo -e "${GREEN}✓ 数据库文件已存在${NC}"
fi
echo ""

# 检查并安装Go依赖
echo "[3/5] 检查Go依赖..."
if [ ! -f "go.sum" ]; then
    echo "正在下载Go依赖..."
    go mod download
    if [ $? -ne 0 ]; then
        echo -e "${RED}✗ 依赖下载失败${NC}"
        exit 1
    fi
fi
echo -e "${GREEN}✓ Go依赖已就绪${NC}"
echo ""

# 构建后端
echo "[4/5] 构建后端应用..."
go build -o llm-fusion-engine ./cmd/server
if [ $? -ne 0 ]; then
    echo -e "${RED}✗ 后端构建失败${NC}"
    exit 1
fi
echo -e "${GREEN}✓ 后端构建成功${NC}"
echo ""

# 启动应用
echo "[5/5] 启动应用..."
echo ""
echo "========================================"
echo "应用正在启动..."
echo "后端地址: http://localhost:8080"
echo "前端地址: http://localhost:5173"
echo "按 Ctrl+C 停止服务"
echo "========================================"
echo ""

# 启动后端（后台运行）
./llm-fusion-engine &
BACKEND_PID=$!
echo "后端进程 PID: $BACKEND_PID"

# 等待后端启动
sleep 3

# 检查Node.js环境
if ! command -v node &> /dev/null; then
    echo -e "${YELLOW}⚠ 未检测到Node.js环境，无法启动前端${NC}"
    echo "请安装Node.js后手动启动前端："
    echo "  cd web"
    echo "  npm install"
    echo "  npm run dev"
    echo ""
    echo "后端已启动，按 Ctrl+C 停止..."
    wait $BACKEND_PID
    exit 0
fi

# 启动前端
cd web
if [ ! -d "node_modules" ]; then
    echo "正在安装前端依赖..."
    npm install
    if [ $? -ne 0 ]; then
        echo -e "${RED}✗ 前端依赖安装失败${NC}"
        kill $BACKEND_PID
        exit 1
    fi
fi

npm run dev &
FRONTEND_PID=$!
echo "前端进程 PID: $FRONTEND_PID"
cd ..

echo ""
echo -e "${GREEN}✓ 应用启动完成！${NC}"
echo ""

# 捕获退出信号，清理进程
trap "echo ''; echo '正在停止服务...'; kill $BACKEND_PID $FRONTEND_PID 2>/dev/null; exit" INT TERM

# 等待
wait