package plugin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/devsebastianops/watt-tf/internal/logger"
)

type DispatchConfig struct {
	Event       Event
	Registry    *Registry
	Input       map[string]interface{}
	Environment map[string]interface{}
	BasePath    string
}

func DispatchEvents(config DispatchConfig) (Context, error) {

	event := config.Event
	registry := config.Registry
	input := config.Input
	environment := config.Environment
	basePath := config.BasePath

	plugins := registry.GetPluginsForEvent(event)
	context := newContext(input, environment)

	logger.Debug("dispatching event", "event", event, "plugin_count", len(plugins))

	var response *PluginResponse
	var err error
	for _, plugin := range plugins {
		response, err = DispatchEvent(plugin, context, basePath)
		if err != nil {
			return context, fmt.Errorf("Error dispatching event '%s' to plugin '%s': %v", event, plugin.Name, err)
		}

		// Partial Update context with the result from the plugin
		for k, v := range response.Data.Result {
			context.Result[k] = v
		}

		for k, v := range response.Data.Env {
			context.Env[k] = v
		}

		for k, v := range response.Data.Input {
			context.Input[k] = v
		}
	}

	return context, nil
}

func DispatchEvent(plugin Plugin, context Context, configPath string) (*PluginResponse, error) {
	pluginCmd := plugin.Cmd

	payload := PluginRequest{
		Version: "1.1.0",
		Event:   plugin.On,
		Data:    context,
	}

	inputBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("Error building plugin request. %s", err)
	}

	args := fixPathes(plugin.Args, configPath)

	cmd := exec.Command(pluginCmd, args...)

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdin = bytes.NewReader(inputBytes)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("Error executing plugin command '%s': %s. Stderr: %s", pluginCmd, err, stderr.String())
	}

	var response PluginResponse
	if err := json.Unmarshal(stdout.Bytes(), &response); err != nil {
		return nil, fmt.Errorf("Plugin %s returned invalid JSON: %v\nRaw output: %s\n", plugin.Name, err, stdout.String())
	}

	if response.Status != "success" {
		return nil, fmt.Errorf("Plugin %s returned an error: %s", plugin.Name, response.Error)
	}

	return &response, nil
}

func fixPathes(args []string, configPath string) []string {
	basePath := filepath.Dir(configPath)

	for i, arg := range args {
		// Wenn das Argument ein relativer Pfad zu einer Datei ist:
		if !filepath.IsAbs(arg) {
			args[i] = filepath.Join(basePath, arg)
		}
	}

	return args
}
