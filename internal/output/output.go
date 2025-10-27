package output

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/TylerBrock/colorjson"
	"golang.org/x/term"
)

func shouldColor(mode string) bool {
	mode = strings.ToLower(strings.TrimSpace(mode))
	// Respect NO_COLOR env (https://no-color.org/)
	if _, ok := os.LookupEnv("NO_COLOR"); ok && mode != "always" {
		return false
	}
	switch mode {
	case "always":
		return true
	case "never":
		return false
	default: // auto
		return term.IsTerminal(int(os.Stdout.Fd()))
	}
}

// FormatJSON pretty-prints and optionally colorizes JSON.
// If body is not valid JSON, it returns the input as-is.
func FormatJSON(body []byte, colorMode string) string {
	// Try to parse JSON once
	var obj any
	if err := json.Unmarshal(body, &obj); err != nil {
		// Not JSON: print raw
		return string(body)
	}

	// Pretty JSON without colors
	pretty, _ := json.MarshalIndent(obj, "", "  ")

	if !shouldColor(colorMode) {
		return string(pretty)
	}

	// Colorized JSON
	f := colorjson.NewFormatter()
	f.Indent = 2
	colored, err := f.Marshal(obj)
	if err != nil {
		// Fallback to non-colored pretty
		return string(pretty)
	}
	return string(colored)
}
