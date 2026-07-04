# Feature: Support for Multiple Template Engines

## Overview

Extend wtf beyond Jinja2 to support multiple template engines, allowing users to choose the best-fit templating language for their infrastructure needs.

### User Story
> As an DevOps engineer working with diverse teams, I want to use template engines familiar to my team (Go, Liquid, Handlebars) without learning new syntax, so that template development is faster and more collaborative.

---

## Feature Specification

### Core Capability

Support a pluggable template engine system that allows users to select different engines per-template:

```yaml
# .wtf.yaml
defaults:
  template_engine: "jinja2"  # Global default

transform:
  # Uses global default (jinja2)
  - target: resource.aws_instance.main
    template: templates/instance.j2
    context:
      instances: ${input.instances}
  
  # Overrides with Go templates
  - target: resource.aws_vpc.main
    engine: "go"
    template: templates/vpc.gotml
    context:
      vpc_config: ${input.vpc}
  
  # Uses Liquid for safety-critical configs
  - target: resource.storage.main
    engine: "liquid"
    template: templates/storage.liquid
    context:
      buckets: ${input.buckets}
```

### Supported Engines

#### 1. **Go Templates** (Recommended Default)

**Engine Name:** `go`

```yaml
engine: "go"
template: |
  {
    {{- range .instances }}
    "{{ .name }}": {
      "instance_type": "{{ .type }}"
    }
    {{- end }}
  }
```

**Pros:**
- ✅ Built-in to Go (zero dependencies)
- ✅ Familiar to Terraform users
- ✅ Highest performance
- ✅ Industry standard for IaC

**Cons:**
- ❌ Fewer built-in filters
- ❌ Syntax less intuitive for non-Go developers

**Dependencies:** None (stdlib)

**Example:**
```go
// internal/template/engines/go_engine.go
type GoEngine struct{}

func (ge *GoEngine) Name() string { return "go" }

func (ge *GoEngine) Render(template string, context map[string]interface{}) (string, error) {
    tmpl, err := text.New("tmpl").Parse(template)
    if err != nil {
        return "", fmt.Errorf("go template parse error: %w", err)
    }
    
    var buf bytes.Buffer
    if err := tmpl.Execute(&buf, context); err != nil {
        return "", fmt.Errorf("go template render error: %w", err)
    }
    
    return buf.String(), nil
}
```

---

#### 2. **Jinja2** (Feature-Rich)

**Engine Name:** `jinja2`

```yaml
engine: "jinja2"
template: |
  {
    {% for instance in instances %}
    "{{ instance.name }}": {
      "instance_type": "{{ instance.type }}"
    }
    {{- if not loop.last }},{{ endif }}
    {% endfor %}
  }
```

**Pros:**
- ✅ Feature-rich (filters, macros, inheritance)
- ✅ Familiar to Python developers
- ✅ Large ecosystem & documentation
- ✅ Powerful for complex transformations

**Cons:**
- ❌ External dependency (CGO)
- ❌ Slower than Go templates
- ❌ Complex filtering of unsafe expressions

**Dependencies:** `pybind11`, `jinja2` (system or embedded)

**Implementation Options:**
- Option A: Python subprocess (simplest, slower)
- Option B: pybind11 with embedded Python
- Option C: Go port like `estebanpdl/pyjinja2` (if available)

**Example (Python Subprocess):**
```go
// internal/template/engines/jinja2_engine.go
type Jinja2Engine struct{}

func (je *Jinja2Engine) Render(template string, context map[string]interface{}) (string, error) {
    contextJSON, _ := json.Marshal(context)
    
    cmd := exec.Command("python", "-c", fmt.Sprintf(`
import jinja2
context = %s
tmpl = jinja2.Template(%q)
print(tmpl.render(context))
    `, string(contextJSON), template))
    
    output, err := cmd.Output()
    return string(output), err
}
```

---

#### 3. **Liquid** (Simple & Safe)

**Engine Name:** `liquid`

