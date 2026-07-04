# Example 16: For-Each Loops

This example demonstrates the `for_each` directive for iterating over arrays in the input and generating multiple resources.

## What It Does

- Reads an array of job definitions from input
- Uses `for_each` to iterate over each job
- Generates multiple job resources with dynamic target names (based on job name)
- Each resource gets a unique ID based on the loop index

## Key Features

- **for_each directive:** Iterates over input arrays
- **item variable:** Represents the current element (e.g., `item.name`)
- **item_index variable:** 0-based index of current element
- **Dynamic target:** Target can use `{{ item.* }}` to create unique resource names
- **CEL-based:** Uses CEL expressions for full flexibility

## Configuration Syntax

```yaml
transform:
  - target: "resource.type.{{ item.field | filter }}"
    for_each: input.array_path
    value:
      field1: "{{ item.prop1 }}"
      field2: "{{ item_index + 1 }}"
      # ... other fields
```

## Input

```json
{
  "jobs": [
    { "name": "my-job" },
    { "name": "another-job" },
    { "name": "third-job" }
  ]
}
```

## Configuration

```yaml
transform:
  - target: "resource.job.{{ item.name | snake_case }}"
    for_each: input.jobs
    value:
      name: "{{ item.name }}"
      enabled: true
      id: "job-{{ item_index + 1 }}"
```

**Breakdown:**
- `for_each: input.jobs` - Iterate over the `jobs` array
- `target: "resource.job.{{ item.name | snake_case }}"` - Dynamic target using item data
  - First iteration: `resource.job.my_job`
  - Second iteration: `resource.job.another_job`
  - Third iteration: `resource.job.third_job`
- `item` - Current job object (has `.name` property)
- `item_index` - 0, 1, 2 for each iteration
- `id: "job-{{ item_index + 1 }}"` - Results in job-1, job-2, job-3

## Output

Generates 3 job resources under `resource.job` with:
- Key: Job name in snake_case (my_job, another_job, third_job)
- ID: Auto-incremented (job-1, job-2, job-3)
- Name: Original name from input
- Enabled: Always true

## Available Variables in for_each

Inside a `for_each` block, these variables are available:

- `item` - The current element from the array (type depends on array content)
- `item_index` - The 0-based index of the current element (integer)
- `input` - The full input data (unchanged)
- `env` - Environment variables (unchanged)

## Advanced: Conditional Items

You can also use `if` conditions with `for_each`:

```yaml
transform:
  - target: "resource.job.{{ item.name | snake_case }}"
    for_each: input.jobs
    if: "item.enabled == true"  # Only process enabled jobs
    value:
      name: "{{ item.name }}"
      id: "job-{{ item_index + 1 }}"
```

## How to Test

```bash
go test ./tests/e2e -v -run 16-for-each-loops
```

Or run all E2E tests:
```bash
go test ./tests/e2e -v
```

## Manual Test

```bash
./bin/wtf build \
  --config example/16-for-each-loops/.wtf.yaml \
  --input example/16-for-each-loops/input.json \
  --output output.tf.json
```
