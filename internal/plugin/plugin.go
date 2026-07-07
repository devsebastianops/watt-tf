package plugin

type Plugin struct {
	Name    string
	Version string
	On      string
	Cmd     string
	Args    []string
}
