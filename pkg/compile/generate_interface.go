package compile

import (
	"fmt"
	"sort"
	"strings"
	"unicode"

	"github.com/kalo-build/plugin-morpherepo-go/pkg/compile/cfg"
	"github.com/kalo-build/plugin-morpherepo-go/pkg/repo"
)

// GenerateInterface generates a Go interface file from a RepoSpec.
func GenerateInterface(spec repo.RepoSpec, config cfg.CompileConfig) string {
	var b strings.Builder

	repoPkg := cfg.PackageNameFromPath(config.RepoConfig.PackagePath)
	modelsPkg := cfg.PackageNameFromPath(config.ModelsConfig.PackagePath)

	// Package declaration
	b.WriteString(fmt.Sprintf("package %s\n\n", repoPkg))

	// Imports
	b.WriteString("import (\n")
	b.WriteString("\t\"context\"\n\n")
	b.WriteString(fmt.Sprintf("\t\"%s\"\n", config.ModelsConfig.PackagePath))
	b.WriteString(")\n\n")

	// Interface declaration
	b.WriteString(fmt.Sprintf("// %s defines the data access contract for %s.\n", spec.Name, spec.Model))
	b.WriteString(fmt.Sprintf("type %s interface {\n", spec.Name))

	// GetAll (list)
	if spec.Operations.List {
		params := buildGetAllParams(spec.Filters)
		b.WriteString(fmt.Sprintf("\tGetAll(%s) ([]%s.%s, error)\n", params, modelsPkg, spec.Model))
	}

	// Get by identifier
	if spec.Operations.Get {
		// Sort identifiers: primary first, then alphabetical
		idNames := sortedIdentifierNames(spec.Identifiers)
		for _, idName := range idNames {
			id := spec.Identifiers[idName]
			methodName := getByMethodName(idName, id)
			methodParams := getByMethodParams(id)
			b.WriteString(fmt.Sprintf("\t%s(%s) (*%s.%s, error)\n", methodName, methodParams, modelsPkg, spec.Model))
		}
	}

	// Create
	if spec.Operations.Create {
		b.WriteString(fmt.Sprintf("\tCreate(ctx context.Context, input *%s.%s) (*%s.%s, error)\n", modelsPkg, spec.Model, modelsPkg, spec.Model))
	}

	// Update
	if spec.Operations.Update {
		b.WriteString(fmt.Sprintf("\tUpdate(ctx context.Context, id string, input *%s.%s) (*%s.%s, error)\n", modelsPkg, spec.Model, modelsPkg, spec.Model))
	}

	// Delete
	if spec.Operations.Delete {
		b.WriteString("\tDelete(ctx context.Context, id string) error\n")
	}

	// Query-by-example (always generated)
	b.WriteString(fmt.Sprintf("\tQuery(ctx context.Context, example *%s.%s) ([]%s.%s, error)\n", modelsPkg, spec.Model, modelsPkg, spec.Model))
	b.WriteString(fmt.Sprintf("\tQueryOne(ctx context.Context, example *%s.%s) (*%s.%s, error)\n", modelsPkg, spec.Model, modelsPkg, spec.Model))

	b.WriteString("}\n")

	return b.String()
}

func buildGetAllParams(filters []repo.Filter) string {
	parts := []string{"ctx context.Context"}
	for _, f := range filters {
		goType := morpheTypeToGoType(f.Type)
		parts = append(parts, fmt.Sprintf("%s *%s", f.Name, goType))
	}
	return strings.Join(parts, ", ")
}

func getByMethodName(idName string, id repo.Identifier) string {
	if idName == "primary" {
		return "GetByID"
	}
	// Capitalize the identifier name
	return "GetBy" + upperFirst(idName)
}

func getByMethodParams(id repo.Identifier) string {
	parts := []string{"ctx context.Context"}
	for _, f := range id.Fields {
		goType := morpheTypeToGoType(f.Type)
		paramName := toParamName(f.Name)
		parts = append(parts, fmt.Sprintf("%s %s", paramName, goType))
	}
	return strings.Join(parts, ", ")
}

// toParamName converts a field name to a Go parameter name.
// Handles abbreviations like "ID" -> "id" properly.
func toParamName(name string) string {
	if name == "" {
		return name
	}
	// If the entire name is uppercase (e.g., "ID"), lowercase it all
	if strings.ToUpper(name) == name {
		return strings.ToLower(name)
	}
	return lowerFirst(name)
}

func morpheTypeToGoType(morpheType string) string {
	switch morpheType {
	case "UUID":
		return "string"
	case "String":
		return "string"
	case "Integer":
		return "int"
	case "Float":
		return "float64"
	case "Boolean":
		return "bool"
	case "AutoIncrement":
		return "int64"
	default:
		return "string"
	}
}

func sortedIdentifierNames(ids map[string]repo.Identifier) []string {
	names := make([]string, 0, len(ids))
	for name := range ids {
		names = append(names, name)
	}
	sort.SliceStable(names, func(i, j int) bool {
		if names[i] == "primary" {
			return true
		}
		if names[j] == "primary" {
			return false
		}
		return names[i] < names[j]
	})
	return names
}

func upperFirst(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

func lowerFirst(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}
