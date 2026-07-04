package config

type Config struct {
	Transform []Transformable `yaml:"transform"`
}

type Transformable struct {
	Target string                 `yaml:"target"`
	If     string                 `yaml:"if,omitempty"`
	Value  map[string]interface{} `yaml:"value"`
}
