package client

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Jeomhps/projet-iac-cli/internal/securestore"
)

type Config struct {
	APIBase               string
	APIPrefix             string
	VerifyTLS             bool
	TokenFile             string
	RewriteLocalhost      bool
	DockerHostGatewayName string
	KeychainMode          string // "auto" (default), "on", "off"
}

type HTTPResponse struct {
	StatusCode int
	Body       []byte
	Header     http.Header
}

type Client struct {
	cfg         Config
	client      *http.Client
	tokenStore  securestore.Store
	usingSecret bool
}

func New(cfg Config) *Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: !cfg.VerifyTLS, // dev: allow self-signed
		},
	}
	// Determine store
	mode := securestore.Mode(strings.ToLower(strings.TrimSpace(cfg.KeychainMode)))
	if mode == "" {
		mode = securestore.ModeAuto
	}
	key := securestore.KeyNameFor(cfg.APIBase, cfg.APIPrefix)
	file := cfg.TokenFile
	if file == "" {
		home, _ := os.UserHomeDir()
		file = filepath.Join(home, ".projet-iac", "token.json")
	}
	store, usingSecret := securestore.New(mode, key, file)

	return &Client{
		cfg:         cfg,
		client:      &http.Client{Transport: tr, Timeout: 60 * time.Second},
		tokenStore:  store,
		usingSecret: usingSecret,
	}
}

func (c *Client) url(path string) string {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return c.cfg.APIBase + c.cfg.APIPrefix + path
}

func (c *Client) ShouldRewrite(host string) bool {
	if !c.cfg.RewriteLocalhost {
		return false
	}
	h := strings.TrimSpace(strings.ToLower(host))
	return h == "localhost" || h == "127.0.0.1"
}

func (c *Client) do(req *http.Request) (*HTTPResponse, error) {
	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	b, _ := io.ReadAll(res.Body)
	return &HTTPResponse{StatusCode: res.StatusCode, Body: b, Header: res.Header.Clone()}, nil
}

func (c *Client) Get(path string, token string) (*HTTPResponse, error) {
	req, _ := http.NewRequest(http.MethodGet, c.url(path), nil)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	return c.do(req)
}

func (c *Client) PostJSON(path string, token string, body any) (*HTTPResponse, error) {
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, c.url(path), bytes.NewReader(b))
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	req.Header.Set("Content-Type", "application/json")
	return c.do(req)
}

func (c *Client) Delete(path string, token string) (*HTTPResponse, error) {
	req, _ := http.NewRequest(http.MethodDelete, c.url(path), nil)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	return c.do(req)
}

// Login posts username/password to /login and returns access token and expiry.
// If API doesn't return expires_in, we try to read exp from the JWT; fallback to 60m.
func (c *Client) Login(username, password string) (token string, expiresAt *time.Time, err error) {
	payload := map[string]string{"username": username, "password": password}
	res, err := c.PostJSON("/login", "", payload)
	if err != nil {
		return "", nil, err
	}
	if res.StatusCode != http.StatusOK {
		return "", nil, fmt.Errorf("login failed: %d %s", res.StatusCode, string(res.Body))
	}
	var data struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in,omitempty"`
	}
	if err := json.Unmarshal(res.Body, &data); err != nil {
		return "", nil, err
	}
	if data.AccessToken == "" {
		return "", nil, errors.New("login: empty access_token")
	}
	var exp *time.Time
	if data.ExpiresIn > 0 {
		t := time.Now().Add(time.Duration(data.ExpiresIn) * time.Second)
		exp = &t
	} else if t2, err := parseJWTExp(data.AccessToken); err == nil {
		exp = &t2
	} else {
		t := time.Now().Add(60 * time.Minute)
		exp = &t
	}
	return data.AccessToken, exp, nil
}

// parseJWTExp reads the "exp" claim from a JWT without verification.
func parseJWTExp(tok string) (time.Time, error) {
	parts := strings.Split(tok, ".")
	if len(parts) < 2 {
		return time.Time{}, errors.New("invalid JWT format")
	}
	payloadB, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return time.Time{}, err
	}
	var claims struct {
		Exp int64 `json:"exp"`
	}
	if err := json.Unmarshal(payloadB, &claims); err != nil {
		return time.Time{}, err
	}
	if claims.Exp == 0 {
		return time.Time{}, errors.New("no exp in JWT")
	}
	return time.Unix(claims.Exp, 0), nil
}

func (c *Client) SaveToken(token string, exp *time.Time) error {
	rec := securestore.Record{AccessToken: token}
	if exp != nil {
		rec.ExpiresAt = *exp
	}
	return c.tokenStore.Save(rec)
}

func (c *Client) LoadToken() (string, *time.Time, error) {
	rec, err := c.tokenStore.Load()
	if err != nil {
		return "", nil, err
	}
	if rec.AccessToken == "" {
		return "", nil, errors.New("no token in cache")
	}
	if !rec.ExpiresAt.IsZero() && time.Now().After(rec.ExpiresAt) {
		return "", &rec.ExpiresAt, errors.New("cached token expired")
	}
	return rec.AccessToken, &rec.ExpiresAt, nil
}

// GetToken returns a valid token or an error instructing to login first.
func (c *Client) GetToken() (string, error) {
	token, _, err := c.LoadToken()
	if err != nil {
		return "", errors.New("no valid token found. Please run: projet-iac-cli login")
	}
	return token, nil
}

func (c *Client) DeleteToken() error {
	return c.tokenStore.Delete()
}