```yaml
engine: "liquid"
template: |
  {
    {% for instance in instances %}
    "{{ instance.name }}": {
      "instance_type": "{{ instance.type }}"
    }
    {%- unless forloop.last %},{% endunless %}
    {% endfor %}
  }
```

**Pros:**
- ✅ Simpler syntax than Jinja2
- ✅ Designed for safety (restricted operations)
- ✅ Go library available (Shopify/liquid)
- ✅ Good for non-developer operators

**Cons:**
- ❌ Fewer features than Jinja2
- ❌ Less documentation
- ❌ Smaller ecosystem

**Dependencies:** `github.com/Shopify/liquid`

**Example:**
```go
// internal/template/engines/liquid_engine.go
import "github.com/Shopify/liquid"

type LiquidEngine struct{}

func (le *LiquidEngine) Render(template string, context map[string]interface{}) (string, error) {
    engine := liquid.NewEngine()
    tpl, err := engine.ParseTemplate([]byte(template), nil)
    if err != nil {
        return "", err
    }
    
    output, err := tpl.Render(context)
    return string(output), err
}
```

---

#### 4. **Handlebars** (JavaScript-Style)

**Engine Name:** `handlebars`

```yaml
engine: "handlebars"
template: |
  {
    {{#each instances}}
    "{{this.name}}": {
      "instance_type": "{{this.type}}"
    }
    {{#unless @last}},{{/unless}}
    {{/each}}
  }
```

**Pros:**
- ✅ Familiar to JavaScript developers
- ✅ Good balance of features & simplicity
- ✅ Go library available (aymerick/raymond)
- ✅ Growing adoption in DevOps

**Cons:**
- ❌ Less powerful than Jinja2
- ❌ Different syntax from other engines
- ❌ Smaller Go ecosystem

**Dependencies:** `github.com/aymerick/raymond`

**Example:**
```go
// internal/template/engines/handlebars_engine.go
import "github.com/aymerick/raymond"

type HandlebarsEngine struct{}

func (he *HandlebarsEngine) Render(template string, context map[string]interface{}) (string, error) {
    result, err := raymond.Render(template, context)
    return result, err
}
```

---

### Engine Factory Architecture

```go
// internal/template/engine.go
package template

// Engine defines the interface for all template engines
type Engine interface {
    Name() string
    Render(template string, context map[string]interface{}) (string, error)
}

// Registry maps engine names to implementations
type Registry struct {
    engines map[string]Engine
    mu      sync.RWMutex
}

// NewRegistry initializes the engine registry
func NewRegistry() *Registry {
    return &Registry{
        engines: map[string]Engine{
            "go":         engines.NewGoEngine(),
            "jinja2":     engines.NewJinja2Engine(),
            "liquid":     engines.NewLiquidEngine(),
            "handlebars": engines.NewHandlebarsEngine(),
        },
    }
}

// GetEngine retrieves an engine by name
func (r *Registry) GetEngine(name string) (Engine, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    if engine, ok := r.engines[name]; ok {
        return engine, nil
    }
    return nil, fmt.Errorf("unknown template engine: %s", name)
}

// Register adds a custom engine to the registry
func (r *Registry) Register(name string, engine Engine) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    
    if _, exists := r.engines[name]; exists {
        return fmt.Errorf("engine %s already registered", name)
    }
    r.engines[name] = engine
    return nil
}

// List returns all available engines
func (r *Registry) List() []string {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    engines := make([]string, 0, len(r.engines))
    for name := range r.engines {
        engines = append(engines, name)
    }
    sort.Strings(engines)
    return engines
}
```

---

### Configuration Integration

#### Extended Config Structure

```go
// internal/config/types.go
type TemplateConfig struct {
    Engine  string                 `yaml:"engine"`
    Source  string                 `yaml:"template"` // Inline or file path
    Context map[string]interface{} `yaml:"context"`
}

type Transformable struct {
    Target       string            `yaml:"target"`
    If           string            `yaml:"if"`
    Value        map[string]interface{} `yaml:"value"`
    Template     TemplateConfig    `yaml:"template"` // NEW
}

type Config struct {
    Defaults struct {
        TemplateEngine string `yaml:"template_engine"`
    } `yaml:"defaults"`
    Transform []Transformable `yaml:"transform"`
}
```

