#!/usr/bin/env bash
# 完整本地验证：构建、架构、CLI 黑盒、Go 测试、race 和静态检查。

set -euo pipefail

SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
ROOT_DIR=$(cd "$SCRIPT_DIR/.." && pwd)
cd "$ROOT_DIR"

VERIFY_DIR=$(mktemp -d "${TMPDIR:-/tmp}/ttl-verify.XXXXXX")
VERIFY_HOME="$VERIFY_DIR/home"
VERIFY_GOPATH=$(go env GOPATH)
BINARY="$VERIFY_DIR/ttl"

cleanup() {
    rm -rf "$VERIFY_DIR"
}
trap cleanup EXIT

mkdir -p "$VERIFY_HOME"

echo "========================================="
echo "  TTL 项目完整验证"
echo "========================================="

echo ""
echo "[1/8] 源码与脚本格式检查..."
git diff --check
UNFORMATTED=$(gofmt -s -l .)
if [[ -n "$UNFORMATTED" ]]; then
    echo "以下 Go 文件需要格式化:" >&2
    echo "$UNFORMATTED" >&2
    exit 1
fi
bash -n "$SCRIPT_DIR/regression.sh" "$SCRIPT_DIR/verify.sh" "$SCRIPT_DIR/cli-composability.sh"
echo "✅ 源码与脚本格式检查通过"

echo ""
echo "[2/8] 客户端构建检查..."
go build -o "$BINARY" ./cmd/ttl
echo "✅ 构建成功"

echo ""
echo "[3/8] 架构依赖检查..."
HOME="$VERIFY_HOME" GOPATH="$VERIFY_GOPATH" go test ./internal/architecture
echo "✅ 架构依赖检查通过"

echo ""
echo "[4/8] CLI 黑盒回归..."
"$SCRIPT_DIR/regression.sh" "$BINARY"

echo ""
echo "[5/8] 全部 Go 测试..."
HOME="$VERIFY_HOME" GOPATH="$VERIFY_GOPATH" go test ./...
echo "✅ 全部 Go 测试通过"

echo ""
echo "[6/8] 集成测试..."
HOME="$VERIFY_HOME" GOPATH="$VERIFY_GOPATH" go test ./integration_test/...
echo "✅ 集成测试通过"

echo ""
echo "[7/8] Race 检查..."
HOME="$VERIFY_HOME" GOPATH="$VERIFY_GOPATH" go test -race ./...
echo "✅ Race 检查通过"

echo ""
echo "[8/8] 静态检查..."
HOME="$VERIFY_HOME" GOPATH="$VERIFY_GOPATH" go vet ./...
echo "✅ 静态检查通过"

echo ""
echo "========================================="
echo "✅ 所有验证通过！"
echo "========================================="
