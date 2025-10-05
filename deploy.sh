#!/bin/bash

# LLM Fusion Engine 快速部署脚本
# 用于重新构建和部署 Docker 镜像

set -e  # 遇到错误立即退出

echo "=========================================="
echo "  LLM Fusion Engine 部署脚本"
echo "=========================================="
echo ""

# 颜色定义
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# 检查 Docker 是否安装
if ! command -v docker &> /dev/null; then
    echo -e "${RED}错误: Docker 未安装${NC}"
    echo "请先安装 Docker: https://docs.docker.com/get-docker/"
    exit 1
fi

# 检查 Docker Compose 是否安装
if ! command -v docker-compose &> /dev/null; then
    echo -e "${YELLOW}警告: docker-compose 未安装，尝试使用 docker compose${NC}"
    DOCKER_COMPOSE="docker compose"
else
    DOCKER_COMPOSE="docker-compose"
fi

# 显示当前 Git 信息
if command -v git &> /dev/null && [ -d .git ]; then
    echo -e "${GREEN}当前分支:${NC} $(git branch --show-current)"
    echo -e "${GREEN}最新提交:${NC} $(git log -1 --oneline)"
    echo ""
fi

# 询问是否清理旧构建
echo -e "${YELLOW}是否清理旧的 Docker 镜像和容器? (y/N):${NC}"
read -r clean_old

if [[ "$clean_old" =~ ^[Yy]$ ]]; then
    echo -e "${GREEN}正在清理旧构建...${NC}"
    $DOCKER_COMPOSE down --rmi local --volumes --remove-orphans 2>/dev/null || true
    docker system prune -f
    echo -e "${GREEN}清理完成!${NC}"
    echo ""
fi

# 构建镜像
echo -e "${GREEN}开始构建 Docker 镜像...${NC}"
$DOCKER_COMPOSE build --no-cache

if [ $? -ne 0 ]; then
    echo -e "${RED}构建失败!${NC}"
    exit 1
fi

echo -e "${GREEN}构建成功!${NC}"
echo ""

# 启动服务
echo -e "${GREEN}启动服务...${NC}"
$DOCKER_COMPOSE up -d

if [ $? -ne 0 ]; then
    echo -e "${RED}启动失败!${NC}"
    exit 1
fi

echo ""
echo -e "${GREEN}部署完成!${NC}"
echo ""

# 等待服务启动
echo -e "${YELLOW}等待服务启动...${NC}"
sleep 5

# 检查服务状态
echo ""
echo -e "${GREEN}服务状态:${NC}"
$DOCKER_COMPOSE ps

echo ""
echo -e "${GREEN}查看实时日志:${NC}"
echo "  $DOCKER_COMPOSE logs -f"
echo ""
echo -e "${GREEN}停止服务:${NC}"
echo "  $DOCKER_COMPOSE down"
echo ""
echo -e "${GREEN}访问应用:${NC}"
echo "  http://localhost:8080"
echo ""

# 询问是否查看日志
echo -e "${YELLOW}是否查看实时日志? (Y/n):${NC}"
read -r show_logs

if [[ ! "$show_logs" =~ ^[Nn]$ ]]; then
    echo ""
    echo -e "${GREEN}按 Ctrl+C 退出日志查看${NC}"
    echo ""
    $DOCKER_COMPOSE logs -f
fi