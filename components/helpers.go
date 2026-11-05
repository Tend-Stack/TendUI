package components

import "github.com/a-h/templ"

func fallback(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}

func htmxAttrs(get, post, target, swap, indicator string) templ.Attributes {
	attrs := templ.Attributes{}
	for key, value := range map[string]string{
		"hx-get": get, "hx-post": post, "hx-target": target,
		"hx-swap": swap, "hx-indicator": indicator,
	} {
		if value != "" {
			attrs[key] = value
		}
	}
	return attrs
}
