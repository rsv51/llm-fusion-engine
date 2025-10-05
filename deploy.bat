@echo off
REM LLM Fusion Engine 快速部署脚本 (Windows)
REM 用于重新构建和部署 Docker 镜像

echo ==========================================
echo   LLM Fusion Engine 部署脚本
echo ==========================================
echo.

REM 检查 Docker 是否安装
docker --version >nul 2>&1
if errorlevel 1 (
    echo [错误] Docker 未安装
    echo 请先安装 Docker Desktop: https://www.docker.com/products/docker-desktop
    pause
    exit /b 1
)

REM 检查 docker-compose 是否可用
docker-compose --version >nul 2>&1
if errorlevel 1 (
    echo [警告] docker-compose 未安装，尝试使用 docker compose
    set DOCKER_COMPOSE=docker compose
) else (
    set DOCKER_COMPOSE=docker-compose
)

REM 显示当前 Git 信息
git --version >nul 2>&1
if not errorlevel 1 (
    if exist .git (
        echo [信息] 当前分支: 
        git branch --show-current
        echo [信息] 最新提交:
        git log -1 --oneline
        echo.
    )
)

REM 询问是否清理旧构建
set /p clean_old="是否清理旧的 Docker 镜像和容器? (y/N): "
if /i "%clean_old%"=="y" (
    echo [执行] 正在清理旧构建...
    %DOCKER_COMPOSE% down --rmi local --volumes --remove-orphans 2>nul
    docker system prune -f
    echo [完成] 清理完成!
    echo.
)

REM 构建镜像
echo [执行] 开始构建 Docker 镜像...
%DOCKER_COMPOSE% build --no-cache

if errorlevel 1 (
    echo [错误] 构建失败!
    pause
    exit /b 1
)

echo [完成] 构建成功!
echo.

REM 启动服务
echo [执行] 启动服务...
%DOCKER_COMPOSE% up -d

if errorlevel 1 (
    echo [错误] 启动失败!
    pause
    exit /b 1
)

echo.
echo [完成] 部署完成!
echo.

REM 等待服务启动
echo [信息] 等待服务启动...
timeout /t 5 /nobreak >nul

REM 检查服务状态
echo.
echo [信息] 服务状态:
%DOCKER_COMPOSE% ps

echo.
echo [信息] 查看实时日志:
echo   %DOCKER_COMPOSE% logs -f
echo.
echo [信息] 停止服务:
echo   %DOCKER_COMPOSE% down
echo.
echo [信息] 访问应用:
echo   http://localhost:8080
echo.

REM 询问是否查看日志
set /p show_logs="是否查看实时日志? (Y/n): "
if /i not "%show_logs%"=="n" (
    echo.
    echo [信息] 按 Ctrl+C 退出日志查看
    echo.
    %DOCKER_COMPOSE% logs -f
)