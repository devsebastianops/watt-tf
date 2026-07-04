# Simple Test Example

This example demonstrates how to use Watt-TF to transform a simple input JSON file.

```sh
wtf build --input example/01-simple-json/input.json \
    --output example/01-simple-json/watt.tf.json \
    --config example/01-simple-json/.wtf.yaml
```

Or as a one-liner:

```sh
wtf build --input example/01-simple-json/input.json --output example/01-simple-json/watt.tf.json --config example/01-simple-json/.wtf.yaml
```