package postgresx

import (
	"time"
)

// Config contains explicit PostgreSQL connection settings.
type Config struct {
	Host            string        `json:"host"`
	Port            int           `json:"port"`
	Database        string        `json:"database"`
	User            string        `json:"user"`
	Password        SecretString  `json:"password"`
	SSLMode         string        `json:"sslmode"`
	MaxOpenConns    int32         `json:"max_open_conns"`
	MinIdleConns    int32         `json:"min_idle_conns"`
	MaxConnLifetime time.Duration `json:"max_conn_lifetime"`
	MaxConnIdleTime time.Duration `json:"max_conn_idle_time"`
	ConnectTimeout  time.Duration `json:"connect_timeout"`
	HealthTimeout   time.Duration `json:"health_timeout"`
	ApplicationName string        `json:"application_name"`
}

// SanitizedConfig is safe for logs, health metadata, and evidence files.
type SanitizedConfig struct {
	Host            string `json:"host"`
	Port            int    `json:"port"`
	Database        string `json:"database"`
	User            string `json:"user"`
	Password        string `json:"password"`
	SSLMode         string `json:"sslmode"`
	MaxOpenConns    int32  `json:"max_open_conns"`
	MinIdleConns    int32  `json:"min_idle_conns"`
	MaxConnLifetime string `json:"max_conn_lifetime"`
	MaxConnIdleTime string `json:"max_conn_idle_time"`
	ConnectTimeout  string `json:"connect_timeout"`
	HealthTimeout   string `json:"health_timeout"`
	ApplicationName string `json:"application_name"`
}

// DefaultConfig returns non-sensitive default connection settings.
func DefaultConfig() Config {
	return Config{
		Port:            5432,
		SSLMode:         "disable",
		MaxOpenConns:    10,
		MinIdleConns:    1,
		MaxConnLifetime: time.Hour,
		MaxConnIdleTime: 30 * time.Minute,
		ConnectTimeout:  5 * time.Second,
		HealthTimeout:   2 * time.Second,
	}
}

// Validate checks the explicit connection contract without reading env.
func (c Config) Validate() error {
	const op = "postgresx.Config.Validate"
	if c.Host == "" {
		return NewError(ErrorKindConfig, op, "host is required")
	}
	if c.Port <= 0 || c.Port > 65535 {
		return NewError(ErrorKindConfig, op, "port must be between 1 and 65535")
	}
	if c.Database == "" {
		return NewError(ErrorKindConfig, op, "database is required")
	}
	if c.User == "" {
		return NewError(ErrorKindConfig, op, "user is required")
	}
	if c.Password.IsZero() {
		return NewError(ErrorKindConfig, op, "password is required")
	}
	if c.SSLMode == "" {
		return NewError(ErrorKindConfig, op, "sslmode is required")
	}
	if c.MaxOpenConns < 0 {
		return NewError(ErrorKindConfig, op, "max open connections must be non-negative")
	}
	if c.MinIdleConns < 0 {
		return NewError(ErrorKindConfig, op, "min idle connections must be non-negative")
	}
	if c.MaxOpenConns > 0 && c.MinIdleConns > c.MaxOpenConns {
		return NewError(ErrorKindConfig, op, "min idle connections cannot exceed max open connections")
	}
	if c.MaxConnLifetime < 0 {
		return NewError(ErrorKindConfig, op, "max connection lifetime must be non-negative")
	}
	if c.MaxConnIdleTime < 0 {
		return NewError(ErrorKindConfig, op, "max connection idle time must be non-negative")
	}
	if c.ConnectTimeout < 0 {
		return NewError(ErrorKindConfig, op, "connect timeout must be non-negative")
	}
	if c.HealthTimeout < 0 {
		return NewError(ErrorKindConfig, op, "health timeout must be non-negative")
	}
	return nil
}

// Sanitize returns a representation that does not expose the password.
func (c Config) Sanitize() SanitizedConfig {
	return SanitizedConfig{
		Host:            c.Host,
		Port:            c.Port,
		Database:        c.Database,
		User:            c.User,
		Password:        c.Password.String(),
		SSLMode:         c.SSLMode,
		MaxOpenConns:    c.MaxOpenConns,
		MinIdleConns:    c.MinIdleConns,
		MaxConnLifetime: c.MaxConnLifetime.String(),
		MaxConnIdleTime: c.MaxConnIdleTime.String(),
		ConnectTimeout:  c.ConnectTimeout.String(),
		HealthTimeout:   c.HealthTimeout.String(),
		ApplicationName: c.ApplicationName,
	}
}

func (c Config) withDefaults() Config {
	defaults := DefaultConfig()
	if c.Port == 0 {
		c.Port = defaults.Port
	}
	if c.SSLMode == "" {
		c.SSLMode = defaults.SSLMode
	}
	if c.MaxOpenConns == 0 {
		c.MaxOpenConns = defaults.MaxOpenConns
	}
	if c.MinIdleConns == 0 {
		c.MinIdleConns = defaults.MinIdleConns
	}
	if c.MaxConnLifetime == 0 {
		c.MaxConnLifetime = defaults.MaxConnLifetime
	}
	if c.MaxConnIdleTime == 0 {
		c.MaxConnIdleTime = defaults.MaxConnIdleTime
	}
	if c.ConnectTimeout == 0 {
		c.ConnectTimeout = defaults.ConnectTimeout
	}
	if c.HealthTimeout == 0 {
		c.HealthTimeout = defaults.HealthTimeout
	}
	return c
}
