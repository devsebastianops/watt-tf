# Transform

The `transform` block is the fundamental building block of any **Watt TF** blueprint. While the `target` defines *where* data goes, the `value` field within a transformation defines *what* is injected.

This section explains how the transformation engine processes payloads, supports various data types, and evaluates dynamic code blocks.

## Processing Order and Behavior

Watt TF processes the `transform` array inside your blueprint file in a strict **sequential top-down order**. 

1. **Evaluation:** The engine takes the `value` object and runs the variable interpolation engine over it, evaluating all `${...}` strings and CEL expressions.
2. **Placement:** The evaluated payload is placed directly at the destination specified by the `target`.
3. **Overwriting:** If a later transformation block writes to the exact same `target` path, it will replace or merge into the existing data structure (depending on your blueprint).

## Supported Data Types

The `value` field accepts any valid JSON/YAML data structure. You are not limited to objects; you can inject primitives or structural arrays depending on what your final Terraform configuration requires.

### 1. Primitive Values
You can inject strings, integers, floats, or booleans directly into a target.

```yaml
transform:
  - target: resource.aws_s3_bucket.main.force_destroy
    value: true # Boolean type preserved

  - target: resource.aws_instance.web.cpu_threads_per_core
    value: 2 # Integer type preserved
```

### 2. Complex Objects (Maps)
The most common use case is passing a full map of parameters to a resource.

```yaml
transform:
  - target: resource.aws_vpc.main
    value:
      cidr_block: "10.0.0.0/16"
      enable_dns_hostnames: true
```

### 3. Arrays

If your target points to an array block (such as an array of security group blocks or container definitions), you can pass a list structure.

```yaml
transform:
  - target: resource.aws_security_group.allow_tls.ingress
    value:
      - description: "TLS from VPC"
        from_port: 443
        to_port: 443
        protocol: "tcp"
        cidr_blocks: ["10.0.0.0/16"]
```

## Dynamic Values inside Payloads
You can seamlessly mix static properties with dynamic inline calculations using the `${...}` syntax. The entire block is evaluated in memory before being translated into the final structural layer.

```yaml
transform:
  - target: resource.aws_instance.app
    value:
      ami: ${input.ami_id}
      # Mixing hardcoded text with variables:
      instance_type: "t3.${input.size}" 
      # Dynamic object definitions based on environment variables:
      tags:
        Environment: ${env.STAGE}
        CostCenter: ${input.billing.code.upper()}
```

When interpolating variables, Watt TF is smart about types and preserves the original data type. This is crucial for Terraform, which will reject the wrong type.

