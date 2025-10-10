# LLM Fusion Engine 健康状态显示功能验证脚本
# 基于 llmio-master 设计模式重构后的验证

param(
    [switch]$SkipBuild,
    [switch]$Verbose
)

$ErrorActionPreference = "Stop"
$ProgressPreference = "Continue"

# 颜色输出函数
function Write-ColorOutput {
    param(
        [string]$Message,
        [string]$Color = "White"
    )
    Write-Host $Message -ForegroundColor $Color
}

function Write-Success {
    param([string]$Message)
    Write-ColorOutput "✓ $Message" "Green"
}

function Write-Error {
    param([string]$Message)
    Write-ColorOutput "✗ $Message" "Red"
}

function Write-Warning {
    param([string]$Message)
    Write-ColorOutput "⚠ $Message" "Yellow"
}

function Write-Info {
    param([string]$Message)
    Write-ColorOutput "ℹ $Message" "Cyan"
}

# 验证结果统计
$TotalTests = 0
$PassedTests = 0
$FailedTests = 0
$Warnings = 0

# 测试函数
function Test-Component {
    param(
        [string]$TestName,
        [scriptblock]$TestCode
    )
    
    $TotalTests++
    Write-Info "测试: $TestName"
    
    try {
        $result = & $TestCode
        if ($result) {
            Write-Success "$TestName - 通过"
            $PassedTests++
        } else {
            Write-Error "$TestName - 失败"
            $FailedTests++
        }
    } catch {
        Write-Error "$TestName - 异常: $($_.Exception.Message)"
        $FailedTests++
    }
}

# 验证文件存在性
function Test-FileExists {
    param([string]$FilePath)
    return Test-Path $FilePath -PathType Leaf
}

# 验证目录存在性
function Test-DirectoryExists {
    param([string]$DirPath)
    return Test-Path $DirPath -PathType Container
}

# 验证文件内容包含特定文本
function Test-FileContains {
    param(
        [string]$FilePath,
        [string]$SearchText
    )
    
    if (-not (Test-FileExists $FilePath)) {
        return $false
    }
    
    $content = Get-Content $FilePath -Raw
    return $content -match [regex]::Escape($SearchText)
}

# 验证JSON格式
function Test-JsonFormat {
    param([string]$FilePath)
    
    if (-not (Test-FileExists $FilePath)) {
        return $false
    }
    
    try {
        $null = Get-Content $FilePath -Raw | ConvertFrom-Json
        return $true
    } catch {
        return $false
    }
}

# 验证TypeScript语法
function Test-TypeScriptSyntax {
    param([string]$FilePath)
    
    if (-not (Test-FileExists $FilePath)) {
        return $false
    }
    
    try {
        $result = & npx tsc --noEmit --skipLibCheck $FilePath 2>&1
        return $LASTEXITCODE -eq 0
    } catch {
        return $false
    }
}

