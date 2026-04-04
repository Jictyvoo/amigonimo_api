package sqltest

import (
	"unicode"
)

// dbNameNormalizer converts a test name into a valid MySQL identifier, prefixed
// with "test_". Non-alphanumeric characters are replaced with underscores.
func dbNameNormalizer(name string) string {
	runes := []rune(name)
	for i, ch := range runes {
		ch = unicode.ToLower(ch)
		if !unicode.IsLetter(ch) && !unicode.IsNumber(ch) {
			ch = '_'
		}
		runes[i] = ch
	}
	return "sqlrepo_" + string(runes)
}
