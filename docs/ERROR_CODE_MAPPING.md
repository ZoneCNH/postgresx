# Error Code Mapping

`postgresx` 将数据库、上下文和传输错误归一化到当前实现暴露的
`foundationx` error-kind 词表。SQLSTATE 映射优先于文本启发式；未命中的错误默认归一化为
`internal`，且默认不可重试。

| 来源 | ErrorKind | Retryable |
| --- | --- | --- |
| `pgx.ErrNoRows` | `not_found` | no |
| `context.DeadlineExceeded` | `timeout` | yes |
| `context.Canceled` | `canceled` | no |
| `net.Error` timeout | `timeout` | yes |
| non-timeout `net.Error` | `connection` | yes |
| SQLSTATE `42601` | `validation` | no |
| SQLSTATE `42P01` | `not_found` | no |
| SQLSTATE `23505` | `already_exists` | no |
| SQLSTATE `23503` | `conflict` | no |
| SQLSTATE `23502` | `validation` | no |
| SQLSTATE `23514` | `validation` | no |
| SQLSTATE `40001` | `conflict` | yes |
| SQLSTATE `40P01` | `conflict` | yes |
| SQLSTATE `55P03` | `conflict` | yes |
| SQLSTATE `57014` | `timeout` | yes |
| SQLSTATE class `08` | `connection` | yes |
| SQLSTATE class `53` | `unavailable` | yes |
| SQLSTATE class `57` | `unavailable` | yes |
| text containing `password authentication failed` | `auth` | no |
| text containing `connection refused`, `cannot connect`, or `connection reset` | `connection` | yes |

`foundationx` 当前没有 `resource_exhausted` kind，因此 PostgreSQL 资源耗尽类 SQLSTATE
`53*` 按 `unavailable` 归一化并标记为可重试。