#### Loader Enhancement

```go
// internal/config/loader.go
func (c *Config) GetTemplateEngine(t *Transformable) string {
    if t.Template.Engine != "" {
        return t.Template.Engine
    }
    if c.Defaults.TemplateEngine != "" {
        return c.Defaults.TemplateEngine
    }
    return "go" // Fallback default
}

func (c *Config) ResolveTemplate(t *Transformable) (string, error) {
    source := t.Template.Source
    
    // Check if it's a file path
    if stat, err := os.Stat(source); err == nil && !stat.IsDir() {
        data, err := ioutil.ReadFile(source)
        return string(data), err
    }
    
    // Otherwise treat as inline template
    return source, nil
}
```

---

### Transformer Integration

```go
// internal/transformer/transformer.go
func (t *Transformer) Transform(
    input map[string]interface{},
    config *config.Config,
) (map[string]interface{}, error) {
    result := make(map[string]interface{})
    registry := template.NewRegistry()
    
    for _, transformable := range config.Transform {
        // Existing: interpolate, condition check, value
        
        // NEW: Template rendering
        if transformable.Template.Source != "" {
            tmplStr, err := config.ResolveTemplate(&transformable)
            if err != nil {
                return nil, fmt.Errorf("resolve template: %w", err)
            }
            
            engineName := config.GetTemplateEngine(&transformable)
            engine, err := registry.GetEngine(engineName)
            if err != nil {
                return nil, err
            }
            
            // Render template with context
            context := t.interpolateMap(transformable.Template.Context, input)
            rendered, err := engine.Render(tmplStr, context)
            if err != nil {
                return nil, fmt.Errorf("render %s template: %w", engineName, err)
            }
            
            // Parse rendered JSON/YAML and merge
            var templateResult map[string]interface{}
            if err := json.Unmarshal([]byte(rendered), &templateResult); err != nil {
                return nil, fmt.Errorf("parse template output: %w", err)
            }
            
            mergedValue := mergeDeep(transformable.Value, templateResult)
            target := t.interpolateString(transformable.Target, input)
            result = unflatten(result, target, mergedValue)
        } else {
            // Existing transformation pipeline
        }
    }
    
    return result, nil
}
```

---

### CLI Enhancements

#### New Commands

```bash
# List available engines
wtf template engines
Output:
  go         - Go Standard Templates (built-in)
  jinja2     - Jinja2 Python Templates
  liquid     - Liquid Templates (Shopify)
  handlebars - Handlebars Templates

# Validate template with specific engine
wtf template validate --engine jinja2 templates/vpc.j2
wtf template validate --engine liquid templates/storage.liquid

# Render template (useful for debugging)
wtf template render \
  --engine jinja2 \
  --input input.json \
  templates/vpc.j2

# Check template syntax across all engines
wtf template check templates/
  vpc.gotml           ✓ go
  instance.j2         ✓ jinja2
  storage.liquid      ✓ liquid
  cache.handlebars    ✓ handlebars
```

#### Flag Extensions

```bash
# Build with engine selection
wtf build \
  --config .wtf.yaml \
  --input terraform.json \
  --template-engine jinja2 \
  --output output.tf.json
```

---

## File Structure

```
internal/
├── template/
│   ├── engine.go          # Engine interface & registry
│   ├── resolver.go        # Template file resolution
│   └── engines/
│       ├── go.go          # Go template engine
│       ├── jinja2.go      # Jinja2 engine
│       ├── liquid.go      # Liquid engine
│       └── handlebars.go  # Handlebars engine
├── config/
│   └── types.go           # Extended with TemplateConfig
└── transformer/
    └── transformer.go     # Integration with templates

tests/
├── template/
│   ├── engines_test.go     # Test all engines
│   ├── registry_test.go    # Registry functionality
│   └── fixtures/
│       ├── instance.gotml
│       ├── instance.j2
│       ├── instance.liquid
│       └── instance.handlebars
└── e2e/
    └── e2e_test.go        # Include template E2E tests

docs/
├── TEMPLATE_ENGINES.md    # User guide for template engines
└── examples/
    ├── go-templates/
    ├── jinja2-templates/
    ├── liquid-templates/
    └── handlebars-templates/
```

