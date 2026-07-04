package config

import (
	"fmt"

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

	transformList, ok := configMap["transform"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid transform list in config")
	}

	for _, transformable := range transformList {
		transformableMap := transformable.(map[string]interface{})
		target, ok := transformableMap["target"].(string)
		if !ok {
			return nil, fmt.Errorf("missing or invalid 'target' field")
		}

		// Parse condition if present
		condition, _ := transformableMap["if"].(string)

		// Parse value (if template is not used)
		var value map[string]interface{}
		if val, ok := transformableMap["value"].(map[string]interface{}); ok {
			value = val
		}

		// Parse for_each (for iteration)
		forEach, _ := transformableMap["for_each"].(string)

		// Validate: either value or template must be present
		if value == nil {
			return nil, fmt.Errorf("transform at target '%s' must have 'value' defined", target)
		}

		config.Transform = append(config.Transform, Transformable{
			Target:  target,
			If:      condition,
			Value:   value,
			ForEach: forEach,
		})
	}

	return &config, nil
}
