package config

import (
	"github.com/devsebastianops/watt-tf/internal/parser"
)

func LoadConfig(filePath string) (*Config, error) {
	configMap, err := parser.ParseYAML(filePath)
	if err != nil {
		return nil, err
	}

	config := Config{
		Transform: []Transformable{},
	}

	for _, transformable := range configMap["transform"].([]interface{}) {
		transformableMap := transformable.(map[string]interface{})
		target := transformableMap["target"].(string)
		value := transformableMap["value"].(map[string]interface{})
		condition, ok := transformableMap["if"].(string)
		if ok {
			config.Transform = append(config.Transform, Transformable{
				Target: target,
				If:     condition,
				Value:  value,
			})
		} else {
			config.Transform = append(config.Transform, Transformable{
				Target: target,
				Value:  value,
			})
		}
	}

	return &config, nil
}
