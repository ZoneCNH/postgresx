# postgresx L2 基础设施适配层标准工厂 Goal 可执行方案

> 生成日期：2026-06-04 Asia/Tokyo  
> Goal ID：`GOAL-20260604-POSTGRESX-L2-FACTORY-001`  
> 目标仓库：`github.com/ZoneCNH/postgresx`  
> 标准源：`github.com/ZoneCNH/xlib-standard`  
> L0 依赖：`github.com/ZoneCNH/kernel`  
> L1 依赖：`configx / observex / testkitx / resiliencx / schedulex`  
> 目标层级：L2 Infrastructure Adapter / PostgreSQL 基础设施适配库  
> 执行协议：Goal Runtime Prompt v3.1  
> 默认模式：Full Mode  
> 目标版本：`v0.1.0` MVA → `v0.2.0` Contract Hardening → `v1.0.0` Stable API  
> 完成声明格式：只能使用 `DONE with evidence:`  
> 核心约束：禁止 main 直接开发；必须使用 git worktree；不得依赖 `x.go`；不得包含业务 repository / 业务表结构 / 应用事务编排；不得泄露 `/home/k8s/secrets/env/*` 真实内容、DSN、密码、token、SQL 参数或生产日志。

---

## 0. 执行结论

`postgresx` 不应该只是对 `pgx` 或 `database/sql` 的薄包装。

`postgresx` 的正确定位是：

```text
PostgreSQL infrastructure adapter contract layer
```

它的真实价值不是“帮业务执行 SQL”，而是把 PostgreSQL 这种外部基础设施的不确定性统一收敛为：

```text
显式配置
+ 可关闭连接生命周期
+ 连接池治理
+ context-first 查询与事务
+ PostgreSQL 错误分类
+ DSN / SQL 参数脱敏
+ health contract
+ metrics / logs / traces contract
+ fake / mockable testkit
+ opt-in integration evidence
+ release manifest
+ downstream adoption proof
+ self-improving patch
```

最终推荐路径：

```text
xlib-standard
  = Standard Source / Generator / Harness / Evidence Runtime
        ↓ governs / generates / audits
kernel
  = L0 primitive: error / lifecycle / clock / context / shutdown / validation / health
        ↓ allowed dependency
configx / observex / testkitx / resiliencx / schedulex
  = L1 cross-cutting capabilities
        ↓ allowed dependency
postgresx
  = L2 PostgreSQL infrastructure adapter
        ↓ consumed by
x.go / market-data / macro-data / engines / services
```

第一阶段不要追求 ORM、migration engine、业务 repository 或分库分表。`v0.1.0` 的正确 MVA 是：

```text
module identity fixed
+ xlib-standard adoption artifacts
+ config schema
+ pgxpool client lifecycle
+ Ping / Health / Close
+ Exec / Query / QueryRow minimal contract
+ WithinTx transaction helper
+ pg error classification
+ pool/query/tx metrics
+ DSN and SQL params redaction
+ fake Queryer/Execer/Tx
+ integration opt-in with PostgreSQL
+ release manifest + checksum
+ downstream smoke adoption proof
+ retrospective patches
```

---

## 1. 当前事实基线

### 1.1 当前 `postgresx` 仓库事实

当前 `ZoneCNH/postgresx` 已存在为独立公开仓库，默认分支为 `main`。仓库当前体量极小，README 内容只有最小标题和项目名：

```text
# postgresx
postgresx
```

这意味着本阶段不是“修复一个成熟库”，而是从一个空骨架开始，把它纳入 `xlib-standard` 标准源控制下的 L2 基础设施适配层标准工厂。

### 1.2 当前 `xlib-standard` 事实

`xlib-standard` 已经被定义为基础库标准与交付运行时仓库，职责包括：

```text
Standard Source
Go Reference Template
Generator
Harness
Evidence Runtime
```

`xlib-standard` 还明确登记了 `postgresx` 作为生成库之一，并把 `postgresx` 列为 L2 目标库。

### 1.3 当前下游采纳事实

当前 `xlib-standard` 下游矩阵中，`postgresx` 的状态仍应按：

```text
adoption_status = not_adopted
evidence_state = not_run
```

解释。登记不等于采纳，骨架不等于实现，dry-run 不等于 release usable。没有当前 `postgresx` 仓库内的 gate 输出、manifest、checksum、artifact 和 adoption proof，不允许写：

```text
postgresx adopted
postgresx release ready
postgresx DONE
postgresx usable by x.go
```

### 1.4 当前角色裁决

`postgresx` 的角色固定为：

```text
L2 数据库适配器库
```

必须包含：

```text
PostgreSQL profile
连接配置
连接池 lifecycle
健康检查
错误分类
测试夹具
release evidence
```

禁止包含：

```text
业务 repository
业务表结构
业务 schema migration 默认执行
应用 transaction 编排
x.go 业务模型
x.go 反向依赖
生产密钥读取
```

---

## 2. 问题的底层本质

`postgresx` 的底层问题不是“如何把 pgx 封装得更顺手”，而是：

> 如何把 PostgreSQL 这种有状态、会失败、带认证、带事务、带连接池、带 SQL 注入/日志泄密风险的外部系统，转化成一个可配置、可测试、可观测、可降级、可发布证明、可被下游安全采纳的基础设施契约。

如果没有 `postgresx` 这一层，上层业务会反复发明：

```text
DSN 配置
连接池默认值
statement timeout
事务 begin/commit/rollback
deadlock / serialization failure 重试判断
unique violation 映射
health 输出格式
pool metrics
SQL 参数脱敏
integration test 连接方式
release evidence
```

这会形成 9 类结构债：

| 债务 | 表现 | 后果 |
|---|---|---|
| 配置债 | 每个服务自定义 DSN/env/pool 字段 | 无法统一运维与审计 |
| 生命周期债 | 隐藏全局 db、init() 连接、Close 不幂等 | 泄露连接与 goroutine |
| 事务债 | 应用层散落 begin/commit/rollback | 数据一致性风险 |
| 错误债 | pgx/pgconn 原始错误泄漏 | 重试、告警、熔断无法统一 |
| 观测债 | 每个业务服务 metrics 名称不同 | dashboard/SLO 不可复用 |
| 安全债 | DSN、password、SQL 参数进入日志或 manifest | 密钥泄漏 |
| 测试债 | 单元测试依赖真实 PostgreSQL | CI 不稳定、不可复现 |
| 发布债 | 没有 manifest/checksum/gate 输出 | 下游无法判断是否可采纳 |
| 治理债 | 业务 repository 下沉到基础库 | L2 被业务污染 |

所以 `postgresx` 的本质是：

```text
把 PostgreSQL 的失败语义、连接语义、事务语义和观测语义标准化。
```

---

## 3. 不可再拆解的基本真理

### 3.1 分层真理

```text
TRUTH-POSTGRESX-001  postgresx 是 L2 基础设施适配库，不是业务 repository 层。
TRUTH-POSTGRESX-002  postgresx 可以依赖 L0 kernel 与 L1 configx/observex/testkitx/resiliencx/schedulex。
TRUTH-POSTGRESX-003  postgresx 不得依赖其他 L2，例如 redisx/kafkax/taosx/ossx/clickhousex/natsx。
TRUTH-POSTGRESX-004  postgresx 不得依赖 x.go、market-data、macro-data、regime-engine 或任何业务系统。
TRUTH-POSTGRESX-005  postgresx 不得定义业务表、业务 SQL、业务 schema、业务 repository。
TRUTH-POSTGRESX-006  postgresx 只提供 PostgreSQL 基础设施能力和契约。
TRUTH-POSTGRESX-007  postgresx public API 必须 context-first。
TRUTH-POSTGRESX-008  postgresx 不得隐藏全局 client，不得在 init() 中连接数据库。
TRUTH-POSTGRESX-009  postgresx release 必须有 gate output、manifest、checksum、contract hash。
TRUTH-POSTGRESX-010  postgresx downstream adoption 只能由当前下游命令证据证明。
```

### 3.2 PostgreSQL 失败语义真理

```text
TRUTH-POSTGRESX-011  PostgreSQL 一定会失败：网络断开、认证失败、连接池耗尽、锁等待、deadlock、serialization failure、statement timeout、constraint violation 都必须被建模。
TRUTH-POSTGRESX-012  pgconn.PgError 必须映射到统一 ErrorKind，不得原样泄漏给上层作为唯一判断依据。
TRUTH-POSTGRESX-013  所有操作必须接受 context.Context。
TRUTH-POSTGRESX-014  所有 query / exec / tx 必须有显式超时或依赖 caller context deadline。
TRUTH-POSTGRESX-015  默认不自动重试写事务；事务重试必须显式开启并要求幂等语义。
TRUTH-POSTGRESX-016  Close 必须幂等；Close 后操作必须返回稳定错误。
TRUTH-POSTGRESX-017  连接池状态必须可观测。
```

### 3.3 安全与脱敏真理

```text
TRUTH-POSTGRESX-018  DSN、password、token、TLS key、SQL 参数不得进入 logs/errors/traces/manifest。
TRUTH-POSTGRESX-019  SQL 文本也可能包含业务敏感信息；默认只记录 operation name / query name / hash，不记录 raw SQL。
TRUTH-POSTGRESX-020  config.Sanitize() 必须是 release manifest、health、debug 输出的唯一配置出口。
TRUTH-POSTGRESX-021  integration evidence 只能记录 sanitized endpoint、server version、database hash 或 safe name。
```

### 3.4 测试与发布真理

```text
TRUTH-POSTGRESX-022  单元测试不得依赖真实 PostgreSQL。
TRUTH-POSTGRESX-023  fake Queryer/Execer/Tx 是 P0，不是 P1。
TRUTH-POSTGRESX-024  integration test 必须 opt-in，并输出 pass/skip/fail evidence。
TRUTH-POSTGRESX-025  没有 release-final-check，不允许 tag。
TRUTH-POSTGRESX-026  没有 downstream adoption proof，不允许宣称可被 x.go 采纳。
TRUTH-POSTGRESX-027  没有 retrospective patch，不允许宣称 self-improving 成立。
```

