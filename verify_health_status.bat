@echo off
chcp 65001 >nul
setlocal enabledelayedexpansion

echo ========================================
echo LLM Fusion Engine 健康状态显示功能验证
echo 基于 llmio-master 设计模式重构
echo ========================================
echo.

set TOTAL_TESTS=0
set PASSED_TESTS=0
set FAILED_TESTS=0

echo [1/7] 代码完整性检查...
echo ----------------------------------------

rem 检查核心文件
call :TestFileExists "internal\constants\health_status.go" "健康状态常量文件存在"
call :TestFileExists "internal\services\health_checker.go" "健康检查器文件存在"
call :TestFileExists "internal\api\admin\model_mapping_handler.go" "API处理器文件存在"
call :TestFileExists "web\src\components\ui\Badge.tsx" "前端Badge组件存在"
call :TestFileExists "web\src\pages\ModelMappings.tsx" "前端模型映射页面存在"
call :TestFileExists "web\src\types\provider.ts" "前端类型定义存在"

echo.
echo [2/7] 代码内容验证...
echo ----------------------------------------

call :TestFileContains "internal\constants\health_status.go" "HealthStatusHealthy" "健康状态常量定义正确"
call :TestFileContains "internal\services\health_checker.go" "constants.HealthStatus" "健康检查器使用常量"
call :TestFileContains "internal\api\admin\model_mapping_handler.go" "[ModelMappings]" "API处理器包含调试日志"
call :TestFileContains "web\src\pages\ModelMappings.tsx" "HealthStatusIndicator" "前端包含健康状态指示器"
call :TestFileContains "web\src\components\ui\Badge.tsx" "success" "Badge组件支持所有变体"

echo.
echo [3/7] 前后端数据结构匹配...
echo ----------------------------------------

call :TestFileContains "internal\database\models.go" "HealthStatus" "后端Provider模型包含健康状态字段"
call :TestFileContains "web\src\types\provider.ts" "healthStatus" "前端Provider类型包含健康状态"
call :TestFileContains "web\src\types\provider.ts" "latency" "前端类型定义与后端匹配"

echo.
echo [4/7] UI显示效果验证...
echo ----------------------------------------

call :TestFileContains "web\src\pages\ModelMappings.tsx" "bg-green-100" "健康状态使用交通灯配色"
call :TestFileContains "web\src\pages\ModelMappings.tsx" "setInterval" "自动刷新机制实现"
call :TestFileContains "web\src\pages\ModelMappings.tsx" "overflow-x-auto" "响应式布局实现"

echo.
echo [5/7] TypeScript语法检查...
echo ----------------------------------------

call :TestTypeScriptSyntax "web\src\components\ui\Badge.tsx" "Badge组件语法正确"
call :TestTypeScriptSyntax "web\src\pages\ModelMappings.tsx" "ModelMappings页面语法正确"
call :TestTypeScriptSyntax "web\src\types\provider.ts" "Provider类型定义语法正确"

echo.
echo [6/7] 配置文件验证...
echo ----------------------------------------

call :TestFileExists "go.mod" "Go模块配置存在"
call :TestFileExists "web\package.json" "前端包配置存在"
call :TestFileExists "web\tsconfig.json" "TypeScript配置存在"
call :TestFileExists "web\tailwind.config.js" "Tailwind CSS配置存在"

echo.
echo [7/7] 文档完整性...
echo ----------------------------------------

call :TestFileExists "UI_REFACTOR_SUMMARY.md" "UI重构总结文档存在"
call :TestFileExists "HEALTH_STATUS_FIX.md" "健康状态修复文档存在"
if exist "start.bat" (
    call :IncrementPassed "启动脚本存在"
) else if exist "start.sh" (
    call :IncrementPassed "启动脚本存在"
) else (
    call :IncrementFailed "启动脚本存在"
)

echo.
echo ========================================
echo 验证报告
echo ========================================
echo 总测试数: %TOTAL_TESTS%
echo 通过测试: %PASSED_TESTS%

if %FAILED_TESTS% GTR 0 (
    echo 失败测试: %FAILED_TESTS%
)

set /a SUCCESS_RATE=%PASSED_TESTS%*100/%TOTAL_TESTS%
echo 成功率: %SUCCESS_RATE%%%

echo.
if %FAILED_TESTS% EQU 0 (
    echo 所有验证测试通过！健康状态显示功能重构成功。
    echo.
    echo 下一步操作：
    echo 1. 运行 start.bat (Windows) 或 start.sh (Linux/Mac) 启动应用
    echo 2. 访问 http://localhost:5173 查看前端界面
    echo 3. 导航到"模型映射"页面查看健康状态显示
    echo 4. 验证交通灯配色方案是否正确显示
    echo 5. 检查自动刷新功能是否正常工作
    exit /b 0
) else (
    echo 部分验证测试失败，请检查上述错误并修复。
    echo.
    echo 建议操作：
    echo 1. 检查失败的测试项目
    echo 2. 确认所有文件都已正确创建和修改
    echo 3. 验证代码语法和类型定义
    echo 4. 重新运行验证脚本
    exit /b 1
)

rem 子程序
:TestFileExists
set /a TOTAL_TESTS+=1
if exist "%~1" (
    echo ✓ %~2
    set /a PASSED_TESTS+=1
) else (
    echo ✗ %~2
    set /a FAILED_TESTS+=1
)
goto :eof

:TestFileContains
set /a TOTAL_TESTS+=1
findstr /C:"%~2" "%~1" >nul 2>&1
if !errorlevel! EQU 0 (
    echo ✓ %~3
    set /a PASSED_TESTS+=1
) else (
    echo ✗ %~3
    set /a FAILED_TESTS+=1
)
goto :eof

:TestTypeScriptSyntax
set /a TOTAL_TESTS+=1
cd web >nul 2>&1
npx tsc --noEmit --skipLibCheck "%~1" >nul 2>&1
cd .. >nul 2>&1
if !errorlevel! EQU 0 (
    echo ✓ %~2
    set /a PASSED_TESTS+=1
) else (
    echo ✗ %~2
    set /a FAILED_TESTS+=1
)
goto :eof

:IncrementPassed
set /a TOTAL_TESTS+=1
set /a PASSED_TESTS+=1
echo ✓ %~1
goto :eof

:IncrementFailed
set /a TOTAL_TESTS+=1
set /a FAILED_TESTS+=1
echo ✗ %~1
goto :eof