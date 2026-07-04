# Jinja2 Template Engine for Complex Scenarios

## Overview

A template engine integration that allows users to define complex transformation logic using Jinja2 templates, supporting loops, conditionals, filters, and macros for sophisticated infrastructure generation.

## User Experience

### Basic Template Usage

```yaml
# .wtf.yaml
transform:
  - target: resource.aws_subnet.subnets
    template: templates/subnet.json.j2
    context:
      vpc_cidr: "${input.vpc_cidr}"
      availability_zones: "${input.azs}"
      enable_nat: "${input.enable_nat}"
```

### Template Example

```jinja2
{# templates/subnet.json.j2 #}
{
  {% for az in availability_zones %}
  "subnet_{{ loop.index }}": {
    "vpc_id": "{{ vpc_id }}",
    "cidr_block": "{{ subnet_cidrs[loop.index0] }}",
    "availability_zone": "{{ az }}",
    "tags": {
      "Name": "subnet-{{ az }}",
      "Index": {{ loop.index }}
    }
    {% if enable_nat %}
    ,
    "nat_gateway_id": "${module.nat_gateway_{{ loop.index }}.id}"
    {% endif %}
  }
  {%- if not loop.last %},{% endif %}
  {% endfor %}
}
```

### Template with Filters & Macros

```jinja2
{# templates/instance-config.json.j2 #}
{% macro instance_tags(name, env, cost_center) %}
{
  "Name": "{{ name }}",
  "Environment": "{{ env | upper }}",
  "CostCenter": "{{ cost_center }}",
  "ManagedBy": "wtf"
}
{% endmacro %}

{
  {% for instance in instances %}
  "instance_{{ loop.index }}": {
    "ami": "{{ ami_ids[instance.type] }}",
    "instance_type": "{{ instance.type }}",
    "tags": {{ instance_tags(instance.name, env, cost_center) | indent(4) }}
  }
  {%- if not loop.last %},{% endif %}
  {% endfor %}
}
```

## Technical Implementation

### 1. Template Engine Abstraction

```go
// internal/template/engine.go
package template

import (
    "github.com/jinja2cpp/jinja2cpp-go"
)

type Engine interface {
    Render(tmpl string, context map[string]interface{}) (string, error)
    RenderFile(filepath string, context map[string]interface{}) (string, error)
}

type Jinja2Engine struct {
    env *jinja2.Environment
}

func NewJinja2Engine() *Jinja2Engine {
    env := jinja2.NewEnvironment()
    // Add custom filters
    env.AddFilter("upper", filterUpper)
    env.AddFilter("kebab_case", filterKebabCase)
    return &Jinja2Engine{env: env}
}

func (e *Jinja2Engine) Render(tmpl string, context map[string]interface{}) (string, error) {
    template, err := e.env.FromString(tmpl)
    if err != nil {
        return "", err
    }
    return template.Render(context)
}
```

### 2. Template Integration in Transformer

```go
// internal/transformer/template_transformer.go
package transformer

type TemplateTransform struct {
    Target   string
    Template string
    Context  map[string]interface{}
}

func (t *TemplateTransform) Transform(engine template.Engine) (map[string]interface{}, error) {
    // Interpolate context values
    interpolatedContext := interpolate(t.Context)
    
    // Render template
    output, err := engine.RenderFile(t.Template, interpolatedContext)
    if err != nil {
        return nil, fmt.Errorf("template rendering failed: %w", err)
    }
    
    // Parse JSON output
    var result map[string]interface{}
    if err := json.Unmarshal([]byte(output), &result); err != nil {
        return nil, fmt.Errorf("template output is not valid JSON: %w", err)
    }
    
    return result, nil
}
```

### 3. Config Loading

```go
// internal/config/loader.go (updated)
type Config struct {
    Transform []Transformable `yaml:"transform"`
    Templates map[string]Template `yaml:"templates"`
}

type Template struct {
    Path    string                 `yaml:"path"`
    Engine  string                 `yaml:"engine"` // "jinja2", "golang", etc
    Context map[string]interface{} `yaml:"context"`
}
```

