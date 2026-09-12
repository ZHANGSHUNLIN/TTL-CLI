#!/usr/bin/env bash
# CLI machine-contract black-box tests against the built ttl entry point.

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

TEST_DIR=$(mktemp -d "${TMPDIR:-/tmp}/ttl-composability.XXXXXX")
TEST_HOME="$TEST_DIR/home"
TEST_CONF="$TEST_DIR/test.conf"

cleanup() {
    rm -rf "$TEST_DIR"
}
trap cleanup EXIT

mkdir -p "$TEST_HOME"

if [[ -z "$BINARY" ]]; then
    BINARY="$TEST_DIR/ttl"
    go build -o "$BINARY" ./cmd/ttl
elif [[ ! -x "$BINARY" ]]; then
    echo "CLI binary does not exist or is not executable: $BINARY" >&2
    exit 1
fi

cat > "$TEST_CONF" << EOF
db_path = $TEST_DIR/data.bbolt
storage_type = bbolt
EOF

run_cli() {
    HOME="$TEST_HOME" LC_ALL=C LANG=C "$BINARY" --conf "$TEST_CONF" "$@"
}

json_assert() {
    local expression="$1"
    python3 -c "import json,sys; data=json.load(sys.stdin); assert $expression, data"
}

assert_error() {
    local expected_exit="$1"
    local expected_code="$2"
    shift 2
    local stdout_file="$TEST_DIR/stdout"
    local stderr_file="$TEST_DIR/stderr"
    local status
    set +e
    run_cli "$@" >"$stdout_file" 2>"$stderr_file"
    status=$?
    set -e
    [[ "$status" -eq "$expected_exit" ]]
    [[ ! -s "$stdout_file" ]]
    if [[ "$expected_code" == "text" ]]; then
        grep -Fq "multiple resources matched" "$stderr_file"
    else
        json_assert "data['schema_version'] == 1 and data['ok'] is False and data['error']['code'] == '$expected_code'" <"$stderr_file"
        python3 -c "import json,sys; decoder=json.JSONDecoder(); text=sys.stdin.read(); _, end=decoder.raw_decode(text); assert not text[end:].strip(), text" <"$stderr_file"
    fi
}

assert_json_success() {
    local stdout_file="$TEST_DIR/stdout"
    local stderr_file="$TEST_DIR/stderr"
    run_cli "$@" >"$stdout_file" 2>"$stderr_file"
    [[ ! -s "$stderr_file" ]]
    python3 -c "import json,sys; decoder=json.JSONDecoder(); text=sys.stdin.read(); data,end=decoder.raw_decode(text); assert data['schema_version'] == 1 and data['ok'] is True; assert not text[end:].strip(), text" <"$stdout_file"
}

run_cli add note value --tag work --json |
    json_assert "data['ok'] is True and data['data']['resource']['key'] == 'note' and data['data']['resource']['tags'] == ['work']"

run_cli get note --json |
    json_assert "data['data']['resource']['value'] == 'value'"
assert_json_success --debug get note --json

printf 'line 1\nline 2\n' | run_cli update note - --json |
    json_assert "data['data']['resource']['value'] == 'line 1\\nline 2\\n' and data['data']['resource']['tags'] == ['work']"

run_cli tag note ci --json |
    json_assert "data['data']['resource']['tags'] == ['work', 'ci']"

run_cli dtag note work --json |
    json_assert "data['data']['resource']['tags'] == ['ci']"

run_cli get --json |
    json_assert "len(data['data']['resources']) == 1 and data['data']['resources'][0]['key'] == 'note'"

assert_error 3 not_found get missing --json
assert_error 4 conflict add note duplicate --json

run_cli add note-two second --tag ci --json >/dev/null
assert_error 4 ambiguous get ci --json
assert_error 4 text get ci --non-interactive
assert_error 2 invalid_argument version --json
assert_error 2 invalid_argument get --unknown-flag --json
assert_error 2 invalid_argument get --json=maybe

run_cli add empty - --json </dev/null |
    json_assert "data['data']['resource']['value'] == ''"
run_cli get empty --json |
    json_assert "data['data']['resource']['value'] == ''"

run_cli del note --json |
    json_assert "data['data'] == {'key': 'note', 'deleted': True}"
assert_error 3 not_found del note --json

echo "CLI composability contract passed"
