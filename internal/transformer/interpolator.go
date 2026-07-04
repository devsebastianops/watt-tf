package transformer

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/google/cel-go/cel"
)

var interpRegex = regexp.MustCompile(`\${([^}]+)}`)

func interpolate(val any, env *cel.Env, inputData map[string]any, envVars map[string]string) (any, error) {
	// Check the type of val and handle accordingly
	switch v := val.(type) {

	case string:
		// Is it a string?
		matches := interpRegex.FindAllStringSubmatch(v, -1)
		if len(matches) == 0 {
			return v, nil // Just a regular string, return as is
		}

		// There is only one match and the whole string is the interpolation, we can evaluate it directly
		if len(matches) == 1 && matches[0][0] == v {
			return evalCelExpression(matches[0][1], env, inputData, envVars)
		}

		// There are multiple matches or the interpolation is part of a larger string, we need to replace them
		resultStr := v
		for _, match := range matches {
			fullMatch := match[0]  // ${input.env}
			expression := match[1] // input.env

			celVal, err := evalCelExpression(expression, env, inputData, envVars)
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

			interpItem, err := interpolate(item, env, inputData, envVars)
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

			interpItem, err := interpolate(item, env, inputData, envVars)
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

func evalCelExpression(expr string, env *cel.Env, inputData map[string]any, envVars map[string]string) (any, error) {
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
