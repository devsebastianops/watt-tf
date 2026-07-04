package parser

import (
	"errors"
	"strings"
)

func ParseInput(input string) (map[string]interface{}, error) {
	if isJson(input) {
		return parseJSON(input)
	}

	if isYaml(input) {
		return parseYAML(input)
	}

	return nil, errors.New("unsupported input format")
}

func isJson(input string) bool {
	// Input ends with .json
	if strings.HasSuffix(input, ".json") {
		return true
	}

	return false
}

func isYaml(input string) bool {
	// Input ends with .yaml or .yml or starts with ---
	if strings.HasSuffix(input, ".yaml") || strings.HasSuffix(input, ".yml") {
		return true
	}

	return false
}
