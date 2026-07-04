# Interpolate Targets Example

Sometimes you may want to interpolate target strings with values from the input file. This example demonstrates how to use Watt-TF to transform an input JSON file and interpolate target strings.

```sh
wtf build --input example/03-json-interpolate-target/input.json \
    --output example/03-json-interpolate-target/watt.tf.json \
    --config example/03-json-interpolate-target/.wtf.yaml
```

Or as a one-liner:

```sh
wtf build --input example/03-json-interpolate-target/input.json --output example/03-json-interpolate-target/watt.tf.json --config example/03-json-interpolate-target/.wtf.yaml
```