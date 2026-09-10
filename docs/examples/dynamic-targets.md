# Dynamic Targets

Demonstrates how to use dynamic targets within Watt TF blueprints.

## Scenario

We want to create a blueprint to generate terraform resources for multiple microservices. Each microservice should have its own module identifier, to prevent terraform from encountering duplicate module names.

## Solution

### Input

As described in the scenario, we have multiple microservices, each with its own configuration.

```yaml
microservices:
-   name: orders
    port: 8080
    image: myregistry/orders:latest
-   name: payments
    port: 9090
    image: myregistry/payments:latest
```

### Blueprint

We iterate over each microservice in the input and create a corresponding terraform module for it.

```yaml
transform:
  - target: "module.${item.name}-service"
    for_each: "${input.microservices}"
    value:
        source: "git::https://github.com/myorg/myrepo.git?ref=v1.2.3"
        name: "${item.name}-service"
        port: "${item.port}"
        image: "${item.image}"
```

### Output

As expected we build the following terraform module structure:

```json
{
    "module": {
        "orders-service": {
            "source": "git::https://github.com/myorg/myrepo.git?ref=v1.2.3",
            "name": "orders-service",
            "port": 8080,
            "image": "myregistry/orders:latest"
        },
        "payments-service": {
            "source": "git::https://github.com/myorg/myrepo.git?ref=v1.2.3",
            "name": "payments-service",
            "port": 9090,
            "image": "myregistry/payments:latest"
        }
    }
}
```