# postgresx 完整可执行 Goal Prompt v1.1

> 文件名：`postgresx_goal_executable_prompt_v1_1.md`  
> 版本定位：Template-bound + Foundation-bound 实战执行版  
> 目标模块：`github.com/ZoneCNH/postgresx`  
> 模块定位：PostgreSQL 独立公共基础开发库 / L2 基础设施适配层  
> 上游模板：`github.com/ZoneCNH/baselib-template`  
> 上游契约：`github.com/ZoneCNH/foundationx`  
> 上游包路径：`github.com/ZoneCNH/foundationx/pkg/foundationx`  
> 适用项目：x.go、Market Data、Macro Data、Regime Engine、Trading Server、基础库体系  
> 执行方法：Goal Runtime Prompt v3.1 + baselib-template + foundationx + Harness + AutoResearch + Self-improving + Evidence Protocol  
> 生成日期：2026-06-01  
> 时区：Asia/Tokyo  

---

# 0. v1.1 更新目的

`postgresx_goal_executable_prompt_v1_0.md` 是独立设计版。

本文件是 v1.1，必须绑定当前已经完成的两个事实前置：

```text
1. github.com/ZoneCNH/baselib-template 已完成，作为基础库模板事实标准。
2. github.com/ZoneCNH/foundationx 已完成，作为 L0 契约事实标准。
```

因此 v1.1 不再从零手工创建全部目录，而是必须优先从 `baselib-template` 渲染 `postgresx` 骨架，再在该骨架上实现 PostgreSQL 专属能力。

---

# 1. 事实锚点

## 1.1 baselib-template 事实锚点

当前模板事实标准：

```text
Repository: github.com/ZoneCNH/baselib-template
Module:     github.com/ZoneCNH/baselib-template
Go:         1.23
```

`baselib-template` 提供：

```text
pkg/{{PACKAGE_NAME}}
internal/
testkit/
examples/
contracts/
docs/
scripts/
.agent/
release/manifest/
Makefile
CI/Harness gate
Release Evidence
```

生成具体基础库时，必须优先使用：

```bash
scripts/render_template.sh \
  --module-name postgresx \
  --module-path github.com/ZoneCNH/postgresx \
  --package-name postgresx \
  --out ../postgresx
```

## 1.2 foundationx 事实锚点

当前 foundationx 事实标准：

```text
Repository: github.com/ZoneCNH/foundationx
Module:     github.com/ZoneCNH/foundationx
Package:    github.com/ZoneCNH/foundationx/pkg/foundationx
Go:         1.23
Layer:      L0 基础设施契约层
```

foundationx 已提供：

```text
ErrorKind
Error
NewError
WrapError
HealthStatus
HealthChecker
Lifecycle
RetryPolicy
Sanitizer
SecretString
Clock
VersionInfo
```

postgresx 必须复用 foundationx，不得重新发明这些 L0 契约。

---

# 2. 使用方式

将本文完整交给 Agent Teams / Codex / Claude Code / Cursor Agent / GitHub Copilot Workspace 执行。

执行前必须确认：

```text
1. 当前目标是创建或完善独立 Go module：github.com/ZoneCNH/postgresx
2. postgresx 是 PostgreSQL L2 基础设施适配库
3. postgresx 必须从 github.com/ZoneCNH/baselib-template 渲染骨架
4. postgresx 必须依赖 github.com/ZoneCNH/foundationx
5. postgresx 必须 import github.com/ZoneCNH/foundationx/pkg/foundationx
6. postgresx 可以依赖 PostgreSQL driver，例如 pgx/pgxpool
7. postgresx 不允许依赖 x.go
8. postgresx 不允许包含 x.go 业务表结构、业务模型、业务 topic、业务 key
9. postgresx 不允许隐式读取 /home/k8s/secrets/env/*
10. postgresx 不允许自动读取 .env、production.yaml、config.local.yaml
11. postgresx 不允许强依赖 configx；configx 只能在 app bootstrap 或文档示例中作为可选组合
12. postgresx 不允许定义全局 DB / 默认 Client / 单例连接池
13. postgresx 不允许在日志、错误、Evidence、Release Manifest 中输出明文密码或 DSN
14. 所有独立性验证必须包含 GOWORK=off
15. 所有完成声明必须使用 DONE with evidence:
```

---

# 3. Master Goal

```text
GOAL-20260601-POSTGRESX-001

基于已完成的 github.com/ZoneCNH/baselib-template 和 github.com/ZoneCNH/foundationx，建立 github.com/ZoneCNH/postgresx 独立 PostgreSQL 公共基础开发库，为 x.go 及其未来服务提供可复用、可测试、可观测、可发布、可审计的 PostgreSQL 访问基础能力。

postgresx 必须封装 PostgreSQL 连接池、配置校验、DSN 脱敏、Ping、HealthCheck、事务执行、Migration Runner、错误归一化、Metrics Hook、TestKit、Examples、CI/Harness/Evidence/Release 流程。

postgresx 必须作为独立 Go module 发布，不得依赖 x.go，不得包含 Market Data / Macro Data / Regime / Trading 等业务语义，不得内置任何 x.go 业务 schema。
```

---

# 4. v1.1 与 v1.0 的关键差异

```text
1. module path 从 github.com/bytechainx/postgresx 改为 github.com/ZoneCNH/postgresx
2. foundationx import 从 github.com/bytechainx/foundationx 改为 github.com/ZoneCNH/foundationx/pkg/foundationx
3. 不再手工创建 skeleton，必须通过 baselib-template/scripts/render_template.sh 渲染
4. foundationx 不再是假设依赖，而是已完成的 L0 事实依赖
5. 增加 GOWORK=off 独立模块验证
6. 增加 make ci-extended
7. 增加 make release-preflight VERSION=v0.1.0
8. 增加 make release-evidence-check
9. 增加 make release-final-check
10. 增加 lint / govulncheck 必需 gate，不得伪造 skipped
11. release/manifest/latest.json 是生成 artifact，不提交源码历史
12. 增加 baselib-template contract alignment gate
13. 增加 foundationx API compatibility gate
14. 增加 configx boundary：postgresx core 不依赖 configx，但 docs 展示组合方式
```

---

# 5. 问题底层本质

postgresx 不是“把 pgxpool 包一层”。

postgresx 的底层本质是：

```text
把 PostgreSQL 作为基础设施能力，抽象成稳定、可治理、可复用、可验证、可发布的工程资产。
```

它解决的是以下结构性问题：

```text
1. 避免 x.go 各服务重复创建 PostgreSQL 连接池
2. 避免各模块各自拼接 DSN 并泄露密码
3. 避免事务处理、rollback、panic、context timeout 逻辑散落在业务代码中
4. 避免 migration runner、health check、metrics 口径不统一
5. 避免测试依赖本地手动 PostgreSQL
6. 避免基础设施错误直接暴露 driver 细节给业务层
7. 避免基础库绕过 baselib-template 的统一 Evidence 和 Release Gate
8. 为 Agent Teams 提供可执行、可验证、可发布的基础库任务边界
```

核心价值：

```text
PostgreSQL 能力标准化 + foundationx 契约复用 + baselib-template Gate 继承 + Release Evidence 可证明。
```

---

# 6. 不可再拆解的基本真理

## 6.1 postgresx 是 L2，不是 L0

允许依赖：

```text
github.com/ZoneCNH/foundationx
PostgreSQL driver
测试辅助库
```

禁止依赖：

```text
x.go
market_data
macro_data
regime_engine
trading_server
业务 schema
业务 repository
```

## 6.2 postgresx 只理解 PostgreSQL，不理解业务

postgresx 可以知道：

```text
Connection
Pool
Transaction
Migration
Query
Exec
Ping
Health
Timeout
Retryable Error
```

postgresx 不应该知道：

```text
Kline
BTCUSDT
MacroRegime
M1-M7
S1-S7
TradingSignal
Order
Position
RiskGate
```

## 6.3 postgresx 只接收 Config，不加载 Config

postgresx 不允许自动读取：

```text
/home/k8s/secrets/env/postgres.env
.env
config.local.yaml
production.yaml
```

x.go 或 app bootstrap 可以使用 configx 显式读取配置，然后构造 `postgresx.Config`。

正确链路：

```text
x.go bootstrap -> configx optional -> postgresx.Config -> postgresx.New(...)
```

错误链路：

```text
postgresx -> configx -> /home/k8s/secrets/env/postgres.env
```

## 6.4 没有 Evidence 不得声称完成

完成声明必须是：

```text
DONE with evidence:
- GOWORK=off go test ./...
- GOWORK=off go test -race ./...
- GOWORK=off make ci
- GOWORK=off make ci-extended
- GOWORK=off make release-check
- GOWORK=off make release-preflight VERSION=v0.1.0
- GOWORK=off make release-evidence-check
- GOWORK=off make release-final-check
- boundary gate passed
- secret gate passed
- integration PostgreSQL test passed
- release manifest generated
```

