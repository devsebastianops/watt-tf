# Error Handling & Robustness

This example demonstrates how Watt-TF handles missing or incomplete input gracefully.

Watt-TF is designed to work with partial configurations and missing keys. When a key is referenced in your configuration but doesn't exist in the input, it is replaced with `null` instead of causing an error. This robustness allows you to build flexible infrastructure templates that work with varying input configurations without failing.

## Lenient Mode (Default)

By default, Watt-TF runs in **lenient mode**, where missing keys are:
- Logged as **warnings**
- Replaced with **null** in the output
- Do not break the build process

This is perfect for optional configuration values.

## Strict Mode

You can enable **strict mode** with the `--strict` flag, which treats missing keys as errors and fails the build:

```bash
wtf build --input example/12-missing-keys/input.json \
    --config example/12-missing-keys/.wtf.yaml \
    --output example/12-missing-keys/watt.tf.json \
    --strict
```

In strict mode, the same configuration will fail if any referenced keys are missing.

## Example

**Configuration (.wtf.yaml):**
```yaml
transform:
  - target: resource.server.main
    value:
      name: "my-server"
      existing_key: "${input.valid_input}"
      not_existing: "${input.missing_input}"
```

**Input (input.json):**
```json
{
  "valid_input": "value"
}
```

**Output (lenient mode, default):**
```json
{
  "resource": {
    "server": {
      "main": {
        "name": "my-server",
        "existing_key": "value",
        "not_existing": null
      }
    }
  }
}
```

Notice that `input.missing_input` doesn't exist, so `not_existing` is set to `null`.

## One-Liner

Lenient mode (default):

```sh
wtf build --input example/12-missing-keys/input.json --output example/12-missing-keys/watt.tf.json --config example/12-missing-keys/.wtf.yaml
```

Strict mode:

```sh
wtf build --input example/12-missing-keys/input.json --output example/12-missing-keys/watt.tf.json --config example/12-missing-keys/.wtf.yaml --strict
```
