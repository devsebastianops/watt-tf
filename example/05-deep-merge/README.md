# Deep Merge

This example demonstrates how to merge multiple transformations that target the same resource.

When you apply multiple transformations to the same target path, Watt-TF intelligently merges their values together. This is useful when you want to build up a resource definition in stages or from different configuration sources.

```sh
wtf build --input example/05-deep-merge/input.json \
    --output example/05-deep-merge/watt.tf.json \
    --config example/05-deep-merge/.wtf.yaml
```

Or as a one-liner:

```sh
wtf build --input example/05-deep-merge/input.json --output example/05-deep-merge/watt.tf.json --config example/05-deep-merge/.wtf.yaml
```
