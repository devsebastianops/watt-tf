# Interactive CLI Mode

## Overview

An interactive CLI mode that guides users through building transformations step-by-step, with real-time validation, suggestions, and preview of generated configurations.

## User Experience

### Starting Interactive Mode

```bash
$ wtf interactive
# or
$ wtf build -i

Welcome to Watt TF Interactive Builder!
Let's build your Terraform configuration step by step.

Step 1: Input File
────────────────────────────────────────────────────
? Select input file: (Tab to browse) ./input.json
✓ Found input.json
? Preview input? [Y/n]: Y

{
  "project": "myapp",
  "environment": "prod",
  "region": "us-west-2"
}
```

### Building Transforms Interactively

```
Step 2: Create Transform
────────────────────────────────────────────────────
? Transform #1 target path: resource.aws_s3_bucket.data
? Add value? [Y/n]: Y

Resource Type: aws_s3_bucket
Fields to configure:
  ✓ bucket (required)
  ⚪ acl (optional)
  ⚪ versioning (optional)
  ⚪ tags (optional)

? Configure 'bucket' [required]:
  Input type: [string/number/boolean/array/object]
  └─ string
? Value (supports ${input.path}):
  $ Enter value: ${input.project}-${input.environment}-data
  
? Add another field? [Y/n]: Y
? Configure 'tags' [optional]:
  └─ object
? Number of tags: 2
  ? Tag 1 key: Environment
  ? Tag 1 value: ${input.environment}
  ? Tag 2 key: Project
  ? Tag 2 value: ${input.project}

Preview of this resource:
{
  "resource": {
    "aws_s3_bucket": {
      "data": {
        "bucket": "myapp-prod-data",
        "tags": {
          "Environment": "prod",
          "Project": "myapp"
        }
      }
    }
  }
}

? Looks good? [Y/n]: Y
```

### Adding Conditions

```
Step 3: Add Condition (Optional)
────────────────────────────────────────────────────
? Add condition to this transform? [Y/n]: Y
? Condition expression:
  Available variables: input.*
  $ ${input.environment} == 'prod'

? Preview with condition:
  Will this transform run?
  Input: {"environment": "prod"} → YES ✓
  Input: {"environment": "dev"} → NO ✗

? Condition looks good? [Y/n]: Y
```

### Review & Save

```
Step 4: Review & Save
────────────────────────────────────────────────────
Summary of your configuration:
  Transforms: 3
  With conditions: 2
  Output file: watt.tf.json

? Save config as: .wtf.yaml
? Validate before saving? [Y/n]: Y

✓ Configuration validated successfully!
✓ Saved to .wtf.yaml

Generated .wtf.yaml:
───────────────────
transform:
  - target: resource.aws_s3_bucket.data
    value:
      bucket: "${input.project}-${input.environment}-data"
      tags:
        Environment: "${input.environment}"
        Project: "${input.project}"
...

? Build configuration now? [Y/n]: Y
✓ Generated watt.tf.json (324 bytes)
```

## Technical Implementation

### 1. Interactive CLI Framework

```go
// internal/cli/interactive.go
package cli

import (
    "github.com/charmbracelet/bubbles/textinput"
    "github.com/charmbracelet/bubbles/list"
    "github.com/charmbracelet/bubbletea"
)

type InteractiveBuilder struct {
    input      map[string]interface{}
    transforms []Transformable
    currentStep int
}

func (ib *InteractiveBuilder) Run() error {
    model := NewModel(ib)
    _, err := bubbletea.NewProgram(model).Run()
    return err
}
```

### 2. Step-by-Step Wizard

```go
// internal/cli/wizard/steps.go
package wizard

type Step interface {
    Title() string
    Execute(ctx Context) (StepResult, error)
    Validate() error
}

type SelectInputStep struct{}
type SelectTargetStep struct{}
type ConfigureValueStep struct{}
type AddConditionStep struct{}
type ReviewStep struct{}
type SaveStep struct{}

func (s *SelectInputStep) Execute(ctx Context) (StepResult, error) {
    // File picker UI
    // Load and preview input
}
```