---

## 4. 被误认为真理的常见假设

| 常见假设 | 为什么错 | 正确裁决 |
|---|---|---|
| postgresx 就是 pgx 的薄包装 | 薄包装不能统一配置、错误、观测、测试、发布证据 | postgresx 是 PostgreSQL contract adapter |
| 直接暴露 pgxpool.Pool 最灵活 | 会把所有下游绑定到 pgx 类型和实现细节 | public API 提供稳定 interface；pgx 放 internal/driver 或 explicit escape hatch |
| L2 应该顺手做 repository pattern | repository 是业务边界，不是基础设施边界 | repository 留给 L3/L4/L5/L6 |
| 自动 migration 很方便 | 隐式生产启动动作风险极高 | migration 只能是可选 helper / explicit command，不进入 Open 默认路径 |
| 事务 helper 应该自动重试所有错误 | 写事务自动重试可能重复副作用 | 只对明确幂等、明确错误类、明确 policy 开启 |
| Ping 通过就表示数据库健康 | Ping 不代表 query、pool、权限、statement timeout 全部正常 | health 输出分多个 check |
| 日志记录 SQL 和 args 方便排障 | SQL/args 可能包含敏感业务数据 | 默认记录 query name/hash，不记录 raw params |
| 单元测试连一个本地 PostgreSQL 就够 | 不可复现、CI 慢、会引入环境耦合 | 单元 fake，integration opt-in |
| 配置可以直接读 `/home/k8s/secrets/env/*` | 基础库不应隐式读取生产环境 | 只接受 caller 显式传入 config/secret source |
| release 只要 go test 通过 | 无 manifest/checksum/contracts/evidence 不可审计 | release 必须 Full Gate |

---

## 5. 可以被打破的限制

```text
LIMIT-POSTGRESX-001  不需要 Day 1 做 ORM。
LIMIT-POSTGRESX-002  不需要 Day 1 做 migration engine。
LIMIT-POSTGRESX-003  不需要 Day 1 支持读写分离、分库分表、sharding。
LIMIT-POSTGRESX-004  不需要 Day 1 完整封装 PostgreSQL 所有高级特性。
LIMIT-POSTGRESX-005  不需要暴露 provider SDK 类型才能好用。
LIMIT-POSTGRESX-006  不需要 integration test 默认跑真实数据库。
LIMIT-POSTGRESX-007  不需要把 x.go 的表结构放进 postgresx。
LIMIT-POSTGRESX-008  不需要引入应用框架、ORM、DI 框架。
LIMIT-POSTGRESX-009  不需要一次性 v1 稳定 API；先 v0.1 MVA，保留 ADR 与 API diff gate。
```

---

## 6. 从零设计的新方案

### 6.1 系统结构

```text
[Consumer: x.go / services / engines]
        |
        | imports
        v
[postgresx public package]
        |
        | uses interfaces/options/contracts
        v
[internal/driver/pgx]
        |
        | depends on selected provider SDK
        v
[PostgreSQL]
```

横切依赖：

```text
postgresx
  -> kernel: error / lifecycle / context / health / validation primitives
  -> configx: explicit config loading and redaction contract
  -> observex: metrics/log/trace/health contract
  -> testkitx: test-only fixtures and contract helpers
  -> resiliencx: optional retry/timeout/breaker policy
  -> schedulex: optional scheduled maintenance / background jobs, not P0
```

### 6.2 统一目录结构

```text
.
├── .agent/
│   ├── goal.md
│   ├── spec.md
│   ├── design.md
│   ├── plan.md
│   ├── traceability-matrix.md
│   ├── risk-register.md
│   ├── decision-log.md
│   ├── evidence/
│   ├── reviews/
│   ├── release/
│   └── retrospectives/
├── .github/workflows/
│   ├── ci.yml
│   ├── release-check.yml
│   └── security.yml
├── contracts/
│   ├── api.contract.yaml
│   ├── config.schema.json
│   ├── health.schema.json
│   ├── metrics.contract.yaml
│   ├── errors.contract.yaml
│   └── release-manifest.contract.yaml
├── docs/
│   ├── spec.md
│   ├── design.md
│   ├── api.md
│   ├── config.md
│   ├── health.md
│   ├── observability.md
│   ├── resilience.md
│   ├── testing.md
│   ├── integration.md
│   ├── security.md
│   ├── release.md
│   ├── migration-boundary.md
│   └── research/
│       └── dependency-research.md
├── examples/
│   ├── basic/
│   ├── health/
│   ├── observability/
│   ├── tx/
│   └── integration/
├── internal/
│   ├── driver/pgx/
│   ├── config/
│   ├── errors/
│   ├── health/
│   ├── metrics/
│   └── testutil/
├── pkg/postgresx/
│   ├── client.go
│   ├── config.go
│   ├── errors.go
│   ├── health.go
│   ├── metrics.go
│   ├── options.go
│   ├── query.go
│   ├── tx.go
│   └── doc.go
├── testkit/
│   ├── fake.go
│   ├── recorder.go
│   ├── assertions.go
│   └── fixtures/
├── release/manifest/
│   └── .gitkeep
├── scripts/
│   ├── boundary_check.sh
│   ├── contract_check.sh
│   ├── docs_check.sh
│   ├── integration_check.sh
│   ├── secret_scan.sh
│   ├── generate_manifest.sh
│   └── release_final_check.sh
├── Makefile
├── go.mod
├── go.sum
├── AGENTS.md
├── CONSTITUTION.md
├── README.md
├── CHANGELOG.md
├── LICENSE
└── renovate.json
```

### 6.3 Public API 设计

P0 public API 不追求“漂亮”，只追求稳定、可测试、可观测、可证明。

```go
package postgresx

import "context"

type Client interface {
    Name() string
    Ping(ctx context.Context) error
    Health(ctx context.Context) Health
    Stats() PoolStats
    Exec(ctx context.Context, query Query, args ...any) (CommandResult, error)
    Query(ctx context.Context, query Query, args ...any) (Rows, error)
    QueryRow(ctx context.Context, query Query, args ...any) Row
    WithinTx(ctx context.Context, fn func(context.Context, Tx) error, opts ...TxOption) error
    Close(ctx context.Context) error
}

type Query struct {
    Name string // stable operation/query name
    SQL  string // raw SQL allowed for execution, not logs by default
}

type Tx interface {
    Exec(ctx context.Context, query Query, args ...any) (CommandResult, error)
    Query(ctx context.Context, query Query, args ...any) (Rows, error)
    QueryRow(ctx context.Context, query Query, args ...any) Row
}

type Queryer interface {
    Query(ctx context.Context, query Query, args ...any) (Rows, error)
    QueryRow(ctx context.Context, query Query, args ...any) Row
}

type Execer interface {
    Exec(ctx context.Context, query Query, args ...any) (CommandResult, error)
}
```

原则：

```text
- Query.Name 是日志、metrics、trace 的主标识。
- Query.SQL 只用于执行；默认不进入日志。
- args 默认永不记录。
- Tx helper 不暴露业务 repository。
- provider SDK 类型默认不出现在 public API。
- 如必须提供 escape hatch，必须单独 ADR：ADR-POSTGRESX-PGX-ESCAPE-HATCH。
```

### 6.4 Config Contract

P0 配置字段：

```yaml
name: postgresx
provider: postgres
database:
  dsn: optional explicit DSN
  host: localhost
  port: 5432
  database: app
  username: app
  password: explicit secret value or secret source reference
  sslmode: disable|prefer|require|verify-ca|verify-full
  application_name: postgresx
pool:
  min_conns: 0
  max_conns: 10
  max_conn_lifetime: 1h
  max_conn_idle_time: 30m
  health_check_period: 1m
timeouts:
  connect_timeout: 5s
  query_timeout: 30s
  tx_timeout: 60s
  statement_timeout: 30s
resilience:
  retry:
    enabled: false
    max_attempts: 2
    idempotent_only: true
  circuit_breaker:
    enabled: false
observability:
  metrics_enabled: true
  traces_enabled: true
  log_sql: false
  log_args: false
  query_name_required: true
integration:
  enabled: false
  allow_external_dsn: false
```

硬约束：

```text
- DefaultConfig() 不得指向生产 endpoint。
- Validate() 必须拒绝空 endpoint、非法 pool、非法 timeout、危险 logging 配置。
- Sanitize() 必须屏蔽 password、DSN credentials、TLS key、raw SQL args。
- config.schema.json 必须覆盖所有 P0 字段。
- examples/config 必须通过 schema validation。
```

### 6.5 Error Contract

统一错误维度：

```text
ErrorKind:
  UNKNOWN
  CONFIG_INVALID
  AUTH_FAILED
  NETWORK_UNREACHABLE
  CONNECTION_FAILED
  TIMEOUT
  CONTEXT_CANCELED
  POOL_EXHAUSTED
  QUERY_FAILED
  TRANSACTION_FAILED
  SERIALIZATION_FAILED
  DEADLOCK_DETECTED
  LOCK_NOT_AVAILABLE
  UNIQUE_VIOLATION
  FOREIGN_KEY_VIOLATION
  CHECK_VIOLATION
  NOT_NULL_VIOLATION
  RESOURCE_NOT_FOUND
  PERMISSION_DENIED
  UNSUPPORTED_OPERATION
  SHUTDOWN
```

PostgreSQL SQLSTATE 映射 P0：

