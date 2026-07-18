# CEL Functions

Watt TF extends the standard Common Expression Language (CEL) with custom utility functions to handle data serialization, encoding, and type inspection directly within your blueprints.

## Serialization & Encoding

These functions allow you to transform data formats dynamically, which is particularly useful when injecting complex objects into Terraform tags or resource attributes.

| Function | Signature | Description |
| :--- | :--- | :--- |
| `toJSON` | `any -> string` | Converts any object into a JSON-formatted string. |
| `fromJSON` | `string -> any` | Parses a JSON-encoded string back into a structural object. |
| `toBase64` | `string -> string` | Encodes a string to Base64. |
| `fromBase64`| `string -> string` | Decodes a Base64-encoded string. |

## Type Inspection

These functions provide a safe way to check the underlying type of a variable. They are especially useful when working with dynamic inputs where you need to branch your logic based on the data structure.

| Function | Signature | Description |
| :--- | :--- | :--- |
| `isMap` | `any -> bool` | Returns `true` if the value is a map/object. |
| `isString` | `any -> bool` | Returns `true` if the value is a string. |
| `isArray` | `any -> bool` | Returns `true` if the value is a list/array. |
| `isNumber` | `any -> bool` | Returns `true` if the value is an integer or float. |
| `isBoolean`| `any -> bool` | Returns `true` if the value is a boolean. |

## Standard CEL Resources

Beyond these custom extensions, Watt TF supports the full standard CEL feature set. If you need to perform more advanced operations, we recommend these resources to get started:

- [CEL Overview](https://cel.dev/overview/cel-overview): The official spec and language guide.
- [CEL by Example](https://celbyexample.com/): A practical guide with examples for common use cases.

<br>

::: info Watt TF CEL libraries
Watt TF automatically includes common CEL libraries. You can use standard functions like `size()`, `has()`, or collection macros like `.exists()` and `.map()` out of the box in your blueprints.
::: 