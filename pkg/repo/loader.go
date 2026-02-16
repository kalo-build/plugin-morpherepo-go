package repo

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// LoadRepoSpecs loads all .repo YAML files from the specified directory.
func LoadRepoSpecs(dirPath string) ([]RepoSpec, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read repo directory %s: %w", dirPath, err)
	}

	var fileNames []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if strings.HasSuffix(entry.Name(), ".repo") {
			fileNames = append(fileNames, entry.Name())
		}
	}
	sort.Strings(fileNames)

	var specs []RepoSpec
	for _, fileName := range fileNames {
		filePath := filepath.Join(dirPath, fileName)
		data, readErr := os.ReadFile(filePath)
		if readErr != nil {
			return nil, fmt.Errorf("failed to read repo file %s: %w", filePath, readErr)
		}

		var spec RepoSpec
		if parseErr := yaml.Unmarshal(data, &spec); parseErr != nil {
			return nil, fmt.Errorf("failed to parse repo file %s: %w", filePath, parseErr)
		}

		specs = append(specs, spec)
	}

	return specs, nil
}
