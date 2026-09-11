#!/usr/bin/env bash
# 完整本地验证：构建产物、CLI 黑盒回归、单元测试、集成测试和静态检查。

set -euo pipefail

SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
ROOT_DIR=$(cd "$SCRIPT_DIR/.." && pwd)
cd "$ROOT_DIR"

VERIFY_DIR=$(mktemp -d "${TMPDIR:-/tmp}/ttl-verify.XXXXXX")
VERIFY_HOME="$VERIFY_DIR/home"
VERIFY_GOPATH=$(go env GOPATH)
BINARY="$VERIFY_DIR/ttl"
SERVER_BINARY="$VERIFY_DIR/ttl-server"

cleanup() {
    rm -rf "$VERIFY_DIR"
}
trap cleanup EXIT

mkdir -p "$VERIFY_HOME"

echo "========================================="
echo "  TTL 项目完整验证"
echo "========================================="

echo ""
echo "[1/5] 客户端与服务端构建检查..."
go build -o "$BINARY" ./cmd/ttl
go build -o "$SERVER_BINARY" ./cmd/ttl-server
echo "✅ 构建成功"

echo ""
echo "[2/5] CLI 黑盒回归..."
"$SCRIPT_DIR/regression.sh" "$BINARY"

echo ""
echo "[3/5] 单元测试..."
HOME="$VERIFY_HOME" GOPATH="$VERIFY_GOPATH" go test ./...
echo "✅ 单元测试通过"

echo ""
echo "[4/5] 集成测试..."
HOME="$VERIFY_HOME" GOPATH="$VERIFY_GOPATH" go test ./integration_test/...
echo "✅ 集成测试通过"

echo ""
echo "[5/5] 静态检查..."
HOME="$VERIFY_HOME" GOPATH="$VERIFY_GOPATH" go vet ./...
echo "✅ 静态检查通过"

echo ""
echo "========================================="
echo "✅ 所有验证通过！"
echo "========================================="
