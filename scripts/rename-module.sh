#!/bin/bash

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 脚本说明
echo -e "${GREEN}╔════════════════════════════════════════════╗"
echo -e "║   Go 项目模块名称重命名工具 v1.0         ║"
echo -e "╚════════════════════════════════════════════╝${NC}\n"

# 检查参数
if [ $# -eq 0 ]; then
    echo -e "${RED}❌ 错误: 请提供新的模块名称${NC}"
    echo ""
    echo -e "${BLUE}用法:${NC}"
    echo "  $0 <new-module-name>"
    echo ""
    echo -e "${BLUE}示例:${NC}"
    echo "  $0 github.com/yourname/yourproject"
    echo "  $0 github.com/mycompany/awesome-api"
    echo "  $0 mycompany.com/internal/api-service"
    echo ""
    echo -e "${YELLOW}💡 提示:${NC}"
    echo "  - 建议使用 github.com/username/projectname 格式"
    echo "  - 模块名称应该是小写，用连字符分隔单词"
    echo "  - 不要包含空格或特殊字符"
    echo ""
    exit 1
fi

NEW_MODULE=$1
OLD_MODULE="gin_demo"

# 提取项目名称（用于文件名等）
PROJECT_NAME=$(basename "$NEW_MODULE")

echo -e "${YELLOW}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${YELLOW}📦 旧模块名: ${OLD_MODULE}${NC}"
echo -e "${GREEN}📦 新模块名: ${NEW_MODULE}${NC}"
echo -e "${GREEN}📝 项目名称: ${PROJECT_NAME}${NC}"
echo -e "${YELLOW}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

# 确认操作
echo -e "${YELLOW}⚠️  此操作将修改以下内容:${NC}"
echo "  • go.mod 模块名"
echo "  • 所有 Go 文件的 import 路径"
echo "  • Makefile 配置"
echo "  • README.md 和文档"
echo ""
read -p "确认要执行重命名吗？(y/n) " -n 1 -r
echo ""
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo -e "${RED}❌ 操作已取消${NC}"
    exit 1
fi

echo ""
echo -e "${GREEN}🚀 开始重命名...${NC}"
echo ""

# 1. 更新 go.mod
echo -e "${BLUE}[1/8]${NC} 📝 更新 go.mod..."
if [ -f "go.mod" ]; then
    sed -i.bak "s|^module ${OLD_MODULE}|module ${NEW_MODULE}|g" go.mod
    rm -f go.mod.bak
    echo -e "${GREEN}      ✓ go.mod 更新完成${NC}"
else
    echo -e "${RED}      ✗ go.mod 文件不存在${NC}"
    exit 1
fi

# 2. 更新所有 .go 文件中的导入路径
echo -e "${BLUE}[2/8]${NC} 📝 更新 Go 文件导入路径..."
GO_FILES=$(find . -type f -name "*.go" \
    ! -path "./vendor/*" \
    ! -path "./.git/*" \
    ! -path "./bin/*")

COUNT=0
for file in $GO_FILES; do
    if grep -q "${OLD_MODULE}" "$file"; then
        sed -i.bak "s|${OLD_MODULE}|${NEW_MODULE}|g" "$file"
        rm -f "${file}.bak"
        COUNT=$((COUNT + 1))
    fi
done
echo -e "${GREEN}      ✓ 更新了 ${COUNT} 个 Go 文件${NC}"

# 3. 更新 Makefile
echo -e "${BLUE}[3/8]${NC} 📝 更新 Makefile..."
if [ -f "Makefile" ]; then
    # 更新模块名
    sed -i.bak "s|${OLD_MODULE}|${NEW_MODULE}|g" Makefile
    # 更新项目名称（如 BINARY_NAME）
    sed -i.bak "s|BINARY_NAME=gin-demo|BINARY_NAME=${PROJECT_NAME}|g" Makefile
    rm -f Makefile.bak
    echo -e "${GREEN}      ✓ Makefile 更新完成${NC}"
fi

# 4. 更新 README.md
echo -e "${BLUE}[4/8]${NC} 📝 更新 README.md..."
if [ -f "README.md" ]; then
    sed -i.bak "s|${OLD_MODULE}|${NEW_MODULE}|g" README.md
    # 更新项目标题（可选）
    sed -i.bak "s|Gin Demo|${PROJECT_NAME}|g" README.md
    rm -f README.md.bak
    echo -e "${GREEN}      ✓ README.md 更新完成${NC}"
fi

# 5. 更新 TEMPLATE_USAGE.md
echo -e "${BLUE}[5/8]${NC} 📝 更新 TEMPLATE_USAGE.md..."
if [ -f "TEMPLATE_USAGE.md" ]; then
    sed -i.bak "s|${OLD_MODULE}|${NEW_MODULE}|g" TEMPLATE_USAGE.md
    rm -f TEMPLATE_USAGE.md.bak
    echo -e "${GREEN}      ✓ TEMPLATE_USAGE.md 更新完成${NC}"
fi

# 6. 更新 docs 目录下的文档
echo -e "${BLUE}[6/8]${NC} 📝 更新文档..."
if [ -d "docs" ]; then
    DOC_COUNT=0
    DOC_FILES=$(find docs -type f -name "*.md")
    for file in $DOC_FILES; do
        if grep -q "${OLD_MODULE}" "$file"; then
            sed -i.bak "s|${OLD_MODULE}|${NEW_MODULE}|g" "$file"
            rm -f "${file}.bak"
            DOC_COUNT=$((DOC_COUNT + 1))
        fi
    done
    echo -e "${GREEN}      ✓ 更新了 ${DOC_COUNT} 个文档文件${NC}"
fi

# 7. 清理并重新整理依赖
echo -e "${BLUE}[7/8]${NC} 🔧 整理依赖..."
if go mod tidy > /dev/null 2>&1; then
    echo -e "${GREEN}      ✓ 依赖整理完成${NC}"
else
    echo -e "${YELLOW}      ⚠ 依赖整理出现警告${NC}"
fi

# 8. 验证编译
echo -e "${BLUE}[8/8]${NC} ✅ 验证编译..."
if go build -o /dev/null ./... 2>/dev/null; then
    echo -e "${GREEN}      ✓ 编译成功！${NC}"
else
    echo -e "${YELLOW}      ⚠ 编译出现警告（可能需要运行 make generate）${NC}"
fi

# 完成
echo ""
echo -e "${GREEN}╔════════════════════════════════════════════╗"
echo -e "║          ✅ 重命名完成！                  ║"
echo -e "╚════════════════════════════════════════════╝${NC}"
echo ""
echo -e "${BLUE}📋 已完成的操作:${NC}"
echo "  ✅ 更新 go.mod 模块名"
echo "  ✅ 更新 Go 文件导入路径"
echo "  ✅ 更新 Makefile"
echo "  ✅ 更新文档"
echo "  ✅ 整理依赖"
echo "  ✅ 验证编译"
echo ""
echo -e "${YELLOW}🎯 下一步操作:${NC}"
echo ""
echo "  1️⃣  重新生成代码"
echo "      ${GREEN}make generate${NC}"
echo ""
echo "  2️⃣  运行测试"
echo "      ${GREEN}make test${NC}"
echo ""
echo "  3️⃣  启动项目"
echo "      ${GREEN}make init && make run${NC}"
echo ""
echo "  4️⃣  提交更改"
echo "      ${GREEN}git add .${NC}"
echo "      ${GREEN}git commit -m \"chore: rename module to ${NEW_MODULE}\"${NC}"
echo ""
echo -e "${YELLOW}💡 提示:${NC}"
echo "  • 项目名称已改为: ${GREEN}${PROJECT_NAME}${NC}"
echo "  • 模块路径已改为: ${GREEN}${NEW_MODULE}${NC}"
echo "  • 建议查看并修改 README.md 中的项目描述"
echo ""
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""
