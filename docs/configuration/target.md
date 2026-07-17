# Target Path Resolution

The `target` property in a transformation block tells **Watt TF** exactly where to place your evaluated payload in the final Terraform JSON structure. 

By utilizing a flexible dot-notated string syntax, Watt TF acts as an automated structural architect—building deeply nested objects on the fly without requiring you to pre-define the parent object trees.

---

## Dynamic Object Creation

When Watt TF processes a target path like `resource.aws_s3_bucket.main`, it evaluates the string from left to right:

1. It checks if the top-level object `resource` exists. If not, it creates it.
2. It moves to the next segment (`aws_s3_bucket`) and inserts it as a nested map inside `resource`.
3. It resolves the final leaf node (`main`) and assigns your `value` payload directly to it.

Because this evaluation happens dynamically in memory, you can safely write to any deep path without worrying about initializing empty parent blocks.

---

## Dynamic Targets (String Interpolation)

The `target` path itself fully supports variable interpolation using the `${...}` syntax. This allows you to namespace resources dynamically based on developer input or loop variables (`item`).

```yaml
transform:
  # Stamping out a resource named exactly after the input service
  - target: resource.aws_instance.${input.service_name}
    value:
      ami: "ami-123456"

  # Using a loop variable to generate multiple unique resource keys
  - for_each: input.topics
    target: "resource.aws_sns_topic.${item}"
    value:
      name: ${item}
```

## Handling Special Characters (Escaping Dots)

Terraform frequently utilizes dots or slashes inside property keys—especially when dealing with resource tags, Kubernetes labels, or AWS integration configurations (e.g., tags.kubernetes.io/cluster/main).

If you write a path like resource.aws_instance.web.tags.company.com, Watt TF will interpret it as a deeply nested structure: tags -> company -> com.

To enforce a literal dot inside a single path segment, wrap that specific segment in backticks (`) or double quotes (depending on your preference/engine version):


```yaml
transform:
  # The backticks ensure "company.com" remains a single JSON key
  - target: resource.aws_instance.web.tags.`[company.com/managed-by](https://company.com/managed-by)`
    value: "WattTF"
    
  # Useful for complex cloud annotations
  - target: resource.kubernetes_deployment.api.metadata.annotations.`cert-manager.io/cluster-issuer`
    value: "letsencrypt-prod"
```

The compiled output correctly preserves the dot-notation inside the key name:

```json
{
  "resource": {
    "aws_instance": {
      "web": {
        "tags": {
          "[company.com/managed-by](https://company.com/managed-by)": "WattTF"
        }
      }
    }
  }
}
```

