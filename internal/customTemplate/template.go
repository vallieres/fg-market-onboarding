package customtemplate

import (
	"html/template"
)

// Unescape returns the unescaped.
func Unescape(s string) template.HTML {
	// #nosec G203
	return template.HTML(template.HTMLEscapeString(s))
}

// Inc is a template function to increment a variable
// usage: {{ inc $index }}
func Inc(i int) int {
	return i + 1
}
