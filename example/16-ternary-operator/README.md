# Example 16: Ternary Operator

This example demonstrates the use of ternary operators (conditional expressions) in CEL to create dynamic configurations based on input values.

Syntax: `condition ? true_value : false_value`

## Input
```yaml
env: prod
region: us-west-2
debug_enabled: true
```

## Config
The transformation uses ternary operators to set values based on conditions:
```yaml
transform:
  - target: resource.server.default
    value:
      port: "${input.env == 'prod' ? 443 : 80}"
      protocol: "${input.env == 'prod' ? 'https' : 'http'}"
      ssl_enabled: "${input.env == 'prod' ? true : false}"
  
  - target: resource.config.logging
    value:
      level: "${input.debug_enabled ? 'DEBUG' : 'INFO'}"
      output: "${input.debug_enabled ? 'stdout' : 'file'}"
  
  - target: resource.location
    value:
      region: "${input.region}"
      environment: "${input.env}"
      backup_region: "${input.region == 'us-west-2' ? 'us-east-1' : 'eu-central-1'}"
```

## Output
```json
{
  "resource": {
    "server": {
      "default": {
        "port": 443,
        "protocol": "https",
        "ssl_enabled": true
      }
    },
    "config": {
      "logging": {
        "level": "DEBUG",
        "output": "stdout"
      }
    },
    "location": {
      "region": "us-west-2",
      "environment": "prod",
      "backup_region": "us-east-1"
    }
  }
}
```

## Key Features
- Simple equality comparisons: `input.env == 'prod'`
- Boolean value selection: `true ? 'https' : 'http'`
- Nested conditions can be combined with other expressions
- Works with any type (strings, numbers, booleans)

## Use Cases
- Environment-specific configuration (dev/staging/prod)
- Feature flags and conditional resource creation
- Dynamic resource properties based on input conditions
- Safe default values based on conditions
