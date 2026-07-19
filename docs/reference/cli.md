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
| `--blueprint` <br> `-b` | `string` | Path to your blueprint YAML file. |
| `--input`<br> `-i` | `string` | Path to your input JSON or YAML file containing variables. |
| `--output`<br> `-o` | `string` | Path where the generated Terraform JSON file will be saved. |
| `--schema`<br> `-s` | `string` | Path/URL to a JSON Schema to validate the input file before transformation. |
| `--strict` | `flag` | Enables strict mode (halts on missing keys or syntax errors). |
| `--strip-nulls` | `flag` | Removes any keys with `null` values from the final output. |

::: warning Deprecation Note
Previously, it was possible to pass blueprints to Watt TF using `--config` or `-c`. This is still possible, but will be deprecated in the next major version (v2.0.0).
:::

## Example Pipeline Integration

In a typical CI/CD scenario (like GitHub Actions or GitLab CI), you would use the following pattern to ensure your infrastructure configuration is valid before deployment:

```bash
wtf build \
  --blueprint blueprint.yaml \
  --input env.prod.yaml \
  --output main.tf.json \
  --schema schemas/input-schema.json
```
