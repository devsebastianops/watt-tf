# E2E Tests

End-to-end tests for the `watt-tf` CLI tool.

## Running Tests

Run all E2E tests:
```bash
go test ./tests/e2e -v
```

Run a specific example test:
```bash
go test ./tests/e2e -v -run TestE2EExamples/01-simple-json
```

Run with detailed output:
```bash
go test ./tests/e2e -v -count=1
```

## Test Structure

Each test:
1. Reads the input file (`.json`, `.yaml`, or `.yml`)
2. Loads the `.wtf.yaml` configuration
3. Runs the transformer
4. Compares the result with `expected.tf.json`

For a test to pass:
- The example directory must contain:
  - `input.json` or `input.yaml` (or `.yml`)
  - `.wtf.yaml` configuration
  - `expected.tf.json` expected output

## Adding a New Test

1. Create a new directory under `example/XX-description/`
2. Add:
   - `input.json` or `input.yaml`
   - `.wtf.yaml`
   - `expected.tf.json` (the expected output)
   - `README.md` (optional documentation)

3. Run tests:
   ```bash
   go test ./tests/e2e -v
   ```

## Test Coverage

The tests currently cover:
- JSON parsing
- YAML parsing
- Configuration loading
- Transformations
- String interpolation
- Conditional execution
- Path unflattening

## Debugging Tests

If a test fails, the output will show:
```
Expected:
{...actual expected JSON...}

Actual:
{...what was generated...}
```

Use this to identify differences and adjust either:
- The `.wtf.yaml` configuration
- The `expected.tf.json` file
