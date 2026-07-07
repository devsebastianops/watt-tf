package config

import (
	"fmt"
	"path/filepath"

	"github.com/devsebastianops/watt-tf/internal/parser"
	p "github.com/devsebastianops/watt-tf/internal/plugin"
)

func LoadConfig(filePath string) (*Config, error) {
	config := &Config{
		Transform: []Transformable{},
	}

	// Load the main config file
	configMap, err := parser.ParseYAML(filePath)
	if err != nil {
		return nil, err
	}

	// Get the directory of the config file for resolving relative paths
	configDir := filepath.Dir(filePath)

	// Parse includes if present
	if includes, ok := configMap["include"].([]interface{}); ok {
		for _, include := range includes {
			if includePath, ok := include.(string); ok {
				// Resolve relative paths from config directory
				if !filepath.IsAbs(includePath) {
					includePath = filepath.Join(configDir, includePath)
				}

				// Recursively load included config
				includedConfig, err := loadConfigWithoutIncludes(includePath)
				if err != nil {
					return nil, fmt.Errorf("failed to load included config '%s': %w", include, err)
				}

				// Append transforms from included config
				config.Transform = append(config.Transform, includedConfig.Transform...)
				// Append plugins from included config
				config.Plugins = append(config.Plugins, includedConfig.Plugins...)
			}
		}
	}

	// Parse main config transforms
	mainConfig, err := loadConfigWithoutIncludes(filePath)
	if err != nil {
		return nil, err
	}

	// Append main config transforms
	config.Transform = append(config.Transform, mainConfig.Transform...)
	// Append main config plugins
	config.Plugins = append(config.Plugins, mainConfig.Plugins...)

	return config, nil
}

// loadConfigWithoutIncludes loads a single config file without processing includes
func loadConfigWithoutIncludes(filePath string) (*Config, error) {
	configMap, err := parser.ParseYAML(filePath)
	if err != nil {
		return nil, err
	}

	config := &Config{
		Transform: []Transformable{},
		Plugins:   []p.Plugin{},
	}

	// Parse transforms
	transformList, ok := configMap["transform"].([]interface{})
	if !ok {
		// If no transform list, return empty config
		return config, nil
	}

	for _, transformable := range transformList {
		transformableMap, ok := transformable.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid transform entry")
		}

		target, ok := transformableMap["target"].(string)
		if !ok {
			return nil, fmt.Errorf("missing or invalid 'target' field")
		}

		// Parse condition if present
		condition, _ := transformableMap["if"].(string)

		// Parse value (if template is not used)
		var value map[string]interface{}
		if val, ok := transformableMap["value"].(map[string]interface{}); ok {
			value = val
		}

		// Parse for_each (for iteration)
		forEach, _ := transformableMap["for_each"].(string)

		// Validate: either value or template must be present
		if value == nil {
			return nil, fmt.Errorf("transform at target '%s' must have 'value' defined", target)
		}

		config.Transform = append(config.Transform, Transformable{
			Target:  target,
			If:      condition,
			Value:   value,
			ForEach: forEach,
		})
	}

	// Parse plugins
	pluginList, ok := configMap["plugins"].([]interface{})
	if ok {
		for _, plugin := range pluginList {
			pluginMap, ok := plugin.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("invalid plugin entry")
			}

			name, ok := pluginMap["name"].(string)
			if !ok {
				return nil, fmt.Errorf("missing or invalid 'name' field in plugin")
			}

			cmd, ok := pluginMap["cmd"].(string)
			if !ok {
				return nil, fmt.Errorf("missing or invalid 'cmd' field in plugin '%s'", name)
			}

			on, ok := pluginMap["on"].(string)
			if !ok {
				return nil, fmt.Errorf("missing or invalid 'on' field in plugin '%s'", name)
			}

			version, ok := pluginMap["version"].(string)
			if !ok {
				return nil, fmt.Errorf("missing or invalid 'version' field in plugin '%s'", name)
			}

			argsInterface, _ := pluginMap["args"].([]interface{})
			args := []string{}
			for _, arg := range argsInterface {
				if argStr, ok := arg.(string); ok {
					args = append(args, argStr)
				} else {
					return nil, fmt.Errorf("invalid argument in 'args' for plugin '%s'", name)
				}
			}

			config.Plugins = append(config.Plugins, p.Plugin{
				Name:    name,
				Version: version,
				On:      on,
				Cmd:     cmd,
				Args:    args,
			})
		}
	}

	return config, nil
}
