package transformer

import (
	"fmt"

	"github.com/google/cel-go/cel"
)

func evalCelCondition(expr string, env *cel.Env, inputData map[string]any, envVars map[string]string) (bool, error) {
	ast, iss := env.Compile(expr)
	if iss.Err() != nil {
		return false, fmt.Errorf("syntax error: %v", iss.Err())
	}
	program, err := env.Program(ast)
	if err != nil {
		return false, err
	}
	out, _, err := program.Eval(map[string]any{"input": inputData, "env": envVars})
	if err != nil {
		return false, err
	}

	boolVal, ok := out.Value().(bool)
	if !ok {
		return false, fmt.Errorf("expression did not return a boolean value")
	}
	return boolVal, nil
}
