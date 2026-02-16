package cfg

import "path"

// CompileConfig holds all configuration for the morpherepo-to-go compilation pipeline.
type CompileConfig struct {
	// InputDirPath is the directory containing .repo YAML files.
	InputDirPath string

	// OutputDirPath is the directory for generated Go interface files.
	OutputDirPath string

	// ModelsConfig defines the Go models package.
	ModelsConfig ModelsConfig

	// RepoConfig defines the Go repository package.
	RepoConfig RepoConfig
}

// ModelsConfig holds configuration for the models package.
type ModelsConfig struct {
	PackagePath string
}

// RepoConfig holds configuration for the repository package.
type RepoConfig struct {
	PackagePath string
}

// PackageNameFromPath extracts the package name from a full Go import path.
func PackageNameFromPath(packagePath string) string {
	return path.Base(packagePath)
}
