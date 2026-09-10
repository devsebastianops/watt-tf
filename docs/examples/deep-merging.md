# Deep Merging

Demonstrates how to perform deep merging of configuration values within Watt TF blueprints.

## Scenario

We want to create a blueprint to generate terraform resources for a microservice. In our input we allow optional labels. If there are any labels, the terraform module should receive them. If there are no labels given, the module should not include the labels section.

## Solution

### Input

This is an example of the input configuration for the microservice with labels.

```yaml
microservice:
  name: orders
  port: 8080
  image: myregistry/orders:latest
  labels:
    environment: production
```

### Blueprint

The blueprint demonstrates how to conditionally merge labels into the module configuration based on the input provided.

```yaml
transform:
  - target: "module.${input.microservice.name}-service"
    value:
        source: "git::https://github.com/myorg/myrepo.git?ref=v1.2.3"
        name: "${input.microservice.name}-service"
        port: "${input.microservice.port}"
        image: "${input.microservice.image}"
  - target: "module.${input.microservice.name}-service"
    if: has(input.microservice.labels) && length(input.microservice.labels) > 0
    value:
        labels: "${input.microservice.labels}"
```

### Command

```bash
wtf build --input input.yaml --blueprint blueprint.yaml --output build.tf.json
```

### Output

Because we provided labels in our input, they are included in the output configuration for the module.

```json
{
    "module": {
        "orders-service": {
            "source": "git::https://github.com/myorg/myrepo.git?ref=v1.2.3",
            "name": "orders-service",
            "port": 8080,
            "image": "myregistry/orders:latest",
            "labels": {
                "environment": "production"
            }
        }
    }
}
```