| SQLSTATE / Class | PostgreSQL 语义 | ErrorKind | Retry |
|---|---|---|---|
| `23505` | unique_violation | UNIQUE_VIOLATION | no |
| `23503` | foreign_key_violation | FOREIGN_KEY_VIOLATION | no |
| `23514` | check_violation | CHECK_VIOLATION | no |
| `23502` | not_null_violation | NOT_NULL_VIOLATION | no |
| `40001` | serialization_failure | SERIALIZATION_FAILED | conditional |
| `40P01` | deadlock_detected | DEADLOCK_DETECTED | conditional |
| `55P03` | lock_not_available | LOCK_NOT_AVAILABLE | conditional |
| `57014` | query_canceled | TIMEOUT or CONTEXT_CANCELED | conditional |
| `42P01` | undefined_table | RESOURCE_NOT_FOUND | no |
| `42501` | insufficient_privilege | PERMISSION_DENIED | no |
| `28P01` | invalid_password | AUTH_FAILED | no |
| `08xxx` | connection exception class | CONNECTION_FAILED | yes |
| pool acquire timeout | pool exhausted | POOL_EXHAUSTED | conditional |

错误输出必须包含：

```text
module=postgresx
provider=postgres
operation
query_name
error_kind
retryable
sqlstate if safe
redacted_message
```

禁止包含：

```text
password
raw DSN
SQL args
raw SQL by default
TLS key
secret path contents
```

### 6.6 Health Contract

Health 输出必须匹配 `contracts/health.schema.json`：

```json
{
  "name": "postgresx",
  "status": "pass|warn|fail",
  "provider": "postgres",
  "version": "optional",
  "checks": [
    {
      "name": "ping",
      "status": "pass|warn|fail",
      "latency_ms": 3,
      "last_success_at": "2026-06-04T00:00:00Z",
      "error_kind": "TIMEOUT",
      "message": "redacted"
    },
    {
      "name": "pool",
      "status": "pass|warn|fail",
      "latency_ms": 0,
      "message": "acquired=1 idle=4 total=5"
    }
  ],
  "observed_at": "2026-06-04T00:00:00Z"
}
```

P0 checks：

```text
ping
pool_stats
last_error
version optional
```

P1 checks：

```text
readiness_query with query_name
replication role optional
server parameter sanity optional
```

### 6.7 Metrics Contract

公共 metrics：

```text
l2_operation_total{module="postgresx",provider="postgres",operation,result,error_kind}
l2_operation_duration_seconds{module="postgresx",provider="postgres",operation,result}
l2_retry_total{module="postgresx",provider="postgres",operation,error_kind}
l2_health_status{module="postgresx",provider="postgres",status}
```

postgresx 专属 metrics：

```text
postgresx_pool_acquired{pool}
postgresx_pool_idle{pool}
postgresx_pool_total{pool}
postgresx_pool_max{pool}
postgresx_pool_acquire_count_total{pool,result}
postgresx_pool_acquire_duration_seconds{pool,result}
postgresx_query_total{query_name,result,error_kind}
postgresx_query_duration_seconds{query_name,result}
postgresx_tx_total{result,error_kind}
postgresx_tx_duration_seconds{result}
postgresx_tx_rollback_total{reason}
```

label 基数限制：

```text
- query_name 必须是受控枚举或显式命名，不允许 raw SQL。
- database name 默认 hash 或 sanitized safe name。
- user、password、dsn、host with credentials 不允许作为 label。
```

### 6.8 Trace / Log Attributes Contract

标准 attrs：

```text
l2.module = postgresx
l2.provider = postgres
db.system = postgresql
db.operation = query|exec|tx|ping|health
db.query_name = <safe name>
db.sql_hash = optional sha256
db.result = success|error
db.error_kind = <ErrorKind>
db.retry_count = <n>
db.timeout_ms = <n>
postgresx.pool.acquired = <n>
postgresx.pool.idle = <n>
```

禁止 attrs：

```text
db.statement.raw
db.params.raw
password
dsn.raw
token
secret
tls.key
```

### 6.9 Transaction Contract

`WithinTx` P0 语义：

```text
1. BeginTx 成功后执行 fn。
2. fn 返回 nil：commit。
3. fn 返回 error：rollback，并返回原 error；rollback error 进入 joined/suppressed safe error。
4. fn panic：rollback，然后 re-panic；不得吞 panic。
5. commit error：返回 TRANSACTION_FAILED 或具体 pg error mapping。
6. rollback error：不得覆盖主错误，必须保留可审计信息。
7. context canceled：尝试 rollback，返回 CONTEXT_CANCELED 或 TIMEOUT。
8. 默认不做自动 retry。
9. 只有 caller 显式选择 retry policy，且 Tx 标记 idempotent，才允许对 SERIALIZATION_FAILED / DEADLOCK_DETECTED 做有限 retry。
```

P0 Tx options：

```text
IsolationLevel
ReadOnly
Deferrable
Timeout
QueryNamePrefix
```

### 6.10 Migration Boundary

P0 只允许：

```text
- docs/migration-boundary.md 说明边界。
- examples/integration 可以演示创建临时测试表。
- integration test 可以在临时数据库/临时 schema 中创建测试表。
```

禁止：

```text
- Open() 默认执行 migration。
- package init 自动 migration。
- 在 core 包中内置 x.go 业务表。
- 提供业务 schema。
```

P1 可选：

```text
- migration runner adapter interface
- migration lock helper
- advisory lock helper
```

必须通过 ADR 才能进入 public API。

---

## 7. Goal Runtime v3.1 对象模型

### 7.1 Master Goal

```text
GOAL-20260604-POSTGRESX-L2-FACTORY-001
Title: Upgrade postgresx into xlib-standard governed L2 PostgreSQL adapter factory
Mode: Full
Owner: ZoneCNH
Layer: L2
Target Repo: github.com/ZoneCNH/postgresx
Standard Source: github.com/ZoneCNH/xlib-standard
State Machine:
  INIT → CONTEXT_READY → GOAL_READY → SPEC_READY → DESIGN_READY → PLAN_READY
  → TASKS_READY → EXECUTING → VERIFYING → REVIEWING → RELEASING
  → RETROSPECTING → DONE
Exception States:
  BLOCKED / FAILED / NEEDS_RESEARCH / NEEDS_DECISION / NEEDS_REPLAN
  / NEEDS_ROLLBACK / NEEDS_HUMAN_APPROVAL / INCONSISTENT_STATE
```

### 7.2 Specs

```text
SPEC-POSTGRESX-L2-v1.0
SPEC-POSTGRESX-CONFIG-v1.0
SPEC-POSTGRESX-ERRORS-v1.0
SPEC-POSTGRESX-HEALTH-v1.0
SPEC-POSTGRESX-METRICS-v1.0
SPEC-POSTGRESX-TX-v1.0
SPEC-POSTGRESX-TESTKIT-v1.0
SPEC-POSTGRESX-RELEASE-v1.0
```

### 7.3 Designs

```text
DESIGN-POSTGRESX-L2-v1.0
DESIGN-POSTGRESX-PGX-DRIVER-v1.0
DESIGN-POSTGRESX-TX-v1.0
DESIGN-POSTGRESX-EVIDENCE-v1.0
```

### 7.4 ADRs

```text
ADR-20260604-POSTGRESX-001  postgresx is L2 adapter, not ORM or repository.
ADR-20260604-POSTGRESX-002  Use pgx/v5 as P0 provider SDK after AutoResearch.
ADR-20260604-POSTGRESX-003  Provider SDK types remain internal by default.
ADR-20260604-POSTGRESX-004  Query.Name required for observability; raw SQL not logged by default.
ADR-20260604-POSTGRESX-005  Unit tests use fake Queryer/Execer; real PostgreSQL only opt-in.
ADR-20260604-POSTGRESX-006  Tx retry is explicit and idempotency-gated.
ADR-20260604-POSTGRESX-007  Migration is boundary/helper, not startup behavior.
ADR-20260604-POSTGRESX-008  Release requires downstream adoption proof before “usable by x.go” claim.
```

---

## 8. Requirements 与 Acceptance Criteria

### 8.1 P0 Requirements

| ID | Requirement | Priority |
|---|---|---|
| REQ-POSTGRESX-001 | module path 必须为 `github.com/ZoneCNH/postgresx` | P0 |
| REQ-POSTGRESX-002 | README/AGENTS/CONSTITUTION/docs 必须声明 L2 PostgreSQL adapter 身份 | P0 |
| REQ-POSTGRESX-003 | 禁止 `x.go`、业务系统、其他 L2 imports | P0 |
| REQ-POSTGRESX-004 | 必须实现 Config/DefaultConfig/Validate/Sanitize | P0 |
| REQ-POSTGRESX-005 | 必须提供 `contracts/config.schema.json` | P0 |
| REQ-POSTGRESX-006 | 必须提供 pgxpool lifecycle：Open/Ping/Health/Stats/Close | P0 |
| REQ-POSTGRESX-007 | Close 必须幂等且可测试 | P0 |
| REQ-POSTGRESX-008 | 必须提供 Exec/Query/QueryRow minimal wrapper | P0 |
| REQ-POSTGRESX-009 | 必须提供 Queryer/Execer/Tx interface | P0 |
| REQ-POSTGRESX-010 | 必须提供 WithinTx helper 并覆盖 commit/rollback/panic/context 语义 | P0 |
| REQ-POSTGRESX-011 | 必须映射常见 pg SQLSTATE 到 ErrorKind | P0 |
| REQ-POSTGRESX-012 | DSN/password/SQL args 必须脱敏 | P0 |
| REQ-POSTGRESX-013 | 必须提供 health schema + golden tests | P0 |
| REQ-POSTGRESX-014 | 必须提供 metrics contract + emission tests | P0 |
| REQ-POSTGRESX-015 | 单元测试必须不依赖真实 PostgreSQL | P0 |
| REQ-POSTGRESX-016 | integration test 必须 opt-in 且记录 pass/skip/fail evidence | P0 |
| REQ-POSTGRESX-017 | 必须生成 release manifest + sha256 | P0 |
| REQ-POSTGRESX-018 | 必须有 downstream smoke adoption proof | P0 for stable, P1 for v0.1 |
| REQ-POSTGRESX-019 | 必须输出 retrospective + Prompt/Harness/Rule Patch candidates | P0 |
| REQ-POSTGRESX-020 | 完成声明必须使用 `DONE with evidence:` | P0 |

