# Plugin After Transform

This example demonstrates how to hook into the watt tf transformation process using a plugin that modifies the configuration after the transformation occurs.
The plugin is a simple python script, that reads the input configuration and adds a new field to it after the transformation is applied.

```sh
wtf build --input example/23-plugin-after-transform/input.json \
    --output example/23-plugin-after-transform/watt.tf.json \
    --config example/23-plugin-after-transform/.wtf.yaml
```

Or as a one-liner:

```sh
wtf build --input example/23-plugin-after-transform/input.json --output example/23-plugin-after-transform/watt.tf.json --config example/23-plugin-after-transform/.wtf.yaml
```
