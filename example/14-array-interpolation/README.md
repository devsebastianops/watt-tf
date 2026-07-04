# Example 14: Array Interpolation

This example demonstrates how to interpolate arrays directly into configuration objects.

## What it does

- Uses array interpolation to insert an array from the input into a configuration object
- Array is `${input.tags}` where input.tags is `[1, 2, 3]`

## Key Features

- **Direct array interpolation:** `tags: ${input.tags}`
- The entire value is replaced with the interpolated array result
- Arrays can be nested within objects
- Type is preserved (numbers remain numbers, not converted to strings)

## Expected Output

The array `[1, 2, 3]` from the input is directly interpolated into the `tags` field within a configuration object.
