package config

import (
	"fmt"
	"path/filepath"

	"github.com/devsebastianops/watt-tf/internal/data"
	p "github.com/devsebastianops/watt-tf/internal/plugin"
	"github.com/devsebastianops/x/parser"
)

func LoadConfig(filePath string) (*Config, error) {
	config := &Config{
		Transform:  []Transformable{},
		Plugins:    []p.Plugin{},
		Variables:  map[string]any{},
		Data:       map[string]data.Data{},
		LoadedData: map[string]any{},
	}

	// Load the main config file
	configMap, err := parser.ParseInput(filePath)
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
				// Append variables from included config
				config.Variables = mergeMaps(config.Variables, includedConfig.Variables)
				config.Data = mergeDataMaps(config.Data, includedConfig.Data)
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
	// Append main config variables
	config.Variables = mergeMaps(config.Variables, mainConfig.Variables)
	config.Data = mergeDataMaps(config.Data, mainConfig.Data)

	return config, nil
}

// loadConfigWithoutIncludes loads a single config file without processing includes
func loadConfigWithoutIncludes(filePath string) (*Config, error) {
	configMap, err := parser.ParseInput(filePath)
	if err != nil {
		return nil, err
	}

	config := &Config{
		Transform:  []Transformable{},
		Plugins:    []p.Plugin{},
		Variables:  map[string]any{},
		Data:       map[string]data.Data{},
		LoadedData: map[string]any{},
	}

	// Parse variables if present
	if variables, ok := configMap["variables"].(map[string]any); ok {
		config.Variables = variables
	}

	// Parse additional data sources if present
	if dataMap, ok := configMap["data"].(map[string]any); ok {
		for dataKey, rawData := range dataMap {
			dataConfig, ok := rawData.(map[string]any)
			if !ok {
				return nil, fmt.Errorf("invalid data entry '%s'", dataKey)
			}

			dataType, _ := dataConfig["type"].(string)
			dataPath, _ := dataConfig["path"].(string)
			if dataType == "" || dataPath == "" {
				return nil, fmt.Errorf("data entry '%s' must include 'type' and 'path'", dataKey)
			}
			if !filepath.IsAbs(dataPath) {
				dataPath = filepath.Join(filepath.Dir(filePath), dataPath)
			}

			config.Data[dataKey] = data.NewData(dataType, dataPath)
		}
	}

	// Parse transforms
	transformList, ok := configMap["transform"].([]any)
	if !ok {
		// If no transform list, return empty config
		return config, nil
	}

	for _, transformable := range transformList {
		transformableMap, ok := transformable.(map[string]any)
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
		var value any
		if val, exists := transformableMap["value"]; exists {
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
	pluginList, ok := configMap["plugins"].([]any)
	if ok {
		for _, plugin := range pluginList {
			pluginMap, ok := plugin.(map[string]any)
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

			argsInterface, _ := pluginMap["args"].([]any)
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

func mergeMaps(dest, src map[string]any) map[string]any {
	if dest == nil {
		dest = make(map[string]any)
	}
	for k, v := range src {
		dest[k] = v
	}
	return dest
}

func mergeDataMaps(dest, src map[string]data.Data) map[string]data.Data {
	if dest == nil {
		dest = make(map[string]data.Data)
	}
	for k, v := range src {
		dest[k] = v
	}
	return dest
}
