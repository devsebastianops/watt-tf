package config

type Config struct {
	Include   []string        `yaml:"include,omitempty"`
	Transform []Transformable `yaml:"transform"`
}

type Transformable struct {
	Target  string                 `yaml:"target"`
	If      string                 `yaml:"if,omitempty"`
	Value   map[string]interface{} `yaml:"value,omitempty"`
	ForEach string                 `yaml:"for_each,omitempty"`
}
