#!/usr/bin/env bash
# 验证版本目录、current 指针和失败升级回滚流程。

set -euo pipefail

SERVER_BINARY=${1:?用法：server-upgrade-rollback-smoke.sh <ttl-server-binary>}
if [[ ! -x "$SERVER_BINARY" ]]; then
    echo "服务端二进制不可执行：$SERVER_BINARY" >&2
    exit 1
fi

WORK_DIR=$(mktemp -d "${TMPDIR:-/tmp}/ttl-server-upgrade.XXXXXX")
DATA_DIR="$WORK_DIR/data"
RELEASES_DIR="$WORK_DIR/releases"
CURRENT="$WORK_DIR/current"
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

mkdir -p "$RELEASES_DIR/v1" "$RELEASES_DIR/v2"
cp "$SERVER_BINARY" "$RELEASES_DIR/v1/ttl-server"
cp "$SERVER_BINARY" "$RELEASES_DIR/v2/ttl-server"
chmod 755 "$RELEASES_DIR/v1/ttl-server" "$RELEASES_DIR/v2/ttl-server"
ln -s "$RELEASES_DIR/v1" "$CURRENT"

"$CURRENT/ttl-server" --data-dir "$DATA_DIR" user add --id rollback --name Rollback >/dev/null
"$CURRENT/ttl-server" --listen 127.0.0.1 --port "$PORT" --data-dir "$DATA_DIR" serve >"$WORK_DIR/server.log" 2>&1 &
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
    exit 1
fi

cp -R "$DATA_DIR" "$WORK_DIR/data-backup"
kill -TERM "$SERVER_PID"
wait "$SERVER_PID"
SERVER_PID=""

cat > "$RELEASES_DIR/v2/ttl-server" <<'EOF'
#!/usr/bin/env sh
exit 1
EOF
chmod 755 "$RELEASES_DIR/v2/ttl-server"
ln -sfn "$RELEASES_DIR/v2" "$CURRENT"
if "$CURRENT/ttl-server" serve --data-dir "$DATA_DIR" >/dev/null 2>&1; then
    echo "损坏的 v2 release 不应启动成功" >&2
    exit 1
fi

ln -sfn "$RELEASES_DIR/v1" "$CURRENT"
"$CURRENT/ttl-server" --listen 127.0.0.1 --port "$PORT" --data-dir "$DATA_DIR" serve >"$WORK_DIR/server.log" 2>&1 &
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
    exit 1
fi

kill -TERM "$SERVER_PID"
wait "$SERVER_PID"
SERVER_PID=""
test -s "$WORK_DIR/data-backup/users.json"
echo "服务端升级回滚冒烟通过"
