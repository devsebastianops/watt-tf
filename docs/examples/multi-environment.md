# Multi Environment

Managing multiple environments (like `staging` and `production`) is one of the most common challenges in infrastructure management.

With Watt TF, your infrastructure footprint remains a single, static blueprint, while your environments are defined as pure, minimalist data files.

## 1. Environment Definitions (`inputs/`)

Instead of duplicating infrastructure code, you simply maintain flat data files for your environments.

::: code-group
```yaml [staging.yaml]
environment: staging
domain: staging.api.mycompany.com
db_replica_count: 1
instance_type: t4g.micro
backup_retention_days: 1
```

```yaml [production.yaml]
environment: production
domain: api.mycompany.com
db_replica_count: 3
instance_type: m6g.large
backup_retention_days: 30
```
:::

## 2. The Shared Blueprint ( `blueprint.yaml` )

Your blueprint stays completely independent of the environment. It uses Google's Common Expression Language (CEL) to dynamically adapt values based on the input data, ensuring strict architectural standards across all environments.

```yaml
transform:
  # 1. Create the application server
  - target: resource.aws_instance.api_server
    value:
      ami: "ami-0c55b159cbfafe1f0"
      instance_type: ${input.instance_type}
      
      tags:
        Name: "api-server-${input.environment}"
        Environment: ${input.environment}

  # 2. Create the primary database
  - target: resource.aws_db_instance.database
    value:
      allocated_storage: 20
      engine: "postgres"
      instance_class: ${input.instance_type == "m6g.large" ? "db.m6g.large" : "db.t4g.micro"}
      backup_retention_period: ${input.backup_retention_days}

  # 3. Scale database replicas dynamically
  # Watt TF uses the built-in CEL range() or integer values to scale blocks on the fly
  - for_each: range(0, input.db_replica_count)
    target: "resource.aws_db_instance.database_replica_${item}"
    value:
      replicate_source_db: "aws_db_instance.database.id"
      instance_class: ${input.instance_type == "m6g.large" ? "db.m6g.large" : "db.t4g.micro"}
```

## The deployment process

Because the output is pure TF JSON, your deployment pipeline becomes incredibly predictable. You just feed the specific environment file into the compiler.

```bash
# 1. Compile Staging
wtf build -i inputs/staging.yaml -b blueprint.yaml -o dist/staging.tf.json

# 2. Compile Production
wtf build -i inputs/production.yaml -b blueprint.yaml -o dist/production.tf.json
```

### CI / CD Example

Inside your pipeline, you can dynamically switch your Terraform target directory or use the compiled backend configurations:

```bash
# Example workflow deployment for a specific environment
ENV="staging"

wtf build -i inputs/${ENV}.yaml -b blueprint.yaml -o current.tf.json

terraform init
terraform plan -out=${ENV}.tfplan
terraform apply ${ENV}.tfplan
```

