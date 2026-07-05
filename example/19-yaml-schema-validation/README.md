# Example 19: YAML Input with Schema Validation

This example demonstrates that JSON Schema validation works seamlessly with YAML inputs.

## What It Does

- Loads a YAML configuration file (instead of JSON)
- Validates the YAML input against a JSON Schema
- Transforms the YAML data into Terraform format
- Shows that schema validation is format-agnostic (works with both JSON and YAML)

## Key Features

- **YAML + Schema:** Schema validation works with YAML inputs
- **Format Agnostic:** Same validation logic for JSON and YAML
- **Complex Structures:** Validates nested objects and required fields
- **Version Validation:** Uses regex patterns (e.g., semantic versioning)
- **Enum Validation:** Restricts environment to valid values

## Input Structure (YAML)

```yaml
app:
  name: my-application
  version: "1.0.0"
  environment: production
  replicas: 3
  image: my-app:latest
  resources:
    cpu: "500m"
    memory: "512Mi"
  health_check:
    enabled: true
    path: /health
    port: 8080
```

## Schema Requirements

```json
{
  "app": {
    "name": "string (required)",
    "version": "string matching semantic versioning (required)",
    "environment": "enum: development|staging|production (required)",
    "replicas": "integer 1-100 (required)",
    "image": "string (required)",
    "resources": {
      "cpu": "string (required)",
      "memory": "string (required)"
    },
    "health_check": {
      "enabled": "boolean (required)",
      "path": "string (required)",
      "port": "integer 1-65535 (required)"
    }
  }
}
```

## Configuration

```yaml
transform:
  - target: resource.kubernetes_deployment.app
    value:
      name: "${input.app.name}"
      namespace: "${input.app.environment}"
      image: "${input.app.image}"
      replicas: ${input.app.replicas}
      resources:
        requests:
          cpu: "${input.app.resources.cpu}"
          memory: "${input.app.resources.memory}"
      health_check:
        enabled: ${input.app.health_check.enabled}
        path: "${input.app.health_check.path}"
        port: ${input.app.health_check.port}
```

## Output

```json
{
  "resource": {
    "kubernetes_deployment": {
      "app": {
        "name": "my-application",
        "namespace": "production",
        "image": "my-app:latest",
        "replicas": 3,
        "resources": {
          "requests": {
            "cpu": "500m",
            "memory": "512Mi"
          }
        },
        "health_check": {
          "enabled": true,
          "path": "/health",
          "port": 8080
        }
      }
    }
  }
}
```

## How to Test

With YAML input and schema validation:

```bash
./bin/wtf build \
  --input example/19-yaml-schema-validation/input.yaml \
  --config example/19-yaml-schema-validation/.wtf.yaml \
  --schema example/19-yaml-schema-validation/schema.json \
  --output example/19-yaml-schema-validation/watt.tf.json
```

Without schema (YAML + transformation only):

```bash
./bin/wtf build \
  --input example/19-yaml-schema-validation/input.yaml \
  --config example/19-yaml-schema-validation/.wtf.yaml \
  --output example/19-yaml-schema-validation/watt.tf.json
```

## One-Liner

With schema validation:

```sh
./bin/wtf build --input example/19-yaml-schema-validation/input.yaml --config example/19-yaml-schema-validation/.wtf.yaml --schema example/19-yaml-schema-validation/schema.json --output example/19-yaml-schema-validation/watt.tf.json
```

## Format Agnostic Validation

The schema validation works the same way for both JSON and YAML inputs because:

1. **Parser:** The input file is parsed to `map[string]any` (JSON or YAML)
2. **Validation:** The map is converted to JSON and validated against the schema
3. **Transformation:** The map is then transformed using the configuration

This means the same schema works for both formats:

```bash
# Works with JSON
./bin/wtf build --input config.json --schema schema.json --config .wtf.yaml

# Works with YAML
./bin/wtf build --input config.yaml --schema schema.json --config .wtf.yaml

# Works with any format as long as the structure matches
./bin/wtf build --input config.yml --schema schema.json --config .wtf.yaml
```

## Use Cases

- **Kubernetes Deployments:** Validate deployment configurations
- **Multi-format Support:** Accept both JSON and YAML inputs with same validation
- **Environment-specific Configs:** Restrict environments to known values
- **Type Safety:** Ensure numeric fields are numbers, not strings
- **Version Validation:** Enforce semantic versioning or custom formats
