# Array Values

This example demonstrates how to work with arrays and list interpolations.

Watt-TF can interpolate values within arrays, including accessing specific array elements using index notation like `${input.tags.0}`. You can also mix interpolated values with literal values in the same array. This is crucial for defining Terraform configurations that require lists of resources or tags.

```sh
wtf build --input example/08-array-values/input.json \
    --output example/08-array-values/watt.tf.json \
    --config example/08-array-values/.wtf.yaml
```

Or as a one-liner:

```sh
wtf build --input example/08-array-values/input.json --output example/08-array-values/watt.tf.json --config example/08-array-values/.wtf.yaml
```
