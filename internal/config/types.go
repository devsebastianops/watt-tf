package config

import "github.com/devsebastianops/watt-tf/internal/plugin"

type Config struct {
	Variables map[string]any  `yaml:"variables,omitempty"`
	Plugins   []plugin.Plugin `yaml:"plugins,omitempty"`
	Include   []string        `yaml:"include,omitempty"`
	Transform []Transformable `yaml:"transform"`
}

type Transformable struct {
	Target  string `yaml:"target"`
	If      string `yaml:"if,omitempty"`
	Value   any    `yaml:"value,omitempty"`
	ForEach string `yaml:"for_each,omitempty"`
}
