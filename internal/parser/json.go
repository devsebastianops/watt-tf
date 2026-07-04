package parser

import (
	"encoding/json"
	"fmt"
	"os"
)

// parseJSON reads a JSON file and returns its contents as a map
func parseJSON(filePath string) (map[string]interface{}, error) {
	// Read file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read JSON file: %w", err)
	}

	// Unmarshal JSON
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return result, nil
}
