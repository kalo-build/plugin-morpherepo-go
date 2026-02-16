package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kalo-build/plugin-morpherepo-go/pkg/compile"
	"github.com/kalo-build/plugin-morpherepo-go/pkg/compile/cfg"
)

type CompileConfigEntry struct {
	PackagePath string `json:"PackagePath"`
}

type PluginConfig struct {
	InputPath  string             `json:"inputPath"`
	OutputPath string             `json:"outputPath"`
	Config     PluginConfigFields `json:"config"`
	Verbose    bool               `json:"verbose,omitempty"`
}

type PluginConfigFields struct {
	Models CompileConfigEntry `json:"models"`
	Repo   CompileConfigEntry `json:"repo"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: plugin-morpherepo-go <config>")
		os.Exit(3)
	}

	var pluginConfig PluginConfig
	if err := json.Unmarshal([]byte(os.Args[1]), &pluginConfig); err != nil {
		fmt.Fprintln(os.Stderr, "Error parsing config JSON:", err)
		os.Exit(4)
	}

	if pluginConfig.InputPath == "" {
		fmt.Fprintln(os.Stderr, "Error: Input path is required")
		os.Exit(12)
	}
	if pluginConfig.OutputPath == "" {
		fmt.Fprintln(os.Stderr, "Error: Output path is required")
		os.Exit(13)
	}
	if pluginConfig.Config.Models.PackagePath == "" {
		fmt.Fprintln(os.Stderr, "Error: Models package path is required (config.models.PackagePath)")
		os.Exit(14)
	}
	if pluginConfig.Config.Repo.PackagePath == "" {
		fmt.Fprintln(os.Stderr, "Error: Repo package path is required (config.repo.PackagePath)")
		os.Exit(14)
	}

	inputAbs, err := filepath.Abs(pluginConfig.InputPath)
	if err == nil {
		pluginConfig.InputPath = inputAbs
	}
	outputAbs, err := filepath.Abs(pluginConfig.OutputPath)
	if err == nil {
		pluginConfig.OutputPath = outputAbs
	}

	compileConfig := cfg.CompileConfig{
		InputDirPath:  pluginConfig.InputPath,
		OutputDirPath: pluginConfig.OutputPath,
		ModelsConfig: cfg.ModelsConfig{
			PackagePath: pluginConfig.Config.Models.PackagePath,
		},
		RepoConfig: cfg.RepoConfig{
			PackagePath: pluginConfig.Config.Repo.PackagePath,
		},
	}

	if compileErr := compile.MorpheRepoToGo(compileConfig); compileErr != nil {
		fmt.Fprintln(os.Stderr, "Compilation failed:", compileErr)
		os.Exit(1)
	}

	os.Exit(0)
}
