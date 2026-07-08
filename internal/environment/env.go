package environment

import (
	"os"
	"strings"
)

func GetEnvVars() map[string]string {
	envMap := make(map[string]string)
	for _, envVar := range os.Environ() {
		key, value, _ := strings.Cut(envVar, "=")
		envMap[key] = value
	}
	return envMap
}
