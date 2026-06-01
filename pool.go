package postgresx

// PoolStats is a driver-neutral snapshot of connection pool state.
type PoolStats struct {
	TotalConns        int32 `json:"total_conns"`
	IdleConns         int32 `json:"idle_conns"`
	AcquiredConns     int32 `json:"acquired_conns"`
	ConstructingConns int32 `json:"constructing_conns"`
	MaxConns          int32 `json:"max_conns"`
}
