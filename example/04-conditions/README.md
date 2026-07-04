# Conditions

In some cases you may want to apply a transformation only if a certain condition is met. This can be done using the `if` field in the transformation configuration.

Watt TF uses the [Google CEL](https://github.com/google/cel-spec) library to evaluate these conditions.

```sh
wtf build --input example/04-conditions/input.json \
    --output example/04-conditions/watt.tf.json \
    --config example/04-conditions/.wtf.yaml
```

Or as a one-liner:

```sh
wtf build --input example/04-conditions/input.json --output example/04-conditions/watt.tf.json --config example/04-conditions/.wtf.yaml
```