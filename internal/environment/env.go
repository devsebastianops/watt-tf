package environment

import (
	"os"
	"strings"
)

// GetEnvVars collects all environment variables into a map
func GetEnvVars() map[string]string {
	envMap := make(map[string]string)
	for _, envVar := range os.Environ() {
		key, value, _ := strings.Cut(envVar, "=")
		envMap[key] = value
	}
	return envMap
}
