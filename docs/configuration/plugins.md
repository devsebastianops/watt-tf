# Extending Watt TF with Plugins

::: danger Attention
Please be aware that you are going to be writing code to extend Watt TF. That means you are responsible for maintenance, further development and security of your plugin. We recommend plugins only when they are absolutely necessary.
:::

Plugins allow you to extend the functionality of Watt TF by hooking into the transformation process at specific points. This can be useful for adding custom logic, modifying data, or integrating with external systems.

## Plugin Lifecycle Events

| Event | Description |
|-------|------------|
| beforeTransform | Triggered before the transformation process begins. You can modify the input data or configuration here. |
| afterTransform | Triggered after the transformation process is complete. You can modify the output data here. |

## How Plugins Work

Plugins can be written in any language, as long as they can be executed from the command line. Watt TF will call the plugin script with the event name and all relevant data passed as JSON via standard input. The plugin should read the input, perform any necessary modifications, and output the modified data as JSON to standard output.

### Incoming Data 

When a plugin is triggered, it receives the following JSON structure as input:

```json
{
    "version": "1.0",
    "event": "beforeTransform",
    "data": {
        "input": {
            "enabled": true,
            "size": 10
        },
        "env": {
            "API_KEY": "my-secret-key"
        },
        "result": {

        }
    }
}
```

#### Field descriptions

| Field | Description | Value |
|-------|------------|-------|
| version | The version of the plugin API. | "1.0" |
| event | The lifecycle event that triggered the plugin. | `beforeTransform` or `afterTransform` |
| data | The data relevant to the event. | See below |
| data.input | The input data provided to Watt TF. | JSON object, coming directly from your `--input` file |
| data.env | The environment variables available to the plugin. | JSON object, containing all environment variables |
| data.result | The result of the transformation process. | JSON object, initially empty for `beforeTransform`, populated for `afterTransform` |

### Expected Outgoing Data  

Watt TF expects the plugin to return a JSON object with the following structure:

```json
{
    "status": "success",
    "data": {
        "input": {
            "enabled": false,
            "size": 20
        },
        "env": {
            "API_KEY": "my-secret-key"
        },
        "result": {

        }
    }
}
```

#### Field descriptions

| Field | Description | Value |
|-------|------------|-------|
| status | The status of the plugin execution. | `success` or `error` |
| data | The modified data after the plugin has executed. | See below |
| data.input | The modified input data. | JSON object, can be modified by the plugin |
| data.env | The environment variables available to the plugin. | JSON object, can be modified by the plugin |
| data.result | The modified result of the transformation process. | JSON object, can be modified by the plugin |


### Error Handling
In case your plugin encounters an error you can return a JSON object with the following structure:

```json
{
    "status": "error",
    "error": "An error occurred while processing the plugin."
}
```

#### Field descriptions

| Field | Description | Value |
|-------|------------|-------|
| status | The status of the plugin execution. | `error` |
| error | A descriptive error message. | String describing the error |


## Example Plugin

```python
#!/usr/bin/env python3
import json
import os
import sys

def handle(version, event, data):
    if event == "beforeTransform":
        data["input"]["value"] = "This field was added by the beforeTransform plugin."
    return data

def main():
    for line in sys.stdin:
        if not line.strip():
            continue

        try:
            req = json.loads(line)
            version = req.get("version", "")
            data = req.get("data", {})
            event = req.get("event", "")

            if event == "beforeTransform":
                updatedData = handle(version, event, data)

            response = {"status": "success", "data": updatedData}
        except Exception as e:
            response = {"status": "error", "error": str(e)}

        sys.stdout.write(json.dumps(response) + "\n")
        sys.stdout.flush()  # Verhindert, dass Go blockiert
        os._exit(0)  # Beendet das Plugin nach der Verarbeitung


if __name__ == "__main__":
    main()
```

This example plugin modifies the input data before the transformation process begins. It adds a new field `value` to the input data with a custom message. You can adapt this example to suit your specific needs, such as modifying the output data after transformation or integrating with external systems via API calls.

## How to Use a Plugin

To use a plugin with Watt TF, you need to specify the plugin in your `blueprint.yaml` configuration file. Here's an example configuration that uses the above plugin:

```yaml
plugins:
- name: setValuePlugin
  version: 1.0.0
  on: beforeTransform
  cmd: python
  args:
    - plugin.py
```

This configuration tells Watt TF to execute the `plugin.py` script before the transformation process begins. The plugin will receive the input data, modify it, and return the modified data to Watt TF for further processing.

All paths in arguments are relative to the `blueprint.yaml` file. You can also use absolute paths if needed.