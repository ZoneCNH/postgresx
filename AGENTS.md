# 仓库贡献指南

## 项目概述

本仓库是 Go 1.25 PostgreSQL 基础库 `postgresx`，模块路径为 `github.com/ZoneCNH/postgresx`，提供 PostgreSQL 运行时基础设施：pgx pool 生命周期、配置校验、sqlc 兼容执行、显式事务、重试策略、迁移运行器、错误归一化、健康检查、连接池统计、可选 metrics/tracing 适配器、DSN 脱敏和集成测试辅助。

## 项目结构与模块组织

- `pkg/postgresx`：公共 API（client、pool、tx、migration、errors、query、health、config、metrics、options、version）。
- `internal/secretmask`：内部辅助——配置脱敏。
- `contracts/`：JSON schema 定义和 `contracts_test.go` 验证映射。
- `testkit/`：可复用测试夹具和断言工具。
- `examples/basic`、`examples/transaction`、`examples/migration`、`examples/sqlc`：最小示例。
- `scripts/`：Harness gate shell 脚本。
- `.agent/`：Goal Runtime 工件。
- `docs/`：设计、ADR、版本矩阵和发布证据文档。
- `release/manifest/`：manifest 模板；`latest.json` 是生成产物，不提交到源码历史。

## 构建、测试与开发命令

### 基础开发

- `make fmt`：执行 `gofmt -w`。
- `make vet`：执行 `go vet ./...`。
- `make test`：运行全部单元测试。
- `make race`：使用 race detector 运行测试。
- `make lint`：执行 `golangci-lint run ./...`；缺少 `golangci-lint` 时必须失败。
- `make security`：执行密钥扫描，并在 `XLIB_ENABLE_VULNCHECK=1` 或 `XLIB_FORCE_VULNCHECK=1` 时运行 `govulncheck`。

### 运行单个测试

```bash
go test ./pkg/postgresx/ -run TestConfigValidate
go test ./contracts/ -run TestContracts
go test ./... -run 'Test.*Property|Test.*Invariant'   # 属性测试
go test ./... -run 'Test.*Golden|Test.*Snapshot'       # golden 测试
```

### CI 与 Gate

- `make ci`：fmt + vet + test + race + boundary + contracts + secret-scan。
- `make ci-extended`：ci + foundationx-api + template-alignment。
- `make boundary`：模块边界检查。
- `make contracts`：JSON schema 契约检查。
- `POSTGRES_TEST_DSN=... make integration`：集成测试。

### 发布（必须 GOWORK=off）

所有发布和验证命令必须使用 `GOWORK=off`，避免本地 `go.work` 改写 module 解析：

```bash
GOWORK=off make release-check
GOWORK=off make release-final-check
make evidence
```

## 编码风格与命名约定

使用标准 Go 风格：交给 `gofmt` 处理缩进，包名保持简短，导出标识符要清晰表达用途。公共库能力放入 `pkg/postgresx`，私有辅助逻辑放入 `internal/`。golangci-lint 启用的 linter 见 `.golangci.yml`。

## 测试规范

测试使用 Go 标准 `testing` 包，命名遵循 `TestXxx`；场景较多时优先使用表驱动测试。必须覆盖配置校验、客户端创建、连接池生命周期、事务、迁移、错误归一化和健康检查。

- 小改动：至少 `go test ./...`
- 并发相关：`make race`
- 发布流程：`make integration`
- 属性/不变量：`make property`
- Golden/快照：`make golden`

## 提交与 Pull Request 规范

提交信息必须遵循 Lore protocol：第一行说明变更意图，正文使用 `Constraint:`、`Rejected:`、`Confidence:`、`Scope-risk:`、`Directive:`、`Tested:` 和 `Not-tested:` 等 trailer 记录决策和验证。PR 需要说明对库行为的影响，关联相关 issue，列出已运行命令。

## 关键约束

- **GOWORK=off**：所有发布和验证命令必须使用，避免 `go.work` 改写 module 解析。
- **无真实凭据**：不得提交真实 DSN、密码或令牌，`scripts/check_secrets.sh` 会扫描。
- **Evidence 完成**：最终完成声明必须包含 `DONE with evidence:`。
- **Release manifest 测试**：必须在临时 fixture 仓库中构造所需 `.omc` state，不得依赖当前工作区的 Agent 运行态文件。

## 文档语言规则

所有仓库文档必须默认使用中文叙述，包括 `README.md`、`docs/`、`.agent/`、`contracts/*.md`、变更日志、发布说明、PR 描述模板和贡献指南。专业术语、代码标识符、命令、路径、包名、外部专有名词、协议固定短语和提交标题必须保留项目惯用原文。

## Agent 专用说明

仓库协作、代码评审和进度更新默认使用中文，除非用户明确要求其他语言。代码、命令、路径、包名和提交标题保留项目惯用语言。
