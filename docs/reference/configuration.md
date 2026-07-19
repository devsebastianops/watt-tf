# Configuration

The Watt TF blueprint file (typically `blueprint.yaml`) is the core definition of your infrastructure transformations. It follows a declarative structure that maps input data to Terraform resources.

## Top Level Structure

| Key | Type | Description |
| :--- | :--- | :--- |
| `plugins` | `array` | A list of external plugins to hook into the transformation lifecycle. |
| `transform` | `array` | The sequence of transformation rules applied to your resources. |

### The transform Block

| Key | Type | Description |
| :--- | :--- | :--- |
| `target` | `string` | The path where the value will be placed (e.g., `resource.aws_s3_bucket.main`). |
| `value` | `any` | The payload or configuration data to set at the target. |
| `if` | `string` | An optional CEL expression that determines if this block should execute. |
| `for_each`| `string` | An optional CEL expression that iterates over a collection. |

### The plugins Block

| Key | Type | Description |
| :--- | :--- | :--- |
| `name` | `string` | A unique identifier for the plugin. |
| `version` | `string` | The version string of the plugin. |
| `on` | `string` | The lifecycle event: `beforeTransform` or `afterTransform`. |
| `cmd` | `string` | The base command to execute (e.g., `python`, `node`, `./script.sh`). |
| `args` | `array` | A list of arguments passed to the command (e.g., `["plugin.py"]`). |

## Example Blueprint

```yaml
plugins:
  - name: my-custom-plugin
    version: 1.0.0
    on: beforeTransform
    cmd: python
    args: ["plugin.py"]

transform:
  # Example of conditional resource creation
  - if: input.environment == "prod"
    target: resource.aws_instance.web
    value:
      instance_type: "t3.large"
      monitoring: true

  # Example of loop-based resource creation
  - for_each: input.subnets
    target: resource.aws_subnet.${item.name}
    value:
      cidr_block: "${item.cidr}"

```

In case you want to learn more about the key concepts of Watt TF, check out the [Core Concepts](../guide/concepts.md) guide.