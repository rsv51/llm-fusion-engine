@echo off
REM LLM Fusion Engine - Windows 快速启动脚本

setlocal enabledelayedexpansion

REM 颜色定义 (使用 PowerShell 实现彩色输出)
set "INFO=[94m[INFO][0m"
set "SUCCESS=[92m[SUCCESS][0m"
set "WARNING=[93m[WARNING][0m"
set "ERROR=[91m[ERROR][0m"

REM 显示欢迎信息
echo.
echo ========================================
echo    LLM Fusion Engine - 快速启动工具
echo ========================================
echo.

REM 解析命令行参数
set "COMMAND=%1"
if "%COMMAND%"=="" set "COMMAND=start"

if "%COMMAND%"=="start" goto START
if "%COMMAND%"=="stop" goto STOP
if "%COMMAND%"=="restart" goto RESTART
if "%COMMAND%"=="logs" goto LOGS
if "%COMMAND%"=="status" goto STATUS
if "%COMMAND%"=="cleanup" goto CLEANUP
if "%COMMAND%"=="help" goto HELP
if "%COMMAND%"=="--help" goto HELP
if "%COMMAND%"=="-h" goto HELP

echo %ERROR% 未知命令: %COMMAND%
goto HELP

:START
call :CHECK_DEPS
if errorlevel 1 exit /b 1

call :SETUP_ENV
call :START_SERVICES
goto END

:STOP
echo %INFO% 停止 LLM Fusion Engine...
docker-compose stop
echo %SUCCESS% 服务已停止
goto END

:RESTART
echo %INFO% 重启 LLM Fusion Engine...
docker-compose restart
echo %SUCCESS% 服务已重启
goto END

:LOGS
echo %INFO% 显示服务日志 (按 Ctrl+C 退出)...
docker-compose logs -f
goto END

:STATUS
echo %INFO% 服务状态:
docker-compose ps
goto END

:CLEANUP
set /p "CONFIRM=确定要删除所有数据吗? (yes/no): "
if /i "%CONFIRM%"=="yes" (
    echo %WARNING% 删除所有容器和数据...
    docker-compose down -v
    if exist .env del .env
    echo %SUCCESS% 清理完成
) else (
    echo %INFO% 取消清理操作
)
goto END

:HELP
echo LLM Fusion Engine - 快速启动脚本
echo.
echo 用法: start.bat [命令]
echo.
echo 可用命令:
echo   start    - 启动服务 (默认)
echo   stop     - 停止服务
echo   restart  - 重启服务
echo   logs     - 查看日志
echo   status   - 显示服务状态
echo   cleanup  - 清理所有数据
echo   help     - 显示帮助信息
echo.
echo 示例:
echo   start.bat           # 启动服务
echo   start.bat stop      # 停止服务
echo   start.bat logs      # 查看日志
goto END

:CHECK_DEPS
echo %INFO% 检查系统依赖...

docker --version >nul 2>&1
if errorlevel 1 (
    echo %ERROR% 未找到 Docker,请先安装 Docker Desktop
    exit /b 1
)

docker-compose --version >nul 2>&1
if errorlevel 1 (
    echo %ERROR% 未找到 Docker Compose,请先安装 Docker Compose
    exit /b 1
)

echo %SUCCESS% 所有依赖检查通过
exit /b 0

:SETUP_ENV
if not exist .env (
    echo %INFO% 创建环境配置文件...
    
    REM 生成随机密码 (简单版本)
    set "DB_PASS=%RANDOM%%RANDOM%"
    set "ADMIN_PASS=%RANDOM%%RANDOM%"
    
    (
        echo # 数据库配置
        echo DB_HOST=postgres
        echo DB_PORT=5432
        echo DB_USER=llmfusion
        echo DB_PASSWORD=!DB_PASS!
        echo DB_NAME=llmfusion
        echo.
        echo # 应用配置
        echo APP_PORT=8080
        echo APP_ENV=production
        echo LOG_LEVEL=info
        echo.
        echo # 管理员凭证
        echo ADMIN_USERNAME=admin
        echo ADMIN_PASSWORD=!ADMIN_PASS!
    ) > .env
    
    echo %SUCCESS% 环境配置文件已创建: .env
    echo %WARNING% 请查看 .env 文件并根据需要修改配置
    echo %INFO% 管理员密码: !ADMIN_PASS!
    echo %WARNING% 请妥善保存此密码,首次登录后建议修改
) else (
    echo %INFO% 环境配置文件已存在,跳过创建
)
exit /b 0

:START_SERVICES
echo %INFO% 启动 LLM Fusion Engine...

echo %INFO% 构建 Docker 镜像...
docker-compose build --no-cache
if errorlevel 1 (
    echo %ERROR% 镜像构建失败
    exit /b 1
)

echo %INFO% 启动服务容器...
docker-compose up -d
if errorlevel 1 (
    echo %ERROR% 服务启动失败
    exit /b 1
)

echo %INFO% 等待服务启动...
timeout /t 5 /nobreak >nul

docker-compose ps | findstr "Up" >nul
if errorlevel 1 (
    echo %ERROR% 服务启动失败,请查看日志
    docker-compose logs
    exit /b 1
)

echo %SUCCESS% 服务启动成功!
echo.
echo %INFO% 访问地址:
echo   Web 管理界面: http://localhost:8080
echo   API 端点:     http://localhost:8080/v1
echo   管理 API:     http://localhost:8080/api
echo.
echo %INFO% 查看日志: docker-compose logs -f
echo %INFO% 停止服务: docker-compose stop
exit /b 0

:END
endlocal