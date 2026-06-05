# Error Code Mapping

postgresx maps database and transport failures onto the foundationx error-kind
vocabulary exposed by the current implementation.

| Source | ErrorKind | Retryable |
| --- | --- | --- |
| `pgx.ErrNoRows` | `not_found` | no |
| `context.DeadlineExceeded` | `timeout` | yes |
| `context.Canceled` | `canceled` | no |
| `net.Error` timeout | `timeout` | yes |
| non-timeout `net.Error` | `connection` | yes |
| SQLSTATE class `08` | `connection` | yes |
| SQLSTATE `57014` | `timeout` | yes |
| SQLSTATE `40001` | `conflict` | yes |
| SQLSTATE `40P01` | `conflict` | yes |
| SQLSTATE `53300` | `conflict` | yes |
| SQLSTATE `55P03` | `conflict` | yes |
| SQLSTATE class `23` | `conflict` | no |
| text containing `connection refused`, `cannot connect`, or `connection reset` | `connection` | yes |

All other errors normalize to `internal` and are not retryable by default.
