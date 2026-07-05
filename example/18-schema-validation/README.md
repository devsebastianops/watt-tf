# Example 18: Schema Validation

This example demonstrates how to validate input against a JSON Schema before transformation.

## What It Does

- Loads a JSON Schema that defines the input structure
- Validates the input against the schema before transformation
- Only proceeds with transformation if validation passes
- Fails immediately if the input doesn't match the schema

## Key Features

- **Optional schema validation:** Use `--schema` flag to validate input
- **Early validation:** Errors are caught before transformation
- **Clear error messages:** Validation failures show which fields are invalid
- **Standard JSON Schema:** Uses JSON Schema Draft 7

## Input Validation Rules

The schema in this example requires:

```json
{
  "server": {
    "name": "string (required, non-empty)",
    "port": "integer (required, 1-65535)",
    "enabled": "boolean (optional)",
    "tags": "array of strings (optional)"
  }
}
```

## Configuration

```yaml
transform:
  - target: resource.compute_instance.server
    value:
      name: "${input.server.name}"
      port: ${input.server.port}
      enabled: ${input.server.enabled}
      tags: ${input.server.tags}
```

## Input

```json
{
  "server": {
    "name": "prod-server",
    "port": 8080,
    "enabled": true,
    "tags": ["production", "critical"]
  }
}
```

## Output

```json
{
  "resource": {
    "compute_instance": {
      "server": {
        "name": "prod-server",
        "port": 8080,
        "enabled": true,
        "tags": ["production", "critical"]
      }
    }
  }
}
```

## Validation Examples

### Valid Input (Passes)

```json
{
  "server": {
    "name": "prod-server",
    "port": 8080
  }
}
```

### Invalid Input (Fails)

```json
{
  "server": {
    "name": "prod-server"
    // Missing required "port" field
  }
}
```

Validation error:
```
Schema validation errors:
  1. server: port is required
```

### Invalid Type (Fails)

```json
{
  "server": {
    "name": "prod-server",
    "port": "8080"  // Should be integer, not string
  }
}
```

Validation error:
```
Schema validation errors:
  1. server.port: Invalid type. Expected: integer, given: string
```

## How to Test

With schema validation:

```bash
./bin/wtf build \
  --input example/18-schema-validation/input.json \
  --config example/18-schema-validation/.wtf.yaml \
  --schema example/18-schema-validation/schema.json \
  --output example/18-schema-validation/watt.tf.json
```

Without schema validation (schema is optional):

```bash
./bin/wtf build \
  --input example/18-schema-validation/input.json \
  --config example/18-schema-validation/.wtf.yaml \
  --output example/18-schema-validation/watt.tf.json
```

### Testing Validation Errors

To test validation errors, use the invalid input files:

Missing required field:

```bash
./bin/wtf build \
  --input example/18-schema-validation/input-invalid.json \
  --schema example/18-schema-validation/schema.json \
  --output /tmp/test.tf.json
```

Expected error: `server: port is required`

Wrong type:

```bash
./bin/wtf build \
  --input example/18-schema-validation/input-invalid-type.json \
  --schema example/18-schema-validation/schema.json \
  --output /tmp/test.tf.json
```

Expected error: `server.port: Invalid type. Expected: integer, given: string`

## One-Liner

With schema validation:

```sh
./bin/wtf build --input example/18-schema-validation/input.json --config example/18-schema-validation/.wtf.yaml --schema example/18-schema-validation/schema.json --output example/18-schema-validation/watt.tf.json
```

## Benefits

1. **Fail Fast:** Catch input errors before expensive transformations
2. **Clear Contracts:** Schema documents what input structure is required
3. **Type Safety:** Ensure inputs have correct types (string, integer, boolean, etc.)
4. **Validation Rules:** Enforce ranges (port: 1-65535), required fields, and more
5. **Developer Experience:** Clear error messages help fix issues quickly
