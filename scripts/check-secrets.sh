#!/bin/bash
# 安全检查脚本 - 检测代码中可能泄露的敏感信息

set -e

echo "🔍 开始安全检查..."
echo ""

# 检查是否存在敏感配置文件
echo "📋 检查敏感配置文件..."
if [ -f "config/conf.yaml" ]; then
    echo "⚠️  警告: config/conf.yaml 存在 - 确保已添加到 .gitignore"
    if git ls-files --error-unmatch config/conf.yaml 2>/dev/null; then
        echo "❌ 错误: config/conf.yaml 已被Git跟踪，请立即移除！"
        echo "   运行: git rm --cached config/conf.yaml"
        exit 1
    else
        echo "✅ 已添加到 .gitignore"
    fi
else
    echo "✅ config/conf.yaml 不存在（正常）"
fi

# 检查.env文件
echo ""
echo "📋 检查环境变量文件..."
if [ -f ".env" ]; then
    echo "⚠️  警告: .env 存在 - 确保已添加到 .gitignore"
    if git ls-files --error-unmatch .env 2>/dev/null; then
        echo "❌ 错误: .env 已被Git跟踪，请立即移除！"
        echo "   运行: git rm --cached .env"
        exit 1
    else
        echo "✅ 已添加到 .gitignore"
    fi
else
    echo "⚠️  .env 文件不存在，请从 .env.example 创建"
fi

# 检查代码中的敏感模式
echo ""
echo "📋 扫描代码中的敏感信息..."

FOUND_ISSUES=0

# 检查硬编码的密码
if grep -r "password.*=.*\"[^$]" --include="*.go" --include="*.yaml" . 2>/dev/null | grep -v ".example" | grep -v "binding:" | grep -v "json:"; then
    echo "❌ 发现硬编码的密码"
    FOUND_ISSUES=1
fi

# 检查硬编码的密钥
if grep -r "secret.*=.*\"[^$]" --include="*.go" --include="*.yaml" . 2>/dev/null | grep -v ".example" | grep -v "jwt" | grep -v "binding:" | grep -v "json:"; then
    echo "❌ 发现硬编码的密钥"
    FOUND_ISSUES=1
fi

# 检查可疑的API密钥格式
if grep -rE "['\"][A-Za-z0-9]{20,}['\"]" --include="*.go" . 2>/dev/null | grep -i "key\|token\|secret" | grep -v "jwt.SigningMethod"; then
    echo "⚠️  发现可能的API密钥"
    FOUND_ISSUES=1
fi

# 检查IP地址（可能是内部IP）
echo ""
echo "📋 检查硬编码的IP地址..."
if grep -rE "[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}" --include="*.go" . 2>/dev/null | grep -v "127.0.0.1\|localhost\|0.0.0.0"; then
    echo "⚠️  发现硬编码的IP地址，建议使用配置文件"
fi

# 检查TODO和FIXME
echo ""
echo "📋 检查未完成的TODO..."
TODO_COUNT=$(grep -r "TODO\|FIXME" --include="*.go" . 2>/dev/null | wc -l)
if [ $TODO_COUNT -gt 0 ]; then
    echo "⚠️  发现 $TODO_COUNT 个 TODO/FIXME 注释"
    grep -rn "TODO\|FIXME" --include="*.go" . 2>/dev/null | head -5
    if [ $TODO_COUNT -gt 5 ]; then
        echo "   ... 还有 $((TODO_COUNT - 5)) 个"
    fi
fi

# 检查print/println调试代码
echo ""
echo "📋 检查调试代码..."
DEBUG_COUNT=$(grep -r "print(\|println(" --include="*.go" . 2>/dev/null | grep -v "fmt.Print\|log.Print" | wc -l)
if [ $DEBUG_COUNT -gt 0 ]; then
    echo "⚠️  发现 $DEBUG_COUNT 个 print/println 调试语句"
    grep -rn "print(\|println(" --include="*.go" . 2>/dev/null | grep -v "fmt.Print\|log.Print"
    FOUND_ISSUES=1
fi

# 总结
echo ""
echo "================================"
if [ $FOUND_ISSUES -eq 0 ]; then
    echo "✅ 安全检查通过！"
    exit 0
else
    echo "❌ 发现安全问题，请修复后再提交"
    exit 1
fi