### 8.2 P1 Requirements

| ID | Requirement | Priority |
|---|---|---|
| REQ-POSTGRESX-021 | advisory lock helper | P1 |
| REQ-POSTGRESX-022 | prepared statement cache config | P1 |
| REQ-POSTGRESX-023 | tx retry helper with idempotency guard | P1 |
| REQ-POSTGRESX-024 | migration lock boundary helper | P1 |
| REQ-POSTGRESX-025 | query classifier / named query registry | P1 |
| REQ-POSTGRESX-026 | benchmark smoke for pool/query/tx overhead | P1 |
| REQ-POSTGRESX-027 | API diff gate before v0.2.0 | P1 |

### 8.3 P2 Requirements

| ID | Requirement | Priority |
|---|---|---|
| REQ-POSTGRESX-028 | read/write split optional package | P2 |
| REQ-POSTGRESX-029 | logical replication helper only after ADR | P2 |
| REQ-POSTGRESX-030 | transactional outbox optional package, not core | P2 |
| REQ-POSTGRESX-031 | migration runner adapter interface | P2 |
| REQ-POSTGRESX-032 | pool autotuning research | P2 |

### 8.4 Acceptance Criteria

```text
AC-POSTGRESX-001  `go.mod` contains `module github.com/ZoneCNH/postgresx`.
AC-POSTGRESX-002  README says `postgresx` is L2 PostgreSQL infrastructure adapter governed by xlib-standard.
AC-POSTGRESX-003  `rg 'github.com/bytechainx/x.go|github.com/ZoneCNH/(redisx|kafkax|natsx|taosx|ossx|clickhousex)' --glob '*.go'` returns no forbidden production import.
AC-POSTGRESX-004  `contracts/config.schema.json` validates examples.
AC-POSTGRESX-005  Config.Sanitize redacts DSN credentials, password, token, TLS key.
AC-POSTGRESX-006  `go test ./...` passes without PostgreSQL.
AC-POSTGRESX-007  fake Queryer/Execer/Tx covers success/error/timeout paths.
AC-POSTGRESX-008  WithinTx tests cover commit, rollback, panic rollback, context canceled, commit error, rollback error.
AC-POSTGRESX-009  pg error mapping tests cover 23505, 40001, 40P01, 55P03, 57014, 42P01, 28P01, 08xxx.
AC-POSTGRESX-010  health golden output validates against health schema.
AC-POSTGRESX-011  metrics tests prove pool/query/tx metrics are emitted and declared.
AC-POSTGRESX-012  secret scan and redaction tests pass.
AC-POSTGRESX-013  `L2_INTEGRATION=postgres make integration-check` passes or records explicit skip reason.
AC-POSTGRESX-014  release manifest and sha256 generated and excluded from source history when required by protocol.
AC-POSTGRESX-015  downstream smoke example imports postgresx, uses fake/sanitized config, runs health, closes client.
AC-POSTGRESX-016  retrospective generates Prompt Patch, Harness Patch, Rule Patch, New Issue Candidates.
```

---

## 9. Traceability Matrix

| Requirement | AC | Design | Task | Test | Evidence | Status |
|---|---|---|---|---|---|---|
| REQ-POSTGRESX-001 | AC-POSTGRESX-001 | DESIGN-POSTGRESX-L2 | TASK-POSTGRESX-001 | module check | EVID-POSTGRESX-001 | TODO |
| REQ-POSTGRESX-002 | AC-POSTGRESX-002 | DESIGN-POSTGRESX-L2 | TASK-POSTGRESX-002 | docs-check | EVID-POSTGRESX-002 | TODO |
| REQ-POSTGRESX-003 | AC-POSTGRESX-003 | Boundary Design | TASK-POSTGRESX-003 | boundary-check | EVID-POSTGRESX-003 | TODO |
| REQ-POSTGRESX-004/005 | AC-POSTGRESX-004/005 | Config Design | TASK-POSTGRESX-004 | config schema/redaction tests | EVID-POSTGRESX-004 | TODO |
| REQ-POSTGRESX-006/007 | AC-POSTGRESX-006 | Lifecycle Design | TASK-POSTGRESX-005 | lifecycle tests | EVID-POSTGRESX-005 | TODO |
| REQ-POSTGRESX-008/009 | AC-POSTGRESX-007 | API Design | TASK-POSTGRESX-006 | fake contract tests | EVID-POSTGRESX-006 | TODO |
| REQ-POSTGRESX-010 | AC-POSTGRESX-008 | Tx Design | TASK-POSTGRESX-007 | tx tests | EVID-POSTGRESX-007 | TODO |
| REQ-POSTGRESX-011 | AC-POSTGRESX-009 | Error Design | TASK-POSTGRESX-008 | error mapping tests | EVID-POSTGRESX-008 | TODO |
| REQ-POSTGRESX-012 | AC-POSTGRESX-005/012 | Security Design | TASK-POSTGRESX-009 | redaction/secret tests | EVID-POSTGRESX-009 | TODO |
| REQ-POSTGRESX-013/014 | AC-POSTGRESX-010/011 | Observability Design | TASK-POSTGRESX-010 | health/metrics tests | EVID-POSTGRESX-010 | TODO |
| REQ-POSTGRESX-016 | AC-POSTGRESX-013 | Integration Design | TASK-POSTGRESX-011 | integration-check | EVID-POSTGRESX-011 | TODO |
| REQ-POSTGRESX-017 | AC-POSTGRESX-014 | Release Design | TASK-POSTGRESX-012 | release-final-check | EVID-POSTGRESX-012 | TODO |
| REQ-POSTGRESX-018 | AC-POSTGRESX-015 | Adoption Design | TASK-POSTGRESX-013 | downstream smoke | EVID-POSTGRESX-013 | TODO |
| REQ-POSTGRESX-019 | AC-POSTGRESX-016 | Retro Design | TASK-POSTGRESX-014 | retro check | EVID-POSTGRESX-014 | TODO |

---

## 10. Harness Gates

### 10.1 Required Gates

```bash
GOWORK=off make fmt
GOWORK=off make vet
GOWORK=off make lint
GOWORK=off make test
GOWORK=off make boundary-check
GOWORK=off make contract-check
GOWORK=off make docs-check
GOWORK=off make security
GOWORK=off make evidence
```

### 10.2 Extended Gates

```bash
GOWORK=off make race
GOWORK=off make examples
GOWORK=off make golden
GOWORK=off make fuzz-smoke
GOWORK=off make benchmark-smoke
GOWORK=off make integration-check
GOWORK=off make ci-extended
```

### 10.3 Release Gates

```bash
XLIB_CONTEXT=release_verify GOWORK=off make release-check
XLIB_CONTEXT=release_verify GOWORK=off make release-final-check
XLIB_CONTEXT=release_verify GOWORK=off make release-preflight VERSION=v0.1.0
GOWORK=off go run ./cmd/goalcli score --min 9.8
```

如果 `postgresx` 仓库没有本地 `cmd/goalcli`，必须采用以下裁决之一：

```text
1. 从 xlib-standard 生成/复制标准化 goalcli toolchain；
2. 通过 xlib-standard 的 goalcli 对 postgresx 执行 repo path 参数化验证；
3. 明确标记 BLOCKED，不得把 score gate 记录为 skipped passed。
```

### 10.4 Boundary Gate

必须失败的情况：

```text
- production code import github.com/bytechainx/x.go
- production code import market-data/macro-data/regime-engine
- production code import 其他 L2 adapter
- production code import testkitx
- core package import provider observability exporter vendor
- README/docs/examples/release manifest 包含真实 password/token/DSN
- package 定义隐藏 global client
- init() 中连接数据库
```

建议脚本：

```bash
#!/usr/bin/env bash
set -euo pipefail

rg 'github.com/bytechainx/x.go|github.com/ZoneCNH/(redisx|kafkax|kafka|natsx|taosx|ossx|clickhousex)' --glob '*.go' && {
  echo "forbidden dependency found"; exit 1;
} || true

rg 'var\s+Default(Client|DB|Pool)|func\s+init\(\).*Open|pgxpool\.Connect\(context\.Background\(\)' --glob '*.go' && {
  echo "hidden global client or init connection found"; exit 1;
} || true

rg 'password=|postgres://[^ ]+:[^ ]+@|/home/k8s/secrets/env/.*=' README.md docs examples .agent release scripts --glob '!**/.git/**' && {
  echo "possible secret leak found"; exit 1;
} || true
```

### 10.5 Contract Gate

必须验证：

```text
contracts/config.schema.json is valid
contracts/health.schema.json is valid
contracts/metrics.contract.yaml is parseable
contracts/errors.contract.yaml is parseable
examples configs match config schema
health golden outputs match health schema
metrics emitted by tests are declared
error kinds in code are declared
```

### 10.6 Security Gate

必须验证：

```text
secret scan
redaction tests
optional govulncheck when enabled
no raw DSN/password/token in logs/errors/manifest
no SQL args in logs/traces
GitHub Actions pinned SHA when workflows added
```

### 10.7 Integration Gate

Integration 必须 opt-in：

```bash
L2_INTEGRATION=postgres make integration-check
```

允许 skip 的原因：

```text
docker unavailable
testcontainers unavailable
POSTGRESX_TEST_DSN absent
network unavailable
provider image unavailable
```

不允许：

```text
silent skip
默认连接生产数据库
把真实 DSN 写入 evidence
把 skip 写成 passed
```

Integration evidence JSON：

```json
{
  "module": "postgresx",
  "provider": "postgres",
  "enabled": true,
  "status": "pass|skip|fail",
  "reason": "docker unavailable",
  "server_version": "redacted-or-version",
  "dsn": "redacted",
  "checks": ["ping", "tx_commit", "tx_rollback"],
  "generated_at": "2026-06-04T00:00:00Z"
}
```

---

## 11. Evidence Protocol

