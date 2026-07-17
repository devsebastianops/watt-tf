package transformer

import (
	"fmt"
	"strings"

	"github.com/devsebastianops/watt-tf/internal/config"
	"github.com/devsebastianops/watt-tf/internal/logger"
	"github.com/google/cel-go/cel"
)

func Transform(input map[string]interface{}, envVars map[string]string, config *config.Config, strict bool) (map[string]interface{}, error) {

	transformables := config.Transform
	result := map[string]interface{}{}

	if strict {
		logger.Info("running in strict mode: missing keys will cause errors")
	} else {
		logger.Info("running in lenient mode: missing keys will be replaced with null")
	}

	baseEnv, _ := cel.NewEnv(
		cel.Variable("input", cel.MapType(cel.StringType, cel.AnyType)),
		cel.Variable("env", cel.MapType(cel.StringType, cel.StringType)),
		cel.Variable("item", cel.AnyType),
		cel.Variable("item_index", cel.IntType),
		cel.Macros(cel.StandardMacros...),
		cel.OptionalTypes(),
	)

	// Register custom wtf functions
	logger.Info("registering custom wtf functions")
	env, err := RegisterWtfFunctions(baseEnv)
	if err != nil {
		logger.Error("failed to register wtf functions", "error", err.Error())
		return nil, fmt.Errorf("failed to register wtf functions: %w", err)
	}
	logger.Info("custom wtf functions registered successfully")

	for _, transformable := range transformables {
		target := transformable.Target
		value := transformable.Value
		condition := transformable.If
		forEach := transformable.ForEach

		// Check if this is a for_each transformation
		if forEach != "" {
			logger.Debug("processing for_each transformation", "target", target, "for_each", forEach)

			// Evaluate the for_each expression to get the array
			forEachCompiled, iss := env.Compile(forEach)
			if iss.Err() != nil {
				return nil, fmt.Errorf("failed to compile for_each expression '%s': %w", forEach, iss.Err())
			}

			// Evaluate with base context (no item yet)
			forEachProgram, err := env.Program(forEachCompiled)
			if err != nil {
				return nil, fmt.Errorf("failed to create for_each program '%s': %w", forEach, err)
			}

			evalResult, _, err := forEachProgram.Eval(map[string]any{
				"input":      input,
				"env":        envVars,
				"item":       nil,
				"item_index": 0,
			})
			if err != nil {
				return nil, fmt.Errorf("failed to evaluate for_each expression '%s': %w", forEach, err)
			}

			// Convert items to array
			var itemArray []interface{}
			switch v := evalResult.Value().(type) {
			case []interface{}:
				itemArray = v
			default:
				return nil, fmt.Errorf("for_each expression must return an array, got %T", v)
			}

			logger.Debug("for_each resolved to array", "count", len(itemArray))

			// Process each item in the array
			for idx, item := range itemArray {
				logger.Debug("processing for_each item", "index", idx, "item_type", fmt.Sprintf("%T", item))

				// Interpolate target with item context
				interpolatedTarget, err := interpolateWithItem(target, env, input, envVars, item, idx, strict)
				if err != nil {
					return nil, fmt.Errorf("failed to interpolate target '%s' for item %d: %w", target, idx, err)
				}
				targetStr := interpolatedTarget.(string)

				// Evaluate condition with item context
				if condition != "" {
					logger.Debug("evaluating condition with item", "condition", condition, "index", idx)
					shouldExecute, err := evalCelConditionWithItem(condition, env, input, envVars, item, idx, strict)
					if err != nil {
						return nil, fmt.Errorf("failed to evaluate condition '%s' for item %d: %w", condition, idx, err)
					}
					if !shouldExecute {
						logger.Debug("condition not met for item, skipping", "index", idx)
						continue
					}
				}

				// Interpolate value with item context
				interpolatedValue, err := interpolateWithItem(value, env, input, envVars, item, idx, strict)
				if err != nil {
					return nil, fmt.Errorf("failed to interpolate value for item %d: %w", idx, err)
				}

				// Unflatten the result
				if mapVal, ok := interpolatedValue.(map[string]any); ok {
					unflatten(result, targetStr, mapVal)
				} else {
					unflatten(result, targetStr, interpolatedValue)
				}
			}

		} else {
			// Standard transformation (no for_each)
			// 0. Interpolate target (supports dynamic names like: resource.aws_s3.${input.name})
			interpolatedTarget, err := interpolate(target, env, input, envVars, strict)
			if err != nil {
				return nil, fmt.Errorf("failed to interpolate target '%s': %w", target, err)
			}
			target = interpolatedTarget.(string) // Target must result in a string

			logger.Debug("processing transformation", "target", target, "has_condition", condition != "")

			// 1. Evaluate condition (if specified)
			if condition != "" {
				logger.Debug("evaluating condition", "target", target, "condition", condition)
				shouldExecute, err := evalCelCondition(condition, env, input, envVars, strict)
				if err != nil {
					return nil, fmt.Errorf("failed to evaluate condition '%s': %w", condition, err)
				}
				if !shouldExecute {
					logger.Debug("condition not met, skipping", "target", target)
					continue
				}
			}

			// 2. Interpolate
			interpolatedValue, err := interpolate(value, env, input, envVars, strict)
			if err != nil {
				return nil, err
			}

			// 3. Span up the result (Unflattening)
			logger.Debug("interpolation completed", "target", target)
			if mapVal, ok := interpolatedValue.(map[string]any); ok {
				unflatten(result, target, mapVal)
			} else {
				unflatten(result, target, interpolatedValue)
			}
		}
	}

	// Return the result map (no file writing here)
	logger.Debug("all transformations completed successfully")
	return result, nil
}

// parsePath zerlegt den Pfad in Segmente und ignoriert Punkte innerhalb von Backticks
func parsePath(path string) []string {
	var parts []string
	var current strings.Builder
	inBackticks := false

	for i := 0; i < len(path); i++ {
		char := path[i]

		switch char {
		case '`':
			inBackticks = !inBackticks
			// Die Backticks selbst wollen wir nicht im finalen JSON-Key haben,
			// darum schreiben wir sie nicht in den current-Buffer.
		case '.':
			if inBackticks {
				// Wenn wir in Backticks sind, ist der Punkt Teil des Keys!
				current.WriteByte(char)
			} else {
				// Außerhalb von Backticks trennt der Punkt das Segment
				if current.Len() > 0 {
					parts = append(parts, current.String())
					current.Reset()
				}
			}
		default:
			current.WriteByte(char)
		}
	}

	if current.Len() > 0 {
		parts = append(parts, current.String())
	}

	return parts
}

func unflatten(result map[string]any, path string, value any) {
	// Nutze den intelligenten Parser statt strings.Split
	parts := parsePath(path)
	if len(parts) == 0 {
		return
	}

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
			// Fallback bei Pfadkonflikten (z.B. target: a und target: a.b)
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
