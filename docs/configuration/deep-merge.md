# Deep Merge

By default, when a transformation target points to an existing map or object, `watt-tf` performs a **deep merge** instead of completely overwriting the structure. This allows you to inject or update specific nested keys while preserving the rest of the existing data structure.


## How It Works

When a `transform` block is executed, the engine checks the data type of the resolved value:
1. **Scalar Values / Arrays:** Overwrite the target path.
2. **Maps (`map[string]any`):** The engine recursively traverses the keys. If a key exists in both the original data and the transformation value, it drills deeper. If it encounters a conflict between a map and a scalar, the transformation value takes precedence.

## Concrete Example

Imagine you have a Terraform configuration with existing tags, and you want to inject a globally managed environment tag, plus some team-specific infrastructure details.

### 1. Input Data 

```yaml
environment: "Production"
ami: "ami-0c55b159cbfafe1f0"
instance_type: "t2.micro"
tags:
  Name: "web-server"
  Project: "Phoenix"
```


### 2. Configuration (blueprint.yaml)

```yaml
transform:
  - target: "resource.aws_instance.web"
    value:
      ami: ${input.ami}
      instance_type: ${input.instance_type}
      tags:
        Name: ${input.tags.Name}
        Project: ${input.tags.Project}
  
  - target: "resource.aws_instance.web.tags"
    if: input.environment == "Production"
    value:
      Environment: "Production"
      `[company.com/managed-by](https://company.com/managed-by)`: "Watt TF"
      Project: "Phoenix-Prod" # Overwrites the existing Project key
```

### 3. Resulting Output
The existing keys (Name) are preserved, the conflicting key (Project) is updated, and new keys (Environment, company.com/managed-by) are seamlessly injected:

```json
{
  "resource": {
    "aws_instance": {
      "web": {
        "ami": "ami-0c55b159cbfafe1f0",
        "instance_type": "t2.micro",
        "tags": {
          "Name": "web-server",
          "Project": "Phoenix-Prod",
          "Environment": "Production",
          "[company.com/managed-by](https://company.com/managed-by)": "Watt TF"
        }
      }
    }
  }
}
```

::: info Nested Deep Merging
Deep merging works across multiple nested layers. If your value block mirrors the structure of the input target, Watt TF will safely navigate down the tree.

:::