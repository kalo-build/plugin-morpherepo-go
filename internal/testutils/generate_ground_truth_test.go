package testutils_test

import (
	"path/filepath"
	"testing"

	"github.com/kalo-build/plugin-morpherepo-go/internal/testutils"
	"github.com/kalo-build/plugin-morpherepo-go/pkg/compile"
	"github.com/kalo-build/plugin-morpherepo-go/pkg/compile/cfg"
)

func TestGenerateGroundTruth(t *testing.T) {
	t.Skip("Only run manually to regenerate ground truth files")

	testDirPath := testutils.GetTestDirPath()

	inputPath := filepath.Join(testDirPath, "repo", "minimal")
	outputPath := filepath.Join(testDirPath, "ground-truth", "compile-minimal")

	config := cfg.CompileConfig{
		InputDirPath:  inputPath,
		OutputDirPath: outputPath,
		ModelsConfig: cfg.ModelsConfig{
			PackagePath: "github.com/test/app/internal/types/models",
		},
		RepoConfig: cfg.RepoConfig{
			PackagePath: "github.com/test/app/internal/generated/repo",
		},
	}

	if err := compile.MorpheRepoToGo(config); err != nil {
		t.Fatal(err)
	}

	t.Log("Ground truth files regenerated at:", outputPath)
}
