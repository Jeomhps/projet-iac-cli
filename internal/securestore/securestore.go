package securestore

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	keyring "github.com/zalando/go-keyring"
)

const serviceName = "projet-iac-cli"

// Record is the shared token record type used across the app.
type Record struct {
	AccessToken string    `json:"access_token"`
	ExpiresAt   time.Time `json:"expires_at,omitempty"`
}

type Store interface {
	Save(Record) error
	Load() (Record, error)
	Delete() error
}

// Keyring-backed store (macOS Keychain, Windows Credential Manager, Linux Secret Service)
type keyringStore struct {
	key string
}

func (s keyringStore) Save(rec Record) error {
	b, _ := json.Marshal(rec)
	return keyring.Set(serviceName, s.key, string(b))
}

func (s keyringStore) Load() (Record, error) {
	secret, err := keyring.Get(serviceName, s.key)
	if err != nil {
		return Record{}, err
	}
	var rec Record
	if err := json.Unmarshal([]byte(secret), &rec); err != nil {
		return Record{}, err
	}
	return rec, nil
}

func (s keyringStore) Delete() error {
	return keyring.Delete(serviceName, s.key)
}

// File fallback (~/.projet-iac/token.json), 0700 dir / 0600 file
type fileStore struct {
	path string
}

func (s fileStore) Save(rec Record) error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o700); err != nil {
		return err
	}
	b, _ := json.MarshalIndent(rec, "", "  ")
	return os.WriteFile(s.path, b, 0o600)
}

func (s fileStore) Load() (Record, error) {
	b, err := os.ReadFile(s.path)
	if err != nil {
		return Record{}, err
	}
	var rec Record
	if err := json.Unmarshal(b, &rec); err != nil {
		return Record{}, err
	}
	return rec, nil
}

func (s fileStore) Delete() error {
	if err := os.Remove(s.path); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return nil
}

type Mode string

const (
	ModeAuto Mode = "auto"
	ModeOn   Mode = "on"
	ModeOff  Mode = "off"
)

// New returns a Store and a boolean indicating whether the keyring backend is being used.
func New(mode Mode, keyName, filePath string) (Store, bool) {
	switch mode {
	case ModeOn:
		if ok := keyringAvailable(); ok {
			return keyringStore{key: keyName}, true
		}
		// requested "on" but unavailable: fall back to file to remain usable
		return fileStore{path: filePath}, false
	case ModeOff:
		return fileStore{path: filePath}, false
	case ModeAuto:
		if keyringAvailable() {
			return keyringStore{key: keyName}, true
		}
		return fileStore{path: filePath}, false
	default:
		return fileStore{path: filePath}, false
	}
}

// keyringAvailable probes by setting/deleting a temporary secret.
func keyringAvailable() bool {
	const probeKey = "__projet-iac-cli_probe__"
	if err := keyring.Set(serviceName, probeKey, "ok"); err != nil {
		return false
	}
	_ = keyring.Delete(serviceName, probeKey)
	return true
}

// KeyNameFor builds a stable key name per API base+prefix.
func KeyNameFor(base, prefix string) string {
	return fmt.Sprintf("api:%s%s", base, prefix)
}
