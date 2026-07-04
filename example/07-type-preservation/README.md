# Type Preservation

This example demonstrates how Watt-TF preserves data types during interpolation.

When you interpolate a value that is the entire content of a field (like `port: ${input.port}`), Watt-TF preserves its original type. Numbers stay numbers, booleans stay booleans. However, when interpolation is part of a larger string, the result becomes a string. This ensures your generated Terraform configurations have correct types for validation and execution.

```sh
wtf build --input example/07-type-preservation/input.json \
    --output example/07-type-preservation/watt.tf.json \
    --config example/07-type-preservation/.wtf.yaml
```

Or as a one-liner:

```sh
wtf build --input example/07-type-preservation/input.json --output example/07-type-preservation/watt.tf.json --config example/07-type-preservation/.wtf.yaml
```
