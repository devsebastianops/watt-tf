# Nested Paths

This example demonstrates how Watt-TF handles deeply nested resource structures.

Watt-TF can work with arbitrarily deep dot-notation paths (like `resource.aws_vpc.main.subnets.primary.tags`). It automatically creates all intermediate nested objects needed to reach your target value. This is essential for complex infrastructure configurations where resources have multiple levels of configuration.

```sh
wtf build --input example/06-nested-paths/input.json \
    --output example/06-nested-paths/watt.tf.json \
    --config example/06-nested-paths/.wtf.yaml
```

Or as a one-liner:

```sh
wtf build --input example/06-nested-paths/input.json --output example/06-nested-paths/watt.tf.json --config example/06-nested-paths/.wtf.yaml
```