---

# 7. Scope

## 7.1 In Scope

```text
Template rendering from baselib-template
go.mod module github.com/ZoneCNH/postgresx
foundationx integration
Config / Validate / Sanitize
DSN builder / RedactedDSN
pgxpool-based Client
Ping
Close idempotency
Pool stats
Exec / Query / QueryRow wrappers
DBTX / Queryer interfaces
WithTx transaction helper
Transaction commit / rollback semantics
Migration Runner
schema_migrations table
HealthCheck using foundationx.HealthStatus
Error mapping to foundationx.Error
Retryable classification
Metrics hooks
TestKit
Docker / env-driven integration tests
Examples
CI / Harness / Release Manifest
GOWORK=off verification
Template contract alignment
Foundationx API compatibility
Configx boundary documentation
```

## 7.2 Out of Scope

```text
x.go business schema
market_data tables
macro_data tables
regime tables
trading tables
SQLC generated code
Business repository
ORM abstraction
Read/write splitting
Distributed transaction
PostgreSQL HA orchestration
Connection proxy
Secret manager implementation
Automatic production env loading
Mandatory configx dependency
```

## 7.3 Optional / Deferred

```text
Prepared statement cache policy
Read/write split
Advisory lock
COPY bulk import
LISTEN/NOTIFY
Replica lag check
observex direct dependency
configx compiled example
```

默认裁决：

```text
v0.1.0 不做 read/write split、不做 COPY、不做 LISTEN/NOTIFY、不强依赖 configx 或 observex。
```

---

# 8. 目标仓库与模块

## 8.1 Repository

```text
github.com/ZoneCNH/postgresx
```

## 8.2 go.mod

```go
module github.com/ZoneCNH/postgresx

go 1.23
```

## 8.3 必需依赖

```text
github.com/ZoneCNH/foundationx
```

## 8.4 PostgreSQL driver 推荐

推荐：

```text
github.com/jackc/pgx/v5
github.com/jackc/pgx/v5/pgxpool
```

要求：

```text
具体版本必须由执行时 AutoResearch 确认，并写入 docs/adr/ADR-20260601-001-driver-selection.md。
```

## 8.5 configx 边界

postgresx core 不依赖：

```text
github.com/ZoneCNH/configx
```

docs 可以说明组合方式：

```text
x.go bootstrap -> configx.LoadEnvFile -> configx.Decode -> postgresx.Config -> postgresx.New
```

如果未来加入 `examples/with_configx`：

```text
1. 必须放入独立文档或单独 module
2. 不得让 postgresx core go.mod 必需依赖 configx
3. 不得破坏 boundary gate
```

---

# 9. Template Render Phase

## 9.1 生成方式

执行者必须先获取或引用已完成模板：

```bash
git clone https://github.com/ZoneCNH/baselib-template.git
cd baselib-template
```

渲染：

```bash
scripts/render_template.sh \
  --module-name postgresx \
  --module-path github.com/ZoneCNH/postgresx \
  --package-name postgresx \
  --out ../postgresx
```

进入目标仓库：

```bash
cd ../postgresx
```

验证初始模板：

```bash
GOWORK=off go mod tidy
GOWORK=off go test ./...
GOWORK=off make ci
```

## 9.2 渲染后必须替换的模板语义

必须确认：

```text
{{MODULE_NAME}} 已替换为 postgresx
{{MODULE_PATH}} 已替换为 github.com/ZoneCNH/postgresx
{{PACKAGE_NAME}} 已替换为 postgresx
pkg/{{PACKAGE_NAME}} 已变为 pkg/postgresx
README / docs / contracts / scripts / .agent 中不再出现未渲染占位符
```

## 9.3 Template Alignment Gate

新增 Gate：

```bash
grep -R "{{MODULE_NAME}}\|{{MODULE_PATH}}\|{{PACKAGE_NAME}}" . \
  --exclude-dir=.git \
  --exclude-dir=vendor && exit 1 || true
```

并检查：

```text
Makefile 继承模板 gate
scripts/check_boundary.sh 存在
scripts/check_secrets.sh 存在
scripts/check_contracts.sh 存在
scripts/generate_manifest.sh 存在
.agent/ 目录存在
contracts/ 目录存在
release/manifest/ 目录存在
```

---

# 10. 标准目录结构

渲染后目标结构应为：

```text
postgresx/
├── go.mod
├── go.sum
├── README.md
├── CHANGELOG.md
├── LICENSE
├── Makefile
├── .gitignore
├── .golangci.yml
├── pkg/
│   └── postgresx/
│       ├── doc.go
│       ├── config.go
│       ├── dsn.go
│       ├── client.go
│       ├── pool.go
│       ├── query.go
│       ├── tx.go
│       ├── migration.go
│       ├── health.go
│       ├── metrics.go
│       ├── errors.go
│       ├── options.go
│       ├── version.go
│       └── *_test.go
├── internal/
│   ├── sanitize/
│   ├── validation/
│   └── pgtest/
├── testkit/
│   ├── container.go
│   ├── config.go
│   ├── fixture.go
│   └── assert.go
├── examples/
│   ├── basic/
│   ├── transaction/
│   ├── migration/
│   └── health/
├── contracts/
│   ├── config.schema.json
│   ├── error.schema.json
│   ├── health.schema.json
│   ├── metrics.md
│   └── public_api.md
├── docs/
│   ├── spec.md
│   ├── design.md
│   ├── api.md
│   ├── config.md
│   ├── dsn.md
│   ├── transactions.md
│   ├── migrations.md
│   ├── errors.md
│   ├── health.md
│   ├── observability.md
│   ├── testing.md
│   ├── xgo-integration.md
│   ├── configx-boundary.md
│   ├── release.md
│   └── adr/
│       ├── ADR-20260601-001-driver-selection.md
│       ├── ADR-20260601-002-transaction-semantics.md
│       ├── ADR-20260601-003-migration-runner-scope.md
│       ├── ADR-20260601-004-template-bound-skeleton.md
│       └── ADR-20260601-005-config-loading-boundary.md
├── scripts/
│   ├── check_boundary.sh
│   ├── check_secrets.sh
│   ├── check_contracts.sh
│   ├── check_template_alignment.sh
│   ├── check_foundationx_api.sh
│   ├── run_integration.sh
│   └── generate_manifest.sh
├── release/
│   └── manifest/
│       ├── .gitkeep
│       └── latest.json        # generated artifact, normally not committed
├── .github/
│   └── workflows/
│       ├── ci.yml
│       ├── integration.yml
│       ├── security.yml
│       └── release.yml
└── .agent/
    ├── goal.md
    ├── spec.md
    ├── design.md
    ├── plan.md
    ├── tasks.md
    ├── harness.md
    ├── gates.md
    ├── evidence.md
    ├── review.md
    ├── release.md
    └── retrospective.md
```

说明：

```text
release/manifest/latest.json 与 release/manifest/v*.json 是 release gate 生成物。
是否提交版本化 manifest 由仓库 release 规范裁决；不得手工伪造。
```

---

# 11. Public API 设计

## 11.1 Config

文件：

```text
pkg/postgresx/config.go
```

目标 API：

```go
package postgresx

import (
    "time"

    "github.com/ZoneCNH/foundationx/pkg/foundationx"
)

type Config struct {
    Host            string
    Port            int
    Database        string
    User            string
    Password        foundationx.SecretString
    SSLMode         string
    MaxOpenConns    int32
    MinIdleConns    int32
    MaxConnLifetime time.Duration
    MaxConnIdleTime  time.Duration
    ConnectTimeout  time.Duration
    HealthTimeout   time.Duration
    ApplicationName string
}

type SanitizedConfig struct {
    Host            string
    Port            int
    Database        string
    User            string
    Password        string
    SSLMode         string
    MaxOpenConns    int32
    MinIdleConns    int32
    MaxConnLifetime string
    MaxConnIdleTime  string
    ConnectTimeout  string
    HealthTimeout   string
    ApplicationName string
}

func DefaultConfig() Config
func (c Config) Validate() error
func (c Config) Sanitize() SanitizedConfig
```

要求：

```text
1. zero-value Config 必须 Validate 失败
2. DefaultConfig 只给非敏感默认值
3. Password 必须使用 foundationx.SecretString
4. Validate 返回 foundationx.Error
5. Sanitize 绝不返回明文密码
6. 不允许 Config 自行读取 env
7. 不允许 Config 自行读取 /home/k8s/secrets/env/*
```

## 11.2 DSN Builder

文件：

```text
pkg/postgresx/dsn.go
```

目标 API：

