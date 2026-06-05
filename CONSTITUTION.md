# postgresx 宪章

本仓库是 PostgreSQL 基础库，遵循 xlib-standard 规范。标准正文、门禁、发布证据和下游采纳约束必须保持一致，任何变更都应能通过仓库内的治理命令复核。

## 权威顺序

1. `docs/` 中的设计文档和 ADR 描述目标、边界和标准条款。
2. `.agent/` 描述 Goal Runtime 全链路规则。
3. `Makefile` 和 `scripts/` 描述机器可执行的门禁契约。
4. `release/manifest/` 保存发布证据模板或占位；`latest.json` 等运行时产物按 `.gitignore` 重新生成，不直接提交。

## 修改原则

- 先更新标准与门禁，再更新下游说明。
- 不用占位文件伪造通过证据。
- 所有发布、治理和证据相关命令默认使用 `GOWORK=off`。
