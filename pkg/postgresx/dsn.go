package postgresx

import (
	"net"
	"net/url"
	"strconv"
)

// DSN builds a PostgreSQL URL for the driver. It intentionally includes the
// password and should not be logged.
func (c Config) DSN() string {
	u := c.postgresURL(c.Password.Reveal())
	return u.String()
}

// RedactedDSN builds a PostgreSQL URL with the password masked for diagnostics.
func (c Config) RedactedDSN() string {
	u := c.postgresURL(c.Password.String())
	return u.String()
}

func (c Config) postgresURL(password string) url.URL {
	values := url.Values{}
	if c.SSLMode != "" {
		values.Set("sslmode", c.SSLMode)
	}
	if c.ApplicationName != "" {
		values.Set("application_name", c.ApplicationName)
	}
	if c.ConnectTimeout > 0 {
		values.Set("connect_timeout", strconv.FormatInt(int64(c.ConnectTimeout.Seconds()), 10))
	}

	return url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(c.User, password),
		Host:     net.JoinHostPort(c.Host, strconv.Itoa(c.Port)),
		Path:     "/" + c.Database,
		RawQuery: values.Encode(),
	}
}
