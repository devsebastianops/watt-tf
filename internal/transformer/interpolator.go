package transformer

import (
	"fmt"
	"strings"

	"github.com/devsebastianops/watt-tf/internal/logger"
	"github.com/google/cel-go/cel"
)

// isMissingKeyError checks if the error is about a missing key
func isMissingKeyError(err error) bool {
	if err == nil {
		return false
	}
	errMsg := err.Error()
	return strings.Contains(errMsg, "no such key") ||
		strings.Contains(errMsg, "no such field") ||
		strings.Contains(errMsg, "undefined reference")
}

func interpolate(val any, env *cel.Env, inputData map[string]any, envVars map[string]string, strict bool) (any, error) {
	// Check the type of val and handle accordingly
	switch v := val.(type) {

	case string:
		// Is it a string?
		hasInterpolation := hasInterpolation(v)
		if !hasInterpolation {
			return v, nil // Just a regular string, return as is
		}
		matches := findInterpolations(v)

		// There is only one match and the whole string is the interpolation, we can evaluate it directly
		if len(matches) == 1 && matches[0].Start == 0 && matches[0].End == len(v) {
			return evalCelExpression(matches[0].Expr, env, inputData, envVars, strict)
		}

		// There are multiple matches or the interpolation is part of a larger string, we need to replace them
		resultStr := v
		for _, match := range matches {
			fullMatch := v[match.Start:match.End] // ${input.env}
			expression := match.Expr              // input.env

			celVal, err := evalCelExpression(expression, env, inputData, envVars, strict)
			if err != nil {
				return nil, err
			}

			// Handle nil values in string interpolation
			var replacement string
			if celVal == nil {
				replacement = "null"
			} else {
				replacement = fmt.Sprintf("%v", celVal)
			}
			resultStr = strings.Replace(resultStr, fullMatch, replacement, 1)
		}
		return resultStr, nil

	case map[string]any:
		// If the value is a map, we recursively interpolate its items
		newMap := make(map[string]any)
		for k, item := range v {
			// Handle nil values explicitly
			if item == nil {
				newMap[k] = nil
				continue
			}

			interpItem, err := interpolate(item, env, inputData, envVars, strict)
			if err != nil {
				return nil, err
			}
			newMap[k] = interpItem
		}
		return newMap, nil

	case []interface{}:
		// If the value is an array, recursively interpolate its items
		newArray := make([]interface{}, 0, len(v))
		for _, item := range v {
			// Handle nil values explicitly
			if item == nil {
				newArray = append(newArray, nil)
				continue
			}

			interpItem, err := interpolate(item, env, inputData, envVars, strict)
			if err != nil {
				return nil, err
			}
			newArray = append(newArray, interpItem)
		}
		return newArray, nil

	default:
		return v, nil // For other types (like int, float, bool), return as is
	}
}

func evalCelExpression(expr string, env *cel.Env, inputData map[string]any, envVars map[string]string, strict bool) (any, error) {
	ast, iss := env.Compile(expr)
	if iss.Err() != nil {
		return nil, fmt.Errorf("syntax error in interpolation '%s': %v", expr, iss.Err())
	}
	program, err := env.Program(ast)
	if err != nil {
		return nil, err
	}
	out, _, err := program.Eval(map[string]any{"input": inputData, "env": envVars})
	if err != nil {
		// Handle missing key errors based on strict mode
		if isMissingKeyError(err) {
			if strict {
				return nil, err
			}
			// In lenient mode, log warning and return nil
			logger.Warn("missing key in CEL expression, using null", "expression", expr, "error", err.Error())
			return nil, nil
		}
		return nil, err
	}

	// Get the value and handle CEL null
	val := out.Value()

	// Handle nil and null-like values
	// CEL represents null as structpb.NullValue with string representation "NULL_VALUE"
	if val == nil {
		return nil, nil
	}

	// CEL might represent null differently - check for NULL_VALUE
	valStr := fmt.Sprintf("%v", val)
	if valStr == "<nil>" || valStr == "NULL_VALUE" {
		return nil, nil
	}

	return val, nil
}

// interpolateWithItem is like interpolate but also includes item and item_index variables
func interpolateWithItem(val any, env *cel.Env, inputData map[string]any, envVars map[string]string, item any, itemIndex int, strict bool) (any, error) {
	switch v := val.(type) {

	case string:
		hasInterpolation := hasInterpolation(v)
		if !hasInterpolation {
			return v, nil
		}

		matches := findInterpolations(v)
		if len(matches) == 1 && matches[0].Start == 0 && matches[0].End == len(v) {
			return evalCelExpressionWithItem(matches[0].Expr, env, inputData, envVars, item, itemIndex, strict)
		}

		resultStr := v
		for _, match := range matches {
			fullMatch := v[match.Start:match.End]
			expression := match.Expr

			celVal, err := evalCelExpressionWithItem(expression, env, inputData, envVars, item, itemIndex, strict)
			if err != nil {
				return nil, err
			}

			var replacement string
			if celVal == nil {
				replacement = "null"
			} else {
				replacement = fmt.Sprintf("%v", celVal)
			}
			resultStr = strings.Replace(resultStr, fullMatch, replacement, 1)
		}
		return resultStr, nil

	case map[string]any:
		newMap := make(map[string]any)
		for k, mapItem := range v {
			if mapItem == nil {
				newMap[k] = nil
				continue
			}

			interpItem, err := interpolateWithItem(mapItem, env, inputData, envVars, item, itemIndex, strict)
			if err != nil {
				return nil, err
			}
			newMap[k] = interpItem
		}
		return newMap, nil

	case []interface{}:
		newArray := make([]interface{}, 0, len(v))
		for _, arrayItem := range v {
			if arrayItem == nil {
				newArray = append(newArray, nil)
				continue
			}

			interpItem, err := interpolateWithItem(arrayItem, env, inputData, envVars, item, itemIndex, strict)
			if err != nil {
				return nil, err
			}
			newArray = append(newArray, interpItem)
		}
		return newArray, nil

	default:
		return v, nil
	}
}

func evalCelExpressionWithItem(expr string, env *cel.Env, inputData map[string]any, envVars map[string]string, item any, itemIndex int, strict bool) (any, error) {
	ast, iss := env.Compile(expr)
	if iss.Err() != nil {
		return nil, fmt.Errorf("syntax error in interpolation '%s': %v", expr, iss.Err())
	}
	program, err := env.Program(ast)
	if err != nil {
		return nil, err
	}
	out, _, err := program.Eval(map[string]any{
		"input":      inputData,
		"env":        envVars,
		"item":       item,
		"item_index": itemIndex,
	})
	if err != nil {
		// Handle missing key errors based on strict mode
		if isMissingKeyError(err) {
			if strict {
				return nil, err
			}
			// In lenient mode, log warning and return nil
			logger.Warn("missing key in CEL expression, using null", "expression", expr, "error", err.Error())
			return nil, nil
		}
		return nil, err
	}

	val := out.Value()

	if val == nil {
		return nil, nil
	}

	valStr := fmt.Sprintf("%v", val)
	if valStr == "<nil>" || valStr == "NULL_VALUE" {
		return nil, nil
	}

	return val, nil
}
