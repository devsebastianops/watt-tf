# Examples

Comprehensive examples demonstrating all features of Watt TF.

## Table of Contents

### Beginner - Core Concepts

1. **[01-simple-json](01-simple-json)** - Basic JSON transformation
   - Learn the fundamentals of transforming JSON to Terraform
   - Understand basic string interpolation
   - Perfect starting point for new users

2. **[02-simple-yaml](02-simple-yaml)** - YAML input support
   - Same as 01-simple-json but with YAML input instead
   - Shows format-agnostic input handling
   - Useful if your configs are already in YAML

### Intermediate - String Interpolation & Conditions

3. **[03-interpolate-target](03-interpolate-target)** - Dynamic target names
   - Create resource names dynamically using `${input.xyz}`
   - Combine string interpolation with target naming
   - Essential for multi-instance deployments

4. **[04-conditions](04-conditions)** - Conditional transforms
   - Apply transformations only when conditions are met
   - Use CEL expressions in `if` fields
   - Perfect for multi-environment setups (prod vs dev)

### Advanced - Complex Logic & Merging

5. **[05-deep-merge](05-deep-merge)** - Multiple transforms per target
   - Apply multiple transformations to the same resource
   - Values are intelligently merged together
   - Great for building up resource configs in stages

6. **[06-nested-paths](06-nested-paths)** - Deep nesting support
   - Create arbitrarily deep resource hierarchies
   - Automatic intermediate map creation
   - Useful for complex infrastructure structures

### Feature Showcase - Advanced Capabilities

7. **[07-type-preservation](07-type-preservation)** - Type preservation
   - Integers, booleans, and arrays maintain their types
   - Understand when interpolation creates strings vs primitives
   - Critical for correct Terraform type checking

8. **[08-array-values](08-array-values)** - Array interpolation
   - Interpolate values within arrays
   - Mix interpolated and literal values in the same array
   - Essential for resource lists and tags

9. **[09-complex-cel](09-complex-cel)** - Complex CEL expressions
   - Combine multiple conditions with `&&` and `||`
   - Use parentheses for grouping logic
   - Build sophisticated conditional rules

10. **[10-string-methods](10-string-methods)** - String method conditions
    - Use `.startsWith()`, `.contains()`, `.endsWith()` for pattern matching
    - Create condition logic based on string patterns
    - Perfect for naming conventions and validations

### Edge Cases & Robustness

11. **[11-null-handling](11-null-handling)** - Null and empty values
    - Properly handle null values from input
    - Work with empty arrays and objects
    - Build optional infrastructure components

12. **[12-missing-keys](12-missing-keys)** - Error handling & robustness
    - Gracefully handle missing or incomplete input
    - Mix required and optional configuration keys
    - Build flexible infrastructure templates

---

## Running Examples

### Run a Single Example

```bash
# Basic syntax
wtf build --input example/XX-name/input.json \
  --config example/XX-name/.wtf.yaml \
  --output example/XX-name/watt.tf.json
```

### Run as One-Liner

Each example has a one-liner command in its README.

### View Generated Output

```bash
cat example/01-simple-json/watt.tf.json
```

### Run All E2E Tests

```bash
go test ./tests/e2e -v
```

---

## Learning Path

**New to Watt TF?** Start here:

1. 01-simple-json → Understand the basics
2. 03-interpolate-target → Learn interpolation
3. 04-conditions → Add conditional logic
4. 05-deep-merge → Handle complex configs
5. 10-string-methods → Advanced conditions
6. Explore others based on your needs

---

## Feature Matrix

| Feature | Example(s) |
|---------|-----------|
| JSON Input | 01, 03, 04, 05, 06, 07, 08, 09, 10, 11, 12 |
| YAML Input | 02 |
| String Interpolation | 01, 02, 03, 04, 05, 06, 07, 08, 09, 10 |
| Target Interpolation | 03, 09, 10 |
| Simple Conditions | 04, 09, 10 |
| Complex Conditions | 09 |
| CEL String Methods | 10 |
| Deep Merge | 05 |
| Nested Paths | 06 |
| Type Preservation | 07 |
| Arrays | 08 |
| Null Handling | 11 |
| Error Handling | 12 |

---

## Tips & Tricks

- **Start Simple**: Begin with 01-simple-json and gradually add complexity
- **Test Locally**: Each example is fully self-contained and testable
- **Reuse Configs**: Copy `.wtf.yaml` from examples as templates for your own projects
- **View Expected Output**: Check `expected.tf.json` in each example directory to see what the output should look like
- **Read READMEs**: Each example's README explains the use case and expected behavior

---

## Contributing

Have an idea for a new example? We'd love to see it!

1. Create a new directory: `example/NN-description`
2. Add:
   - `input.json` or `input.yaml` (your input data)
   - `.wtf.yaml` (transformation config)
   - `expected.tf.json` (expected output)
   - `README.md` (explanation and one-liner)
3. Submit a pull request

---

**Happy Learning!** 🚀
