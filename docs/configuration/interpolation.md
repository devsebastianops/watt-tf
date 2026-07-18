# Interpolation

Watt TF uses Google's CEL (Common Expression Language) to provide fast, safe, and powerful interpolation for your Terraform blueprints. Interpolation can be used dynamically within both target paths and value blocks.

## Context Variables

Every expression evaluated by the engine has access to the following top-level contexts:

- `input`: A map containing all custom configuration variables passed into the transformer.
- `env`: A map of all available environment variables.
- `item`: The current element when iterating via for_each (otherwise null).
- `item_index`: The 0-based index of the current element when iterating via for_each.

## Basic String Interpolation

You can inject variables into strings using the ${...} syntax. Anything inside the braces is evaluated as a CEL expression.

```yaml
transform:
  - target: resource.aws_s3_bucket.${input.environment}_storage
    value:
      bucket: "company-${input.project_name}-bucket"
```

If a string consists only of an interpolation expression, Watt TF preserves the underlying type (e.g., boolean, integer, map, or list) instead of forcing it into a string:

```yaml
transform:
  - target: resource.aws_instance.web
    value:
      monitoring: "${input.enable_monitoring}" # Resolves to a true/false boolean type
```

## Strict vs. Lenient Mode

Watt TF's behavior during interpolation depends heavily on the execution mode.

### Lenient Mode: <Badge text="default" vertical="middle" type="info" />

Missing or failed keys are gracefully replaced with null values. The engine continues processing.

### Strict Mode:

Missing or failed keys will lead to a syntax error and halt the entire transformation process.

You can toggle strict mode using the `--strict` flag in the CLI.

## Stripping null values

It can be useful to strip null values from the final output. You can do so by using the `--strip-null` flag in the CLI. This will remove any keys with null values from the final JSON output, which can be helpful for Terraform configurations that do not accept nulls.

## Common Expression Language (CEL)

CEL provides a rich set of built-in functions for string manipulation, math operations, and collection handling. You can use these functions to perform complex transformations directly within your blueprint.

In case you do not know CEL, here are some useful resources to get you started:
- [CEL Overview](https://cel.dev/overview/cel-overview)
- [CEL by Example](https://celbyexample.com/)

