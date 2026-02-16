package compile

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kalo-build/plugin-morpherepo-go/pkg/compile/cfg"
	"github.com/kalo-build/plugin-morpherepo-go/pkg/repo"
)

// MorpheRepoToGo is the main compilation entrypoint.
// It reads .repo YAML files and generates Go interface files.
func MorpheRepoToGo(config cfg.CompileConfig) error {
	// Load all .repo specs
	specs, loadErr := repo.LoadRepoSpecs(config.InputDirPath)
	if loadErr != nil {
		return fmt.Errorf("failed to load repo specs: %w", loadErr)
	}

	if len(specs) == 0 {
		return fmt.Errorf("no .repo files found in %s", config.InputDirPath)
	}

	// Create output directory
	if err := os.MkdirAll(config.OutputDirPath, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Generate Go interface for each spec
	for _, spec := range specs {
		content := GenerateInterface(spec, config)

		fileName := toSnakeCase(spec.Model) + "_repository.go"
		filePath := filepath.Join(config.OutputDirPath, fileName)
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write interface file for %s: %w", spec.Model, err)
		}
	}

	return nil
}

// toSnakeCase converts PascalCase to snake_case.
func toSnakeCase(s string) string {
	if s == "" {
		return s
	}
	var words []string
	wordStart := 0
	runes := []rune(s)
	for i := 1; i < len(runes); i++ {
		if runes[i] >= 'A' && runes[i] <= 'Z' {
			if runes[i-1] >= 'a' && runes[i-1] <= 'z' {
				words = append(words, string(runes[wordStart:i]))
				wordStart = i
			} else if i+1 < len(runes) && runes[i+1] >= 'a' && runes[i+1] <= 'z' {
				words = append(words, string(runes[wordStart:i]))
				wordStart = i
			}
		}
	}
	words = append(words, string(runes[wordStart:]))
	for i, w := range words {
		words[i] = strings.ToLower(w)
	}
	return strings.Join(words, "_")
}
