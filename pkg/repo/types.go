package repo

// RepoSpec represents a parsed .repo YAML file.
type RepoSpec struct {
	Name        string                `yaml:"name"`
	Model       string                `yaml:"model"`
	Identifiers map[string]Identifier `yaml:"identifiers"`
	Filters     []Filter              `yaml:"filters"`
	Operations  Operations            `yaml:"operations"`
}

// Identifier defines a lookup method with one or more fields.
type Identifier struct {
	Fields []IdentifierField `yaml:"fields"`
}

// IdentifierField describes a single field within an identifier.
type IdentifierField struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
}

// Filter defines a filter parameter for list operations.
type Filter struct {
	Name     string `yaml:"name"`
	Type     string `yaml:"type"`
	Relation string `yaml:"relation"`
}

// Operations defines which CRUD methods are enabled.
type Operations struct {
	List   bool `yaml:"list"`
	Get    bool `yaml:"get"`
	Create bool `yaml:"create"`
	Update bool `yaml:"update"`
	Delete bool `yaml:"delete"`
}
