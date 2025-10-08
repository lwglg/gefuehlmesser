package tooling

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var Normalizer = transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)

func StripAccentsLower(content string) (string, error) {
	result, _, err := transform.String(Normalizer, content)
	if err != nil {
		return "", err
	}

	return strings.ToLower(result), nil
}

func HexDigestFromString(text string) string {
	hash := sha256.Sum256([]byte(text))
	hexString := hex.EncodeToString(hash[:])

	return hexString
}
