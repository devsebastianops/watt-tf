# Null and Empty Values

This example demonstrates how Watt-TF handles null values and empty collections.

Watt-TF properly processes null values and empty arrays/objects from your input configuration. This is important for optional infrastructure components—you can define conditional infrastructure that may or may not be deployed based on whether values are provided.

```sh
wtf build --input example/11-null-handling/input.json \
    --output example/11-null-handling/watt.tf.json \
    --config example/11-null-handling/.wtf.yaml
```

Or as a one-liner:

```sh
wtf build --input example/11-null-handling/input.json --output example/11-null-handling/watt.tf.json --config example/11-null-handling/.wtf.yaml
```
