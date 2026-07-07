# Example 23: Type Checking Functions

This example demonstrates all the type checking custom functions:
- `isString(any)` - Returns `true` if the value is a string, `false` otherwise
- `isNumber(any)` - Returns `true` if the value is a number (int or float), `false` otherwise
- `isArray(any)` - Returns `true` if the value is an array/list, `false` otherwise
- `isMap(any)` - Returns `true` if the value is a map/object, `false` otherwise
- `isBoolean(any)` - Returns `true` if the value is a boolean, `false` otherwise

## Input
```json
{
  "name": "John",
  "age": 30,
  "address": {
    "street": "123 Main St",
    "city": "New York"
  },
  "tags": ["developer", "gopher"]
}
```

## Config
The transformation checks the types of different values:
```yaml
transform:
  - target: output
    value:
      name_is_string: "${isString(input.name)}"
      age_is_number: "${isNumber(input.age)}"
      address_is_map: "${isMap(input.address)}"
      tags_is_array: "${isArray(input.tags)}"
      test_boolean: "${isBoolean(true)}"
      name_value: "${input.name}"
      age_value: "${input.age}"
      tags_count: "${input.tags.size()}"
      address_city: "${input.address.city}"
```

## Output
```json
{
  "output": {
    "name_is_string": true,        // "John" is a string
    "age_is_number": true,         // 30 is a number
    "address_is_map": true,        // address object is a map
    "tags_is_array": true,         // tags is an array
    "test_boolean": true,          // true is a boolean
    "name_value": "John",
    "age_value": 30,
    "tags_count": 2,
    "address_city": "New York"
  }
}
```

## Use Cases
- Type validation before transformation
- Conditional logic based on value types
- Safe navigation (check if map before accessing properties)
- Data quality checks before processing
- Runtime type assertions in CEL expressions