### 11.1 必需 Evidence Artifacts

```text
release/manifest/latest.json
release/manifest/latest.json.sha256
.agent/evidence/GOAL-20260604-POSTGRESX-L2-FACTORY-001/gate-output.txt
.agent/evidence/GOAL-20260604-POSTGRESX-L2-FACTORY-001/test-output.txt
.agent/evidence/GOAL-20260604-POSTGRESX-L2-FACTORY-001/contract-hashes.txt
.agent/evidence/GOAL-20260604-POSTGRESX-L2-FACTORY-001/integration-evidence.json
.agent/evidence/GOAL-20260604-POSTGRESX-L2-FACTORY-001/downstream-adoption.md
.agent/reviews/REV-POSTGRESX-20260604-001.md
.agent/retrospectives/RETRO-20260604-POSTGRESX-001.md
```

### 11.2 Manifest 最小字段

```json
{
  "module": "github.com/ZoneCNH/postgresx",
  "package": "postgresx",
  "layer": "L2",
  "role": "postgresql_infrastructure_adapter",
  "standard_source": "github.com/ZoneCNH/xlib-standard",
  "standard_source_commit": "<sha>",
  "kernel_version": "<version-or-sha>",
  "l1_dependencies": {
    "configx": "<version-or-sha>",
    "observex": "<version-or-sha>",
    "testkitx": "test-only:<version-or-sha>",
    "resiliencx": "optional:<version-or-sha>",
    "schedulex": "optional:<version-or-sha>"
  },
  "provider_dependencies": {
    "pgx": "github.com/jackc/pgx/v5@<version>"
  },
  "commit": "<sha>",
  "tree_sha": "<sha>",
  "source_digest": "sha256:<digest>",
  "contract_hashes": {
    "api": "sha256:<digest>",
    "config": "sha256:<digest>",
    "health": "sha256:<digest>",
    "metrics": "sha256:<digest>",
    "errors": "sha256:<digest>"
  },
  "gates": {
    "fmt": "passed",
    "vet": "passed",
    "lint": "passed",
    "test": "passed",
    "race": "passed",
    "boundary": "passed",
    "contract": "passed",
    "docs": "passed",
    "security": "passed",
    "integration": "passed|skip-with-reason",
    "release_final": "passed"
  },
  "integration": {
    "status": "passed|skip-with-reason",
    "evidence": ".agent/evidence/.../integration-evidence.json"
  },
  "downstream_adoption": {
    "status": "not_claimed|passed",
    "consumer": "examples/downstream/postgresx-consumer",
    "evidence": ".agent/evidence/.../downstream-adoption.md"
  },
  "workflow": {
    "workflow_run_id": "local-or-github-run-id",
    "artifact_name": "postgresx-release-manifest",
    "artifact_url": "local-or-url"
  },
  "generated_at": "2026-06-04T00:00:00Z"
}
```

### 11.3 完成声明模板

```text
DONE with evidence:
- goal_id: GOAL-20260604-POSTGRESX-L2-FACTORY-001
- repo: github.com/ZoneCNH/postgresx
- branch: goal/GOAL-20260604-POSTGRESX-L2-FACTORY-001
- commit: <sha>
- tag: v0.1.0 or not created
- gates:
  - GOWORK=off make ci: passed
  - GOWORK=off make boundary-check: passed
  - GOWORK=off make contract-check: passed
  - GOWORK=off make docs-check: passed
  - GOWORK=off make security: passed
  - L2_INTEGRATION=postgres make integration-check: passed|skip-with-reason
  - XLIB_CONTEXT=release_verify GOWORK=off make release-final-check: passed
- manifest: release/manifest/latest.json
- manifest_sha256: release/manifest/latest.json.sha256
- workflow_artifact: <url-or-local-path>
- downstream_adoption: <proof path>
- known_gaps: <none or explicit>
- retrospective: .agent/retrospectives/RETRO-20260604-POSTGRESX-001.md
```

### 11.4 禁止声明

```text
- 禁止说 “tests pass” 但没有命令输出。
- 禁止把 skipped integration 写成 passed。
- 禁止把 README 更新写成 adoption。
- 禁止 dirty workspace 下写 release ready。
- 禁止把 local fake test 写成真实 PostgreSQL integration。
- 禁止没有 downstream proof 时写 “x.go 可直接使用”。
```

---

## 12. Worktree 执行标准

```bash
# 0. 准备主仓库
mkdir -p ~/code/ZoneCNH
cd ~/code/ZoneCNH

# 1. postgresx
git clone git@github.com:ZoneCNH/postgresx.git postgresx || true
cd ~/code/ZoneCNH/postgresx
git checkout main
git pull --ff-only

# 2. 禁止在 main 开发
git status --short

# 3. 创建独立 worktree
mkdir -p ~/code/ZoneCNH/.worktree
git worktree add ~/code/ZoneCNH/.worktree/postgresx-l2-factory-20260604 \
  -b goal/GOAL-20260604-POSTGRESX-L2-FACTORY-001 main

cd ~/code/ZoneCNH/.worktree/postgresx-l2-factory-20260604

# 4. 绑定标准源版本
git clone git@github.com:ZoneCNH/xlib-standard.git ../xlib-standard-standard-source || true
cd ../xlib-standard-standard-source
git checkout main
git pull --ff-only
STANDARD_SOURCE_SHA=$(git rev-parse HEAD)

cd ~/code/ZoneCNH/.worktree/postgresx-l2-factory-20260604

# 5. 创建 evidence 目录
mkdir -p .agent/evidence/GOAL-20260604-POSTGRESX-L2-FACTORY-001
```

禁止：

```text
- main 上直接 commit
- 多个 Agent 共用一个 worktree
- 把 .worktree/ 运行态提交
- 把真实 secret/env 文件提交
- 未跑 gate 创建 PR
- 未有 evidence 声称 DONE
```

---

## 13. 任务拆解

### 13.1 Phase 0：Context Recovery

```text
TASK-POSTGRESX-000  仓库事实盘点
输出：
- .agent/context/current-state.md
- go.mod status
- README status
- Makefile status
- .agent status
- contracts status
- CI status
- xlib-standard source SHA
验收：
- 不把空仓库写成已实现
- 不把 registered 写成 adopted
```

### 13.2 Phase 1：Identity Correction

```text
TASK-POSTGRESX-001  修正 go.mod
TASK-POSTGRESX-002  修正 README
TASK-POSTGRESX-003  添加 AGENTS.md
TASK-POSTGRESX-004  添加 CONSTITUTION.md
TASK-POSTGRESX-005  添加 docs/spec.md / docs/design.md / docs/api.md
TASK-POSTGRESX-006  添加 .agent Goal Runtime v3.1 工件
```

验收：

```bash
rg 'xlib-standard|baselib-template|foundationx' README.md docs .agent
# 只允许作为标准源或迁移语境出现，不允许作为当前仓库身份。
```

### 13.3 Phase 2：Contracts

```text
TASK-POSTGRESX-010  contracts/config.schema.json
TASK-POSTGRESX-011  contracts/health.schema.json
TASK-POSTGRESX-012  contracts/metrics.contract.yaml
TASK-POSTGRESX-013  contracts/errors.contract.yaml
TASK-POSTGRESX-014  contracts/api.contract.yaml
TASK-POSTGRESX-015  contract-check script
```

验收：

```bash
GOWORK=off make contract-check
```

### 13.4 Phase 3：Core API + pgx Driver Isolation

```text
TASK-POSTGRESX-020  pkg/postgresx Config/Validate/Sanitize
TASK-POSTGRESX-021  pkg/postgresx Client interfaces
TASK-POSTGRESX-022  internal/driver/pgx implementation
TASK-POSTGRESX-023  lifecycle Open/Ping/Health/Stats/Close
TASK-POSTGRESX-024  Query/Exec/QueryRow wrapper
TASK-POSTGRESX-025  provider type isolation check
```

验收：

```bash
GOWORK=off go test ./pkg/... ./internal/...
GOWORK=off make boundary-check
```

### 13.5 Phase 4：Transaction Contract

```text
TASK-POSTGRESX-030  Tx interface
TASK-POSTGRESX-031  WithinTx helper
TASK-POSTGRESX-032  commit / rollback / panic / context tests
TASK-POSTGRESX-033  Tx retry ADR and explicit opt-in guard
```

验收：

```bash
GOWORK=off go test ./... -run 'TestWithinTx'
```

### 13.6 Phase 5：Error Mapping + Redaction

```text
TASK-POSTGRESX-040  ErrorKind enum / classifier
TASK-POSTGRESX-041  pgconn.PgError mapping
TASK-POSTGRESX-042  pool/context/network errors mapping
TASK-POSTGRESX-043  redacted error output
TASK-POSTGRESX-044  SQL args never logged tests
```

验收：

```bash
GOWORK=off go test ./... -run 'TestErrorMapping|TestRedaction'
GOWORK=off make security
```

### 13.7 Phase 6：Health + Observability

```text
TASK-POSTGRESX-050  Health model
TASK-POSTGRESX-051  health golden fixtures
TASK-POSTGRESX-052  pool/query/tx metrics
TASK-POSTGRESX-053  log/trace attrs policy
TASK-POSTGRESX-054  observex memory recorder tests
```

验收：

```bash
GOWORK=off make contract-check
GOWORK=off go test ./... -run 'TestHealth|TestMetrics'
```

### 13.8 Phase 7：Fake/Testkit

```text
TASK-POSTGRESX-060  testkit FakeClient
TASK-POSTGRESX-061  testkit FakeQueryer/FakeExecer
TASK-POSTGRESX-062  testkit FakeTx
TASK-POSTGRESX-063  FailureInjector / LatencyInjector
TASK-POSTGRESX-064  ContractAssertions
```

验收：

```bash
GOWORK=off go test ./... 
# 必须不需要真实 PostgreSQL。
```

### 13.9 Phase 8：Integration Opt-in

