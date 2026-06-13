# postgresx 评分审计 — 2026-06-13

## 结论

`postgresx` 当前可验证评分是 `85/100`，对应本仓库 L2 门禁的
`L2-T3` 状态。这个分数已经满足本地发布许可：
`release_allowed=true`，但还不能声明满分、工厂级或生产采用：
`factory_grade_allowed=false`。

不能把当前状态提升为 `100/100`，因为现有证据仍缺少外部 CI、生产
soak、真实下游采用证明，以及 `v1.0.0` 发布 manifest 与已发布 tag 的
历史一致性闭合。

## 评分依据

| 维度 | 当前状态 | 证据 |
| --- | --- | --- |
| 公共 API 与契约 | 通过 | `go test ./pkg/postgresx ./test/contract` 覆盖客户端、事务、迁移、metrics、错误映射和 `Queryer` 边界 |
| 错误归一化 | 通过 | SQLSTATE `42P01` 已映射为 `not_found`，并纳入单元测试、契约测试和 `docs/ERROR_CODE_MAPPING.md` |
| 真实 PostgreSQL 集成 | 通过 | 使用本地 SRE secret 文档加载的开发 DSN，仅通过环境变量注入，未写入文档、manifest 或 evidence |
| chaos / benchmark smoke | 通过 | `make release-check` 生成 `.agent/evidence/raw/*` 与 `.agent/evidence/normalized/*` |
| 本地下游 smoke | 通过 | 临时 consumer module 验证导入、编译、配置脱敏和 `Queryer` 边界 |
| secret scan / 边界检查 | 通过 | 发布 gate 包含 secret scan、boundary、contracts、foundationx API 和 template alignment |
| 发布证据一致性 | 阻塞 | `release-evidence-check` 仍因 manifest source commit 与已发布 tag commit 不具祖先关系失败 |
| 工厂级 / 满分 | 阻塞 | 缺少外部 CI、生产 soak 和真实 consumer release evidence |

## 已验证命令

以下命令均在 `GOWORK=off` 模式下执行，除特别说明外已经通过：

- `go test ./pkg/postgresx ./test/contract`
- `go test ./...`
- `go vet ./...`
- `make lint`
- `make race`
- `make test-integration`
- `make test-unit test-contract test-chaos benchmark-smoke downstream-smoke`
- `make security boundary contracts foundationx-api template-alignment`
- `VERSION=v1.0.0 make release-check`

`VERSION=v1.0.0 make release-evidence-check` 的当前结果是失败，失败原因是
发布 manifest 记录的 source commit `7fe4cfd` 不是已发布 `v1.0.0` tag
commit `310a249` 的祖先。这是发布历史一致性问题，不是测试覆盖或代码行为
失败。

## 真实配置使用边界

集成测试使用 `/home/ZoneCNH/sre/secrets/env/dev.md` 中的开发 PostgreSQL
配置作为本地证据来源。该配置只用于构造进程环境变量，不得写入：

- `docs/`
- `.agent/evidence/`
- `release/manifest/`
- git commit message
- release note

当前文档只记录“使用真实 PostgreSQL 开发 DSN 验证通过”这个事实，不记录 DSN、
用户名、密码、主机、端口或数据库名。

## v1.0.0 发布状态

GitHub release `v1.0.0` 已存在，已发布 tag commit 是 `310a249`。当前
`postgresx` 分支保存的是 tag 之后补齐的 `L2-T3 / 85` evidence 与文档同步，
其中 release manifest 的 source commit 是 `7fe4cfd`。

在没有显式发布历史授权前，不应 force retag 或重写 `v1.0.0`。可选闭合路径：

1. 保持 `v1.0.0` 不可变，从 tag snapshot 重新生成并校验 manifest。
2. 从当前 evidence 分支切后继版本，例如 `v1.0.1`。
3. 仅在获得明确授权时重写 `v1.0.0` tag/release；这是高风险路径，不作为默认方案。

## 满分差距

要把评分从 `85/100` 提升到 `100/100`，至少需要新增并提交以下证据：

- 外部 CI 成功记录，且不受当前 GitHub 账户 billing lock 影响；
- 生产 soak 记录，包括时间窗口、环境、失败率和回滚条件；
- 真实 consumer checkout 的依赖 pin、编译、测试、导入边界和发布 manifest；
- 与最终 tag 一致的 release manifest、checksum 和 `release-evidence-check` 通过记录。

在这些证据补齐前，`85/100` 是当前最高可信评分。
