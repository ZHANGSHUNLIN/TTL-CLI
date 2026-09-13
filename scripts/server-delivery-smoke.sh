#!/usr/bin/env bash
# 验证服务端二进制可以脱离源码和客户端制品独立运行。

set -euo pipefail

SERVER_BINARY=${1:?用法：server-delivery-smoke.sh <ttl-server-binary>}
if [[ ! -x "$SERVER_BINARY" ]]; then
    echo "服务端二进制不可执行：$SERVER_BINARY" >&2
    exit 1
fi

WORK_DIR=$(mktemp -d "${TMPDIR:-/tmp}/ttl-server-smoke.XXXXXX")
DATA_DIR="$WORK_DIR/data"
PORT=$(python3 -c 'import socket; s=socket.socket(); s.bind(("127.0.0.1", 0)); print(s.getsockname()[1]); s.close()')
SERVER_PID=""

cleanup() {
    if [[ -n "$SERVER_PID" ]] && kill -0 "$SERVER_PID" 2>/dev/null; then
        kill -TERM "$SERVER_PID" 2>/dev/null || true
        wait "$SERVER_PID" 2>/dev/null || true
    fi
    rm -rf "$WORK_DIR"
}
trap cleanup EXIT

USER_OUTPUT=$("$SERVER_BINARY" --data-dir "$DATA_DIR" user add --id smoke --name Smoke)
API_KEY=$(printf '%s\n' "$USER_OUTPUT" | awk '/API Key:/ {print $3; exit}')
if [[ -z "$API_KEY" ]]; then
    echo "无法从 user add 输出读取 API Key" >&2
    exit 1
fi

"$SERVER_BINARY" --listen 127.0.0.1 --port "$PORT" --data-dir "$DATA_DIR" serve >"$WORK_DIR/server.log" 2>&1 &
SERVER_PID=$!

ready=0
for _ in $(seq 1 100); do
    if curl -fsS "http://127.0.0.1:${PORT}/healthz" | grep -q '"status":"ok"'; then
        ready=1
        break
    fi
    sleep 0.05
done
if [[ "$ready" -ne 1 ]]; then
    cat "$WORK_DIR/server.log" >&2
    echo "服务端未在限定时间内就绪" >&2
    exit 1
fi

curl -fsS -H "Authorization: Bearer ${API_KEY}" "http://127.0.0.1:${PORT}/api/v1/resources" >/dev/null
kill -TERM "$SERVER_PID"
wait "$SERVER_PID"
SERVER_PID=""

test -s "$DATA_DIR/users.json"
echo "服务端独立交付冒烟通过"
