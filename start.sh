#!/bin/bash

# LLM Fusion Engine - 快速启动脚本
# 适用于 Linux/macOS 系统

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 打印带颜色的消息
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查依赖
check_dependencies() {
    print_info "检查系统依赖..."
    
    if ! command -v docker &> /dev/null; then
        print_error "未找到 Docker,请先安装 Docker"
        exit 1
    fi
    
    if ! command -v docker-compose &> /dev/null; then
        print_error "未找到 Docker Compose,请先安装 Docker Compose"
        exit 1
    fi
    
    print_success "所有依赖检查通过"
}

# 创建环境配置文件
setup_env() {
    if [ ! -f .env ]; then
        print_info "创建环境配置文件..."
        
        cat > .env << EOF
# 数据库配置
DB_HOST=postgres
DB_PORT=5432
DB_USER=llmfusion
DB_PASSWORD=$(openssl rand -base64 32 | tr -dc 'a-zA-Z0-9' | head -c 16)
DB_NAME=llmfusion

# 应用配置
APP_PORT=8080
APP_ENV=production
LOG_LEVEL=info

# 管理员凭证
ADMIN_USERNAME=admin
ADMIN_PASSWORD=$(openssl rand -base64 32 | tr -dc 'a-zA-Z0-9' | head -c 16)
EOF
        
        print_success "环境配置文件已创建: .env"
        print_warning "请查看 .env 文件并根据需要修改配置"
        
        # 显示管理员密码
        ADMIN_PASS=$(grep ADMIN_PASSWORD .env | cut -d '=' -f2)
        print_info "管理员密码: ${GREEN}${ADMIN_PASS}${NC}"
        print_warning "请妥善保存此密码,首次登录后建议修改"
    else
        print_info "环境配置文件已存在,跳过创建"
    fi
}

# 启动服务
start_services() {
    print_info "启动 LLM Fusion Engine..."
    
    # 构建镜像
    print_info "构建 Docker 镜像..."
    docker-compose build --no-cache
    
    # 启动服务
    print_info "启动服务容器..."
    docker-compose up -d
    
    # 等待服务启动
    print_info "等待服务启动..."
    sleep 5
    
    # 检查服务状态
    if docker-compose ps | grep -q "Up"; then
        print_success "服务启动成功!"
        echo ""
        print_info "访问地址:"
        echo -e "  ${GREEN}Web 管理界面:${NC} http://localhost:8080"
        echo -e "  ${GREEN}API 端点:${NC}     http://localhost:8080/v1"
        echo -e "  ${GREEN}管理 API:${NC}     http://localhost:8080/api"
        echo ""
        print_info "查看日志: docker-compose logs -f"
        print_info "停止服务: docker-compose stop"
    else
        print_error "服务启动失败,请查看日志"
        docker-compose logs
        exit 1
    fi
}

# 停止服务
stop_services() {
    print_info "停止 LLM Fusion Engine..."
    docker-compose stop
    print_success "服务已停止"
}

# 重启服务
restart_services() {
    print_info "重启 LLM Fusion Engine..."
    docker-compose restart
    print_success "服务已重启"
}

# 查看日志
view_logs() {
    print_info "显示服务日志 (按 Ctrl+C 退出)..."
    docker-compose logs -f
}

# 清理数据
cleanup() {
    read -p "确定要删除所有数据吗? (yes/no): " confirm
    if [ "$confirm" = "yes" ]; then
        print_warning "删除所有容器和数据..."
        docker-compose down -v
        rm -f .env
        print_success "清理完成"
    else
        print_info "取消清理操作"
    fi
}

# 显示状态
show_status() {
    print_info "服务状态:"
    docker-compose ps
}

# 显示帮助
show_help() {
    echo "LLM Fusion Engine - 快速启动脚本"
    echo ""
    echo "用法: ./start.sh [命令]"
    echo ""
    echo "可用命令:"
    echo "  start    - 启动服务 (默认)"
    echo "  stop     - 停止服务"
    echo "  restart  - 重启服务"
    echo "  logs     - 查看日志"
    echo "  status   - 显示服务状态"
    echo "  cleanup  - 清理所有数据"
    echo "  help     - 显示帮助信息"
    echo ""
    echo "示例:"
    echo "  ./start.sh           # 启动服务"
    echo "  ./start.sh stop      # 停止服务"
    echo "  ./start.sh logs      # 查看日志"
}

# 主函数
main() {
    echo ""
    echo "╔════════════════════════════════════════╗"
    echo "║   LLM Fusion Engine - 快速启动工具    ║"
    echo "╚════════════════════════════════════════╝"
    echo ""
    
    case "${1:-start}" in
        start)
            check_dependencies
            setup_env
            start_services
            ;;
        stop)
            stop_services
            ;;
        restart)
            restart_services
            ;;
        logs)
            view_logs
            ;;
        status)
            show_status
            ;;
        cleanup)
            cleanup
            ;;
        help|--help|-h)
            show_help
            ;;
        *)
            print_error "未知命令: $1"
            show_help
            exit 1
            ;;
    esac
}

# 运行主函数
main "$@"