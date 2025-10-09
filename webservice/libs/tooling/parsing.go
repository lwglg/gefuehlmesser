package tooling

import (
	"crypto/sha256"
	"encoding/hex"
	"math/big"
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

func HexDigestFromString(text string) (*string, *big.Int, error) {
	hasher := sha256.New()
	_, err := hasher.Write([]byte(text))
	if err != nil {
		return nil, nil, err
	}

	sum := hasher.Sum(nil)
	hashString := hex.EncodeToString(sum)
	hashInt := new(big.Int).SetBytes(sum)

	return &hashString, hashInt, nil
}
