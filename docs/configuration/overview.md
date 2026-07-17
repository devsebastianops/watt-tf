# Configuration Overview

Every **Watt TF** project centers around two structural layers: the **Input Data** (what your applications or developers require) and the **Blueprint Configuration** (how those requirements map to infrastructure).

This section covers the basic syntax, file structure, and top-level fields available in a Watt TF blueprint.

## The Anatomy of a Blueprint

A blueprint is a standard YAML file. Its primary purpose is to hold an array of declarative instructions inside the `transform` key. 

Here is a minimalist but complete blueprint structure:

```yaml
# blueprint.yaml

# The transform block is the core array containing all orchestration steps
transform:
  - target: resource.local_file.welcome
    value:
      filename: "welcome.txt"
      content: "Hello World"
```

## Inside a Transform Block

Each object inside the transform array is a separate operation evaluated in strict sequential order. A transform block accepts the following configuration keys:

```yaml
transform:
  - if: <CEL-Expression>       # (Optional) Conditional execution check
    for_each: <CEL-Expression> # (Optional) Collection to iterate over
    target: <String-Path>      # (Required) The dot-notated destination path
    value: <Any-Type>          # (Required) The data structure or primitive to inject
```

- `target` (The Destination): Points to the exact location in the final Terraform configuration where the value should live (e.g., resource.aws_s3_bucket.main).
- `value` (The Payload): Can be a single string, an integer, a boolean, an array, or a deeply nested map/object. Watt TF evaluates all dynamic variables inside this payload before placement.
- `if` (The Gatekeeper): A common expression language (CEL) string returning a boolean. If it yields false, the execution of this entire block terminates instantly.
- `for_each` (The Multiplier): A CEL expression that evaluates to a list or a range. It causes the block to clone itself and execute once for each element in the collection.

## Next Steps

Now that you know how a blueprint is structured globally, dive deeper into how the individual engines behave:

- [Transformations & Value Mapping](./transform.md) - Master how data is injected into payloads.
- [Path Resolution & Targets](./target.md) - Learn how dot-notation automatically builds nested JSON.
- [Variable Interpolation](./interpolation.md) - Read values from inputs and environment variables effortlessly.