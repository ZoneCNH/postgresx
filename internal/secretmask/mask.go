package secretmask

import (
	"net/url"
	"regexp"
	"strings"
)

const replacement = "<masked>"

var (
	passwordKVPattern = regexp.MustCompile(`(?i)\b(password|pass|pwd)=([^\s]+)`)
	pgPasswordPattern = regexp.MustCompile(`(?i)\b(PG` + `PASSWORD=)([^\s]+)`)
	pgURLPattern      = regexp.MustCompile(`(?i)(postgres(?:ql)?://[^\s:/@]+:)[^\s@]+(@)`)
)

// Mask redacts common PostgreSQL password locations from loggable text.
func Mask(value string) string {
	if value == "" {
		return value
	}
	masked := maskURL(value)
	masked = passwordKVPattern.ReplaceAllString(masked, `$1=`+replacement)
	masked = pgPasswordPattern.ReplaceAllString(masked, `$1`+replacement)
	masked = pgURLPattern.ReplaceAllString(masked, `${1}`+replacement+`${2}`)
	return masked
}

// MaskDSN is a semantic alias for callers that specifically redact DSNs.
func MaskDSN(dsn string) string {
	return Mask(dsn)
}

func maskURL(value string) string {
	if !strings.HasPrefix(strings.ToLower(value), "postgres://") &&
		!strings.HasPrefix(strings.ToLower(value), "postgresql://") {
		return value
	}
	parsed, err := url.Parse(value)
	if err != nil || parsed.User == nil {
		return value
	}
	if _, ok := parsed.User.Password(); !ok {
		return value
	}
	parsed.User = url.UserPassword(parsed.User.Username(), replacement)
	return parsed.String()
}
