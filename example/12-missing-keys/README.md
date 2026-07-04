# Error Handling & Robustness

This example demonstrates how Watt-TF handles missing or incomplete input gracefully.

Watt-TF is designed to work with partial configurations and missing keys. You can define transformations that use some optional keys while others remain as literal values. This robustness allows you to build flexible infrastructure templates that work with varying input configurations without failing.

```sh
wtf build --input example/12-missing-keys/input.json \
    --output example/12-missing-keys/watt.tf.json \
    --config example/12-missing-keys/.wtf.yaml
```

Or as a one-liner:

```sh
wtf build --input example/12-missing-keys/input.json --output example/12-missing-keys/watt.tf.json --config example/12-missing-keys/.wtf.yaml
```