### 4. Custom Filters & Functions

```go
// internal/template/filters.go
package template

func filterKebabCase(value interface{}) string {
    s := fmt.Sprintf("%v", value)
    return strings.ReplaceAll(strings.ToLower(s), "_", "-")
}

func filterSnakeCase(value interface{}) string {
    // Convert to snake_case
}

func filterYAMLToJSON(value interface{}) string {
    // Convert YAML to JSON
}

// Custom functions
func funcCIDRSubnet(cidr string, index int) string {
    // Calculate subnet from CIDR
}

func funcGeneratePassword(length int) string {
    // Generate secure password
}
```

### 5. Template Validation

```go
// internal/template/validator.go
package template

func ValidateTemplate(filepath string, context map[string]interface{}) error {
    // Parse and validate template syntax
    // Check for undefined variables
    // Validate Jinja2 syntax
}
```

## Implementation Steps

### Phase 1: Basic Jinja2 Integration
1. Add jinja2 package dependency
2. Implement template engine abstraction
3. Add template rendering in transformer
4. Update config loader for template support

### Phase 2: Custom Filters & Functions
1. Implement common filters (kebab_case, snake_case, etc)
2. Add infrastructure-specific functions (CIDR calculation, etc)
3. Template validation & error reporting

### Phase 3: Advanced Features
1. Template inheritance & includes
2. Template caching & optimization
3. Template library & reusable components

## File Structure

```
internal/
├── template/
│   ├── engine.go          # Template engine interface
│   ├── jinja2.go          # Jinja2 implementation
│   ├── filters.go         # Custom filters
│   ├── functions.go       # Custom functions
│   └── validator.go       # Template validation
│
└── transformer/
    └── template_transformer.go  # Template transform logic

examples/
└── template-examples/
    ├── vpc-subnets.yaml
    ├── templates/
    │   ├── subnet.json.j2
    │   └── instance.json.j2
    └── input.json
```

## Testing

```go
// internal/template/jinja2_test.go
func TestTemplateRender(t *testing.T) {
    engine := NewJinja2Engine()
    
    tmpl := `
    {
      {% for item in items %}
      "{{ item.name }}": {{ item.value }}
      {%- if not loop.last %},{% endif %}
      {% endfor %}
    }
    `
    
    context := map[string]interface{}{
        "items": []map[string]interface{}{
            {"name": "a", "value": "1"},
            {"name": "b", "value": "2"},
        },
    }
    
    result, err := engine.Render(tmpl, context)
    assert.NoError(t, err)
    
    var output map[string]interface{}
    err = json.Unmarshal([]byte(result), &output)
    assert.NoError(t, err)
}
```

## .wtf.yaml Example

```yaml
transform:
  # Simple template
  - target: resource.aws_instance
    template: templates/instances.json.j2
    context:
      instances: "${input.instances}"
      ami_id: "${input.ami_id}"
  
  # Template with inline content
  - target: resource.aws_vpc.main
    template_inline: |
      {
        "cidr_block": "{{ vpc_cidr }}",
        "availability_zones": [
          {% for az in azs %}
          "{{ az }}"
          {%- if not loop.last %},{% endif %}
          {% endfor %}
        ]
      }
    context:
      vpc_cidr: "10.0.0.0/16"
      azs: "${input.availability_zones}"
```

## Benefits

- ✅ Support for loops and conditionals
- ✅ Reusable template components
- ✅ Powerful filters and functions
- ✅ Cleaner than complex string interpolation
- ✅ Industry-standard template language

## Dependencies

```go
// go.mod
require (
    github.com/jinja2cpp/jinja2cpp-go v0.x.x
)
```

## Estimated Effort

- **Phase 1**: 20-30 hours (basic integration)
- **Phase 2**: 15-20 hours (filters & functions)
- **Phase 3**: 10-15 hours (advanced features)
