package e2e

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/devsebastianops/watt-tf/internal/config"
	"github.com/devsebastianops/watt-tf/internal/parser"
	"github.com/devsebastianops/watt-tf/internal/transformer"
)

// TestE2EExamples runs end-to-end tests for all examples
func TestE2EExamples(t *testing.T) {
	exampleDir := "../../example"

	entries, err := os.ReadDir(exampleDir)
	if err != nil {
		t.Fatalf("failed to read example directory: %v", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		exampleName := entry.Name()
		examplePath := filepath.Join(exampleDir, exampleName)

		t.Run(exampleName, func(t *testing.T) {
			runE2ETest(t, examplePath)
		})
	}
}

// runE2ETest runs a single example through the full pipeline
func runE2ETest(t *testing.T, examplePath string) {
	exampleName := filepath.Base(examplePath)

	// Set up environment variables for specific test cases
	if exampleName == "13-env-variables" {
		os.Setenv("DEPLOYMENT_ENV", "production")
		os.Setenv("DEPLOYMENT_REGION", "us-east-1")
		// Restore original values after test if they existed
		defer func() {
			if original, exists := os.LookupEnv("DEPLOYMENT_ENV"); exists {
				os.Setenv("DEPLOYMENT_ENV", original)
			} else {
				os.Unsetenv("DEPLOYMENT_ENV")
			}
			if original, exists := os.LookupEnv("DEPLOYMENT_REGION"); exists {
				os.Setenv("DEPLOYMENT_REGION", original)
			} else {
				os.Unsetenv("DEPLOYMENT_REGION")
			}
		}()
	}

	// Find input file (either .json or .yaml)
	var inputFile string
	for _, ext := range []string{".json", ".yaml", ".yml"} {
		candidate := filepath.Join(examplePath, "input"+ext)
		if _, err := os.Stat(candidate); err == nil {
			inputFile = candidate
			break
		}
	}

	if inputFile == "" {
		t.Fatalf("no input file found in %s", examplePath)
	}

	configFile := filepath.Join(examplePath, ".wtf.yaml")
	if _, err := os.Stat(configFile); err != nil {
		t.Fatalf("config file not found: %s", configFile)
	}

	expectedFile := filepath.Join(examplePath, "expected.tf.json")
	if _, err := os.Stat(expectedFile); err != nil {
		t.Fatalf("expected output file not found: %s", expectedFile)
	}

	// 1. Parse input
	input, err := parser.ParseInput(inputFile)
	if err != nil {
		t.Fatalf("failed to parse input: %v", err)
	}

	// 2. Load config
	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	// 3. Transform
	result, err := transformer.Transform(input, cfg, false)
	if err != nil {
		t.Fatalf("transformation failed: %v", err)
	}

	// 4. Load expected output
	expectedData, err := os.ReadFile(expectedFile)
	if err != nil {
		t.Fatalf("failed to read expected output: %v", err)
	}

	var expected map[string]interface{}
	if err := json.Unmarshal(expectedData, &expected); err != nil {
		t.Fatalf("failed to parse expected output: %v", err)
	}

	// 5. Compare
	if !mapsEqual(result, expected) {
		actualJSON, _ := json.MarshalIndent(result, "", "  ")
		expectedJSON, _ := json.MarshalIndent(expected, "", "  ")

		t.Errorf("output mismatch:\nExpected:\n%s\n\nActual:\n%s",
			string(expectedJSON), string(actualJSON))
	}

	t.Logf("✓ %s passed", filepath.Base(examplePath))
}

// mapsEqual compares two maps recursively
func mapsEqual(a, b map[string]interface{}) bool {
	if len(a) != len(b) {
		return false
	}

	for key, aVal := range a {
		bVal, exists := b[key]
		if !exists {
			return false
		}

		if !valuesEqual(aVal, bVal) {
			return false
		}
	}

	return true
}

// valuesEqual compares two interface{} values recursively
func valuesEqual(a, b interface{}) bool {
	// Handle nil cases
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	// Compare maps
	aMap, aIsMap := a.(map[string]interface{})
	bMap, bIsMap := b.(map[string]interface{})
	if aIsMap && bIsMap {
		return mapsEqual(aMap, bMap)
	}
	if aIsMap || bIsMap {
		return false
	}

	// Compare slices
	aSlice, aIsSlice := a.([]interface{})
	bSlice, bIsSlice := b.([]interface{})
	if aIsSlice && bIsSlice {
		if len(aSlice) != len(bSlice) {
			return false
		}
		for i := range aSlice {
			if !valuesEqual(aSlice[i], bSlice[i]) {
				return false
			}
		}
		return true
	}
	if aIsSlice || bIsSlice {
		return false
	}

	// Compare primitives
	return fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b)
}
