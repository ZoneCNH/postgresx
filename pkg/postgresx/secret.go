package postgresx

// SecretString holds a sensitive value that masks itself in logs.
type SecretString string

// NewSecretString creates a secret from a raw string.
func NewSecretString(s string) SecretString { return SecretString(s) }

// IsZero reports whether the secret is empty.
func (s SecretString) IsZero() bool { return s == "" }

// String returns a masked representation for diagnostics.
func (s SecretString) String() string { return "***" }

// Reveal returns the raw value for use with the driver.
func (s SecretString) Reveal() string { return string(s) }