```text
TASK-POSTGRESX-070  integration_check.sh
TASK-POSTGRESX-071  Docker/testcontainers or explicit external test DSN support
TASK-POSTGRESX-072  Ping integration
TASK-POSTGRESX-073  Tx commit/rollback integration
TASK-POSTGRESX-074  integration evidence JSON
```

验收：

```bash
L2_INTEGRATION=postgres GOWORK=off make integration-check
```

### 13.10 Phase 9：Release Evidence

```text
TASK-POSTGRESX-080  release manifest template
TASK-POSTGRESX-081  manifest generator
TASK-POSTGRESX-082  manifest sha256
TASK-POSTGRESX-083  release-evidence-check
TASK-POSTGRESX-084  release-final-check
```

验收：

```bash
XLIB_CONTEXT=release_verify GOWORK=off make release-final-check
```

### 13.11 Phase 10：Downstream Adoption

```text
TASK-POSTGRESX-090  examples/downstream/postgresx-consumer
TASK-POSTGRESX-091  compile proof
TASK-POSTGRESX-092  fake health smoke
TASK-POSTGRESX-093  optional x.go compile-only branch
TASK-POSTGRESX-094  adoption evidence
```

验收：

```bash
GOWORK=off make downstream-smoke
```

### 13.12 Phase 11：Retrospective + Self-improving

```text
TASK-POSTGRESX-100  RETRO-20260604-POSTGRESX-001
TASK-POSTGRESX-101  PATCH-PROMPT-20260604-POSTGRESX-001
TASK-POSTGRESX-102  PATCH-HARNESS-20260604-POSTGRESX-001
TASK-POSTGRESX-103  PATCH-RULE-20260604-POSTGRESX-001
TASK-POSTGRESX-104  ISSUE-CANDIDATES-20260604-POSTGRESX
```

---

## 14. AI / 自动化 / 研究增强介入位置

### 14.1 gstack

```text
gstack/postgresx-l2
  G0: Context Recovery
  G1: Identity Correction
  G2: Contracts
  G3: API + Driver Isolation
  G4: Tx Semantics
  G5: Error Mapping + Redaction
  G6: Health + Metrics
  G7: Fake/Testkit
  G8: Integration Opt-in
  G9: Release Evidence
  G10: Downstream Adoption
  G11: Retrospective / Self-improving
```

### 14.2 superpowers

