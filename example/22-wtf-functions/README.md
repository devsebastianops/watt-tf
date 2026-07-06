# Example 22: Custom wtf Functions

This example demonstrates the custom CEL functions available under the `wtf` namespace:
- `toJSON(any)` - Converts any value to a JSON string
- `fromJSON(string)` - Parses a JSON string into an object
- `toBase64(string)` - Encodes a string to base64
- `fromBase64(string)` - Decodes a base64 string

## Input
```json
{
  "message": "Hello World",
  "data": {
    "key": "value",
    "nested": [1, 2, 3]
  },
  "encoded": "SGVsbG8gQmFzZTY0"
}
```

## Config
The transformation uses all four custom functions:
- `toJSON()` converts the data object to a JSON string
- `fromJSON()` parses the JSON string back to an object
- `toBase64()` encodes the message string to base64
- `fromBase64()` decodes the base64 string
- Complex combinations like encoding JSON as base64

## Output
```json
{
  "output": {
    "json_string": "{\"key\":\"value\",\"nested\":[1,2,3]}",
    "parsed_back": {
      "key": "value",
      "nested": [1, 2, 3]
    },
    "base64_encoded": "SGVsbG8gV29ybGQ=",
    "base64_decoded": "Hello Base64",
    "json_as_base64": "eyJrZXkiOiJ2YWx1ZSIsIm5lc3RlZCI6WzEsMiwzXX0="
  }
}
```

## Key Features
- Full JSON serialization/deserialization support
- Base64 encoding/decoding for strings
- Composable functions (e.g., `toBase64(toJSON(object))`)
- Type-safe conversion between different representations
