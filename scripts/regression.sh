#!/usr/bin/env bash
# CLI 黑盒回归测试：只验证构建后的用户入口和可观察结果。

set -euo pipefail

SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
ROOT_DIR=$(cd "$SCRIPT_DIR/.." && pwd)
INPUT_BINARY="${1:-}"

if [[ -n "$INPUT_BINARY" ]]; then
    case "$INPUT_BINARY" in
        /*) BINARY="$INPUT_BINARY" ;;
        *) BINARY="$(pwd)/$INPUT_BINARY" ;;
    esac
else
    BINARY=""
fi

cd "$ROOT_DIR"

TEST_DIR=$(mktemp -d "${TMPDIR:-/tmp}/ttl-regression.XXXXXX")
TEST_HOME="$TEST_DIR/home"
TEST_CONF="$TEST_DIR/test.conf"

cleanup() {
    rm -rf "$TEST_DIR"
}
trap cleanup EXIT

mkdir -p "$TEST_HOME"

if [[ -z "$BINARY" ]]; then
    BINARY="$TEST_DIR/ttl"
    go build -o "$BINARY" .
elif [[ ! -x "$BINARY" ]]; then
    echo "回归测试二进制不存在或不可执行: $BINARY" >&2
    exit 1
fi

cat > "$TEST_CONF" << EOF
db_path = $TEST_DIR/data.bbolt
storage_type = bbolt
EOF

run_cli() {
    HOME="$TEST_HOME" LC_ALL=zh_CN.UTF-8 LANG=zh_CN.UTF-8 "$BINARY" --conf "$TEST_CONF" "$@"
}

assert_cli_contains() {
    local expected="$1"
    shift
    local output
    output=$(run_cli "$@")
    grep -Fq -- "$expected" <<< "$output"
}

assert_cli_matches() {
    local pattern="$1"
    shift
    local output
    output=$(run_cli "$@")
    grep -Eq -- "$pattern" <<< "$output"
}

assert_file_exists() {
    if [[ ! -f "$1" ]]; then
        echo "预期文件不存在: $1" >&2
        exit 1
    fi
}

assert_file_not_exists() {
    if [[ -e "$1" ]]; then
        echo "预期文件仍然存在: $1" >&2
        exit 1
    fi
}

echo "========================================="
echo "  TTL CLI 黑盒回归测试"
echo "========================================="
echo "测试目录: $TEST_DIR"

echo ""
echo "[1/1] 用户入口与数据行为..."

echo "   - add/get"
run_cli add "test-resource" "https://example.com" > /dev/null
assert_file_exists "$TEST_DIR/data.bbolt"
assert_cli_contains "example.com" get test-resource

echo "   - tag/export/dtag"
run_cli tag test-resource ci automated > /dev/null
run_cli export -t resources > "$TEST_DIR/export.csv"
grep -q "key,value,tags" "$TEST_DIR/export.csv"
grep -q "test-resource" "$TEST_DIR/export.csv"
run_cli dtag test-resource automated > /dev/null

echo "   - update/rename/delete"
run_cli update test-resource "https://updated.example.com" > /dev/null
run_cli rename test-resource "renamed-resource" > /dev/null
assert_cli_contains "updated.example.com" get renamed-resource
run_cli del renamed-resource > /dev/null
if run_cli get renamed-resource > /dev/null 2>&1; then
    echo "删除后的资源仍然可以读取" >&2
    exit 1
fi

echo "   - import/log/history/audit"
run_cli add "test-resource-2" "https://example2.com" > /dev/null
run_cli import "$TEST_DIR/export.csv" > /dev/null
run_cli log "CLI 测试日志" > /dev/null
assert_cli_contains "CLI 测试日志" log -l
assert_cli_contains "test-resource" history 5
assert_cli_contains "总操作次数" audit

echo "   - version/config/tags"
assert_cli_matches "[0-9]+\.[0-9]+" version
assert_cli_contains "$TEST_DIR/data.bbolt" config
assert_cli_contains "$TEST_CONF" config
assert_cli_contains "ci" tags
run_cli add "tag-test-1" "value1" -t work > /dev/null
run_cli add "tag-test-2" "value2" -t work > /dev/null
assert_cli_contains "tag-test-1" tags work

echo "   - legacy server command compatibility"
assert_cli_contains "user" server --help

echo "   - encrypt/key verify/decrypt"
run_cli encrypt --migrate > /dev/null
assert_file_exists "$TEST_HOME/.ttl/.key"
assert_cli_contains "example2.com" get test-resource-2
run_cli key verify > /dev/null
printf 'y\n' | run_cli decrypt > /dev/null
assert_cli_contains "example2.com" get test-resource-2

echo "   - workspace list/create/switch/isolation/delete"
assert_cli_contains "暂无工作空间" workspace list
run_cli workspace create work > /dev/null
run_cli workspace create life > /dev/null
assert_cli_contains "work" workspace list
assert_cli_contains "life" workspace list
run_cli workspace switch work > /dev/null
assert_cli_contains "work" workspace current
run_cli add "work-only-resource" "work-value" > /dev/null
assert_file_exists "$TEST_DIR/workspaces/work.bbolt"
run_cli workspace switch life > /dev/null
if run_cli get "work-only-resource" > /dev/null 2>&1; then
    echo "工作空间之间读取到了不应存在的资源" >&2
    exit 1
fi
run_cli add "life-only-resource" "life-value" > /dev/null
assert_file_exists "$TEST_DIR/workspaces/life.bbolt"
assert_cli_contains "life-value" get "life-only-resource"
assert_cli_contains "Database" workspace show work
assert_cli_contains "Resources:" workspace show work
run_cli ws work > /dev/null
assert_cli_contains "work-value" get "work-only-resource"
run_cli workspace switch work > /dev/null
run_cli workspace delete life > /dev/null
assert_file_not_exists "$TEST_DIR/workspaces/life.bbolt"
assert_cli_contains "work" workspace list
if assert_cli_contains "life" workspace list; then
    echo "删除工作空间后仍然能看到 life" >&2
    exit 1
fi

echo "✅ CLI 黑盒回归通过"
