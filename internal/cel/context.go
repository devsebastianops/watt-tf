package cel

type Context struct {
	Input map[string]any
	Env   map[string]string
	Vars  map[string]any
	Data  map[string]any
	Each  EachContext
}

func NewContext() Context {
	return Context{
		Input: make(map[string]any),
		Env:   make(map[string]string),
		Vars:  make(map[string]any),
		Data:  make(map[string]any),
		Each:  NewEachContext(),
	}
}

func (c *Context) SetInputs(inputs map[string]any) {
	c.Input = inputs
}

func (c *Context) GetInputs() map[string]any {
	return c.Input
}

func (c *Context) SetEnv(env map[string]string) {
	c.Env = env
}

func (c *Context) GetEnv() map[string]string {
	return c.Env
}

func (c *Context) SetVars(vars map[string]any) {
	c.Vars = vars
}

func (c *Context) GetVars() map[string]any {
	return c.Vars
}

func (c *Context) SetData(data map[string]any) {
	c.Data = data
}

func (c *Context) GetData() map[string]any {
	return c.Data
}

func (c *Context) SetEach(each EachContext) {
	c.Each = each
}

func (c *Context) GetEach() EachContext {
	return c.Each
}
