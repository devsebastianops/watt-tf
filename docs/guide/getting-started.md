# Getting Started 

## The Problems of Native Terraform HCL

I worked with terraform and I have a hate love relationship with it. I love the idea of infrastructure as code, but I hate the complexity of HCL, especially because it is **not a programming language**.

Another problem is that developers do not understand Terraform. Explaining how to generate dynamic routing tables, complex IAM bindings, or nested Kubernetes volumes using pure HCL is a nightmare.

Oftentimes, you end up with write-only masterpieces of unreadable nested `for`/`for_each` loops, manual string-splitting, and ternary operators that nobody - not even you in two weeks - will ever understand.

## The "improvement": Fighting fire with fire

The first solution that comes to mind was a custom generator that would take simple JSON input and generate the complex HCL output, based on the rules defined by the platform team. This is a common approach, but it has some drawbacks:

- It is a lot of work to maintain, especially when the platform team is busy with other tasks.
- It is not generic, so it cannot be reused for other projects or teams.
- It is not flexible, so it cannot adapt to changing requirements or new features.

So in the end, you end up with a lot of duplicated tooling, inconsistent deployments, and platform teams spending valuable time maintaining custom generators instead of improving their infrastructure.

## The Solution: Watt TF

Watt TF is a generic transformation engine that sits between your application configuration and Terraform.

Instead of asking developers to write Terraform, platform teams define how structured input maps to reusable Terraform modules and resources.

And instead of asking your platform team to maintain a custom generator, Watt TF provides a generic engine that can be reused for any project or team.

### 1. Input <Badge type="info" text="Typically provided by: Developers" />

Watt TF takes structured input in `JSON` or `YAML`. Imagine you want to enable your developers to deploy a microservice with a database. Instead of asking them to write complex Terraform HCL, you can provide them with a simple JSON input like this:

```json
{
  "service": {
    "name": "ms-orders",
  },
  "app": {
    "image": "ghcr.io/acme/ms-orders",
    "version": "v1.0.0",
    "port": 8080
  },
  "database": {
    "type": "postgres",
    "version": "14",
    "size": "20Gi"
  }
}
```

Also Watt TF will read **environment variables**, so you can use them to provide secrets or other sensitive information that should not be stored in the input file.

### 2. Blueprint <Badge type="info" text="Typically provided by: Platform Team" />

The platform team defines a blueprint that describes how the input maps to terraform modules and resources. The blueprint is written in `YAML`, and it can contain complex logic, loops, and conditionals, empowered by Google CEL. The blueprint is the heart of Watt TF, and it is what makes it generic and reusable.

```yaml
transform:
  - target: module.app.${input.service.name}
    value: 
        source: "git::https://git.acme.com/terraform-modules/app.git?ref=v1.0.0"
        version: "${input.app.version}"
        port: "${input.app.port}"
        min_replicas: ${has(input.app.min_replicas) ? optional(input.app.min_replicas) : optional.none()}
        max_replicas: ${has(input.app.max_replicas) ? optional(input.app.max_replicas) : optional.none()}
        ram: ${env.ENVIRONMENT == "production" ? "2Gi" : "512Mi"}
        database: ${has(input.database) ? "module.database.${input.service.name}.dsn" : optional.none()}

  - target: module.database.${input.service.name}
    if: has(input.database)
    value:
        source: "git::https://git.acme.com/terraform-modules/database.git?ref=v1.0.0"
        type: "${input.database.type}"
        version: "${input.database.version}"
        size: "${input.database.size}"
```

### 3. Generating Terraform JSON 

::: info Local usage 
For testing purposes, you can run the Watt TF CLI locally. But it is highly recommended to follow the credo of "Automate everything" and integrate Watt TF into your CI/CD pipeline.
:::

Now it's time to generate the Terraform JSON. 



```bash
wtf build --input input.json --config blueprint.yaml --output terraform.json
```

The result should look like this:

```json
{
  "module": {
    "app": {
      "ms-orders": {
        "source": "git::https://git.acme.com/terraform-modules/app.git?ref=v1.0.0",
        "version": "v1.0.0",
        "port": 8080,
        "min_replicas": null,
        "max_replicas": null,
        "ram": "512Mi",
        "database": "module.database.ms-orders.dsn"
      }
    },
    "database": {
      "ms-orders": {
        "source": "git::https://git.acme.com/terraform-modules/database.git?ref=v1.0.0",
        "type": "postgres",
        "version": "14",
        "size": "20Gi"
      }
    }
  }
}
```

Please refer to the [reference documentation](../reference/cli.md) for more details on the CLI commands and options, and the [examples](../examples/overview.md) for more detailed use cases and configurations.
