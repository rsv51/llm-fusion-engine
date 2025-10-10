@echo off
chcp 65001 >nul
echo ========================================
echo 健康状态修复验证脚本
echo ========================================
echo.

echo [1/5] 检查修复文件是否存在...
set FILES_OK=1

if not exist "internal\constants\health_status.go" (
    echo ✗ 缺少: internal\constants\health_status.go
    set FILES_OK=0
)

if not exist "internal\services\health_checker.go" (
    echo ✗ 缺少: internal\services\health_checker.go
    set FILES_OK=0
)

if not exist "internal\api\admin\model_mapping_handler.go" (
    echo ✗ 缺少: internal\api\admin\model_mapping_handler.go
    set FILES_OK=0
)

if not exist "web\src\pages\ModelMappings.tsx" (
    echo ✗ 缺少: web\src\pages\ModelMappings.tsx
    set FILES_OK=0
)

if %FILES_OK%==1 (
    echo ✓ 所有修复文件都存在
) else (
    echo.
    echo ✗ 部分文件缺失，请检查修复是否完整
    pause
    exit /b 1
)
echo.

echo [2/5] 检查常量文件内容...
findstr /C:"HealthStatusHealthy" internal\constants\health_status.go >nul 2>&1
if errorlevel 1 (
    echo ✗ health_status.go 中未找到健康状态常量
    set FILES_OK=0
) else (
    echo ✓ 健康状态常量定义正确
)
echo.

echo [3/5] 检查健康检查器更新...
findstr /C:"constants.HealthStatus" internal\services\health_checker.go >nul 2>&1
if errorlevel 1 (
    echo ✗ health_checker.go 未使用健康状态常量
    set FILES_OK=0
) else (
    echo ✓ 健康检查器已使用常量
)
echo.

echo [4/5] 检查API日志...
findstr /C:"[ModelMappings]" internal\api\admin\model_mapping_handler.go >nul 2>&1
if errorlevel 1 (
    echo ✗ model_mapping_handler.go 未添加调试日志
    set FILES_OK=0
) else (
    echo ✓ API调试日志已添加
)
echo.

echo [5/5] 检查前端组件...
findstr /C:"HealthStatusIndicator" web\src\pages\ModelMappings.tsx >nul 2>&1
if errorlevel 1 (
    echo ✗ ModelMappings.tsx 缺少健康状态组件
    set FILES_OK=0
) else (
    echo ✓ 前端健康状态组件存在
)
echo.

echo ========================================
echo 验证结果
echo ========================================
if %FILES_OK%==1 (
    echo ✓ 所有修复验证通过！
    echo.
    echo 下一步：
    echo 1. 运行 start.bat 启动应用
    echo 2. 访问 http://localhost:5173
    echo 3. 查看模型映射页面的健康状态列
    echo 4. 检查后端日志中的调试信息
) else (
    echo ✗ 部分验证未通过，请检查修复内容
)
echo ========================================
echo.
pause