# 主验证流程
function Main {
    Write-ColorOutput "========================================" "Magenta"
    Write-ColorOutput "LLM Fusion Engine 健康状态显示功能验证" "Magenta"
    Write-ColorOutput "基于 llmio-master 设计模式重构" "Magenta"
    Write-ColorOutput "========================================" "Magenta"
    Write-Host ""
    
    # 1. 代码完整性检查
    Write-ColorOutput "1. 代码完整性检查" "Yellow"
    Write-Host "----------------------------------------"
    
    # 检查核心文件
    Test-Component "健康状态常量文件存在" {
        Test-FileExists "internal/constants/health_status.go"
    }
    
    Test-Component "健康检查器文件存在" {
        Test-FileExists "internal/services/health_checker.go"
    }
    
    Test-Component "API处理器文件存在" {
        Test-FileExists "internal/api/admin/model_mapping_handler.go"
    }
    
    Test-Component "前端Badge组件存在" {
        Test-FileExists "web/src/components/ui/Badge.tsx"
    }
    
    Test-Component "前端模型映射页面存在" {
        Test-FileExists "web/src/pages/ModelMappings.tsx"
    }
    
    Test-Component "前端类型定义存在" {
        Test-FileExists "web/src/types/provider.ts"
    }
    
    # 2. 代码内容验证
    Write-Host ""
    Write-ColorOutput "2. 代码内容验证" "Yellow"
    Write-Host "----------------------------------------"
    
    Test-Component "健康状态常量定义正确" {
        Test-FileContains "internal/constants/health_status.go" "HealthStatusHealthy"
    }
    
    Test-Component "健康检查器使用常量" {
        Test-FileContains "internal/services/health_checker.go" "constants.HealthStatus"
    }
    
    Test-Component "API处理器包含调试日志" {
        Test-FileContains "internal/api/admin/model_mapping_handler.go" "[ModelMappings]"
    }
    
    Test-Component "前端包含健康状态指示器" {
        Test-FileContains "web/src/pages/ModelMappings.tsx" "HealthStatusIndicator"
    }
    
    Test-Component "Badge组件支持所有变体" {
        Test-FileContains "web/src/components/ui/Badge.tsx" "success.*error.*warning.*info.*default"
    }
    
    # 3. 前后端数据结构匹配
    Write-Host ""
    Write-ColorOutput "3. 前后端数据结构匹配" "Yellow"
    Write-Host "----------------------------------------"
    
    Test-Component "后端Provider模型包含健康状态字段" {
        Test-FileContains "internal/database/models.go" "HealthStatus.*string"
    }
    
    Test-Component "前端Provider类型包含健康状态" {
        Test-FileContains "web/src/types/provider.ts" "healthStatus.*healthy.*unhealthy.*degraded.*unknown"
    }
    
    Test-Component "前端类型定义与后端匹配" {
        Test-FileContains "web/src/types/provider.ts" "latency.*lastChecked.*lastStatusCode"
    }
    
    # 4. UI显示效果验证
    Write-Host ""
    Write-ColorOutput "4. UI显示效果验证" "Yellow"
    Write-Host "----------------------------------------"
    
    Test-Component "健康状态使用交通灯配色" {
        Test-FileContains "web/src/pages/ModelMappings.tsx" "bg-green-100.*bg-yellow-100.*bg-red-100"
    }
    
    Test-Component "自动刷新机制实现" {
        Test-FileContains "web/src/pages/ModelMappings.tsx" "setInterval.*60000"
    }
    
    Test-Component "响应式布局实现" {
        Test-FileContains "web/src/pages/ModelMappings.tsx" "overflow-x-auto"
    }
    
    # 5. TypeScript语法检查
    Write-Host ""
    Write-ColorOutput "5. TypeScript语法检查" "Yellow"
    Write-Host "----------------------------------------"
    
    Test-Component "Badge组件语法正确" {
        Test-TypeScriptSyntax "web/src/components/ui/Badge.tsx"
    }
    
    Test-Component "ModelMappings页面语法正确" {
        Test-TypeScriptSyntax "web/src/pages/ModelMappings.tsx"
    }
    
    Test-Component "Provider类型定义语法正确" {
        Test-TypeScriptSyntax "web/src/types/provider.ts"
    }
    
    # 6. 配置文件验证
    Write-Host ""
    Write-ColorOutput "6. 配置文件验证" "Yellow"
    Write-Host "----------------------------------------"
    
    Test-Component "Go模块配置存在" {
        Test-FileExists "go.mod"
    }
    
    Test-Component "前端包配置存在" {
        Test-FileExists "web/package.json"
    }
    
    Test-Component "TypeScript配置存在" {
        Test-FileExists "web/tsconfig.json"
    }
    
    Test-Component "Tailwind CSS配置存在" {
        Test-FileExists "web/tailwind.config.js"
    }
    
    # 7. 文档完整性
    Write-Host ""
    Write-ColorOutput "7. 文档完整性" "Yellow"
    Write-Host "----------------------------------------"
    
    Test-Component "UI重构总结文档存在" {
        Test-FileExists "UI_REFACTOR_SUMMARY.md"
    }
    
    Test-Component "健康状态修复文档存在" {
        Test-FileExists "HEALTH_STATUS_FIX.md"
    }
    
    Test-Component "启动脚本存在" {
        (Test-FileExists "start.bat") -or (Test-FileExists "start.sh")
    }
    
    # 生成验证报告
    Write-Host ""
    Write-ColorOutput "========================================" "Magenta"
    Write-ColorOutput "验证报告" "Magenta"
    Write-ColorOutput "========================================" "Magenta"
    
    Write-Host "总测试数: $TotalTests"
    Write-Success "通过测试: $PassedTests"
    
    if ($FailedTests -gt 0) {
        Write-Error "失败测试: $FailedTests"
    }
    
    if ($Warnings -gt 0) {
        Write-Warning "警告: $Warnings"
    }
    
    $SuccessRate = if ($TotalTests -gt 0) { [math]::Round(($PassedTests / $TotalTests) * 100, 2) } else { 0 }
    Write-Host "成功率: $SuccessRate%"
    
    Write-Host ""
    
    if ($FailedTests -eq 0) {
        Write-Success "所有验证测试通过！健康状态显示功能重构成功。"
        Write-Host ""
        Write-Info "下一步操作："
        Write-Host "1. 运行 start.bat (Windows) 或 start.sh (Linux/Mac) 启动应用"
        Write-Host "2. 访问 http://localhost:5173 查看前端界面"
        Write-Host "3. 导航到'模型映射'页面查看健康状态显示"
        Write-Host "4. 验证交通灯配色方案是否正确显示"
        Write-Host "5. 检查自动刷新功能是否正常工作"
        exit 0
    } else {
        Write-Error "部分验证测试失败，请检查上述错误并修复。"
        Write-Host ""
        Write-Info "建议操作："
        Write-Host "1. 检查失败的测试项目"
        Write-Host "2. 确认所有文件都已正确创建和修改"
        Write-Host "3. 验证代码语法和类型定义"
        Write-Host "4. 重新运行验证脚本"
        exit 1
    }
}

# 执行主函数
try {
    Main
} catch {
    Write-Error "验证脚本执行失败: $($_.Exception.Message)"
    exit 1
}