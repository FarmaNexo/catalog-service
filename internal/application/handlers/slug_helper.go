// internal/application/handlers/slug_helper.go
package handlers

import (
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// GenerateSlug genera un slug URL-friendly desde un nombre
// "Paracetamol 500mg Tabletas" → "paracetamol-500mg-tabletas"
func GenerateSlug(name string) string {
	// Normalize unicode and remove accents
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	result, _, _ := transform.String(t, name)

	// Lowercase
	result = strings.ToLower(result)

	// Replace non-alphanumeric with hyphens
	reg := regexp.MustCompile(`[^a-z0-9]+`)
	result = reg.ReplaceAllString(result, "-")

	// Trim leading/trailing hyphens
	result = strings.Trim(result, "-")

	return result
}
