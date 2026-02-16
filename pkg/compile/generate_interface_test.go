package compile_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/kalo-build/plugin-morpherepo-go/pkg/compile"
	"github.com/kalo-build/plugin-morpherepo-go/pkg/compile/cfg"
	"github.com/kalo-build/plugin-morpherepo-go/pkg/repo"
)

type GenerateInterfaceTestSuite struct {
	suite.Suite

	DefaultConfig cfg.CompileConfig
}

func TestGenerateInterfaceTestSuite(t *testing.T) {
	suite.Run(t, new(GenerateInterfaceTestSuite))
}

func (suite *GenerateInterfaceTestSuite) SetupTest() {
	suite.DefaultConfig = cfg.CompileConfig{
		ModelsConfig: cfg.ModelsConfig{
			PackagePath: "github.com/test/app/internal/types/models",
		},
		RepoConfig: cfg.RepoConfig{
			PackagePath: "github.com/test/app/internal/generated/repo",
		},
	}
}

func (suite *GenerateInterfaceTestSuite) TestGenerateInterface_SimpleModel() {
	spec := repo.RepoSpec{
		Name:  "OrganizationRepository",
		Model: "Organization",
		Identifiers: map[string]repo.Identifier{
			"primary": {Fields: []repo.IdentifierField{{Name: "ID", Type: "UUID"}}},
			"code":    {Fields: []repo.IdentifierField{{Name: "Code", Type: "String"}}},
		},
		Filters: nil,
		Operations: repo.Operations{
			List: true, Get: true, Create: true, Update: true, Delete: true,
		},
	}

	result := compile.GenerateInterface(spec, suite.DefaultConfig)

	suite.Contains(result, "package repo")
	suite.Contains(result, "\"context\"")
	suite.Contains(result, "\"github.com/test/app/internal/types/models\"")
	suite.Contains(result, "type OrganizationRepository interface {")
	suite.Contains(result, "GetAll(ctx context.Context) ([]models.Organization, error)")
	suite.Contains(result, "GetByID(ctx context.Context, id string) (*models.Organization, error)")
	suite.Contains(result, "GetByCode(ctx context.Context, code string) (*models.Organization, error)")
	suite.Contains(result, "Create(ctx context.Context, input *models.Organization) (*models.Organization, error)")
	suite.Contains(result, "Update(ctx context.Context, id string, input *models.Organization) (*models.Organization, error)")
	suite.Contains(result, "Delete(ctx context.Context, id string) error")
	suite.Contains(result, "Query(ctx context.Context, example *models.Organization) ([]models.Organization, error)")
	suite.Contains(result, "QueryOne(ctx context.Context, example *models.Organization) (*models.Organization, error)")
}

func (suite *GenerateInterfaceTestSuite) TestGenerateInterface_ModelWithFilters() {
	spec := repo.RepoSpec{
		Name:  "ProjectRepository",
		Model: "Project",
		Identifiers: map[string]repo.Identifier{
			"primary": {Fields: []repo.IdentifierField{{Name: "ID", Type: "UUID"}}},
			"code":    {Fields: []repo.IdentifierField{{Name: "Code", Type: "String"}}},
		},
		Filters: []repo.Filter{
			{Name: "organizationID", Type: "UUID", Relation: "Organization"},
		},
		Operations: repo.Operations{
			List: true, Get: true, Create: true, Update: true, Delete: true,
		},
	}

	result := compile.GenerateInterface(spec, suite.DefaultConfig)

	suite.Contains(result, "GetAll(ctx context.Context, organizationID *string) ([]models.Project, error)")
}

func (suite *GenerateInterfaceTestSuite) TestGenerateInterface_PrimaryOnly() {
	spec := repo.RepoSpec{
		Name:  "TaskRepository",
		Model: "Task",
		Identifiers: map[string]repo.Identifier{
			"primary": {Fields: []repo.IdentifierField{{Name: "ID", Type: "UUID"}}},
		},
		Filters: []repo.Filter{
			{Name: "projectID", Type: "UUID", Relation: "Project"},
		},
		Operations: repo.Operations{
			List: true, Get: true, Create: true, Update: true, Delete: true,
		},
	}

	result := compile.GenerateInterface(spec, suite.DefaultConfig)

	suite.Contains(result, "GetByID(ctx context.Context, id string) (*models.Task, error)")
	suite.NotContains(result, "GetByCode")
}

func (suite *GenerateInterfaceTestSuite) TestGenerateInterface_DisabledOperations() {
	spec := repo.RepoSpec{
		Name:  "ReadOnlyRepository",
		Model: "ReadOnly",
		Identifiers: map[string]repo.Identifier{
			"primary": {Fields: []repo.IdentifierField{{Name: "ID", Type: "UUID"}}},
		},
		Operations: repo.Operations{
			List: true, Get: true, Create: false, Update: false, Delete: false,
		},
	}

	result := compile.GenerateInterface(spec, suite.DefaultConfig)

	suite.Contains(result, "GetAll")
	suite.Contains(result, "GetByID")
	suite.NotContains(result, "Create(")
	suite.NotContains(result, "Update(")
	suite.NotContains(result, "Delete(")
}
