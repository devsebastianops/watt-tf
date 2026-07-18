# Input Schema Validation

When building complex infrastructure topologies, you want to ensure that the input variables (`input.json` or `input.yaml`) passed into Watt TF strictly adhere to your expected data formats. 

Watt TF allows you to validate your inputs against a custom **JSON Schema** before any transformations are executed. This prevents malformed data or missing attributes from breaking your blueprint execution halfway through.

## The `--schema` Flag

You can enforce input validation by passing a local path or a remote URL of a JSON Schema file to the `wtf build` command via the `--schema` flag.

```bash
wtf build \
  --config blueprint.yaml \
  --input inputs.prod.yaml \
  --schema schemas/inputs-schema.json \
  --output main.tf.json
```

If the data inside `inputs.prod.yaml` violates the structural rules defined in `inputs-schema.json`, the execution halts instantly and prints a detailed validation report.

## Example

### Input Schema

```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "title": "Watt TF Input Schema",
  "type": "object",
  "properties": {
    "environment": {
      "type": "string",
      "enum": ["dev", "staging", "prod"]
    },
    "project_name": {
      "type": "string",
      "minLength": 3
    },
    "enable_monitoring": {
      "type": "boolean"
    }
  },
  "required": ["environment", "project_name"]
}
```

### Invalid input

```yaml
environment: "local"      # Error: Not in the allowed enum (dev, staging, prod)
project_name: "tf"        # Error: Fails minLength of 3 characters
# enable_monitoring missing but optional, but environment/project_name fail constraints
```

### Expected Output

```text
❌     input validation failed
Error: Schema validation errors:
  1. environment: environment must be one of the following: "dev", "staging", "prod"
  2. project_name: String length must be greater than or equal to 3
```