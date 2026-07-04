package transformer

import (
	"fmt"
	"os"
	"strings"

	"github.com/devsebastianops/watt-tf/internal/config"
	"github.com/devsebastianops/watt-tf/internal/logger"
	"github.com/google/cel-go/cel"
)

func Transform(input map[string]interface{}, config *config.Config) (map[string]interface{}, error) {

	transformables := config.Transform
	result := map[string]interface{}{}
	envVars := getEnvVars()

	env, _ := cel.NewEnv(
		cel.Variable("input", cel.MapType(cel.StringType, cel.AnyType)),
		cel.Variable("env", cel.MapType(cel.StringType, cel.StringType)),
		cel.Macros(cel.StandardMacros...),
	)

	for _, transformable := range transformables {
		target := transformable.Target
		value := transformable.Value
		condition := transformable.If

		// 0. Interpolate target (supports dynamic names like: resource.aws_s3.${input.name})
		interpolatedTarget, err := interpolate(target, env, input, envVars)
		if err != nil {
			return nil, fmt.Errorf("failed to interpolate target '%s': %w", target, err)
		}
		target = interpolatedTarget.(string) // Target must result in a string

		logger.Debug("processing transformation", "target", target, "has_condition", condition != "")

		// 1. Evaluate condition (if specified)
		if condition != "" {
			logger.Debug("evaluating condition", "target", target, "condition", condition)
			shouldExecute, err := evalCelCondition(condition, env, input, envVars)
			if err != nil {
				return nil, fmt.Errorf("failed to evaluate condition '%s': %w", condition, err)
			}
			if !shouldExecute {
				logger.Debug("condition not met, skipping", "target", target)
				continue
			}
		}

		// 2. Interpolate
		interpolatedValue, err := interpolate(value, env, input, envVars)
		if err != nil {
			return nil, err
		}

		// 3. Span up the result (Unflattening)
		logger.Debug("interpolation completed", "target", target)
		// Wir casten interpolatedValue zu map[string]any, da dein Value-Typ im Struct so definiert ist
		if mapVal, ok := interpolatedValue.(map[string]any); ok {
			unflatten(result, target, mapVal)
		} else {
			// Falls der Nutzer ein primitives Value wie einen nackten String/Int übergibt
			unflatten(result, target, interpolatedValue)
		}
	}

	// Return the result map (no file writing here)
	logger.Debug("all transformations completed successfully")
	return result, nil
}

func unflatten(result map[string]any, path string, value any) {
	parts := strings.Split(path, ".")
	current := result

	for i := 0; i < len(parts)-1; i++ {
		part := parts[i]

		if _, exists := current[part]; !exists {
			current[part] = make(map[string]any)
		}

		// Typ-Absicherung: Falls der Pfad bereits existiert, MUSS es eine Map sein
		if nextMap, ok := current[part].(map[string]any); ok {
			current = nextMap
		} else {
			// Fallback, falls jemand Pfadkonflikte erzeugt (z.B. target: a und target: a.b)
			newMap := make(map[string]any)
			current[part] = newMap
			current = newMap
		}
	}

	lastPart := parts[len(parts)-1]

	// Deep Merge für das finale Value, falls dort bereits Daten liegen
	if existingMap, ok := current[lastPart].(map[string]any); ok {
		if newMap, ok := value.(map[string]any); ok {
			for k, v := range newMap {
				existingMap[k] = v
			}
			return
		}
	}

	current[lastPart] = value
}

// getEnvVars collects all environment variables into a map
func getEnvVars() map[string]string {
	envMap := make(map[string]string)
	for _, envVar := range os.Environ() {
		key, value, _ := strings.Cut(envVar, "=")
		envMap[key] = value
	}
	return envMap
}