---

## Technical Implementation

### Phase 1: Foundation (Week 1-2)
- [ ] Design Engine interface
- [ ] Implement Registry pattern
- [ ] Go template engine (built-in)
- [ ] Unit tests for engine interface
- [ ] Config loader extensions

### Phase 2: Additional Engines (Week 3-4)
- [ ] Liquid engine implementation
- [ ] Handlebars engine implementation
- [ ] Jinja2 engine (Python subprocess)
- [ ] E2E tests for all engines
- [ ] CLI commands for engine management

### Phase 3: Documentation & Polish (Week 5)
- [ ] User guide: Template Engine Comparison
- [ ] Migration guide from single to multiple engines
- [ ] Example templates for each engine
- [ ] Performance benchmarks
- [ ] Troubleshooting guide

---

## Testing Strategy

### Unit Tests

```go
// tests/template/engines_test.go
type TestCase struct {
    Name     string
    Engine   string
    Template string
    Context  map[string]interface{}
    Expected string
}

var testCases = []TestCase{
    {
        Name:     "go-loop",
        Engine:   "go",
        Template: `{{- range . -}}{{ . }},{{- end }}`,
        Context:  map[string]interface{}{"items": []string{"a", "b"}},
        Expected: `a,b,`,
    },
    {
        Name:     "jinja2-loop",
        Engine:   "jinja2",
        Template: `{%- for item in items -%}{{ item }},{% endfor %}`,
        Context:  map[string]interface{}{"items": []string{"a", "b"}},
        Expected: `a,b,`,
    },
    {
        Name:     "liquid-loop",
        Engine:   "liquid",
        Template: `{%- for item in items -%}{{ item }},{% endfor %}`,
        Context:  map[string]interface{}{"items": []string{"a", "b"}},
        Expected: `a,b,`,
    },
}

func TestAllEngines(t *testing.T) {
    for _, tc := range testCases {
        t.Run(tc.Name, func(t *testing.T) {
            engine := registry.GetEngine(tc.Engine)
            result, err := engine.Render(tc.Template, tc.Context)
            
            assert.NoError(t, err)
            assert.Equal(t, tc.Expected, result)
        })
    }
}

func TestEngineComparison(t *testing.T) {
    // Same transformation in all engines should produce identical output
    context := map[string]interface{}{
        "resources": []map[string]string{
            {"name": "api", "type": "t2.micro"},
            {"name": "web", "type": "t2.small"},
        },
    }
    
    templates := map[string]string{
        "go": `{{- range .resources }}{{ .name }},{{ end }}`,
        "jinja2": `{%- for r in resources %}{{ r.name }},{% endfor %}`,
        "liquid": `{%- for r in resources %}{{ r.name }},{% endfor %}`,
        "handlebars": `{{#each resources}}{{this.name}},{{/each}}`,
    }
    
    expected := "api,web,"
    for engineName, tmpl := range templates {
        t.Run(engineName, func(t *testing.T) {
            engine := registry.GetEngine(engineName)
            result, _ := engine.Render(tmpl, context)
            assert.Equal(t, expected, result)
        })
    }
}
```

### E2E Tests

```yaml
# tests/e2e/examples/13-template-engines/
├── .wtf.yaml
├── input.json
└── expected.tf.json

# .wtf.yaml
defaults:
  template_engine: "go"

