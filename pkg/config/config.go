package config

import "os"

// Get returns the environment variable value for the given key.
// Returns empty string if not set.
func Get(key string) string {
	return os.Getenv(key)
}

// MustGet returns the environment variable value, panics if empty.
func MustGet(key string) string {
	v := os.Getenv(key)
	if v == "" {
		panic("required env var not set: " + key)
	}
	return v
}
