package envflags

import (
	"fmt"
	"os"
	"strconv"
)

// GetStringDefault looks up Environment Variable for a key.
// Variable value is returned if it is non-empty.
// If it is empty default value is returned.
func GetStringDefault(envKey, def string) string {
	v := os.Getenv(envKey)
	if v != "" {
		return v
	}

	return def
}

// GetIntDefault looks up Environment Variable for a key.
// Variable value is converted if it is non-empty.
// If it is empty default value is returned.
// Panics when conversion fails.
func GetIntDefault(envKey string, def int) int {
	v := os.Getenv(envKey)
	if v != "" {
		i, err := strconv.Atoi(v)
		if err != nil {
			panic(fmt.Sprintf("Provided value for Environment Variable %q not convertible to int: %s", envKey, v))
		}
		return i
	}
	return def
}
