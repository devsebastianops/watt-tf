package plugin

type PluginRequest struct {
	Version string  `json:"version"`
	Event   string  `json:"event"`
	Data    Context `json:"data"`
}

type PluginResponse struct {
	Status string  `json:"status"`
	Data   Context `json:"data,omitempty"`
	Error  string  `json:"error,omitempty"`
}
