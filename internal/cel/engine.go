package cel

import (
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/ext"
)

type Engine struct {
	env     *cel.Env
	context Context
}

func BuildCelEngine() (Engine, error) {

	context := NewContext()

	env, err := cel.NewEnv(
		cel.Variable("input", cel.MapType(cel.StringType, cel.AnyType)),
		cel.Variable("env", cel.MapType(cel.StringType, cel.StringType)),
		cel.Variable("item", cel.AnyType),
		cel.Variable("item_index", cel.IntType),
		cel.Variable("vars", cel.MapType(cel.StringType, cel.AnyType)),
		cel.Variable("data", cel.MapType(cel.StringType, cel.AnyType)),
		cel.Macros(cel.StandardMacros...),
		cel.OptionalTypes(),
		ext.Encoders(),
		ext.Lists(),
		ext.Strings(),
		ext.Sets(),
		ext.Regex(),
	)

	if err != nil {
		return Engine{}, err
	}

	env, err = registerFunctions(env)
	if err != nil {
		return Engine{}, err
	}

	return Engine{
		env:     env,
		context: context,
	}, nil
}

func (e *Engine) Context() *Context {
	return &e.context
}

func (e *Engine) Evaluate(expression string) (ref.Val, error) {
	ast, iss := e.env.Compile(expression)
	if iss != nil && iss.Err() != nil {
		return nil, iss.Err()
	}

	prg, err := e.env.Program(ast)
	if err != nil {
		return nil, err
	}

	each := e.context.GetEach()
	out, _, err := prg.Eval(map[string]any{
		"input":      e.context.GetInputs(),
		"env":        e.context.GetEnv(),
		"item":       each.GetItem(),
		"item_index": each.GetItemIndex(),
		"vars":       e.context.GetVars(),
		"data":       e.context.GetData(),
	})
	if err != nil {
		return nil, err
	}

	return out, nil
}
