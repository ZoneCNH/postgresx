# Error Code Mapping

| Source | ErrorKind | Retryable |
| --- | --- | --- |
| `pgx.ErrNoRows` | `not_found` | no |
| SQLSTATE `23505` | `unique_violation` | no |
| SQLSTATE `23503` | `foreign_key_violation` | no |
| SQLSTATE class `23` | `constraint` | no |
| SQLSTATE `40001` | `serialization_failure` | yes |
| SQLSTATE `40P01` | `deadlock_detected` | yes |
| SQLSTATE class `08` | `connection` | yes |
| SQLSTATE `57014` | `timeout` | yes |
| `context.DeadlineExceeded` | `timeout` | yes |
| `context.Canceled` | `canceled` | no |
| `net.Error` timeout | `timeout` | yes |
| non-timeout `net.Error` | `connection` | yes |
| other error text containing `connect` | `connection` | yes |

All other errors normalize to `unknown`.
