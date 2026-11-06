// Package recipes extracts the exact self-contained sources served by the demo catalog.
package recipes

import (
	"fmt"
	"regexp"
	"strings"
)

var identifier = regexp.MustCompile(`^[A-Z][A-Za-z0-9_]*$`)
var markers = regexp.MustCompile(`(?m)^// catalog:([A-Z][A-Za-z0-9_]*)$`)

func Names(source string) []string {
	var names []string
	for _, match := range markers.FindAllStringSubmatch(source, -1) {
		names = append(names, match[1])
	}
	return names
}

func Source(source, name string) (string, error) {
	if !identifier.MatchString(name) {
		return "", fmt.Errorf("invalid recipe name %q", name)
	}
	marker := "// catalog:" + name + "\n"
	if strings.Count(source, marker) != 1 {
		return "", fmt.Errorf("missing or duplicate recipe %s", name)
	}
	body := source[strings.Index(source, marker)+len(marker):]
	end := strings.Index(body, "// catalog:end")
	if end < 0 || strings.Contains(body[:end], "// catalog:") {
		return "", fmt.Errorf("unterminated recipe %s", name)
	}
	return "package examples\n\nimport \"github.com/Tend-Stack/TendUI/v2026/components\"\n\n" + strings.TrimSpace(body[:end]) + "\n", nil
}
