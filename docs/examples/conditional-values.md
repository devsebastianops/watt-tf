# Conditional Values

Demonstrates how to define and use conditional values within Watt TF blueprints.

## Scenario

We want to create a blueprint to generate terraform resources for a microservice. Based on the environment, the maximum number of instances should be configured conditionally. You could do so in a terraform module as well, but this example shows you how to do it using a blueprint.

## Solution

### Input

As described in the scenario we want to behave differently based on the environment. In this input it is set to `production`

```yaml
microservice:
  name: orders
  port: 8080
  image: myregistry/orders:latest
  environment: production
```

### Blueprint

The blueprint uses a ternary CEL expression to conditionally set the `max_instances` value based on the environment. If the environment is `production`, it sets the value to 5; otherwise, it sets it to 2.

```yaml
transform:
  - target: "module.${input.microservice.name}-service"
    value:
        source: "git::https://github.com/myorg/myrepo.git?ref=v1.2.3"
        name: "${input.microservice.name}-service"
        port: "${input.microservice.port}"
        image: "${input.microservice.image}"
        max_instances: "${input.microservice.environment == 'production' ? 5 : 2}"
```

### Command

```bash
wtf build --input input.yaml --blueprint blueprint.yaml --output build.tf.json
```

### Output

As expected, the output shows the `max_instances` value set to 5 because the environment is `production`.

```json
{
    "module": {
        "orders-service": {
            "source": "git::https://github.com/myorg/myrepo.git?ref=v1.2.3",
            "name": "orders-service",
            "port": 8080,
            "image": "myregistry/orders:latest",
            "max_instances": 5
        }
    }
}
```

## Alternative: Using Environment Variables

You could potentially use environment variables to control the behavior of the blueprint instead of relying solely on the input values. For example, you could set an environment variable `ENVIRONMENT` to `production` and use it within the blueprint to determine the `max_instances` value.

### Adjusted Blueprint

```yaml
transform:
  - target: "module.${input.microservice.name}-service"
    value:
        source: "git::https://github.com/myorg/myrepo.git?ref=v1.2.3"
        name: "${input.microservice.name}-service"
        port: "${input.microservice.port}"
        image: "${input.microservice.image}"
        max_instances: "${env.ENVIRONMENT == 'production' ? 5 : 2}"
```

### Adjusted Command

```bash
export ENVIRONMENT=production
wtf build --input input.yaml --blueprint blueprint.yaml --output build.tf.json
```


## Alternative: Load from another configuration file

You could also load configuration values from another file, such as a JSON or YAML file, and use them within your blueprint to determine the `max_instances`. This is good for separating environment-specific configurations from the main input file.

### Adjusted Blueprint

```yaml
data:
    project:
        type: "file"
        file: "project.yaml"

transform:
  - target: "module.${input.microservice.name}-service"
    value:
        source: "git::https://github.com/myorg/myrepo.git?ref=v1.2.3"
        name: "${input.microservice.name}-service"
        port: "${input.microservice.port}"
        image: "${input.microservice.image}"
        max_instances: "${data.project.environment == 'production' ? 5 : 2}"
```

### Additional Project Configuration File

```yaml
environment: "production"
```
