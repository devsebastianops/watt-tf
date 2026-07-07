package plugin

type Registry struct {
	Plugins map[string][]Plugin
}

func NewRegistry() *Registry {
	return &Registry{
		Plugins: make(map[string][]Plugin),
	}
}

func (r *Registry) Register(plugin Plugin) {
	r.Plugins[plugin.On] = append(r.Plugins[plugin.On], plugin)
}

func (r *Registry) RegisterPlugins(plugins []Plugin) {
	for _, plugin := range plugins {
		r.Register(plugin)
	}
}

func (r *Registry) GetPluginsForEvent(event Event) []Plugin {
	return r.Plugins[string(event)]
}
