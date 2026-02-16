package testutils_test

import (
	"path/filepath"
	"testing"

	"github.com/kalo-build/plugin-morpherepo-go/pkg/compile"
	"github.com/kalo-build/plugin-morpherepo-go/pkg/compile/cfg"
)

func TestRegenRegistry(t *testing.T) {
	t.Skip("Run manually to regenerate registry Go repo interfaces")

	repoInputPath := filepath.Join("..", "..", "..", "kalo-plugin-registry", "morphe", "repo")
	outputPath := filepath.Join("..", "..", "..", "kalo-plugin-registry", "internal", "generated", "repo")

	config := cfg.CompileConfig{
		InputDirPath:  repoInputPath,
		OutputDirPath: outputPath,
		ModelsConfig: cfg.ModelsConfig{
			PackagePath: "github.com/kalo-build/kalo-plugin-registry/internal/types/models",
		},
		RepoConfig: cfg.RepoConfig{
			PackagePath: "github.com/kalo-build/kalo-plugin-registry/internal/generated/repo",
		},
	}

	if err := compile.MorpheRepoToGo(config); err != nil {
		t.Fatal(err)
	}

	t.Log("Registry Go repo interfaces written to:", outputPath)
}
