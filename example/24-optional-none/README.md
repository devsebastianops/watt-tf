# Optional None Example

This example demonstrates how to use the `optional.none()` function in a watt tf transformation. The input JSON does not include the `ssl_certificate` field, and the transformation uses a conditional expression to set the `ssl-certificate` field to `optional.none()` if the input does not have the `ssl_certificate` key.

```sh
wtf build --input example/24-optional-none/input.json \
    --output example/24-optional-none/watt.tf.json \
    --config example/24-optional-none/.wtf.yaml
```

Or as a one-liner:

```sh
wtf build --input example/24-optional-none/input.json --output example/24-optional-none/watt.tf.json --config example/24-optional-none/.wtf.yaml
```

## Stripping Null Values

watt tf provides an option to strip null values from the output.

```sh
wtf build --input example/24-optional-none/input.json \
    --output example/24-optional-none/watt.tf.json \
    --config example/24-optional-none/.wtf.yaml \
    --strip-nulls
```