```go
func (c Config) DSN() string
func (c Config) RedactedDSN() string
```

要求：

```text
DSN 用于 driver 连接，可以包含密码
RedactedDSN 用于日志、错误、Evidence，绝不包含密码
```

测试必须覆盖：

```text
DSN 包含必要字段
RedactedDSN 不包含原始密码
特殊字符经过正确转义
```

## 11.3 Options

文件：

```text
pkg/postgresx/options.go
```

目标 API：

```go
type Option func(*options)

type options struct {
    logger  Logger
    metrics Metrics
    clock   foundationx.Clock
}

type Logger interface {
    Debug(ctx context.Context, msg string, fields ...Field)
    Info(ctx context.Context, msg string, fields ...Field)
    Warn(ctx context.Context, msg string, fields ...Field)
    Error(ctx context.Context, msg string, fields ...Field)
}

type Field struct {
    Key   string
    Value any
}

type Metrics interface {
    IncCounter(name string, labels map[string]string)
    ObserveHistogram(name string, value float64, labels map[string]string)
    SetGauge(name string, value float64, labels map[string]string)
}

func WithLogger(logger Logger) Option
func WithMetrics(metrics Metrics) Option
func WithClock(clock foundationx.Clock) Option
```

v1.1 裁决：

```text
1. postgresx core 不强依赖 observex。
2. observex 的具体 Logger/Metrics 可以在 x.go bootstrap 中注入，只要满足接口。
3. 默认 logger / metrics 必须是 noop。
4. nil logger / metrics 不得 panic。
```

## 11.4 Client

文件：

```text
pkg/postgresx/client.go
```

目标 API：

```go
type Client struct {
    // private fields
}

func New(ctx context.Context, cfg Config, opts ...Option) (*Client, error)
func (c *Client) Ping(ctx context.Context) error
func (c *Client) Close(ctx context.Context) error
func (c *Client) Stats() PoolStats
func (c *Client) Queryer() Queryer
```

要求：

```text
1. New 必须先 Validate Config
2. New 必须使用 context timeout
3. New 失败必须返回 foundationx.Error
4. Close 必须幂等
5. Close 不得 panic
6. Stats 不暴露 driver 原始类型
7. Queryer 返回最小查询接口
```

## 11.5 Query Interfaces

文件：

```text
pkg/postgresx/query.go
```

目标 API：

```go
type CommandTag interface {
    RowsAffected() int64
}

type Row interface {
    Scan(dest ...any) error
}

type Rows interface {
    Close()
    Err() error
    Next() bool
    Scan(dest ...any) error
}

type Queryer interface {
    Exec(ctx context.Context, sql string, args ...any) (CommandTag, error)
    Query(ctx context.Context, sql string, args ...any) (Rows, error)
    QueryRow(ctx context.Context, sql string, args ...any) Row
}
```

要求：

```text
1. 业务 repository 只依赖 Queryer
2. Queryer 可以由 Client 或 Tx 实现
3. 不强制业务导入 pgx 类型
4. 底层适配 pgx.CommandTag / pgx.Rows / pgx.Row
```

## 11.6 Transaction

文件：

```text
pkg/postgresx/tx.go
```

目标 API：

```go
type Tx interface {
    Queryer
}

type TxFunc func(ctx context.Context, tx Tx) error

type TxOptions struct {
    IsolationLevel string
    ReadOnly       bool
}

func (c *Client) WithTx(ctx context.Context, fn TxFunc) error
func (c *Client) WithTxOptions(ctx context.Context, opts TxOptions, fn TxFunc) error
```

事务语义：

```text
1. Begin 失败返回 error
2. fn 返回 nil -> commit
3. fn 返回 error -> rollback
4. rollback 失败需要被记录，但主错误以 fn error 为主
5. commit 失败返回 commit error
6. fn panic -> rollback，然后继续 panic
7. context cancel -> rollback
8. 不做嵌套事务自动识别
```

Panic 策略必须写入：

```text
docs/adr/ADR-20260601-002-transaction-semantics.md
```

## 11.7 Migration Runner

文件：

```text
pkg/postgresx/migration.go
```

目标 API：

```go
type Migration struct {
    Version int64
    Name    string
    UpSQL   string
    DownSQL string
}

type MigrationSource interface {
    List(ctx context.Context) ([]Migration, error)
}

type MigrationRunner struct {
    client *Client
}

func NewMigrationRunner(client *Client) *MigrationRunner
func (r *MigrationRunner) Up(ctx context.Context, source MigrationSource) error
func (r *MigrationRunner) Applied(ctx context.Context) ([]AppliedMigration, error)

type AppliedMigration struct {
    Version   int64
    Name      string
    AppliedAt time.Time
}
```

Schema migration table：

