package parser

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// parseYAML reads a YAML file and returns its contents as a map
func parseYAML(filePath string) (map[string]interface{}, error) {
	// Read file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read YAML file: %w", err)
	}

	// Unmarshal YAML
	var result map[string]interface{}
	if err := yaml.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	return result, nil
}

func ParseYAML(input string) (map[string]interface{}, error) {
	if isYaml(input) {
		return parseYAML(input)
	}

	return nil, fmt.Errorf("unsupported input format for YAML")
}
