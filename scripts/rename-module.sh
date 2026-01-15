#!/bin/bash

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 脚本说明
echo -e "${GREEN}================================"
echo "  Go 项目模块名称重命名工具"
echo -e "================================${NC}\n"

# 检查参数
if [ $# -eq 0 ]; then
    echo -e "${RED}错误: 请提供新的模块名称${NC}"
    echo ""
    echo "用法:"
    echo "  $0 <new-module-name>"
    echo ""
    echo "示例:"
    echo "  $0 github.com/yourname/yourproject"
    echo "  $0 mycompany.com/internal/api-service"
    echo ""
    exit 1
fi

NEW_MODULE=$1
OLD_MODULE="gin_demo"

echo -e "${YELLOW}旧模块名: ${OLD_MODULE}${NC}"
echo -e "${YELLOW}新模块名: ${NEW_MODULE}${NC}"
echo ""

# 确认操作
read -p "确认要执行重命名吗？(y/n) " -n 1 -r
echo ""
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo -e "${RED}操作已取消${NC}"
    exit 1
fi

echo ""
echo -e "${GREEN}开始重命名...${NC}"
echo ""

# 1. 更新 go.mod
echo -e "📝 更新 go.mod..."
if [ -f "go.mod" ]; then
    sed -i.bak "s|^module ${OLD_MODULE}|module ${NEW_MODULE}|g" go.mod
    rm -f go.mod.bak
    echo -e "${GREEN}✓ go.mod 更新完成${NC}"
else
    echo -e "${RED}✗ go.mod 文件不存在${NC}"
fi

# 2. 更新所有 .go 文件中的导入路径
echo -e "📝 更新 Go 文件导入路径..."
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
echo -e "${GREEN}✓ 更新了 ${COUNT} 个 Go 文件${NC}"

# 3. 更新 Makefile（如果存在）
echo -e "📝 更新 Makefile..."
if [ -f "Makefile" ]; then
    sed -i.bak "s|${OLD_MODULE}|${NEW_MODULE}|g" Makefile
    rm -f Makefile.bak
    echo -e "${GREEN}✓ Makefile 更新完成${NC}"
fi

# 4. 更新 README.md
echo -e "📝 更新 README.md..."
if [ -f "README.md" ]; then
    sed -i.bak "s|${OLD_MODULE}|${NEW_MODULE}|g" README.md
    rm -f README.md.bak
    echo -e "${GREEN}✓ README.md 更新完成${NC}"
fi

# 5. 更新 docs 目录下的文档
echo -e "📝 更新文档..."
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
    echo -e "${GREEN}✓ 更新了 ${DOC_COUNT} 个文档文件${NC}"
fi

# 6. 清理并重新整理依赖
echo ""
echo -e "${GREEN}🔧 整理依赖...${NC}"
go mod tidy

# 7. 验证
echo ""
echo -e "${GREEN}✅ 验证结果...${NC}"
if go build -o /dev/null ./... 2>/dev/null; then
    echo -e "${GREEN}✓ 编译成功！${NC}"
else
    echo -e "${YELLOW}⚠ 编译出现警告，请检查${NC}"
fi

# 完成
echo ""
echo -e "${GREEN}================================"
echo "  ✅ 重命名完成！"
echo -e "================================${NC}"
echo ""
echo "下一步操作:"
echo "  1. 检查 go.mod 文件"
echo "  2. 运行 'make generate' 重新生成代码"
echo "  3. 运行 'make test' 确保测试通过"
echo "  4. 提交更改: git add . && git commit -m 'Rename module to ${NEW_MODULE}'"
echo ""
echo -e "${YELLOW}提示: 如果使用了 Git，建议先创建一个新分支${NC}"
echo ""
