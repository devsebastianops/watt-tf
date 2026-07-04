# Complex CEL Expressions

This example demonstrates advanced conditional logic with complex boolean expressions.

Watt-TF's condition system supports complex CEL expressions with multiple logical operators (`&&`, `||`), parentheses for grouping, and comparison operators. This allows you to create sophisticated rules for when transformations should apply—essential for managing infrastructure across different environments and deployment scenarios.

```sh
wtf build --input example/09-complex-cel/input.json \
    --output example/09-complex-cel/watt.tf.json \
    --config example/09-complex-cel/.wtf.yaml
```

Or as a one-liner:

```sh
wtf build --input example/09-complex-cel/input.json --output example/09-complex-cel/watt.tf.json --config example/09-complex-cel/.wtf.yaml
```
