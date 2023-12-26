package blueprint

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"sync"
)

// Generator defines the interface for a blueprint.
type Generator interface {
	Generate(spec Spec, values map[string]string, workdir string) error
	LoadSpec() (Spec, error)
}

// Spec BlueprintSpec represents the structure of a YAML blueprint spec.
type Spec struct {
	ID          string   `yaml:"id"`
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Options     []Option `yaml:"options"`
	Version     string   `yaml:"version"`
}

func (b Spec) String() string {
	return fmt.Sprintf("Name: %s, Description: %s", b.Name, b.Description)
}

type Option struct {
	ID          string   `yaml:"id"`
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Type        string   `yaml:"type"`
	Default     string   `yaml:"default,omitempty"`
	Choices     []string `yaml:"choices,omitempty"`
}

func (b *Option) String() string {
	return fmt.Sprintf("Name: %s, Description: %s, Type: %s, Choices: %v", b.Name, b.Description, b.Type, b.Choices)
}

// blueprintSpecs holds the specs for each registered blueprint.
var blueprints sync.Map

type Type string

const (
	TypeBuiltIn  Type = "built-in"
	TypePlugin   Type = "plugin"
	TypeRgistrar Type = "registrar"
)

type Blueprint struct {
	Type      Type `yaml:"type" json:"type"`
	Spec      Spec `yaml:"spec" json:"spec"`
	Generator Generator
}

// Register registers the spec for a given blueprint type.
func Register(blueprint Blueprint) error {
	if blueprint.Spec.ID == "" {
		return fmt.Errorf("failed to register blueprint spec: missing id")
	}
	blueprints.Store(blueprint.Spec.ID, blueprint)

	return nil
}

// GetGenerator retrieves a generator from the registry.
func GetGenerator(id string) (Generator, bool) {
	load, ok := blueprints.Load(id)
	if !ok {
		return nil, false
	}

	pb, ok := load.(Blueprint)
	if !ok {
		return nil, false
	}

	generator := pb.Generator
	return generator, true
}

// GetBlueprintSpec retrieves the spec for a given blueprint type.
func GetBlueprintSpec(id string) (*Spec, error) {
	pb, ok := blueprints.Load(id)
	if !ok {
		return nil, fmt.Errorf("no spec found for blueprint id: %s", id)
	}

	bp, ok := pb.(Blueprint)
	if !ok {
		return nil, fmt.Errorf("invalid spec type for blueprint id: %s", id)
	}

	blueprintSpec := &bp.Spec

	return blueprintSpec, nil
}

// LoadBlueprintSpec reads and parses the YAML spec file for a blueprint.
func LoadBlueprintSpec(filePath string) (Spec, error) {
	// sanitize paths
	filePath = filepath.Join(filepath.Clean(filePath))
	data, err := os.ReadFile(filePath)
	if err != nil {
		return Spec{}, err
	}

	var spec Spec
	err = yaml.Unmarshal(data, spec)
	return spec, err
}

// LoadBlueprintSpecFromBytes parses the YAML spec file for a blueprint from a byte array.
func LoadBlueprintSpecFromBytes(data []byte) (Spec, error) {
	var spec Spec
	err := yaml.Unmarshal(data, &spec)
	return spec, err
}

// GetAllBlueprintSpecs returns a map of all registered blueprint specs.
func GetAllBlueprintSpecs() map[string]Spec {
	specs := make(map[string]Spec)
	blueprints.Range(func(key, value interface{}) bool {
		specs[key.(string)] = value.(Spec)
		return true
	})
	return specs
}

// GetAllBlueprints returns a map of all registered blueprints.
func GetAllBlueprints() []Blueprint {
	var bps []Blueprint

	blueprints.Range(func(key, value interface{}) bool {
		bps = append(bps, value.(Blueprint))
		return true
	})

	return bps
}
