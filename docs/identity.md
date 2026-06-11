# postgresx 身份

## 我是谁

`postgresx` 是 FoundationX 的 **PostgreSQL 存储扩展模块**，提供 PostgreSQL 客户端封装和标准数据库操作接口。

## 我做什么

- PostgreSQL 连接池和健康检查
- 查询/事务/迁移接口封装
- 配置模型标准化

## 我不做什么

- 不是业务仓储层 — 业务数据访问由调用方定义
- 不是 ORM — 仅提供轻量封装
- 不是模板源 — 模板生成属于 xlib-standard
- 不依赖其他存储模块

## 宪法合规

| 条款 | 遵循方式 |
|------|----------|
| §3.3 | 存储扩展，可依赖 kernel + observex (interface-only) |
| §3.4 | 不依赖 configx、业务域、其他存储扩展 |
| §1 P13 | 存储扩展之间平级协作 |