transform:
  - target: resource.aws_instance.main
    engine: "go"
    template: |
      {
        "ami": "${input.ami}",
        "instance_type": "${input.type}"
      }
    context:
      ami: ${input.ami}
      type: ${input.type}
  
  - target: resource.aws_security_group.main
    engine: "liquid"
    template: |
      {
        {% for rule in rules %}
        "rule_{{ forloop.index }}": {
          "from_port": {{ rule.from_port }}
        }
        {%- unless forloop.last %},{% endunless %}
        {% endfor %}
      }
    context:
      rules: ${input.rules}
```

---

## User Experience

### Example: Multi-Engine Configuration

```yaml
# .wtf.yaml - Real-world example
defaults:
  template_engine: "jinja2"

transform:
  # Jinja2 for complex instance configuration
  - target: resource.aws_instance.{{ input.environment }}
    engine: jinja2
    template: templates/instance.j2
    context:
      env: ${input.environment}
      instances: ${input.instances}
      tags: ${input.tags}
  
  # Go for simple VPC setup (performance critical)
  - target: resource.aws_vpc.main
    engine: go
    template: |
      {
        "cidr_block": "{{ .vpc_cidr }}",
        "enable_dns_hostnames": true
      }
    context:
      vpc_cidr: ${input.vpc_cidr}
  
  # Liquid for operator-friendly storage config
  - target: resource.aws_s3_bucket.data
    engine: liquid
    template: templates/s3-bucket.liquid
    context:
      bucket_name: ${input.bucket_name}
      region: ${input.region}
```

### Learning Path

**Beginner:**
1. Start with Go templates (familiar to Terraform users)
2. Basic loops and conditionals

**Intermediate:**
1. Jinja2 for advanced filtering
2. Handlebars for team with JavaScript background

**Advanced:**
1. Combine multiple engines in single config
2. Template composition with includes/macros
3. Custom filters and helpers

---

## Rollout Plan

### MVP Release
- Go templates (built-in, default)
- Registry pattern for extensibility
- Config support for engine selection
- E2E tests with simple examples

### Follow-up Releases
- Phase 2: Liquid + Handlebars
- Phase 3: Jinja2 (complex integration)
- Phase 4: Template composition & inheritance
- Phase 5: Custom filters & helpers

---

## Estimated Effort

| Component | Effort | Notes |
|-----------|--------|-------|
| Engine interface & registry | 4 hours | Design & core implementation |
| Go template engine | 2 hours | Simple wrapper around stdlib |
| Config extensions | 3 hours | YAML parsing, type handling |
| Liquid engine | 3 hours | Library integration, error handling |
| Handlebars engine | 3 hours | Library integration |
| Jinja2 engine | 8 hours | Python subprocess or pybind11 |
| CLI commands | 4 hours | validate, render, list-engines |
| Unit tests | 6 hours | Comprehensive test coverage |
| E2E tests | 4 hours | Multi-engine example tests |
| Documentation | 5 hours | User guide, examples, migration |
| **Total** | **42 hours** | ~1.5 weeks for full implementation |

### Phased Approach (Recommended)

**Phase 2a (MVP Extensions):**
- Go + Liquid + Handlebars: 12 hours
- Provides most use cases
- Defers complex Jinja2 integration

**Phase 2b (Later):**
- Jinja2: 8 hours
- Optional for power users

---

## Success Criteria

✅ **Functional:**
- All engines render identically for equivalent templates
- Config supports per-template engine selection
- CLI commands work for all engines

✅ **Testing:**
- 100% test coverage for engine interface
- E2E tests for each engine
- Performance benchmarks show Go templates are 3x+ faster

✅ **Documentation:**
- User guide with syntax comparison
- 3 example configs using different engines
- Migration guide from single to multiple engines

✅ **UX:**
- Clear error messages with engine name
- Helpful validation messages for syntax errors
- List of available engines on `wtf template engines`

---

## Related Features

- **Dependent:** #02 (Jinja2 Templates) - foundation layer
- **Complements:** #01 (Plugin System) - extensible architecture
- **Enables:** #03 (Interactive CLI) - template preview/validation
- **Enhances:** #06 (Cost Estimation) - template-based cost configs

