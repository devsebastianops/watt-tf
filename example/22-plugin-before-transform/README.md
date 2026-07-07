# Plugin Before Transform

This example demonstrates how to hook into the watt tf transformation process using a plugin that modifies the configuration before the transformation occurs.
The plugin is a simple python script, that reads the input configuration and adds a new field to it before the transformation is applied.

```sh
wtf build --input example/22-plugin-before-transform/input.json \
    --output example/22-plugin-before-transform/watt.tf.json \
    --config example/22-plugin-before-transform/.wtf.yaml
```

Or as a one-liner:

```sh
wtf build --input example/22-plugin-before-transform/input.json --output example/22-plugin-before-transform/watt.tf.json --config example/22-plugin-before-transform/.wtf.yaml
```
