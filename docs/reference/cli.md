# CLI Reference

The Watt TF command-line interface is designed for high-performance transformations and seamless integration into CI/CD pipelines.

## `wtf build`

The `build` command is the primary entry point for transforming your blueprints into Terraform JSON files.

### Usage
```bash
wtf build [options]
```

### Options

| Option | Type | Description |
| :--- | :--- | :--- |
| `--config` | `string` | Path to your blueprint YAML file. |
| `--input` | `string` | Path to your input JSON or YAML file containing variables. |
| `--output` | `string` | Path where the generated Terraform JSON file will be saved. |
| `--schema` | `string` | Path/URL to a JSON Schema to validate the input file before transformation. |
| `--strict` | `flag` | Enables strict mode (halts on missing keys or syntax errors). |
| `--strip-null` | `flag` | Removes any keys with `null` values from the final output. |

## Example Pipeline Integration

In a typical CI/CD scenario (like GitHub Actions or GitLab CI), you would use the following pattern to ensure your infrastructure configuration is valid before deployment:

```bash
wtf build \
  --config blueprint.yaml \
  --input env.prod.yaml \
  --output main.tf.json \
  --schema schemas/input-schema.json
```
