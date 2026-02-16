package compile_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/kalo-build/go-util/assertfile"
	"github.com/kalo-build/plugin-morpherepo-go/internal/testutils"
	"github.com/kalo-build/plugin-morpherepo-go/pkg/compile"
	"github.com/kalo-build/plugin-morpherepo-go/pkg/compile/cfg"
)

type CompileTestSuite struct {
	assertfile.FileSuite

	TestDirPath            string
	TestGroundTruthDirPath string

	RepoInputDirPath string
}

func TestCompileTestSuite(t *testing.T) {
	suite.Run(t, new(CompileTestSuite))
}

func (suite *CompileTestSuite) SetupTest() {
	suite.TestDirPath = testutils.GetTestDirPath()
	suite.TestGroundTruthDirPath = filepath.Join(suite.TestDirPath, "ground-truth", "compile-minimal")
	suite.RepoInputDirPath = filepath.Join(suite.TestDirPath, "repo", "minimal")
}

func (suite *CompileTestSuite) TearDownTest() {
	suite.TestDirPath = ""
}

func (suite *CompileTestSuite) TestMorpheRepoToGo() {
	workingDirPath := filepath.Join(suite.TestDirPath, "working")
	suite.Nil(os.Mkdir(workingDirPath, 0755))
	defer os.RemoveAll(workingDirPath)

	config := cfg.CompileConfig{
		InputDirPath:  suite.RepoInputDirPath,
		OutputDirPath: workingDirPath,
		ModelsConfig: cfg.ModelsConfig{
			PackagePath: "github.com/test/app/internal/types/models",
		},
		RepoConfig: cfg.RepoConfig{
			PackagePath: "github.com/test/app/internal/generated/repo",
		},
	}

	compileErr := compile.MorpheRepoToGo(config)

	suite.NoError(compileErr)

	// Verify organization_repository.go
	orgPath := filepath.Join(workingDirPath, "organization_repository.go")
	gtOrgPath := filepath.Join(suite.TestGroundTruthDirPath, "organization_repository.go")
	suite.FileExists(orgPath)
	suite.FileEquals(orgPath, gtOrgPath)

	// Verify project_repository.go
	projectPath := filepath.Join(workingDirPath, "project_repository.go")
	gtProjectPath := filepath.Join(suite.TestGroundTruthDirPath, "project_repository.go")
	suite.FileExists(projectPath)
	suite.FileEquals(projectPath, gtProjectPath)

	// Verify task_repository.go
	taskPath := filepath.Join(workingDirPath, "task_repository.go")
	gtTaskPath := filepath.Join(suite.TestGroundTruthDirPath, "task_repository.go")
	suite.FileExists(taskPath)
	suite.FileEquals(taskPath, gtTaskPath)
}
