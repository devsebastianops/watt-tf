package transformer

import (
	"encoding/base64"
	"encoding/json"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
)

// toJSONImpl converts any value to JSON string
func toJSONImpl(val ref.Val) ref.Val {
	data, err := json.Marshal(val.Value())
	if err != nil {
		return types.NewErr("failed to convert to JSON: %v", err)
	}
	return types.String(string(data))
}

// fromJSONImpl parses JSON string to a value
func fromJSONImpl(val ref.Val) ref.Val {
	jsonStr := val.(types.String).Value().(string)
	var result interface{}
	err := json.Unmarshal([]byte(jsonStr), &result)
	if err != nil {
		return types.NewErr("failed to parse JSON: %v", err)
	}
	return types.DefaultTypeAdapter.NativeToValue(result)
}

// toBase64Impl encodes string to base64
func toBase64Impl(val ref.Val) ref.Val {
	str := val.(types.String).Value().(string)
	encoded := base64.StdEncoding.EncodeToString([]byte(str))
	return types.String(encoded)
}

// fromBase64Impl decodes base64 string
func fromBase64Impl(val ref.Val) ref.Val {
	str := val.(types.String).Value().(string)
	decoded, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return types.NewErr("failed to decode base64: %v", err)
	}
	return types.String(string(decoded))
}

// RegisterWtfFunctions registers custom wtf helper functions in the CEL environment
func RegisterWtfFunctions(env *cel.Env) (*cel.Env, error) {
	var err error

	// Register toJSON function
	env, err = env.Extend(
		cel.Function("toJSON",
			cel.Overload("toJSON_any",
				[]*cel.Type{cel.AnyType},
				cel.StringType,
				cel.UnaryBinding(toJSONImpl),
			),
		),
	)
	if err != nil {
		return nil, err
	}

	// Register fromJSON function
	env, err = env.Extend(
		cel.Function("fromJSON",
			cel.Overload("fromJSON_string",
				[]*cel.Type{cel.StringType},
				cel.AnyType,
				cel.UnaryBinding(fromJSONImpl),
			),
		),
	)
	if err != nil {
		return nil, err
	}

	// Register toBase64 function
	env, err = env.Extend(
		cel.Function("toBase64",
			cel.Overload("toBase64_string",
				[]*cel.Type{cel.StringType},
				cel.StringType,
				cel.UnaryBinding(toBase64Impl),
			),
		),
	)
	if err != nil {
		return nil, err
	}

	// Register fromBase64 function
	env, err = env.Extend(
		cel.Function("fromBase64",
			cel.Overload("fromBase64_string",
				[]*cel.Type{cel.StringType},
				cel.StringType,
				cel.UnaryBinding(fromBase64Impl),
			),
		),
	)
	if err != nil {
		return nil, err
	}

	return env, nil
}
