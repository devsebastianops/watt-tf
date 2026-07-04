package writer

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/devsebastianops/watt-tf/internal/logger"
)

// WriteJSON writes the result map to a JSON file
func WriteJSON(data map[string]interface{}, filePath string) error {
	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal result to JSON: %w", err)
	}

	err = os.WriteFile(filePath, jsonBytes, 0644)
	if err != nil {
		return fmt.Errorf("failed to write JSON to file %s: %w", filePath, err)
	}

	logger.Info("output file written successfully", "path", filePath)
	return nil
}
