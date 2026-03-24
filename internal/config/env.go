package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// env returns the value of the environment variable identified by key.
// If the variable is not set or is empty, it returns the provided fallback value.
func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// envInt returns the integer value of the environment variable identified by key.
// If the variable is unset or empty, it returns the fallback value.
// It panics if the value cannot be parsed as an integer.
func envInt(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}

	n, err := strconv.Atoi(v)
	if err != nil {
		panic(fmt.Sprintf("env %q  must be an integer, got %q", key, v))
	}

	return n
}

// envBool retrieves an environment variable by key and attempts to parse it as a boolean.
// If the variable is not set or is empty, the provided fallback value is returned.
// If the value cannot be parsed as a boolean, the function panics with a descriptive error message.
// Supported boolean values are those accepted by strconv.ParseBool (e.g., "true", "false", "1", "0").
func envBool(key string, fallback bool) bool {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}

	b, err := strconv.ParseBool(v)
	if err != nil {
		panic(fmt.Sprintf("env %q  must be a boolean (e.g. \"true\", \"false\", \"1\", \"0\"), got %q", key, v))
	}

	return b
}

// envDuration returns the duration value of the environment variable identified by key.
// If the variable is unset or empty, it returns the fallback value.
// It panics if the value cannot be parsed as a valid duration
// (e.g. "5s", "1m", "1h").
func envDuration(key string, fallback time.Duration) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}

	d, err := time.ParseDuration(v)
	if err != nil {
		panic(fmt.Sprintf("env %q  must be a duration (e.g. 5s, 1m), got %q", key, v))
	}

	return d
}
