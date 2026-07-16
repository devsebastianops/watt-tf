package schema

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/devsebastianops/watt-tf/internal/logger"
	"github.com/xeipuuv/gojsonschema"
)

func ValidateInputSchema(input map[string]any, schemaFile string) error {
	// Read schema file
	schemaData, err := os.ReadFile(schemaFile)
	if err != nil {
		return fmt.Errorf("failed to read schema file '%s': %w", schemaFile, err)
	}

	// Parse schema
	schemaLoader := gojsonschema.NewBytesLoader(schemaData)
	schema, err := gojsonschema.NewSchema(schemaLoader)
	if err != nil {
		return fmt.Errorf("failed to parse schema file '%s': %w", schemaFile, err)
	}

	// Convert input to JSON for validation
	inputJSON, err := json.Marshal(input)
	if err != nil {
		return fmt.Errorf("failed to marshal input for validation: %w", err)
	}

	// Validate input against schema
	inputLoader := gojsonschema.NewBytesLoader(inputJSON)
	result, err := schema.Validate(inputLoader)
	if err != nil {
		return fmt.Errorf("schema validation error: %w", err)
	}

	// Check validation result
	if !result.Valid() {
		logger.Error("input validation failed")
		errorMsg := "Schema validation errors:\n"
		for i, desc := range result.Errors() {
			errorMsg += fmt.Sprintf("  %d. %s\n", i+1, desc.String())
		}
		return fmt.Errorf("%s", errorMsg)
	}

	return nil
}
