# Plugin System for Custom Transformers

## Overview

A flexible plugin system that allows users to extend Watt TF with custom transformation logic, validators, and output formatters without modifying the core codebase.

## User Experience

### Installation

```bash
# Install a plugin from registry
wtf plugin install aws-provider@1.2.0

# Install from local path
wtf plugin install ./plugins/custom-transformer

# List installed plugins
wtf plugin list

# Show plugin details
wtf plugin info aws-provider
```

### Usage in .wtf.yaml

```yaml
plugins:
  - name: "aws-validator"
    version: "1.0.0"
    config:
      auto_tag: true
      validate: strict

transform:
  - target: resource.aws_instance.main
    value:
      instance_type: "${input.instance_type}"
    use_plugin: "aws-validator"
    plugin_config:
      required_tags:
        - Environment
        - CostCenter
```

### Plugin Development

```bash
# Create a new plugin scaffold
wtf plugin create my-plugin --type transformer

# This generates:
# - plugins/my-plugin/
#   - plugin.go (interface implementation)
#   - plugin_test.go
#   - go.mod
```

## Technical Implementation

### 1. Plugin Interface

```go
// internal/plugin/interface.go
package plugin

import "context"

type PluginKind string

const (
    KindTransformer PluginKind = "transformer"
    KindValidator   PluginKind = "validator"
    KindOutput      PluginKind = "output"
)

type Plugin interface {
    Name() string
    Version() string
    Kind() PluginKind
    Init(config map[string]interface{}) error
}

type Transformer interface {
    Plugin
    Transform(ctx context.Context, value map[string]interface{}, input map[string]interface{}) (map[string]interface{}, error)
}

type Validator interface {
    Plugin
    Validate(ctx context.Context, target string, value map[string]interface{}) error
}

type OutputFormatter interface {
    Plugin
    Format(ctx context.Context, data map[string]interface{}) ([]byte, error)
}
```

### 2. Plugin Registry

```go
// internal/plugin/registry.go
package plugin

type Registry struct {
    plugins map[string]Plugin
    loaders map[string]Loader
}

func (r *Registry) Register(plugin Plugin) error
func (r *Registry) Load(name string, version string) (Plugin, error)
func (r *Registry) List() []PluginInfo
```

### 3. Plugin Loader

```go
// internal/plugin/loader.go
package plugin

type Loader interface {
    Load(path string) (Plugin, error)
}

// Go plugin loader (compile-time)
type GoPluginLoader struct{}

// WASM plugin loader (runtime)
type WasmPluginLoader struct{}

// Script plugin loader (Lua/JS)
type ScriptPluginLoader struct{}
```

### 4. Plugin Config

```yaml
# ~/.wtf/plugins.yaml
registry:
  default: "https://plugins.watt-tf.io"
  
plugins:
  aws-provider:
    source: "https://plugins.watt-tf.io/aws-provider@1.2.0"
    kind: "transformer"
    enabled: true
    cache_dir: "~/.wtf/plugin-cache"
  
  custom-validator:
    source: "file:///home/user/plugins/custom-validator"
    kind: "validator"
    enabled: true
```

### 5. Build & Distribution

```makefile
# plugins/aws-validator/Makefile
.PHONY: build test install

build:
	go build -o dist/aws-validator ./cmd

test:
	go test ./...

install:
	mkdir -p ~/.wtf/plugins/aws-validator
	cp dist/aws-validator ~/.wtf/plugins/aws-validator/

publish:
	# Publish to plugin registry
	wtf plugin publish --registry https://plugins.watt-tf.io
```

## Implementation Steps

### Phase 1: Core Plugin System
1. Define plugin interfaces
2. Implement Go plugin loader (using Go's `plugin` package)
3. Create plugin registry & loader
4. Add `wtf plugin` CLI commands

### Phase 2: Plugin Marketplace
1. Build plugin registry server
2. Implement plugin discovery
3. Add version management & checksums

### Phase 3: Multiple Formats
1. Add WebAssembly support (security isolation)
2. Add Lua/JavaScript scripting support
3. Plugin package manager (`wtf plugin install`)

## Testing

```go
// plugins/aws-validator/plugin_test.go
func TestTransform(t *testing.T) {
    plugin := &AWSValidator{}
    plugin.Init(map[string]interface{}{})
    
    result, err := plugin.Transform(ctx, value, input)
    assert.NoError(t, err)
    assert.Equal(t, expected, result)
}
```

## Files to Create

- `internal/plugin/interface.go` - Plugin interface definitions
- `internal/plugin/registry.go` - Plugin registry
- `internal/plugin/loader.go` - Plugin loaders (Go, WASM, Script)
- `internal/cli/plugin.go` - Plugin CLI commands
- `cmd/wtf-plugin/main.go` - Plugin scaffold generator
- `plugins/example-transformer/` - Example plugin

## Estimated Effort

- **Phase 1**: 40-50 hours (core system)
- **Phase 2**: 20-30 hours (marketplace)
- **Phase 3**: 30-40 hours (multiple formats)

## Benefits

- ✅ Extensibility without core modifications
- ✅ Community ecosystem
- ✅ Custom business logic
- ✅ Provider-specific extensions
- ✅ Reusable components
