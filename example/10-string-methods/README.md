# String Methods

This example demonstrates using string manipulation methods in conditions.

Watt-TF supports string methods like `startsWith()`, `endsWith()`, and `contains()` for pattern matching in conditions. This enables powerful filtering logic based on string patterns—perfect for conditional resource creation based on naming conventions, email domains, or environment identifiers.

```sh
wtf build --input example/10-string-methods/input.json \
    --output example/10-string-methods/watt.tf.json \
    --config example/10-string-methods/.wtf.yaml
```

Or as a one-liner:

```sh
wtf build --input example/10-string-methods/input.json --output example/10-string-methods/watt.tf.json --config example/10-string-methods/.wtf.yaml
```
