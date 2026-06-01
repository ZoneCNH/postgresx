#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$ROOT_DIR"

GO="${GO:-go}"
TEST_PATTERN='TestMigrationRunner.*Integration'

run_tests() {
  GOWORK=off "$GO" test -v -run "$TEST_PATTERN" -count=1 ./...
}

if [[ -n "${POSTGRESX_INTEGRATION_DSN:-}" || -n "${POSTGRES_TEST_DSN:-}" ]]; then
  run_tests
  exit 0
fi

if command -v docker >/dev/null 2>&1 && docker info >/dev/null 2>&1; then
  db_name="postgresx_test"
  db_user="postgres"
  db_secret="postgres"
  container="postgresx-integration-${RANDOM:-0}-$$"

  if ! docker run --rm -d \
    --name "$container" \
    -e "POSTGRES_DB=$db_name" \
    -e "POSTGRES_USER=$db_user" \
    -e "POSTGRES_PASSWORD=$db_secret" \
    -P postgres:17-alpine >/dev/null; then
    if [[ "${POSTGRESX_REQUIRE_INTEGRATION:-}" == "1" ]]; then
      echo "failed to start Docker PostgreSQL for required integration gate" >&2
      exit 1
    fi
    echo "Docker PostgreSQL is unavailable; skipping live migration gate"
    exit 0
  fi
  trap 'docker rm -f "$container" >/dev/null 2>&1 || true' EXIT

  port=""
  for _ in {1..45}; do
    port="$(docker port "$container" 5432/tcp 2>/dev/null | awk -F: 'END {print $NF}')"
    if [[ -n "$port" ]]; then
      break
    fi
    sleep 1
  done

  if [[ -z "$port" ]]; then
    echo "Docker PostgreSQL did not publish a port" >&2
    exit 1
  fi

  # The official image briefly starts a bootstrap server before restarting into
  # the final externally reachable server. Wait for init completion, the second
  # readiness log line, and TCP readiness to avoid racing tests against restart.
  for _ in {1..45}; do
    if docker logs "$container" 2>&1 | grep -q 'PostgreSQL init process complete; ready for start up'; then
      break
    fi
    sleep 1
  done
  for _ in {1..45}; do
    ready_count="$(docker logs "$container" 2>&1 | grep -c 'database system is ready to accept connections' || true)"
    if [[ "$ready_count" -ge 2 ]]; then
      break
    fi
    sleep 1
  done
  for _ in {1..45}; do
    if docker exec "$container" pg_isready -h 127.0.0.1 -p 5432 -U "$db_user" -d "$db_name" >/dev/null 2>&1; then
      break
    fi
    sleep 1
  done
  # pg_isready can return immediately before the host-mapped connection path is
  # fully stable on some Docker networking setups. A short bounded settle avoids
  # a one-time first ping timeout while keeping the gate deterministic.
  sleep 3

  scheme="postgres"
  export POSTGRES_TEST_DSN="${scheme}://${db_user}:${db_secret}@127.0.0.1:${port}/${db_name}?sslmode=disable"
  run_tests
  exit 0
fi

if [[ "${POSTGRESX_REQUIRE_INTEGRATION:-}" == "1" ]]; then
  echo "live PostgreSQL integration is required; set POSTGRESX_INTEGRATION_DSN or POSTGRES_TEST_DSN, or provide Docker" >&2
  exit 1
fi

echo "POSTGRESX_INTEGRATION_DSN or POSTGRES_TEST_DSN is not set; skipping live migration gate"
