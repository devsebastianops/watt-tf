# Platform Engineering

One of the most powerful use cases for Watt TF is building an Internal Developer Platform (IDP), powered by JSON and Terraform.

In a modern platform engineering setup, you want to shield application developers from the vast complexity of cloud infrastructure. Developers shouldn't need to know how to configure AWS VPCs, IAM Roles, or complex Terraform modules. They just want to deploy their application.

With Watt TF, you can let developers provide a minimalist, business-centric configuration file, while your platform team maintains a robust infrastructure blueprint.

In this example, we want to allow developers to provision a **Microservice** with an optional **PostgreSQL Database**. 

## 1. The Developer Experience (`input.yaml`)
Instead of confronting developers with hundreds of lines of HCL, they only submit a simple, clean interface. They don't even know Terraform is running under the hood:

```yaml
service_name: payment-processor
team: billing
tier: critical

compute:
  cpu: "1.0"
  memory: "2Gi"

# The developer just flags that they need a DB, the platform handles the rest
database:
  enabled: true
  storage_gb: 20
```

## 2. The Platform Team Blueprint (`blueprint.yaml`)

Your Platform Team maintains the infrastructure standards. This blueprint takes the simple developer input and generates a secure AWS ECS (Fargate) Service, an IAM Role, and a production-grade Amazon RDS instance—completely dynamically.

```yaml
transform:
  # 1. Generate the ECS Task Definition
  - target: resource.aws_ecs_task_definition.${input.service_name}
    value:
      family: ${input.service_name}
      requires_compatibilities: ["FARGATE"]
      cpu: ${input.compute.cpu}
      memory: ${input.compute.memory}
      container_definitions:
        - name: app
          image: "cloudsmith.io/my-company/${input.service_name}:latest"
          essential: true

  # 2. Generate standard corporate Tags automatically
  - target: resource.aws_ecs_task_definition.${input.service_name}.tags
    value:
      ManagedBy: "PlatformTeam-Watt TF"
      OwnerTeam: ${input.team}
      Criticality: ${input.tier}

  # 3. Conditional Database Provisioning
  # If the developer set database.enabled to true, build a secure RDS instance
  - if: has(input.database) && input.database.enabled
    target: resource.aws_db_instance.${input.service_name}_db
    value:
      allocated_storage: ${input.database.storage_gb}
      engine: "postgres"
      engine_version: "15.4"
      instance_class: ${input.tier == "critical" ? "db.r6g.large" : "db.t4g.micro"}
      db_name: ${input.service_name.replace("-", "_")}
      username: "db_admin"
      password: ${env.DB_MASTER_PASSWORD} # Injected safely at build time
      skip_final_snapshot: true
``` 

## 3. The build process

In your CI/CD pipeline (e.g., GitHub Actions or GitLab CI), the platform team runs Watt TF before executing Terraform.

```bash
# Find all microservice yaml definitions and build them into Terraform JSON
find . -type f -name "microservice-*.yaml" | while read -r file; do
    echo "Processing $file..."
    wtf build -i "$file" -b blueprint.yaml -o "${file}.tf.json"
done

# Run standard Terraform deployment
terraform init
terraform apply -auto-approve
```