### 3. Real-time Validation

```go
// internal/cli/validator/realtime.go
package validator

type RealtimeValidator struct {
    schema SchemaProvider
}

func (rv *RealtimeValidator) ValidateField(fieldName string, value interface{}) ([]Warning, error) {
    // Validate field type
    // Check value format
    // Suggest corrections
}

func (rv *RealtimeValidator) PreviewInterpolation(expr string, input map[string]interface{}) (interface{}, error) {
    // Show what ${input.path} will evaluate to
    // Show type information
}
```

### 4. Schema-aware Configuration

```go
// internal/cli/schema/provider.go
package schema

type ProviderSchema struct {
    Resources map[string]ResourceSchema
}

type ResourceSchema struct {
    Name     string
    Fields   map[string]FieldSchema
    Required []string
}

type FieldSchema struct {
    Type        string // string, number, boolean, object, array
    Description string
    Default     interface{}
    Enum        []interface{}
}
```

### 5. Preview & Diff

```go
// internal/cli/preview/preview.go
package preview

func PreviewTransform(transform Transformable, input map[string]interface{}) (string, error) {
    // Show what this transform will produce
    // Format as pretty JSON
}

func PreviewDiff(old, new map[string]interface{}) string {
    // Show what changed between old and new
    // Highlight additions/deletions/modifications
}
```

### 6. UI Components

```go
// internal/cli/ui/components.go

// File picker
func FilePicker(startPath string) (string, error)

// Key-value input
func KeyValueInput(title string, initialData map[string]interface{}) (map[string]interface{}, error)

// Expression editor with autocomplete
func ExpressionEditor(placeholder string) (string, error)

// JSON preview
func JSONPreview(data interface{}, maxLines int) string

// Spinner for long operations
func ShowSpinner(message string, fn func() error) error
```

## File Structure

```
internal/cli/
├── interactive.go         # Main interactive mode
├── wizard/
│   ├── steps.go          # Step definitions
│   ├── context.go        # Wizard context
│   └── flow.go           # Step flow control
├── ui/
│   ├── components.go     # UI components
│   ├── styles.go         # Styling
│   └── colors.go         # Color scheme
├── preview/
│   ├── preview.go        # Preview generation
│   └── diff.go           # Diff display
└── validator/
    ├── realtime.go       # Real-time validation
    └── suggestions.go    # Suggestions & hints
```

## Dependencies

```go
// go.mod
require (
    github.com/charmbracelet/bubbletea v0.x.x
    github.com/charmbracelet/bubbles v0.x.x
    github.com/charmbracelet/lipgloss v0.x.x
    github.com/manifoldco/promptui v0.x.x
)
```

## User Features

```bash
# Start interactive mode
wtf interactive

# Interactive with existing config
wtf interactive --config .wtf.yaml

# Non-interactive mode (existing)
wtf build --input config.json --config .wtf.yaml

# Hybrid: load template, then edit interactively
wtf interactive --template aws-vpc
```

## Testing

```go
// internal/cli/interactive_test.go
func TestSelectInputStep(t *testing.T) {
    step := &SelectInputStep{}
    ctx := NewTestContext()
    
    result, err := step.Execute(ctx)
    assert.NoError(t, err)
    assert.Equal(t, "input.json", result.Input)
}
```

## Benefits

- ✅ Lower learning curve
- ✅ Real-time feedback
- ✅ Guided experience
- ✅ Auto-completion & suggestions
- ✅ Visual preview before save

## Estimated Effort

- **Phase 1**: 30-40 hours (basic wizard, UI components)
- **Phase 2**: 20-25 hours (validation, suggestions)
- **Phase 3**: 15-20 hours (templates, advanced features)
