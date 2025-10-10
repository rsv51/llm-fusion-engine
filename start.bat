@echo off
chcp 65001 >nul
echo ========================================
echo LLM Fusion Engine 启动脚本
echo ========================================
echo.

:: 检查Go环境
echo [1/5] 检查Go环境...
where go >nul 2>nul
if errorlevel 1 (
    echo ✗ 未检测到Go环境
    echo.
    echo 请先安装Go语言环境：
    echo 下载地址: https://go.dev/dl/
    echo 建议版本: Go 1.21 或更高
    echo.
    pause
    exit /b 1
)

for /f "tokens=3" %%v in ('go version') do set GO_VERSION=%%v
echo ✓ Go环境已安装 %GO_VERSION%
echo.

:: 检查数据库文件
echo [2/5] 检查数据库...
if not exist "llm_fusion.db" (
    echo ⚠ 数据库文件不存在，正在初始化...
    if exist "init_database.sql" (
        sqlite3 llm_fusion.db < init_database.sql
        if errorlevel 1 (
            echo ✗ 数据库初始化失败
            pause
            exit /b 1
        )
        echo ✓ 数据库初始化成功
    ) else (
        echo ✗ 找不到 init_database.sql 文件
        pause
        exit /b 1
    )
) else (
    echo ✓ 数据库文件已存在
)
echo.

:: 检查并安装Go依赖
echo [3/5] 检查Go依赖...
if not exist "go.sum" (
    echo 正在下载Go依赖...
    go mod download
    if errorlevel 1 (
        echo ✗ 依赖下载失败
        pause
        exit /b 1
    )
)
echo ✓ Go依赖已就绪
echo.

:: 构建后端
echo [4/5] 构建后端应用...
go build -o llm-fusion-engine.exe ./cmd/server
if errorlevel 1 (
    echo ✗ 后端构建失败
    pause
    exit /b 1
)
echo ✓ 后端构建成功
echo.

:: 启动应用
echo [5/5] 启动应用...
echo.
echo ========================================
echo 应用正在启动...
echo 后端地址: http://localhost:8080
echo 前端地址: http://localhost:5173
echo 按 Ctrl+C 停止服务
echo ========================================
echo.

start "LLM Fusion Engine - Backend" llm-fusion-engine.exe

:: 等待后端启动
timeout /t 3 /nobreak >nul

:: 检查Node.js环境
where node >nul 2>nul
if errorlevel 1 (
    echo ⚠ 未检测到Node.js环境，无法启动前端
    echo 请安装Node.js后手动启动前端：
    echo   cd web
    echo   npm install
    echo   npm run dev
    echo.
    echo 按任意键关闭...
    pause >nul
    exit /b 0
)

:: 启动前端
cd web
if not exist "node_modules" (
    echo 正在安装前端依赖...
    call npm install
    if errorlevel 1 (
        echo ✗ 前端依赖安装失败
        cd ..
        pause
        exit /b 1
    )
)

start "LLM Fusion Engine - Frontend" cmd /k "npm run dev"
cd ..

echo.
echo ✓ 应用启动完成！
echo.
pause