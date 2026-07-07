package plugin

// Context is used to pass information between plugin calls.
// This enables plugin pipes to be built, where the output of one plugin can be used as the input for another.
type Context struct {
	Input  map[string]interface{} `json:"input"`
	Env    map[string]string      `json:"env"`
	Result map[string]interface{} `json:"result"`
}

func newContext(input map[string]interface{}, env map[string]interface{}) Context {
	envStr := make(map[string]string)
	for k, v := range env {
		if strVal, ok := v.(string); ok {
			envStr[k] = strVal
		}
	}
	return Context{
		Input:  input,
		Env:    envStr,
		Result: make(map[string]interface{}),
	}
}
