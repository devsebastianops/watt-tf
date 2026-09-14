# Environment Variables

This example shows how to use environment variables within Watt TF blueprints.

## Scenario

We want to create a blueprint to generate terraform resources for a microservice and a database. The database credentials will be provided through environment variables to prevent hardcoding sensitive information in the blueprint.

### Input

```yaml
microservice:
  name: orders
  port: 8080
  image: myregistry/orders:latest
```

### Blueprint

```yaml
transform:
  - target: module.database
    value:
        source: "git::https://github.com/myorg/myrepo.git?ref=v1.2.3"
        username: "${env.DATABASE_USERNAME}"
        password: "${env.DATABASE_PASSWORD}"
  - target: module.microservice
    value:
        source: "git::https://github.com/myorg/myrepo.git?ref=v1.2.3"
        name: "${input.microservice.name}"
        port: "${input.microservice.port}"
        image: "${input.microservice.image}"
        database: "${module.database}"
```

### Command

```bash
export DATABASE_USERNAME="myuser"
export DATABASE_PASSWORD="mypassword"
wtf build --input input.yaml --blueprint blueprint.yaml --output build.tf.json
```
