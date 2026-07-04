# Simple Test Example ( YAML )

In case you are a fan of YAML, this example demonstrates how to use Watt-TF to transform a simple input YAML file.

```sh
wtf build --input example/02-simple-yaml/input.yaml \
    --output example/02-simple-yaml/watt.tf.json \
    --config example/02-simple-yaml/.wtf.yaml
```

Or as a one-liner:

```sh
wtf build --input example/02-simple-yaml/input.yaml --output example/02-simple-yaml/watt.tf.json --config example/02-simple-yaml/.wtf.yaml
```