```sql
CREATE TABLE IF NOT EXISTS schema_migrations (
    version BIGINT PRIMARY KEY,
    name TEXT NOT NULL,
    applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

要求：

```text
1. postgresx 不拥有业务 migration 文件
2. migration 文件由调用方传入
3. migration 必须按 version 升序执行
4. 每个 migration 在事务中执行
5. 已执行 version 不重复执行
6. 同 version 不同 name 必须报 conflict
7. migration runner 必须有 integration test
```

v0.1.0 非目标：

```text
down migration 执行
checksum 管理
复杂 migration DSL
自动扫描业务目录
```

## 11.8 Health Check

文件：

```text
pkg/postgresx/health.go
```

目标 API：

```go
func (c *Client) Name() string
func (c *Client) Check(ctx context.Context) foundationx.HealthStatus
```

要求：

```text
1. Client 实现 foundationx.HealthChecker
2. Check 使用 HealthTimeout
3. 成功返回 healthy
4. 超时返回 degraded 或 unhealthy，策略文档化
5. 连接不可用返回 unhealthy
6. Metadata 包含 host/database/pool stats，但不得包含 password/dsn 明文
```

## 11.9 PoolStats

文件：

```text
pkg/postgresx/pool.go
```

目标 API：

```go
type PoolStats struct {
    TotalConns        int32
    IdleConns         int32
    AcquiredConns     int32
    ConstructingConns int32
    MaxConns          int32
}
```

要求：

```text
1. 不直接暴露 pgxpool.Stat
2. 字段命名保持稳定
3. metrics 使用该结构
```

## 11.10 Error Mapping

文件：

```text
pkg/postgresx/errors.go
```

目标 API：

```go
func MapError(op string, err error) error
func IsRetryable(err error) bool
```

映射原则：

```text
context.Canceled -> foundationx.ErrorKindCanceled
context.DeadlineExceeded -> foundationx.ErrorKindTimeout
认证失败 -> foundationx.ErrorKindAuth
连接失败 -> foundationx.ErrorKindConnection
唯一键冲突 -> foundationx.ErrorKindConflict
not found / no rows -> foundationx.ErrorKindNotFound
其他 driver 错误 -> foundationx.ErrorKindInternal
```

要求：

```text
1. 返回 foundationx.Error
2. 保留 Cause
3. 设置 Retryable
4. 不泄露 DSN 和密码
```

---

# 12. Spec

```text
SPEC-postgresx-v1.1
```

## REQ-POSTGRESX-001：Template-bound module

Acceptance Criteria：

```text
AC-REQ-POSTGRESX-001-001: postgresx 由 baselib-template/scripts/render_template.sh 生成
AC-REQ-POSTGRESX-001-002: go.mod module 为 github.com/ZoneCNH/postgresx
AC-REQ-POSTGRESX-001-003: 不存在 {{MODULE_NAME}} / {{MODULE_PATH}} / {{PACKAGE_NAME}} 占位符
AC-REQ-POSTGRESX-001-004: pkg/postgresx 存在
AC-REQ-POSTGRESX-001-005: GOWORK=off go test ./... 通过
```

## REQ-POSTGRESX-002：Foundationx compatibility

Acceptance Criteria：

```text
AC-REQ-POSTGRESX-002-001: go.mod 依赖 github.com/ZoneCNH/foundationx
AC-REQ-POSTGRESX-002-002: import 使用 github.com/ZoneCNH/foundationx/pkg/foundationx
AC-REQ-POSTGRESX-002-003: Config.Password 使用 foundationx.SecretString
AC-REQ-POSTGRESX-002-004: errors 使用 foundationx.Error / ErrorKind
AC-REQ-POSTGRESX-002-005: HealthCheck 返回 foundationx.HealthStatus
AC-REQ-POSTGRESX-002-006: Client 实现 foundationx.HealthChecker
```

## REQ-POSTGRESX-003：依赖边界

Acceptance Criteria：

```text
AC-REQ-POSTGRESX-003-001: 允许依赖 foundationx
AC-REQ-POSTGRESX-003-002: 允许依赖 pgx/pgxpool
AC-REQ-POSTGRESX-003-003: 不允许依赖 x.go
AC-REQ-POSTGRESX-003-004: 不允许出现 Market Data / Macro Data / Regime / Trading 业务模型
AC-REQ-POSTGRESX-003-005: postgresx core 不依赖 configx
AC-REQ-POSTGRESX-003-006: postgresx core 不依赖 observex
```

## REQ-POSTGRESX-004：Config

Acceptance Criteria：

```text
AC-REQ-POSTGRESX-004-001: Config 包含 Host/Port/Database/User/Password/SSLMode/Pool/Timeout 字段
AC-REQ-POSTGRESX-004-002: Config.Validate 覆盖必填字段
AC-REQ-POSTGRESX-004-003: Config.Sanitize 不泄露密码
AC-REQ-POSTGRESX-004-004: DefaultConfig 不包含敏感值
AC-REQ-POSTGRESX-004-005: zero-value Config Validate 失败
AC-REQ-POSTGRESX-004-006: Config 不读取 env/file/secret path
```

## REQ-POSTGRESX-005：DSN

Acceptance Criteria：

```text
AC-REQ-POSTGRESX-005-001: DSN 可用于连接 PostgreSQL
AC-REQ-POSTGRESX-005-002: RedactedDSN 不包含原始密码
AC-REQ-POSTGRESX-005-003: 特殊字符正确转义
AC-REQ-POSTGRESX-005-004: DSN 不被日志默认输出
```

## REQ-POSTGRESX-006：Client / Pool

Acceptance Criteria：

```text
AC-REQ-POSTGRESX-006-001: New(ctx, cfg, opts...) 创建连接池
AC-REQ-POSTGRESX-006-002: New 失败返回 foundationx.Error
AC-REQ-POSTGRESX-006-003: Ping 可验证连接
AC-REQ-POSTGRESX-006-004: Close 幂等
AC-REQ-POSTGRESX-006-005: Stats 返回 PoolStats
AC-REQ-POSTGRESX-006-006: GOWORK=off go test -race ./... 通过
```

## REQ-POSTGRESX-007：Queryer

Acceptance Criteria：

```text
AC-REQ-POSTGRESX-007-001: 定义 Queryer interface
AC-REQ-POSTGRESX-007-002: Client 实现 Queryer
AC-REQ-POSTGRESX-007-003: Tx 实现 Queryer
AC-REQ-POSTGRESX-007-004: 业务 repository 可以只依赖 Queryer
```

## REQ-POSTGRESX-008：Transaction

Acceptance Criteria：

```text
AC-REQ-POSTGRESX-008-001: WithTx 支持成功 commit
AC-REQ-POSTGRESX-008-002: WithTx 在 fn error 时 rollback
AC-REQ-POSTGRESX-008-003: WithTx 在 panic 时 rollback
AC-REQ-POSTGRESX-008-004: WithTx 在 context cancel 时 rollback
AC-REQ-POSTGRESX-008-005: commit 失败返回错误
AC-REQ-POSTGRESX-008-006: transaction semantics 写入 ADR
```

## REQ-POSTGRESX-009：Migration Runner

Acceptance Criteria：

```text
AC-REQ-POSTGRESX-009-001: 创建 schema_migrations 表
AC-REQ-POSTGRESX-009-002: 按 version 升序执行 migration
AC-REQ-POSTGRESX-009-003: 已执行 migration 不重复执行
AC-REQ-POSTGRESX-009-004: 同 version 不同 name 返回 conflict
AC-REQ-POSTGRESX-009-005: migration 在事务中执行
AC-REQ-POSTGRESX-009-006: integration test 覆盖
```

## REQ-POSTGRESX-010：Health

Acceptance Criteria：

```text
AC-REQ-POSTGRESX-010-001: Client 实现 foundationx.HealthChecker
AC-REQ-POSTGRESX-010-002: Check 返回 healthy/degraded/unhealthy
AC-REQ-POSTGRESX-010-003: Check 包含 latency
AC-REQ-POSTGRESX-010-004: Metadata 不包含密码/DSN
AC-REQ-POSTGRESX-010-005: 连接失败时 unhealthy
```

## REQ-POSTGRESX-011：Error Mapping

Acceptance Criteria：

```text
AC-REQ-POSTGRESX-011-001: context.Canceled 映射 canceled
AC-REQ-POSTGRESX-011-002: context.DeadlineExceeded 映射 timeout
AC-REQ-POSTGRESX-011-003: no rows 映射 not_found
AC-REQ-POSTGRESX-011-004: unique violation 映射 conflict
AC-REQ-POSTGRESX-011-005: auth error 映射 auth
AC-REQ-POSTGRESX-011-006: connection error 映射 connection/unavailable
AC-REQ-POSTGRESX-011-007: 保留 Cause
```

## REQ-POSTGRESX-012：Configx Boundary

Acceptance Criteria：

```text
AC-REQ-POSTGRESX-012-001: postgresx core 不 import configx
AC-REQ-POSTGRESX-012-002: docs/configx-boundary.md 说明 configx 只在 app bootstrap 使用
AC-REQ-POSTGRESX-012-003: docs/xgo-integration.md 展示 configx + postgresx 组合
AC-REQ-POSTGRESX-012-004: examples 不引入 configx 编译依赖，除非独立 module 或 build tag 明确隔离
```

## REQ-POSTGRESX-013：Harness / Release

Acceptance Criteria：

```text
AC-REQ-POSTGRESX-013-001: GOWORK=off make ci 通过
AC-REQ-POSTGRESX-013-002: GOWORK=off make ci-extended 通过
AC-REQ-POSTGRESX-013-003: GOWORK=off make release-check 通过
AC-REQ-POSTGRESX-013-004: GOWORK=off make release-preflight VERSION=v0.1.0 通过
AC-REQ-POSTGRESX-013-005: GOWORK=off make release-evidence-check 通过
AC-REQ-POSTGRESX-013-006: GOWORK=off make release-final-check 通过
AC-REQ-POSTGRESX-013-007: lint/security 工具缺失不得伪造 passed
AC-REQ-POSTGRESX-013-008: release manifest 生成且不含 secret
```

---

# 13. Plan

```text
PLAN-GOAL-20260601-POSTGRESX-001-v1.1
```

## Phase 0：Context Recovery

目标：

```text
确认 baselib-template、foundationx、postgresx、configx 边界和当前事实状态。
```

输出：

```text
.agent/context.md
```

必须记录：

```text
baselib-template 已完成
foundationx 已完成
postgresx 从模板渲染
foundationx 是 L0
postgresx 是 L2
configx 是 L1 配置加载库，但 postgresx core 不依赖 configx
x.go 读取 /home/k8s/secrets/env/* 后显式构造 postgresx.Config
```

## Phase 1：Render from baselib-template

目标：

```text
使用 scripts/render_template.sh 渲染 postgresx。
```

输出：

```text
github.com/ZoneCNH/postgresx skeleton
```

## Phase 2：Template Alignment

目标：

```text
确认模板占位符、路径、package、Makefile、scripts、contracts、.agent 完整。
```

## Phase 3：Foundationx Integration

目标：

```text
接入 github.com/ZoneCNH/foundationx/pkg/foundationx。
```

输出：

```text
go.mod
foundationx API compatibility tests
```

## Phase 4：PostgreSQL Driver ADR

目标：

```text
AutoResearch pgx/pgxpool 当前版本与语义，写入 ADR。
```

输出：

```text
docs/adr/ADR-20260601-001-driver-selection.md
```

## Phase 5：Core Config + DSN

目标：

```text
实现 Config / Validate / Sanitize / DSN / RedactedDSN。
```

## Phase 6：Client + Pool + Queryer

目标：

```text
实现 New / Ping / Close / Stats / Queryer。
```

## Phase 7：Transaction

目标：

```text
实现 WithTx / WithTxOptions / Tx interface。
```

## Phase 8：Migration Runner

目标：

```text
实现 schema_migrations 和 Up。
```

## Phase 9：Health + Error + Metrics

目标：

```text
实现 HealthCheck、driver error mapping、metrics hooks。
```

## Phase 10：TestKit + Integration

目标：

```text
实现可重复运行的 PostgreSQL 集成测试。
```

## Phase 11：Docs + ADR + Configx Boundary

目标：

```text
补齐 README、docs、ADR、contracts、configx-boundary、xgo-integration。
```

## Phase 12：Harness + Release

目标：

```text
GOWORK=off make ci
GOWORK=off make ci-extended
GOWORK=off make release-check
GOWORK=off make release-preflight VERSION=v0.1.0
GOWORK=off make release-evidence-check
GOWORK=off make release-final-check
```

## Phase 13：Retrospective

目标：

```text
输出 self-improving patch，供 redisx/kafkax/taosx/configx/observex 复用。
```

---

# 14. Task Breakdown

## TASK-POSTGRESX-001：从 baselib-template 渲染 postgresx

操作：

```bash
git clone https://github.com/ZoneCNH/baselib-template.git
cd baselib-template

scripts/render_template.sh \
  --module-name postgresx \
  --module-path github.com/ZoneCNH/postgresx \
  --package-name postgresx \
  --out ../postgresx

cd ../postgresx
```

验收：

```text
go.mod 存在
module 为 github.com/ZoneCNH/postgresx
pkg/postgresx 存在
README/docs/contracts/scripts/.agent 存在
不存在未渲染占位符
```

证据：

```text
EVID-TASK-POSTGRESX-001-20260601-001: tree output
EVID-TASK-POSTGRESX-001-20260601-002: go env GOMOD
EVID-TASK-POSTGRESX-001-20260601-003: check_template_alignment output
```

## TASK-POSTGRESX-002：初始模板独立验证

命令：

```bash
GOWORK=off go mod tidy
GOWORK=off go test ./...
GOWORK=off make ci
```

证据：

```text
EVID-TASK-POSTGRESX-002-20260601-001: GOWORK=off go test output
EVID-TASK-POSTGRESX-002-20260601-002: GOWORK=off make ci output
```

## TASK-POSTGRESX-003：接入 foundationx

操作：

```bash
GOWORK=off go get github.com/ZoneCNH/foundationx
```

要求：

```text
所有 import 使用 github.com/ZoneCNH/foundationx/pkg/foundationx
不得使用 bytechainx 路径
不得复制 foundationx 代码
```

证据：

```text
EVID-TASK-POSTGRESX-003-20260601-001: go.mod diff
EVID-TASK-POSTGRESX-003-20260601-002: foundationx API compatibility test output
```

## TASK-POSTGRESX-004：AutoResearch pgx/pgxpool 并写 ADR

操作：

```text
确认 pgx/v5 当前稳定版本
确认 pgxpool Config / Stat / error behavior
确认 PostgreSQL SQLSTATE 错误码处理方式
写入 ADR-20260601-001-driver-selection.md
```

证据：

```text
EVID-TASK-POSTGRESX-004-20260601-001: ADR driver selection
EVID-TASK-POSTGRESX-004-20260601-002: go.mod diff
```

## TASK-POSTGRESX-005：实现 Config

文件：

```text
pkg/postgresx/config.go
pkg/postgresx/config_test.go
```

实现：

```text
Config
SanitizedConfig
DefaultConfig
Validate
Sanitize
```

测试：

```text
TestDefaultConfig
TestConfigValidateZeroValueFails
TestConfigValidateMissingHost
TestConfigValidateInvalidPort
TestConfigValidateMissingDatabase
TestConfigValidateMissingUser
TestConfigValidateMissingPassword
TestConfigSanitizeMasksPassword
TestConfigDoesNotReadEnv
```

## TASK-POSTGRESX-006：实现 DSN / RedactedDSN

文件：

```text
pkg/postgresx/dsn.go
pkg/postgresx/dsn_test.go
```

测试：

```text
TestDSNBuildsPostgresURL
TestDSNEscapesSpecialChars
TestRedactedDSNMasksPassword
TestRedactedDSNDoesNotContainRawPassword
```

## TASK-POSTGRESX-007：实现 Options / Noop Logger / Noop Metrics

文件：

```text
pkg/postgresx/options.go
pkg/postgresx/metrics.go
```

测试：

```text
TestOptionsDefaultNoop
TestWithLogger
TestWithMetrics
TestWithClock
TestNilLoggerDoesNotPanic
```

## TASK-POSTGRESX-008：实现 Client / Pool

文件：

```text
pkg/postgresx/client.go
pkg/postgresx/pool.go
pkg/postgresx/client_test.go
```

测试：

```text
TestNewInvalidConfigFails
TestCloseIdempotent
TestStatsZeroSafe
TestNewPingCloseIntegration
```

## TASK-POSTGRESX-009：实现 Queryer 适配

文件：

```text
pkg/postgresx/query.go
pkg/postgresx/query_test.go
```

测试：

```text
TestClientImplementsQueryer
TestQueryerExecIntegration
TestQueryerQueryIntegration
TestQueryerQueryRowIntegration
```

## TASK-POSTGRESX-010：实现 Transaction

文件：

```text
pkg/postgresx/tx.go
pkg/postgresx/tx_test.go
docs/adr/ADR-20260601-002-transaction-semantics.md
```

测试：

```text
TestWithTxCommit
TestWithTxRollbackOnError
TestWithTxRollbackOnPanic
TestWithTxContextCanceled
TestTxImplementsQueryer
```

## TASK-POSTGRESX-011：实现 Migration Runner

文件：

```text
pkg/postgresx/migration.go
pkg/postgresx/migration_test.go
docs/adr/ADR-20260601-003-migration-runner-scope.md
```

测试：

```text
TestMigrationRunnerCreatesSchemaMigrations
TestMigrationRunnerAppliesInOrder
TestMigrationRunnerSkipsApplied
TestMigrationRunnerDetectsVersionNameConflict
TestMigrationRunnerRollbackOnFailure
```

## TASK-POSTGRESX-012：实现 Error Mapping

文件：

```text
pkg/postgresx/errors.go
pkg/postgresx/errors_test.go
```

测试：

```text
TestMapErrorContextCanceled
TestMapErrorContextDeadlineExceeded
TestMapErrorNoRows
TestMapErrorUniqueViolation
TestMapErrorPreservesCause
TestIsRetryable
```

## TASK-POSTGRESX-013：实现 HealthCheck

文件：

```text
pkg/postgresx/health.go
pkg/postgresx/health_test.go
```

测试：

```text
TestClientImplementsFoundationHealthChecker
TestHealthCheckHealthy
TestHealthCheckUnhealthyWhenClosedOrUnavailable
TestHealthCheckMetadataSanitized
TestHealthCheckLatency
```

## TASK-POSTGRESX-014：实现 Metrics Hook

文件：

```text
pkg/postgresx/metrics.go
pkg/postgresx/metrics_test.go
```

测试：

```text
TestMetricsNoopDoesNotPanic
TestMetricsRecordsPing
TestMetricsRecordsTx
TestMetricsLabelsNoSecrets
```

## TASK-POSTGRESX-015：实现 TestKit

文件：

```text
testkit/container.go
testkit/config.go
testkit/fixture.go
testkit/assert.go
```

能力：

```text
StartPostgres
ConfigFromContainer
ConfigFromEnv
CreateTempSchema
DropTempSchema
RequireIntegration
```

要求：

```text
1. 若 Docker 可用，启动临时 PostgreSQL
2. 若设置 POSTGRESX_INTEGRATION_DSN，则使用外部 PostgreSQL
3. 若二者都不可用，普通 PR 可明确 skip；release-final-check 不得 skip
4. 不输出明文密码
```

## TASK-POSTGRESX-016：编写 Examples

目录：

```text
examples/basic
examples/transaction
examples/migration
examples/health
```

要求：

```text
1. examples 不包含真实密码
2. examples 支持通过环境变量读取测试 DSN 或分字段配置
3. examples 文档说明仅用于本地/测试
4. examples 不强依赖 configx
```

## TASK-POSTGRESX-017：增加 Template Alignment Gate

文件：

```text
scripts/check_template_alignment.sh
```

检查：

```text
无模板占位符
核心目录存在
关键 scripts 存在
contracts 存在
.agent 存在
```

## TASK-POSTGRESX-018：增加 Foundationx API Compatibility Gate

文件：

```text
scripts/check_foundationx_api.sh
```

检查：

```text
go list -m github.com/ZoneCNH/foundationx
grep -R "github.com/ZoneCNH/foundationx/pkg/foundationx" pkg internal testkit
grep -R "github.com/bytechainx" . 不得出现
```

## TASK-POSTGRESX-019：更新 Boundary Gate

文件：

```text
scripts/check_boundary.sh
```

必须检查：

```text
不依赖 x.go
不依赖 configx core
不依赖 observex core
不出现业务术语
不出现 bytechainx 旧路径
```

## TASK-POSTGRESX-020：编写 Configx Boundary 文档

文件：

```text
docs/configx-boundary.md
docs/xgo-integration.md
docs/adr/ADR-20260601-005-config-loading-boundary.md
```

必须说明：

```text
1. postgresx 只定义 Config，不加载 Config
2. postgresx 不依赖 configx
3. postgresx 不读取 /home/k8s/secrets/env/*
4. configx 负责显式读取 env/envfile/json/map
5. x.go 或 app bootstrap 负责 Decode RuntimeConfig
6. x.go 或 app bootstrap 再构造 postgresx.Config
```

## TASK-POSTGRESX-021：更新 Release / Evidence Gate

必须支持：

```bash
GOWORK=off make ci
GOWORK=off make ci-extended
GOWORK=off make release-check
GOWORK=off make release-preflight VERSION=v0.1.0
GOWORK=off make evidence
GOWORK=off make release-evidence-check
GOWORK=off make release-final-check
```

要求：

```text
lint/security 工具缺失不得伪造 passed
release manifest 不含 secret
release manifest 与当前仓库事实一致
release-final-check 要求工作区 clean
```

## TASK-POSTGRESX-022：Retrospective

输出：

```text
.agent/retrospective.md
.agent/patch_prompt.md
.agent/patch_harness.md
.agent/patch_rule.md
```

必须回答：

```text
1. postgresx 哪些模式可以复制到 redisx/kafkax/taosx？
2. baselib-template 哪些 gate 需要回补？
3. foundationx API 是否足够支持 L2 库？
4. configx 是否仍应保持在 app bootstrap 层？
5. TestKit 是否可抽象到 testkitx？
6. Integration gate 是否可复用？
```

---

# 15. Harness Gates

## Gate 1：Template Alignment

```bash
GOWORK=off ./scripts/check_template_alignment.sh
```

## Gate 2：Foundationx API Compatibility

```bash
GOWORK=off ./scripts/check_foundationx_api.sh
```

## Gate 3：Format

```bash
GOWORK=off go fmt ./...
```

## Gate 4：Vet

```bash
GOWORK=off go vet ./...
```

## Gate 5：Unit Test

```bash
GOWORK=off go test ./...
```

## Gate 6：Race Test

```bash
GOWORK=off go test -race ./...
```

## Gate 7：Boundary

```bash
GOWORK=off ./scripts/check_boundary.sh
```

必须检查：

```text
不依赖 github.com/ZoneCNH/x.go
不依赖 x.go/internal
不依赖 github.com/ZoneCNH/configx
不依赖 github.com/ZoneCNH/observex
不出现 github.com/bytechainx
不出现业务词汇
```

## Gate 8：Secret

```bash
GOWORK=off ./scripts/check_secrets.sh
```

必须无疑似密钥。

## Gate 9：Contract

```bash
GOWORK=off ./scripts/check_contracts.sh
```

检查：

```text
contracts/config.schema.json
contracts/error.schema.json
contracts/health.schema.json
contracts/metrics.md
contracts/public_api.md
docs/api.md
```

## Gate 10：Integration

```bash
GOWORK=off ./scripts/run_integration.sh
```

必须验证真实 PostgreSQL。

## Gate 11：Examples

```bash
GOWORK=off go run ./examples/basic
GOWORK=off go run ./examples/transaction
GOWORK=off go run ./examples/migration
GOWORK=off go run ./examples/health
```

## Gate 12：CI Extended

```bash
GOWORK=off make ci-extended
```

## Gate 13：Release Preflight

```bash
GOWORK=off make release-preflight VERSION=v0.1.0
```

## Gate 14：Evidence

```bash
GOWORK=off make evidence
GOWORK=off make release-evidence-check
```

## Gate 15：Final Release

```bash
GOWORK=off make release-final-check
```

要求：

```text
工作区 clean
manifest 新鲜
所有必需 gate passed
不得伪造 skipped
```

---

# 16. Boundary Gate 脚本模板

```bash
#!/usr/bin/env bash
set -euo pipefail

echo "checking postgresx boundary..."

FORBIDDEN_DEPS=(
  "github.com/ZoneCNH/x.go"
  "github.com/ZoneCNH/x.go/internal"
  "github.com/bytechainx"
  "github.com/ZoneCNH/configx"
  "github.com/ZoneCNH/observex"
)

DEPS="$(GOWORK=off go list -deps ./...)"

for dep in "${FORBIDDEN_DEPS[@]}"; do
  if echo "$DEPS" | grep -q "$dep"; then
    echo "ERROR: forbidden dependency found: $dep"
    exit 1
  fi
done

FORBIDDEN_TERMS=(
  "BTCUSDT"
  "ETHUSDT"
  "Kline"
  "OrderBook"
  "MarketData"
  "MacroData"
  "MacroRegime"
  "MarketRegime"
  "TradingSignal"
  "Position"
  "RiskGate"
  "M1"
  "M2"
  "S1"
  "S2"
)

for term in "${FORBIDDEN_TERMS[@]}"; do
  if grep -R "$term" ./pkg ./internal ./testkit --exclude-dir=.git; then
    echo "ERROR: forbidden business term found: $term"
    exit 1
  fi
done

echo "postgresx boundary check passed"
```

---

# 17. Template Alignment 脚本模板

```bash
#!/usr/bin/env bash
set -euo pipefail

echo "checking template alignment..."

if grep -R "{{MODULE_NAME}}\|{{MODULE_PATH}}\|{{PACKAGE_NAME}}" .   --exclude-dir=.git   --exclude-dir=vendor; then
  echo "ERROR: unresolved template placeholders found"
  exit 1
fi

required_paths=(
  "pkg/postgresx"
  "internal"
  "testkit"
  "examples"
  "contracts"
  "docs"
  "scripts"
  ".agent"
  "release/manifest"
  "Makefile"
)

for path in "${required_paths[@]}"; do
  if [ ! -e "$path" ]; then
    echo "ERROR: missing template path: $path"
    exit 1
  fi
done

echo "template alignment check passed"
```

---

# 18. Foundationx API 脚本模板

```bash
#!/usr/bin/env bash
set -euo pipefail

echo "checking foundationx API compatibility..."

GOWORK=off go list -m github.com/ZoneCNH/foundationx >/dev/null

if ! grep -R "github.com/ZoneCNH/foundationx/pkg/foundationx" ./pkg ./internal ./testkit --exclude-dir=.git; then
  echo "ERROR: foundationx package import not found"
  exit 1
fi

if grep -R "github.com/bytechainx" . --exclude-dir=.git --exclude-dir=vendor; then
  echo "ERROR: old bytechainx import path found"
  exit 1
fi

echo "foundationx API compatibility check passed"
```

---

# 19. Configx Boundary 文档片段

`docs/configx-boundary.md` 必须包含：

```markdown
# Config Loading Boundary

postgresx core does not load configuration files.

## Correct dependency direction

x.go bootstrap -> configx -> postgresx.Config -> postgresx.New

## Forbidden dependency direction

postgresx -> configx -> /home/k8s/secrets/env/postgres.env

## Rules

1. postgresx defines Config.
2. postgresx validates Config.
3. postgresx sanitizes Config.
4. postgresx does not read env.
5. postgresx does not read env files.
6. postgresx does not import configx.
7. x.go or app bootstrap may use configx.
8. x.go owns business runtime config.
```

示例：

```go
result, err := configx.LoadEnvFile(ctx, "/home/k8s/secrets/env/postgres.env")
if err != nil {
    return err
}

var runtimeCfg PostgresRuntimeConfig
if err := configx.Decode(result, &runtimeCfg); err != nil {
    return err
}

pgCfg := postgresx.DefaultConfig()
pgCfg.Host = runtimeCfg.Host
pgCfg.Port = runtimeCfg.Port
pgCfg.Database = runtimeCfg.Database
pgCfg.User = runtimeCfg.User
pgCfg.Password = runtimeCfg.Password
pgCfg.SSLMode = runtimeCfg.SSLMode

client, err := postgresx.New(ctx, pgCfg)
```

注意：

```text
此示例应位于 docs，不应让 postgresx core go.mod 强依赖 configx。
```

---

# 20. Release Manifest 规则

v1.1 必须继承 baselib-template release Evidence 规则：

```text
release/manifest/latest.json 是生成产物
release manifest 记录 module、commit、tree SHA、源码摘要、contract 指纹、依赖清单、工具版本、生成时间、工作区状态和 gate 结果
make release-evidence-check 验证 manifest 与当前仓库事实一致
make release-final-check 要求工作区 clean
```

Manifest 不得包含：

```text
password
token
secret
raw DSN
authorization
private key
```

Release 声明必须包含：

```text
DONE with evidence:
- release manifest path
- commit
- GOWORK=off gate outputs
- integration evidence
- known risks
```

---

# 21. Traceability Matrix

| Requirement | Acceptance Criteria | Design | Task | Test | Evidence | Status |
|---|---|---|---|---|---|---|
| REQ-POSTGRESX-001 | AC-001-* | Template Design | TASK-001/002 | template alignment | EVID-001/002 | TODO |
| REQ-POSTGRESX-002 | AC-002-* | Foundationx Integration | TASK-003 | API compatibility | EVID-003 | TODO |
| REQ-POSTGRESX-003 | AC-003-* | Boundary | TASK-019 | boundary gate | EVID-019 | TODO |
| REQ-POSTGRESX-004 | AC-004-* | Config | TASK-005 | config_test.go | EVID-005 | TODO |
| REQ-POSTGRESX-005 | AC-005-* | DSN | TASK-006 | dsn_test.go | EVID-006 | TODO |
| REQ-POSTGRESX-006 | AC-006-* | Client/Pool | TASK-008 | client_test.go | EVID-008 | TODO |
| REQ-POSTGRESX-007 | AC-007-* | Queryer | TASK-009 | query_test.go | EVID-009 | TODO |
| REQ-POSTGRESX-008 | AC-008-* | Transaction | TASK-010 | tx_test.go | EVID-010 | TODO |
| REQ-POSTGRESX-009 | AC-009-* | Migration | TASK-011 | migration_test.go | EVID-011 | TODO |
| REQ-POSTGRESX-010 | AC-010-* | Health | TASK-013 | health_test.go | EVID-013 | TODO |
| REQ-POSTGRESX-011 | AC-011-* | Error Map | TASK-012 | errors_test.go | EVID-012 | TODO |
| REQ-POSTGRESX-012 | AC-012-* | Configx Boundary | TASK-020 | docs review | EVID-020 | TODO |
| REQ-POSTGRESX-013 | AC-013-* | Release Gates | TASK-021 | release-final-check | EVID-021 | TODO |

---

# 22. Risk Register

## RISK-POSTGRESX-001：模板漂移

风险：

```text
postgresx 手工修改后偏离 baselib-template 的目录、gate、release evidence 规则。
```

缓解：

```text
check_template_alignment.sh
Retrospective 回补 template patch
```

## RISK-POSTGRESX-002：foundationx API 假设错误

风险：

```text
postgresx 使用了 foundationx 不存在或已变化的 API。
```

缓解：

```text
check_foundationx_api.sh
GOWORK=off go test
ADR 记录 foundationx 版本
```

## RISK-POSTGRESX-003：configx 被错误拉入 core

风险：

```text
postgresx 为方便读取 envfile 直接依赖 configx。
```

缓解：

```text
Config loading boundary ADR
Boundary Gate 禁止 configx core import
docs-only integration example
```

## RISK-POSTGRESX-004：DSN 泄露密码

风险：

```text
错误、日志、Evidence 输出原始 DSN。
```

缓解：

```text
Config.Sanitize
RedactedDSN
Secret Gate
DSN redaction tests
```

## RISK-POSTGRESX-005：Migration Runner 越界

风险：

```text
postgresx 变成业务 migration 管理系统。
```

缓解：

```text
只提供 runner，不拥有 migration 内容。
不扫描 x.go 目录。
```

## RISK-POSTGRESX-006：Integration Gate 被跳过

风险：

```text
PR 中集成测试 skip，release 却误认为通过。
```

缓解：

```text
release-final-check 不允许 integration skip
manifest 明确记录 integration status
```

---

# 23. Decision Log

## DEC-20260601-001：Template-bound skeleton

决策：

```text
postgresx v1.1 必须从 baselib-template 渲染。
```

原因：

```text
复用已经完成的目录结构、Harness、Evidence、Release Gate。
```

## DEC-20260601-002：Foundationx-bound contracts

决策：

```text
postgresx 必须复用 foundationx 的 Error / Health / Secret / Clock 契约。
```

原因：

```text
避免 L2 库重新定义 L0 语义。
```

## DEC-20260601-003：Configx 不进入 postgresx core

决策：

```text
postgresx core 只接收 Config，不加载 Config。
```

原因：

```text
保持 L2 基础设施库边界，防止配置加载副作用进入数据库库。
```

## DEC-20260601-004：GOWORK=off 是强制 gate

决策：

```text
所有独立模块 gate 使用 GOWORK=off。
```

原因：

```text
防止父级 workspace 影响独立 module 真实性。
```

## DEC-20260601-005：Release Evidence 不可伪造

决策：

```text
lint/security/tool missing 不得伪造 passed。
```

原因：

```text
Release Manifest 必须反映真实门禁状态。
```

---

# 24. AutoResearch Protocol

触发条件：

```text
1. pgx/pgxpool 当前推荐版本不确定
2. pgxpool Config 字段行为不确定
3. pgx 错误码结构不确定
4. PostgreSQL unique violation SQLSTATE 不确定
5. testcontainers-go 用法不确定
6. GitHub Actions PostgreSQL service 行为不确定
7. baselib-template render_template.sh 参数或行为变化
8. foundationx API 与 README 不一致
9. configx 集成边界不明确
```

输出必须写入：

```text
docs/adr/ADR-YYYYMMDD-NNN-<topic>.md
```

禁止：

```text
1. 不经 ADR 直接引入新依赖
2. 不经测试直接改变事务语义
3. 不经 Review 直接扩大 Public API
4. 不经 ADR 让 postgresx core 依赖 configx
```

---

# 25. Review Checklist

Review 前必须确认：

```text
[ ] postgresx 从 baselib-template 渲染
[ ] module path 是 github.com/ZoneCNH/postgresx
[ ] 无未替换模板占位符
[ ] 无 github.com/bytechainx 旧路径
[ ] 依赖 github.com/ZoneCNH/foundationx
[ ] import github.com/ZoneCNH/foundationx/pkg/foundationx
[ ] 不依赖 x.go
[ ] 不依赖 configx core
[ ] 不依赖 observex core
[ ] 不包含业务模型
[ ] 不读取 /home/k8s/secrets/env/*
[ ] Config Validate 完整
[ ] Password 使用 foundationx.SecretString
[ ] RedactedDSN 不泄露密码
[ ] Client Close 幂等
[ ] Queryer 最小接口可用
[ ] WithTx commit/rollback/panic 语义清晰
[ ] Migration Runner 不拥有业务 migration
[ ] HealthCheck 使用 foundationx.HealthStatus
[ ] Error Mapping 保留 Cause
[ ] TestKit 可运行真实 PostgreSQL
[ ] docs/configx-boundary.md 完整
[ ] GOWORK=off make ci 通过
[ ] GOWORK=off make ci-extended 通过
[ ] GOWORK=off make release-check 通过
[ ] GOWORK=off make release-preflight VERSION=v0.1.0 通过
[ ] GOWORK=off make release-evidence-check 通过
[ ] GOWORK=off make release-final-check 通过
[ ] release manifest 生成且不含 secret
```

---

# 26. Release Protocol

## 26.1 v0.1.0 发布前

执行：

```bash
GOWORK=off make ci
GOWORK=off make ci-extended
GOWORK=off make release-check
GOWORK=off make release-preflight VERSION=v0.1.0
GOWORK=off make evidence
GOWORK=off make release-evidence-check
GOWORK=off make release-final-check
```

必须通过：

```text
fmt
vet
lint
unit test
race
boundary
security
contracts
docs
examples
integration
manifest generation
manifest evidence check
clean workspace check
```

## 26.2 CHANGELOG

```markdown
## v0.1.0 - 2026-06-01

### Added
- Created postgresx from github.com/ZoneCNH/baselib-template.
- Added foundationx integration.
- Added Config / Validate / Sanitize.
- Added DSN and RedactedDSN builders.
- Added pgxpool-based Client.
- Added Ping, Close, and PoolStats.
- Added Queryer abstraction.
- Added WithTx and WithTxOptions transaction helper.
- Added MigrationRunner.
- Added HealthCheck with foundationx.HealthStatus.
- Added error mapping to foundationx.Error.
- Added metrics hook interface and noop metrics.
- Added testkit for PostgreSQL integration tests.
- Added template alignment gate.
- Added foundationx API compatibility gate.
- Added configx boundary documentation.
- Added GOWORK=off release validation.

### Security
- Password uses foundationx.SecretString.
- RedactedDSN masks passwords.
- Secret Gate added.
- Release Manifest must not contain secrets.

### Boundary
- postgresx core does not depend on x.go.
- postgresx core does not depend on configx.
- postgresx core does not depend on observex.
- postgresx does not own business migrations.

### Breaking Changes
- None.
```

## 26.3 Release 声明

```text
DONE with evidence:
- GOWORK=off go test ./... passed
- GOWORK=off go test -race ./... passed
- GOWORK=off make ci passed
- GOWORK=off make ci-extended passed
- GOWORK=off make release-check passed
- GOWORK=off make release-preflight VERSION=v0.1.0 passed
- GOWORK=off make release-evidence-check passed
- GOWORK=off make release-final-check passed
- integration tests passed against PostgreSQL
- boundary gate passed
- secret gate passed
- release manifest generated
```

---

# 27. x.go 集成规范

x.go 错误方式：

```go
db, err := sql.Open("postgres", dsn)
```

x.go 正确方式：

```go
pgCfg := postgresx.DefaultConfig()
pgCfg.Host = runtimeCfg.Postgres.Host
pgCfg.Port = runtimeCfg.Postgres.Port
pgCfg.Database = runtimeCfg.Postgres.Database
pgCfg.User = runtimeCfg.Postgres.User
pgCfg.Password = runtimeCfg.Postgres.Password
pgCfg.SSLMode = runtimeCfg.Postgres.SSLMode

pgClient, err := postgresx.New(ctx, pgCfg)
if err != nil {
    return err
}
defer pgClient.Close(ctx)
```

业务 repository：

```go
type Repository struct {
    db postgresx.Queryer
}

func NewRepository(db postgresx.Queryer) *Repository {
    return &Repository{db: db}
}
```

## configx 组合方式

```text
configx 不是 postgresx core 依赖。
configx 可以在 x.go bootstrap 使用。
```

示意：

```go
result, err := configx.LoadEnvFile(ctx, "/home/k8s/secrets/env/postgres.env")
if err != nil {
    return err
}

var runtimeCfg PostgresRuntimeConfig
if err := configx.Decode(result, &runtimeCfg); err != nil {
    return err
}

pgCfg := postgresx.DefaultConfig()
pgCfg.Host = runtimeCfg.Host
pgCfg.Port = runtimeCfg.Port
pgCfg.Database = runtimeCfg.Database
pgCfg.User = runtimeCfg.User
pgCfg.Password = runtimeCfg.Password
pgCfg.SSLMode = runtimeCfg.SSLMode

client, err := postgresx.New(ctx, pgCfg)
```

x.go 必须保留：

```text
业务 schema
业务 migration
业务 repository
业务 query
业务配置加载
```

postgresx 只提供：

```text
连接
事务
migration runner
health
metrics hook
error mapping
```

---

# 28. Retrospective Protocol

输出：

```text
.agent/retrospective.md
```

模板：

```markdown
# postgresx Retrospective v1.1

## Release
- Version:
- Commit:
- Date:

## Template alignment
- Was baselib-template sufficient?
- What should be patched back into baselib-template?

## Foundationx compatibility
- Which foundationx APIs were used?
- Any missing L0 contracts?

## Configx boundary
- Did postgresx avoid configx core dependency?
- Is docs-only integration enough?

## What worked
-

## What failed
-

## API stability concerns
-

## Boundary risks
-

## Security findings
-

## Test gaps
-

## Harness improvements
-

## Reusable patterns for other base libs
- redisx:
- kafkax:
- taosx:
- ossx:
- configx:
- observex:
- testkitx:

## Patch outputs
- PATCH-PROMPT:
- PATCH-HARNESS:
- PATCH-RULE:
```

---

# 29. Final DoD

## Task DoD

```text
代码实现完成
单元测试完成
integration 测试完成
无业务语义污染
无 x.go 依赖
无 bytechainx 旧路径
无 configx core 依赖
无 observex core 依赖
无密钥泄露
GOWORK=off go fmt / go vet / go test / go test -race 通过
```

## Module DoD

```text
Template render 完整
Foundationx integration 完整
Config 完整
DSN 完整
Client 完整
Queryer 完整
Tx 完整
Migration Runner 完整
Health 完整
Error Mapping 完整
Metrics Hook 完整
TestKit 完整
Examples 完整
Configx Boundary Docs 完整
Docs 完整
ADR 完整
Harness 完整
Release Manifest 完整
```

## Goal DoD

```text
postgresx 可作为 x.go PostgreSQL 基础库使用
postgresx 从 baselib-template 生成并对齐
postgresx 依赖 foundationx
postgresx 不依赖 x.go
postgresx 不包含业务 schema
postgresx 不读取生产密钥
postgresx 不强依赖 configx
postgresx integration test 可验证真实 PostgreSQL
postgresx v0.1.0 release evidence 完整
retrospective patch 生成
```

完成声明必须是：

```text
DONE with evidence:
- GOWORK=off go test ./... passed
- GOWORK=off go test -race ./... passed
- GOWORK=off make ci passed
- GOWORK=off make ci-extended passed
- GOWORK=off make release-check passed
- GOWORK=off make release-preflight VERSION=v0.1.0 passed
- GOWORK=off make release-evidence-check passed
- GOWORK=off make release-final-check passed
- boundary gate passed
- secret gate passed
- integration PostgreSQL test passed
- release manifest generated
```

---

# 30. 最小可行执行顺序

Agent 执行时按以下顺序，不要跳步：

```text
1. clone baselib-template
2. render postgresx from template
3. run template alignment check
4. run GOWORK=off go test ./...
5. add foundationx dependency
6. add foundationx API compatibility check
7. AutoResearch pgx/pgxpool and write ADR
8. implement Config + tests
9. implement DSN + tests
10. implement Options/Noop + tests
11. implement Client/Pool + tests
12. implement Queryer + integration tests
13. implement WithTx + integration tests
14. write transaction semantics ADR
15. implement Migration Runner + integration tests
16. write migration scope ADR
17. implement Error Mapping + tests
18. implement HealthCheck + tests
19. implement Metrics Hook + tests
20. implement TestKit
21. write Examples
22. write Configx Boundary docs
23. update scripts
24. update Makefile / CI
25. run GOWORK=off make ci
26. run GOWORK=off make ci-extended
27. run GOWORK=off make release-check
28. run GOWORK=off make release-preflight VERSION=v0.1.0
29. run GOWORK=off make evidence
30. run GOWORK=off make release-evidence-check
31. run GOWORK=off make release-final-check
32. write retrospective
33. output DONE with evidence
```

---

# 31. 给 Agent 的最终执行指令

```text
你现在要执行 GOAL-20260601-POSTGRESX-001 v1.1。

请严格按 Goal Runtime Prompt v3.1 执行：
Goal → Context Recovery → Spec → Design → Plan → Tasks → Execution → Verification → Evidence → Review → Release → Retrospective → Self-improving。

你必须创建或完善 github.com/ZoneCNH/postgresx。

硬性约束：
1. postgresx 必须从 github.com/ZoneCNH/baselib-template 渲染。
2. postgresx module path 必须是 github.com/ZoneCNH/postgresx。
3. postgresx 必须依赖 github.com/ZoneCNH/foundationx。
4. postgresx 必须 import github.com/ZoneCNH/foundationx/pkg/foundationx。
5. postgresx 可以依赖 pgx/pgxpool，但具体版本必须通过 AutoResearch 或 ADR 记录。
6. postgresx 不允许依赖 github.com/ZoneCNH/x.go。
7. postgresx 不允许出现 github.com/bytechainx 旧路径。
8. postgresx 不允许包含 x.go 业务语义。
9. postgresx 不允许包含业务表结构。
10. postgresx 不允许隐式读取 /home/k8s/secrets/env/*。
11. postgresx core 不允许依赖 configx。
12. postgresx core 不允许依赖 observex。
13. postgresx 不允许使用全局 DB / 单例 Client。
14. postgresx 不允许在日志、错误、Evidence、Release Manifest 中输出明文密码或 DSN。
15. 所有独立模块验证必须使用 GOWORK=off。
16. 不允许没有 Evidence 就声称 DONE。

必须实现：
1. Template render from baselib-template
2. Template alignment gate
3. Foundationx API compatibility gate
4. Config / Validate / Sanitize
5. DSN / RedactedDSN
6. Client / New / Ping / Close / Stats
7. Queryer abstraction
8. WithTx / WithTxOptions
9. MigrationRunner
10. HealthCheck with foundationx.HealthStatus
11. Error Mapping to foundationx.Error
12. Metrics Hook
13. TestKit
14. Integration Tests
15. Examples
16. Configx Boundary docs
17. Harness scripts
18. Makefile
19. GitHub Actions
20. Docs / ADR
21. Release Manifest
22. Retrospective patches

执行完成后输出：

DONE with evidence:
- 具体命令
- 具体测试结果
- 具体文件路径
- release manifest 路径
- known risks
- next recommended issue
```

---

# 32. 最终推荐路径

postgresx v1.1 必须先做“模板绑定 + foundationx 绑定 + PostgreSQL 核心能力”：

```text
Template
Foundationx
Config
DSN
Pool
Ping
Close
Queryer
Tx
Migration Runner
Health
Error
Metrics Hook
TestKit
Evidence
```

暂不做：

```text
ORM
业务 repository
读写分离
复杂 migration DSL
PostgreSQL HA
业务 schema
强依赖 configx
强依赖 observex
```

最重要的四条红线：

```text
1. 不偏离 baselib-template
2. 不绕过 foundationx
3. 不依赖 x.go / configx core / observex core
4. 不泄露密钥
```

最小交付：

```text
v0.1.0 = Template-bound PostgreSQL 连接池 + 事务 + migration runner + health + error mapping + testkit + release evidence
```
