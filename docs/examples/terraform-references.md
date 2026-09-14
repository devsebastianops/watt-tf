# Terraform References

This section demonstrates how to reference Terraform specifics within Watt TF blueprints.

## Scenario

We want to create a blueprint to generate terraform resources for a HTTP service. The generated configuration will be placed inside a already existing terraform configuration. The existing configuration contains other resources and locals that we need to reference within our blueprint to pass them into modules.

## Solution

### Input

```yaml
service:
  name: my-microservice
  image: my-microservice-image:latest
  port: 8080
```

### Blueprint

```yaml
transform:
    - target: module.service.${input.service.name}
      value:
        image: ${input.service.image}
        port: ${input.service.port}
        project: ${local.project_id}
        bucket: ${resource.google_storage_bucket.project_files.name}
```

### Command

```bash
wtf build --input input.yaml --blueprint blueprint.yaml --output build.tf.json
```

### Output

```json
{
    "module": {
        "service": {
            "my-microservice": {
                "image": "my-microservice-image:latest",
                "port": 8080,
                "project": "${local.project_id}",
                "bucket": "${resource.google_storage_bucket.project_files.name}"
            }
        }
    }
}
```
