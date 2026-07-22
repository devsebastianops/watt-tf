# Loops

Loops are a way to create multiple instances of a resource or module based on a collection of values. In Watt TF you can use the `for_each` keyword in your blueprint to iterate over a list or map.

## Using `for_each` in a Transform Block

```yaml
transform:
  - for_each: input.tags
    target: "resource.my_provider_compute.default.volumes"
    value:
      name: ${item.name}
      size: ${item.size}
      id: ${item_index + 1}
```

When using the `for_each` keyword, Watt TF will iterate over the collection provided via CEL, exposing the current element as `item`, and its 0-indexed position as `item_index`.

Context variables available during iteration:
- `item`: The current element when iterating via for_each (otherwise null).
- `item_index`: The index of the current element when iterating via for_each (otherwise null).

