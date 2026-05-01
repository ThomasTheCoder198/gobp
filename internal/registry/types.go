package registry

// Kind classifies a manifest.
type Kind string

const (
	KindFramework Kind = "framework"
	KindLogger    Kind = "logger"
	KindDB        Kind = "db"
	KindSDK       Kind = "sdk"
	KindAddon     Kind = "addon"
	KindPattern   Kind = "pattern"
)

// Manifest is the YAML-decoded description of one feature.
//
// Each manifest is the single source of truth for what a feature contributes:
// imports, config blocks, env overrides, compose stanzas, fragments, and errors.
type Manifest struct {
	ID          string   `yaml:"id"`
	Kind        Kind     `yaml:"kind"`
	Display     string   `yaml:"display"`
	Description string   `yaml:"description,omitempty"`
	GoImports   []string `yaml:"go_imports,omitempty"`
	GoRequires  []string `yaml:"go_requires,omitempty"` // module paths needed by this feature

	Config       []ConfigBlock `yaml:"config,omitempty"`
	EnvOverrides []EnvOverride `yaml:"env_overrides,omitempty"`
	Compose      *Compose      `yaml:"compose,omitempty"`
	Fragments    []Fragment    `yaml:"fragments,omitempty"`
	Errors       []ErrorEntry  `yaml:"errors,omitempty"`
}

// ConfigBlock describes a top-level entry in the project's Config struct.
type ConfigBlock struct {
	Key    string                 `yaml:"key"`
	Type   string                 `yaml:"type"`
	Fields map[string]ConfigField `yaml:"fields"`
}

// ConfigField describes one field inside a ConfigBlock.
type ConfigField struct {
	Type    string `yaml:"type"`
	Default any    `yaml:"default,omitempty"`
	Secret  bool   `yaml:"secret,omitempty"`
	EnvVar  string `yaml:"env,omitempty"`
}

// EnvOverride maps an env-var name to a dotted config key path.
type EnvOverride struct {
	Env string `yaml:"env"`
	Key string `yaml:"key"`
}

// Compose is one service stanza in the generated docker-compose.yml.
type Compose struct {
	Service     string            `yaml:"service"`
	Image       string            `yaml:"image"`
	Ports       []string          `yaml:"ports,omitempty"`
	Environment map[string]string `yaml:"environment,omitempty"`
	Volumes     []string          `yaml:"volumes,omitempty"`
	Command     string            `yaml:"command,omitempty"`
	NamedVolume string            `yaml:"named_volume,omitempty"`
}

// Fragment is one piece of generated content contributed by a manifest.
//
// Either File (a path inside the embedded templates tree) or Inline (literal
// text) must be set. If Section is set, the fragment is appended to that named
// section of Target rather than being its own file.
type Fragment struct {
	Target  string `yaml:"target"`
	Section string `yaml:"section,omitempty"`
	File    string `yaml:"file,omitempty"`
	Inline  string `yaml:"inline,omitempty"`
}

// ErrorEntry contributes an entry to the generated errors/codes.go catalog.
type ErrorEntry struct {
	Code        int    `yaml:"code"`
	Name        string `yaml:"name"`
	HTTP        int    `yaml:"http"`
	Message     string `yaml:"message"`
	Description string `yaml:"description,omitempty"`
}
