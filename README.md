# plugin-morpherepo-go

Generates typed Go repository interfaces from MorpheRepo data access contract definitions (`.repo` YAML files).

## What it generates

For each `.repo` contract file, this plugin produces a Go interface with methods derived
from the contract's identifiers, filters, and CRUD operations:

| Contract element     | Generated method(s)                                          |
|----------------------|--------------------------------------------------------------|
| `list: true`         | `GetAll(ctx, filters...) ([]models.T, error)`                |
| `get: true` + primary| `GetByID(ctx, id) (*models.T, error)`                        |
| Secondary identifier | `GetBy{Name}(ctx, params...) (*models.T, error)`             |
| `create: true`       | `Create(ctx, input) (*models.T, error)`                      |
| `update: true`       | `Update(ctx, id, input) (*models.T, error)`                  |
| `delete: true`       | `Delete(ctx, id) error`                                      |
| _(always)_           | `Query(ctx, example) ([]models.T, error)`                    |
| _(always)_           | `QueryOne(ctx, example) (*models.T, error)`                  |

Filter parameters from `ForOne` / `ForOnePoly` relationships become optional pointer
arguments on `GetAll`.

### Example

Given a `project.repo` with a `code` secondary identifier and an `organizationID` filter:

```go
package repo

import (
	"context"

	"github.com/myapp/internal/types/models"
)

// ProjectRepository defines the data access contract for Project.
type ProjectRepository interface {
	GetAll(ctx context.Context, organizationID *string) ([]models.Project, error)
	GetByID(ctx context.Context, id string) (*models.Project, error)
	GetByCode(ctx context.Context, code string) (*models.Project, error)
	Create(ctx context.Context, input *models.Project) (*models.Project, error)
	Update(ctx context.Context, id string, input *models.Project) (*models.Project, error)
	Delete(ctx context.Context, id string) error
	Query(ctx context.Context, example *models.Project) ([]models.Project, error)
	QueryOne(ctx context.Context, example *models.Project) (*models.Project, error)
}
```

### Type mappings

| Morphe type     | Go type   |
|-----------------|-----------|
| `UUID`          | `string`  |
| `String`        | `string`  |
| `Integer`       | `int`     |
| `Float`         | `float64` |
| `Boolean`       | `bool`    |
| `AutoIncrement` | `int64`   |

## Input / output

| Direction | Format         | Store suggestion | Description                              |
|-----------|----------------|------------------|------------------------------------------|
| Input     | `KA:MR1:YAML1` | `KA_RE_YAML`    | MorpheRepo `.repo` contract files        |
| Output    | `KA:MR1:GO1`   | `KA_RE_GO`      | Go repository interface files             |

Output files are named `{model_name}_repository.go` in snake_case
(e.g., `Project` becomes `project_repository.go`).

## Configuration

| Key                   | Type   | Required | Description                                        |
|-----------------------|--------|----------|----------------------------------------------------|
| `config.models.PackagePath` | string | yes      | Full Go import path of the models package          |
| `config.repo.PackagePath`   | string | yes      | Full Go import path of the generated repo package  |

## Pipeline context

This plugin sits downstream of
[`plugin-morphe-morpherepo`](../plugin-morphe-morpherepo), which generates the `.repo`
contract files from Morphe model definitions:

```
Morphe Models (.mod)
  └─ plugin-morphe-morpherepo ──▶ .repo YAML contracts
       └─ plugin-morpherepo-go ──▶ Go repository interfaces
```

```yaml
stores:
  KA_MO_YAML:
    format: "KA:MO1:YAML1"
    type: "localFileSystem"
    options:
      path: "./morphe"

  KA_RE_YAML:
    format: "KA:MR1:YAML1"
    type: "localFileSystem"
    options:
      path: "./morphe/repo"

  KA_RE_GO:
    format: "KA:MR1:GO1"
    type: "localFileSystem"
    options:
      path: "./internal/generated/repo"

plugins:
  "@kalo-build/plugin-morphe-morpherepo":
    version: "v1.0.0"
    inputs:
      morphe:
        format: "KA:MO1:YAML1"
        store: "KA_MO_YAML"
    output:
      format: "KA:MR1:YAML1"
      store: "KA_RE_YAML"

  "@kalo-build/plugin-morpherepo-go":
    version: "v1.0.0"
    inputs:
      repo:
        format: "KA:MR1:YAML1"
        store: "KA_RE_YAML"
    output:
      format: "KA:MR1:GO1"
      store: "KA_RE_GO"
    config:
      models:
        PackagePath: "github.com/myapp/internal/types/models"
      repo:
        PackagePath: "github.com/myapp/internal/generated/repo"

pipelines:
  generate:
    stages:
      - name: "repo-contracts"
        steps:
          - "plugin: @kalo-build/plugin-morphe-morpherepo"
      - name: "repo-go"
        steps:
          - "plugin: @kalo-build/plugin-morpherepo-go"
```

## Project structure

```
plugin-morpherepo-go/
├── cmd/plugin/             # WASM entry point
├── pkg/
│   ├── compile/            # Compilation pipeline and interface generation
│   │   ├── compile.go      # MorpheRepoToGo entry point
│   │   ├── generate_interface.go  # Go interface code generation
│   │   └── cfg/            # CompileConfig definition
│   └── repo/               # .repo YAML spec types and loader
├── internal/testutils/     # Test helpers and ground-truth regeneration
├── testdata/
│   ├── repo/minimal/       # Sample .repo input files
│   └── ground-truth/       # Expected Go output for integration tests
└── plugin.yaml             # Kalo plugin manifest
```

## Building

```bash
# Native binary
go build ./cmd/plugin

# WASM (for Kalo CLI)
GOOS=wasip1 GOARCH=wasm go build -o dist/plugin.wasm cmd/plugin/main.go
```

## Testing

```bash
go test ./...
```
