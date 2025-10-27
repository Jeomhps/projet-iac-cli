package configloader

import (
	"errors"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// FileConfig uses pointer fields to detect presence in YAML.
type FileConfig struct {
	APIBase               *string `yaml:"api_base"`
	APIPrefix             *string `yaml:"api_prefix"`
	VerifyTLS             *bool   `yaml:"verify_tls"`
	TokenFile             *string `yaml:"token_file"`
	RewriteLocalhost      *bool   `yaml:"rewrite_localhost"`
	DockerHostGatewayName *string `yaml:"docker_host_gateway_name"`
	KeychainMode          *string `yaml:"keychain"` // "auto" | "on" | "off"
}

// DefaultPath returns ~/.projet-iac/config.yaml
func DefaultPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".projet-iac", "config.yaml")
}

// LoadFile loads YAML config if the file exists.
// Returns (cfg, exists, error).
func LoadFile(path string) (FileConfig, bool, error) {
	if path == "" {
		path = DefaultPath()
	}
	b, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return FileConfig{}, false, nil
		}
		return FileConfig{}, false, err
	}
	var fc FileConfig
	if err := yaml.Unmarshal(b, &fc); err != nil {
		return FileConfig{}, true, err
	}
	return fc, true, nil
}