| Agent | Superpower | 输出 |
|---|---|---|
| Repo Scanner | 扫描仓库结构与当前事实 | current-state.md |
| Standard Sync Agent | 从 xlib-standard 同步标准包 | standard-source.md / adoption notes |
| Boundary Auditor | 发现 forbidden imports、global client、secret 泄漏 | boundary evidence |
| Contract Generator | 生成 config/health/metrics/errors contracts | contracts/** |
| API Architect | 设计 context-first API / Tx contract | docs/api.md / pkg/postgresx |
| Driver Agent | 隔离 pgx driver 实现 | internal/driver/pgx |
| Error Taxonomist | SQLSTATE → ErrorKind | errors.contract.yaml |
| Test Agent | fake/testkit/tx/redaction/race tests | test-output.txt |
| Integration Agent | opt-in PostgreSQL integration | integration-evidence.json |
| Evidence Agent | manifest/checksum/artifact | release/manifest |
| Adoption Agent | downstream compile proof | downstream-adoption.md |
| Retro Agent | failure → patch | Prompt/Harness/Rule Patch |

### 14.3 Harness

Harness 的核心裁决：

```text
PostgreSQL behavior is not trusted unless tests, contracts and evidence prove it.
```

每条规则必须形成：

```text
Rule -> Gate -> Evidence -> Release
```

### 14.4 Compound Engineering

复利路径：

```text
一次 DSN 泄露风险
  -> redaction test
  -> secret gate
  -> xlib-standard rule patch
  -> 所有 L2 复用

一次 Tx rollback bug
  -> tx golden fixture
  -> contract test
  -> postgresx testkit
  -> downstream repository 复用

一次 pg error 未映射
  -> ErrorKind patch
  -> errors.contract.yaml
  -> observability dashboard 更稳定
```

### 14.5 Self-improving

每次失败必须输出：

```text
Failure
  -> Root Cause
  -> Missing Rule / Missing Gate / Missing Test
  -> Patch Candidate
  -> xlib-standard / postgresx rule update
  -> Regression fixture
  -> Next release gate stronger
```

### 14.6 AutoResearch

必须触发 AutoResearch 的问题：

```text
- pgx/v5 当前最新 API、Go version、license、pool behavior 不确定。
- pgconn.PgError SQLSTATE mapping 不确定。
- pgxpool Close/Acquire/Stats 行为不确定。
- Tx rollback / panic recover 行为与 context 交互不确定。
- PostgreSQL Docker image / testcontainers 版本不确定。
- GitHub Actions / services container 版本不确定。
- govulncheck / pgx vulnerability 状态不确定。
```

输出模板：

```md
# RESEARCH-20260604-POSTGRESX-001

## Question

## Sources

## Findings

## Decision Needed

## Proposed Patch

## Evidence
```

---

## 15. 可复利增长的系统架构

```text
xlib-standard
  ├─ L2 postgres adapter standard
  ├─ contract pack
  ├─ harness pack
  ├─ evidence protocol
  ├─ release policy
  └─ generator overlay
        ↓
postgresx
  ├─ config contract
  ├─ client lifecycle
  ├─ transaction helper
  ├─ error classifier
  ├─ observability contract
  ├─ fake/testkit
  ├─ integration evidence
  └─ release manifest
        ↓
downstream adoption
  ├─ x.go compile proof
  ├─ market-data config usage
  ├─ macro-data state storage
  ├─ engine metadata persistence
  └─ service template adoption
        ↓
feedback
  ├─ missing SQLSTATE
  ├─ tx edge case
  ├─ pool metric gap
  ├─ config pain point
  └─ secret redaction failure
        ↓
retrospective patches
        ↓
xlib-standard + postgresx improve
```

复利公式：

```text
Postgresx Leverage =
  reusable_contract
  * reusable_fake
  * reusable_tx_semantics
  * reusable_error_classifier
  * downstream_adoption_count
  * regression_memory
  / manual_pg_boilerplate_cost
```

---

## 16. 最小可行行动 MVA

### 16.1 MVA 目标

```text
MVA-POSTGRESX-001  让 postgresx 从空仓库变成 xlib-standard 可审计的 L2 PostgreSQL adapter。
```

### 16.2 MVA 范围

必须做：

```text
1. go.mod / README / docs / .agent 身份修正。
2. contracts/config.schema.json。
3. contracts/health.schema.json。
4. contracts/metrics.contract.yaml。
5. contracts/errors.contract.yaml。
6. pkg/postgresx Config/Validate/Sanitize。
7. Client interface + pgxpool implementation。
8. Ping/Health/Stats/Close。
9. Exec/Query/QueryRow minimal wrapper。
10. WithinTx helper。
11. Error mapping + redaction。
12. Fake Queryer/Execer/Tx。
13. Required gates。
14. Integration opt-in。
15. Release manifest + sha256。
16. Downstream smoke example。
17. Retrospective patches。
```

不做：

```text
- ORM。
- 业务 repository。
- 业务 migration。
- x.go 生产接入。
- 自动生产连接。
- 默认真实 integration。
- 读写分离。
- 分库分表。
- transactional outbox。
```

### 16.3 MVA 完成标准

```text
DONE with evidence:
- postgresx identity corrected.
- Config/health/metrics/errors contracts exist and pass contract-check.
- Unit tests pass without PostgreSQL.
- fake Queryer/Execer/Tx supports downstream tests.
- WithinTx behavior is tested.
- pg SQLSTATE mapping is tested.
- DSN/password/SQL args redaction tests pass.
- integration opt-in passes or records explicit skip reason.
- release manifest + sha256 generated.
- downstream smoke adoption proof exists.
- retrospective patch generated.
```

---

## 17. 1 天行动计划

### Day 1 目标

把 `postgresx` 从空 README 仓库升级为可执行的 L2 MVA 起点。

### Day 1 步骤

```text
1. 创建 worktree：goal/GOAL-20260604-POSTGRESX-L2-FACTORY-001。
2. 记录 current-state：README 只有项目名、无 go.mod/Makefile/.agent/contracts 时如实记录。
3. 从 xlib-standard 固定 standard source commit。
4. 创建 go.mod：module github.com/ZoneCNH/postgresx。
5. 创建 README：声明 L2 PostgreSQL adapter 身份、边界、禁止项、Quickstart。
6. 创建 AGENTS.md / CONSTITUTION.md：禁止 main、禁止 x.go、禁止 secret、DONE with evidence。
7. 创建 .agent/goal/spec/design/plan/traceability/risk/decision。
8. 创建 contracts/config.schema.json、health.schema.json、metrics.contract.yaml、errors.contract.yaml。
9. 实现 Config/DefaultConfig/Validate/Sanitize。
10. 实现 fake Queryer/Execer/Tx。
11. 实现 WithinTx fake contract tests。
12. 添加 Makefile required targets：fmt/vet/test/boundary/contract/docs/security。
13. 跑 `GOWORK=off go test ./...`。
14. 跑 boundary-check、contract-check、docs-check、security。
15. 生成 Day 1 evidence 草案。
```

### Day 1 不做

```text
- 不连接真实 PostgreSQL。
- 不做 migration。
- 不把 x.go 表结构写入 examples。
- 不写 raw production DSN。
- 不创建 release tag。
```

---

## 18. 7 天行动计划

### Day 1：身份与契约

```text
- repo identity correction
- .agent Goal Runtime v3.1 artifacts
- contracts
- fake/testkit MVA
- required gates
```

### Day 2：pgx driver + lifecycle

```text
- AutoResearch pgx/v5
- internal/driver/pgx
- Open/Ping/Health/Stats/Close
- context timeout
- Close idempotency
```

### Day 3：Query / Exec / Error Mapping

```text
- Exec/Query/QueryRow wrapper
- Query.Name policy
- SQLSTATE mapping
- redaction tests
- query metrics
```

### Day 4：Transaction Contract

```text
- WithinTx real implementation
- commit/rollback/panic/context behavior
- Tx options
- no default retry
- optional explicit retry ADR
```

### Day 5：Observability + Security

```text
- health golden tests
- pool/query/tx metrics
- log/trace attrs
- secret scan
- SQL args redaction
```

### Day 6：Integration Opt-in + Downstream Smoke

```text
- PostgreSQL integration via Docker/testcontainers/shared CI service
- integration evidence JSON
- examples/downstream/postgresx-consumer
- compile proof
```

### Day 7：Release Gate + Retrospective

```text
- release manifest
- checksum
- release-final-check
- score >= 9.8 or blocked with reason
- downstream adoption proof
- retrospective patches
- PR summary with DONE with evidence or NOT DONE yet
```

---

## 19. 30 天行动计划

### Week 1：v0.1.0 MVA

```text
目标：postgresx 可作为 L2 PostgreSQL adapter 的最小可证明版本。

完成：
- identity correction
- config/health/metrics/errors contracts
- pgxpool lifecycle
- fake/testkit
- WithinTx helper
- error mapping
- redaction
- integration opt-in
- release evidence
- downstream smoke adoption
```

### Week 2：v0.2.0 Contract Hardening

```text
目标：把 tx/error/observability 从“可用”强化为“可长期复用”。

完成：
- API diff gate
- richer SQLSTATE mapping
- tx retry explicit policy
- advisory lock helper design
- prepared statement cache config
- benchmark smoke
- mutation fixtures for redaction and error mapping
```

### Week 3：Integration & Reliability Baseline

```text
目标：证明 postgresx 在真实 PostgreSQL 下行为稳定。

完成：
- PostgreSQL version matrix
- Docker/testcontainers integration
- pool exhaustion tests
- context timeout tests
- deadlock/serialization integration fixtures if feasible
- race tests
- benchmark baseline
```

### Week 4：Downstream Adoption + Release Train

```text
目标：让 postgresx 能被 x.go / data services 以 compile-proof 方式采纳，但不引入业务污染。

完成：
- x.go compile-only branch optional
- market-data/macro-data storage adapter smoke optional
- release v0.1.x / v0.2.x tag
- xlib-standard downstream status update only with evidence
- retrospective patch 回写 xlib-standard
```

---

## 20. 衡量指标

### 20.1 工程指标

```text
Required gates pass rate: 100%
Unit tests without PostgreSQL: 100%
Boundary violations: 0
Secret findings: 0
Raw DSN/password/SQL args leak: 0
Release manifest generated: 100%
Manifest checksum verified: 100%
Integration silent skips: 0
Downstream smoke compile proof: >= 1
```

### 20.2 API / Contract 指标

```text
Config schema coverage: 100% P0 fields
Health schema golden coverage: 100%
Metrics declared/emitted mismatch: 0
ErrorKind mapping coverage for P0 SQLSTATE: 100%
Public provider SDK leakage: 0 unless ADR-approved
Breaking API changes before v1: tracked by API diff gate
```

### 20.3 PostgreSQL 行为指标

```text
Ping latency p50/p95/p99
Query latency p50/p95/p99 by query_name
Tx duration p50/p95/p99
Tx rollback count by reason
Pool acquired/idle/total/max
Pool acquire duration
Error rate by ErrorKind
Retry count only when explicit policy enabled
Health pass/warn/fail ratio
```

### 20.4 复利指标

```text
Time to bootstrap next DB adapter: target < 30 minutes
postgresx fake reuse count in downstream tests
Error mapping patch reuse by other SQL adapters
Redaction regression count: decreasing
Manual release checklist items: decreasing
Failure converted to patch candidate: >= 95%
```

---

## 21. 迭代优化机制

每个 PR / release 必须回答：

```text
1. 哪个 gate 最有价值？
2. 哪个问题是人工发现但 gate 没拦住？
3. 哪个 pgx/PostgreSQL 行为与预期不一致？
4. 哪个 SQLSTATE 需要升入 errors contract？
5. 哪个 metrics label 基数过高？
6. 哪个 config 字段应该进入 xlib-standard L2 common contract？
7. 哪个 fake/testkit 能被下游复用？
8. 哪个 integration test 不稳定？
9. 哪个安全/脱敏规则应升级为 P0 gate？
10. 哪些 P1/P2 能力需要新 Issue？
```

输出：

```text
RETRO-20260604-POSTGRESX-001
PATCH-PROMPT-20260604-POSTGRESX-001
PATCH-HARNESS-20260604-POSTGRESX-001
PATCH-RULE-20260604-POSTGRESX-001
PATCH-GENERATOR-20260604-POSTGRESX-001
ISSUE-CANDIDATES-20260604-POSTGRESX.md
```

---

## 22. Change Propagation Matrix

| 变更源 | 必须同步 |
|---|---|
| Config 字段变更 | config.schema.json、docs/config.md、examples、redaction tests、manifest |
| ErrorKind 变更 | errors.contract.yaml、error tests、observability docs、downstream docs |
| Health 输出变更 | health.schema.json、golden tests、docs/health.md、manifest contract |
| Metrics 变更 | metrics.contract.yaml、observex tests、dashboard docs |
| Public API 变更 | docs/api.md、api.contract.yaml、examples、downstream smoke、SemVer |
| Tx 语义变更 | docs/design.md、tx tests、risk register、ADR |
| Provider SDK 变更 | dependency research、go.mod、release manifest、security |
| Integration 策略变更 | docs/integration.md、CI、integration-evidence schema |
| Release manifest 变更 | release docs、evidence check、xlib-standard protocol |
| xlib-standard rule 变更 | postgresx AGENTS/CONSTITUTION/Makefile/CI |
| Downstream adoption 变更 | adoption evidence、downstream status、release notes |

---

## 23. Risk Register

| Risk ID | 风险 | 影响 | 处理 |
|---|---|---|---|
| RISK-POSTGRESX-001 | pgx provider 类型泄漏到 public API | 下游被 SDK 锁死 | internal driver + ADR for escape hatch |
| RISK-POSTGRESX-002 | Tx helper 吞 rollback/commit 错误 | 数据一致性风险 | explicit tx tests + joined error policy |
| RISK-POSTGRESX-003 | 自动重试写事务 | 重复副作用 | 默认禁用，idempotency opt-in |
| RISK-POSTGRESX-004 | DSN/SQL args 泄露 | 严重安全风险 | redaction gate + secret scan |
| RISK-POSTGRESX-005 | integration 连接生产 DB | 数据破坏/泄密 | opt-in + sanitized evidence + allow_external_dsn=false default |
| RISK-POSTGRESX-006 | pool 默认值不合理 | 连接耗尽或 DB 压垮 | conservative defaults + docs |
| RISK-POSTGRESX-007 | Ping 伪健康 | 误判 readiness | multi-check health |
| RISK-POSTGRESX-008 | migration helper 越界 | 变成业务 schema 库 | migration boundary ADR |
| RISK-POSTGRESX-009 | metrics label 高基数 | 观测成本爆炸 | query_name controlled enum / SQL hash |
| RISK-POSTGRESX-010 | release overclaim | 标准信任下降 | status no-overclaim + evidence protocol |

---

## 24. Rollback Protocol

触发 rollback：

```text
- secret leak detected
- release-final-check failed after manifest generation
- public API 误暴露 pgx type
- Tx helper 语义错误
- integration accidentally touches production DB
- manifest checksum mismatch
- downstream smoke breaks
```

回滚步骤：

```bash
git status --short
git log --oneline -n 10

# 如果未提交
git restore --staged .
git restore .
git clean -fdX

# 如果已提交到 PR 分支
git revert <bad_commit_sha>

GOWORK=off make ci
GOWORK=off make boundary-check
GOWORK=off make security
XLIB_CONTEXT=release_verify GOWORK=off make release-final-check
```

Rollback Evidence：

```text
ROLLBACK DONE with evidence:
- bad_commit:
- rollback_commit:
- reason:
- affected_files:
- gates_rerun:
- remaining_gaps:
- preventive_patch:
```

---

## 25. Human Approval Gates

以下必须人工批准：

```text
- 改变 postgresx public API。
- 暴露 pgx/pgconn/pgxpool 类型到 public API。
- 引入 ORM。
- 引入 migration runner。
- 默认启用 tx retry。
- 默认记录 raw SQL。
- 默认记录 SQL args。
- 默认连接真实 integration DB。
- 接入 x.go production path。
- 创建 v1.0.0 tag。
- 修改 release evidence semantics。
```

审批记录：

```text
DEC-20260604-POSTGRESX-001
- decision:
- alternatives:
- risk accepted:
- owner:
- expiry:
- evidence:
```

---

## 26. Failure Budget

```text
P0 gate failure: 0 allowed
secret leakage: 0 allowed
raw DSN leak: 0 allowed
raw SQL args leak: 0 allowed
main direct commit: 0 allowed
x.go reverse dependency: 0 allowed
other L2 dependency: 0 allowed
release overclaim: 0 allowed
skipped required gate marked passed: 0 allowed
Tx semantic regression: 0 allowed before release
docs/code drift: <= 1 release cycle
```

超过 failure budget：

```text
- block release
- create retro
- create rule/harness/prompt patch
- require human decision
```

---

## 27. Issue / PR / Commit / Release 规范

### 27.1 Issue 模板

```markdown
## Goal
Implement postgresx L2 PostgreSQL adapter MVA.

## Scope
- Identity correction
- Config contract
- Client lifecycle
- Query/Exec/Tx contract
- Error mapping
- Health contract
- Observability contract
- Fake/testkit
- Integration opt-in
- Release evidence
- Downstream adoption

## Acceptance Criteria
- [ ] go.mod module is github.com/ZoneCNH/postgresx
- [ ] no x.go / other L2 imports
- [ ] config/health/metrics/errors contracts pass
- [ ] go test ./... passes without PostgreSQL
- [ ] WithinTx behavior tested
- [ ] pg SQLSTATE mapping tested
- [ ] redaction tests pass
- [ ] integration pass or explicit skip evidence
- [ ] release manifest generated
- [ ] downstream adoption proof exists

## Evidence Required
DONE with evidence:
- commit
- gate output
- manifest
- checksum
- integration evidence
- downstream adoption
- retrospective
```

### 27.2 PR 模板

```markdown
## What changed

## Layer boundary
- [ ] L2 only
- [ ] no x.go import
- [ ] no other L2 import
- [ ] no business repository/schema
- [ ] no production secret

## PostgreSQL semantics
- [ ] Tx commit/rollback behavior tested
- [ ] pg errors mapped
- [ ] DSN/SQL args redacted

## Gates
- [ ] make ci
- [ ] make boundary-check
- [ ] make contract-check
- [ ] make docs-check
- [ ] make security
- [ ] make integration-check or skip evidence
- [ ] make release-final-check

## Evidence

## Downstream impact

## Retrospective patch
```

### 27.3 Commit 规范

```text
feat(postgresx): add config contract and sanitizer
feat(postgresx): add pgxpool lifecycle client
feat(postgresx): add transaction helper contract
fix(postgresx): redact dsn in error mapping
test(postgresx): add fake tx rollback fixtures
docs(postgresx): document migration boundary
chore(postgresx): add release evidence manifest
```

### 27.4 Release 规范

```text
v0.1.0  First L2 PostgreSQL adapter MVA
v0.2.0  Contract hardening + API diff gate + benchmark smoke
v0.3.0  Integration matrix + downstream adoption expansion
v1.0.0  Stable public API and compatibility promise
```

---

## 28. 交付清单

### 28.1 xlib-standard 侧交付物

如果 `xlib-standard` 已经有 L2 标准包，则只需记录 source SHA 与 impact。若缺失，必须补：

```text
docs/standard/l2-adapter-standard.md
templates/l2-adapter/postgresx.overlay.yaml
contracts/l2/postgresx/config.schema.json
contracts/l2/postgresx/health.schema.json
contracts/l2/postgresx/metrics.contract.yaml
contracts/l2/postgresx/errors.contract.yaml
.agent/rules/l2-postgresx-boundary.md
.agent/rules/l2-postgresx-evidence.md
```

### 28.2 postgresx 仓库交付物

```text
README.md
AGENTS.md
CONSTITUTION.md
go.mod
Makefile
.agent/goal.md
.agent/spec.md
.agent/design.md
.agent/plan.md
.agent/traceability-matrix.md
.agent/risk-register.md
.agent/decision-log.md
contracts/api.contract.yaml
contracts/config.schema.json
contracts/health.schema.json
contracts/metrics.contract.yaml
contracts/errors.contract.yaml
docs/spec.md
docs/design.md
docs/api.md
docs/config.md
docs/health.md
docs/observability.md
docs/resilience.md
docs/testing.md
docs/integration.md
docs/security.md
docs/release.md
docs/migration-boundary.md
pkg/postgresx/*.go
internal/driver/pgx/*.go
internal/errors/*.go
internal/health/*.go
internal/metrics/*.go
testkit/*.go
examples/basic
examples/health
examples/tx
examples/observability
examples/downstream/postgresx-consumer
scripts/boundary_check.sh
scripts/contract_check.sh
scripts/docs_check.sh
scripts/integration_check.sh
scripts/secret_scan.sh
scripts/generate_manifest.sh
release/manifest/latest.json generated
release/manifest/latest.json.sha256 generated
.agent/evidence/... generated
.agent/retrospectives/... generated
```

---

## 29. Master Goal 可执行 Prompt

下面这段可直接交给 Agent / Codex / goalkit 执行。

```text
You are executing GOAL-20260604-POSTGRESX-L2-FACTORY-001.

Objective:
Upgrade github.com/ZoneCNH/postgresx from a minimal standalone repository into an xlib-standard governed L2 PostgreSQL infrastructure adapter library with independent release, independent Evidence, shared L0/L1 contracts, Harness gates, Release Gate, downstream adoption proof, and Self-improving feedback.

Execution protocol:
Use Goal Runtime Prompt v3.1:
Goal → Context Recovery → Spec → Design → Plan → Tasks → Execution → Verification → Evidence → Review → Release → Retrospective → Self-improving.

Mode:
Full.

Hard constraints:
1. Do not commit on main. Use git worktree.
2. Do not import x.go, market-data, macro-data, regime-engine, or any business module.
3. Do not import other L2 adapters from production code.
4. Do not implement business repositories, business table schemas, or application transaction orchestration.
5. Do not leak secrets, DSNs, passwords, tokens, SQL args, or /home/k8s/secrets/env/* contents.
6. Do not expose pgx provider SDK types in public API unless an ADR explicitly approves it.
7. Unit tests must pass without real PostgreSQL.
8. Integration tests must be opt-in and produce pass/skip/fail evidence.
9. Release requires manifest, checksum, gates, contract hashes, dependency versions, downstream adoption proof, and retrospective patches.
10. Completion must be declared only as "DONE with evidence:".

Current context to verify:
- github.com/ZoneCNH/postgresx exists.
- README may currently be minimal.
- xlib-standard is the standard source.
- xlib-standard downstream matrix registers postgresx as L2 with adoption not yet proven.
- Do not claim adoption until current postgresx/downstream commands produce proof.

Primary sequence:
1. Inspect postgresx and xlib-standard.
2. Produce current-state report.
3. Create worktree branch goal/GOAL-20260604-POSTGRESX-L2-FACTORY-001.
4. Correct repo identity: go.mod, README, AGENTS, CONSTITUTION, docs, .agent.
5. Add contracts: config, health, metrics, errors, API.
6. Implement Config/Validate/Sanitize.
7. Implement Client interface and pgxpool driver isolation.
8. Implement Ping/Health/Stats/Close.
9. Implement Exec/Query/QueryRow minimal wrappers.
10. Implement WithinTx with commit/rollback/panic/context semantics.
11. Implement pg error mapping and redaction.
12. Implement fake Queryer/Execer/Tx and testkit.
13. Add required gates and scripts.
14. Add integration opt-in for PostgreSQL.
15. Generate release manifest and sha256.
16. Add downstream smoke adoption example.
17. Run gates.
18. Produce evidence.
19. Review release readiness.
20. Generate retrospective patches.

Required gates:
- GOWORK=off make fmt
- GOWORK=off make vet
- GOWORK=off make lint
- GOWORK=off make test
- GOWORK=off make boundary-check
- GOWORK=off make contract-check
- GOWORK=off make docs-check
- GOWORK=off make security
- L2_INTEGRATION=postgres GOWORK=off make integration-check or explicit skip evidence
- XLIB_CONTEXT=release_verify GOWORK=off make release-check
- XLIB_CONTEXT=release_verify GOWORK=off make release-final-check
- GOWORK=off go run ./cmd/goalcli score --min 9.8 or equivalent xlib-standard score gate

Deliverables:
- postgresx repo identity corrected.
- Contracts complete.
- Public API and fake/testkit complete.
- pgx driver isolated.
- Tx semantics implemented and tested.
- Redaction and error mapping tested.
- Integration opt-in evidence produced.
- Release manifest and checksum produced.
- Downstream smoke adoption proof produced.
- Retrospective patches produced.

Completion format:
DONE with evidence:
- goal_id:
- repo:
- worktree:
- branch:
- commit:
- tag:
- commands:
- gates:
- manifest:
- manifest_sha256:
- integration:
- downstream_adoption:
- known_gaps:
- retrospective:
```

---

## 30. 最终推荐路径

最优路径不是“先把 pgx 封装完整”，而是：

```text
先身份，后契约；先 fake，后真实；先 evidence，后 release；先 smoke adoption，后 x.go production。
```

推荐顺序：

```text
1. 先修正 postgresx 仓库身份，避免所有 Evidence 证明错对象。
2. 立刻建立 contracts/config/health/metrics/errors，让后续代码有机器可判定的目标。
3. 先实现 fake/testkit 和 Tx contract tests，保证单元测试不依赖真实 PostgreSQL。
4. 再接入 pgxpool driver，并隔离到 internal/driver/pgx。
5. 用 error mapping 和 redaction gate 处理 PostgreSQL adapter 的最大风险。
6. 用 integration opt-in 验证真实 PostgreSQL，但不让真实 DB 成为默认测试前提。
7. 生成 release manifest 和 checksum，禁止无 Evidence 的完成声明。
8. 用 downstream smoke example 证明可被消费。
9. 将所有失败回写 xlib-standard 的 Prompt/Harness/Rule/Generator patch。
10. v0.1.0 只交付 PostgreSQL adapter MVA；ORM、migration、outbox、读写分离全部进入 P1/P2 ADR。
```

最终目标状态：

```text
postgresx 不再是空骨架或零散 pgx 包装。
postgresx 成为 xlib-standard 控制下、可测试、可观测、可发布证明、可下游采纳、可持续自我强化的 L2 PostgreSQL 基础设施适配产品。
```

---

## 31. 最终完成标准

真正完成时，必须能够写出：

```text
DONE with evidence:
- postgresx identity corrected.
- no x.go or other L2 reverse dependency.
- config/health/metrics/errors/API contracts passed.
- unit/fake tests passed without PostgreSQL.
- tx contract tests passed.
- pg SQLSTATE error mapping tests passed.
- DSN/password/SQL args redaction tests passed.
- boundary/security/docs/contracts gates passed.
- integration opt-in passed or explicit skip evidence recorded.
- release manifest generated and checksum verified.
- release-final-check passed.
- downstream smoke adoption proof exists.
- retrospective patches generated.
```

只要其中任一项缺少证据，只能写：

```text
NOT DONE yet:
- blocked_by:
- missing_evidence:
- next_action:
```
