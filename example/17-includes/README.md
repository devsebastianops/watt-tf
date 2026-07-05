# Example 17: Config Includes

This example demonstrates how to split your configuration across multiple files using the `include` directive. Includes are loaded before transformations are evaluated, and their transforms are merged into the main configuration.

## What It Does

- Main config file (`.wtf.yaml`) includes another config file (`compute.yaml`)
- Both config files have `transform` sections
- All transforms are combined in order (includes first, then main)
- The combined transform list is executed normally

## Key Features

- **Include directive:** Specify other config files to load
- **Order:** Included transforms are processed before main transforms
- **Relative paths:** Include paths are relative to the main config file's directory
- **Recursive:** Included configs can also include other files

## Configuration Files

### Main Config (.wtf.yaml)

```yaml
include:
  - compute.yaml

transform:
  - target: resource.google_cloud_run_service.default
    value:
      port: ${input.cloudrun.port}
      image: ${input.cloudrun.image}
```

### Included Config (compute.yaml)

```yaml
transform:
  - target: resource.compute.default
    value:
      port: ${input.compute.port}
```

## Input

```json
{
  "cloudrun": {
    "port": 8080,
    "image": "my-app:latest"
  },
  "compute": {
    "port": 5432
  }
}
```

## Execution Order

1. **Load includes:** `compute.yaml` is loaded
   - Adds transform for `resource.compute.default`
2. **Load main:** `.wtf.yaml` is loaded
   - Adds transform for `resource.google_cloud_run_service.default`
3. **Execute transforms:** All transforms run in order
   - First: compute resources
   - Then: cloud run service

## Output

Both resources are generated:

```json
{
  "resource": {
    "compute": {
      "default": {
        "port": 5432
      }
    },
    "google_cloud_run_service": {
      "default": {
        "image": "my-app:latest",
        "port": 8080
      }
    }
  }
}
```

## Use Cases

- **Modular configs:** Split infrastructure into logical modules
- **Shared transforms:** Common transforms in a shared file
- **Environment-specific:** Include different configs for dev/staging/prod

## Example with Multiple Includes

```yaml
include:
  - base.yaml
  - networking.yaml
  - compute.yaml

transform:
  - target: resource.app
    value:
      version: ${input.version}
```

## How to Test

```bash
go test ./tests/e2e -v -run 17-includes
```

Or run all E2E tests:
```bash
go test ./tests/e2e -v
```

## Manual Test

```bash
./bin/wtf build \
  --config example/17-includes/.wtf.yaml \
  --input example/17-includes/input.json \
  --output output.tf.json
```
