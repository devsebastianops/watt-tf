# Conditions

Watt TF allows you to control whether a transformation should be executed or skipped by using the `if` property. The condition is evaluated as a CEL (Common Expression Language) expression that must resolve to a boolean value (true or false).

Conditional blocks are ideal for applying configurations dynamically based on the environment (e.g., adding heavy monitoring only in prod) or filtering out elements during loop iterations.

## Basic Conditions

To apply a transformation conditionally, add the `if` key to your transform block. The expression has full access to your `input` and `env` contexts.

```yaml
transform:
  # Only create production-grade replica counts if the environment matches
  - if: env.ENVIRONMENT == "prod"
    target: resource.aws_rds_cluster.db
    value:
      backup_retention_period: 30
      skip_final_snapshot: false

  # Enable specific configurations based on flags
  - if: input.enable_monitoring == true
    target: resource.aws_instance.web
    value:
      monitoring: true
```


## Evaluation Behavior

The behavior of evaluating conditional keys relies heavily on your configured execution mode.

Defaultwise, if an if expression references a key that does not exist in the context (e.g., `input.missing_flag == true`), the engine resolves the missing key to null. In CEL, evaluating `null == true` results in `false`. Therefore, the block is simply skipped without crashing the pipeline.

Strict mode, however, will not tolerate missing keys. If a key is referenced in an if expression but does not exist, the engine will throw a syntax error and halt the entire transformation process. You can enable strict mode using the `--strict` flag in the CLI.

## Available context variables

In conditions you have access to the following context variables:

- `input`: The root of your input data (e.g., `${input.nested.property}`).
- `env`: The root of your environment variables (e.g., `${env.STAGE}`).

Because conditions are always evaluated first, before any loops, you cannot access loop variables like `item` or `item_index` in the condition. 