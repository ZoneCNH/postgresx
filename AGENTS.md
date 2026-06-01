# Repository Guidelines

## 项目结构与模块组织

本仓库是 Go 模块 `github.com/bytechainx/postgresx`，核心代码位于根目录，包名为 `postgresx`。
主要入口包括 `client.go`、`pool.go`、`tx.go`、`migrator.go`、`errors.go`、`retry.go` 和 `health.go`。
内部实现放在 `internal/secretmask/`，可复用测试工具放在 `testkit/`。
示例程序在 `examples/basic/`、`examples/transaction/`、`examples/migration/` 和 `examples/sqlc/`。
设计、ADR、版本矩阵和发布证据保存在 `docs/`。

## 构建、测试与开发命令

- `make fmt`：对所有 Go 文件执行 `gofmt -w`。
- `make vet`：运行 `go vet ./...`，对应 CI lint 步骤。
- `make test`：运行 `go test ./...`。
- `make race`：运行 `go test -race ./...`。
- `make secret-scan`：扫描 PostgreSQL DSN、`PGPASSWORD` 和明文 `password=` 模式。
- `make ci`：依次执行格式化、静态检查、单元测试、竞态检测和密钥扫描。
- `POSTGRES_TEST_DSN=... make integration`：运行迁移 up/down/up 集成门禁；未设置 DSN 时脚本会跳过。

## 编码风格与命名约定

使用 Go 1.26.3。提交前保持 `gofmt` 后的制表符缩进和标准 Go import 分组。
导出 API 使用清晰的 PascalCase 名称并补充文档注释；内部 helper 使用 camelCase。
错误码、配置项和公共接口应保持向后兼容，避免让 `postgresx` 拥有业务 schema 或 ORM 行为。

## 测试指南

测试使用标准库 `testing`，文件命名为 `*_test.go`，测试函数命名为 `TestXxx`。
普通单元测试应可通过 `go test ./...` 无外部依赖运行。需要真实 PostgreSQL 的测试放在 integration 测试路径中，并通过 `POSTGRES_TEST_DSN` 控制。
修改迁移、连接池、错误归一化或 secret masking 时，优先补充针对性回归测试。

## 提交与 Pull Request 规则

当前 Git 历史只有 `Initial commit`，尚未形成可观察的长期提交风格。
新提交使用简短祈使句说明意图，并遵循 Lore trailer：至少记录 `Tested:`，有风险时补充 `Constraint:`、`Rejected:`、`Confidence:`、`Scope-risk:` 和 `Not-tested:`。
PR 应包含变更目的、验证命令输出摘要、关联 issue 或 ADR；影响示例、文档或外部 API 时同步更新 `README.md`、`docs/` 或 `examples/`。

## 安全、配置与代理规则

不要提交真实 DSN、密码或令牌；本地使用环境变量，例如 `POSTGRES_TEST_DSN`。
所有协作说明、评审意见和代理回复默认并强制使用中文；代码标识符、命令、错误文本和外部 API 名称保持原文。
