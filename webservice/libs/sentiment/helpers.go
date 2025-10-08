package sentiment

import (
	"crypto/sha256"
	"encoding/hex"
	"math"
	"regexp"
	"strings"
	"time"
	"unicode"

	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"

	"golang.org/x/text/runes"
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

func IsPrime(v int64) bool {
	if v < 2 || v%2 == 0 {
		return false
	}

	if v == 2 {
		return true
	}

	for i := 3; i <= int(math.Sqrt(float64(v)))+1; i += 2 {
		if v%int64(i) == 0 {
			return false
		}
	}

	return true
}

func FilterFuture(messages []FeedMessage, nowUtc time.Time) ([]FeedMessage, error) {
	result := []FeedMessage{}

	for _, message := range messages {
		ts, err := time.Parse(time.RFC3339, message.Timestamp)
		if err != nil {
			return nil, err
		}

		if ts.Before(nowUtc.Add(5*time.Second)) || ts.Equal(nowUtc.Add(5*time.Second)) {
			message.TimeWindow = ts
			result = append(result, message)
		}
	}

	return result, nil
}

func IsMsgWithinWindow(m FeedMessage, anchor time.Time, minutes int) bool {
	dt := m.TimeWindow

	return dt.After(anchor.Add(-time.Duration(minutes)*time.Minute)) && (dt.Equal(anchor) || dt.Before(anchor))
}

func Filter[T any](ss []T, test func(T) bool) (ret []T) {
	for _, s := range ss {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return ret
}

func CheckCandidateAwareneness(content string, referenceStr ...string) bool {
	reference := "teste técnico mbras"

	if len(referenceStr) > 0 {
		reference = referenceStr[0]
	}

	strippedNorm, _ := StripAccentsLower(content)
	strippedRef, _ := StripAccentsLower(reference)

	return strippedNorm == strippedRef
}

func TokenizeContent(content string) []ContentToken {
	tokens := [](ContentToken){}

	re := regexp.MustCompile(TokenizationRegexString)
	matches := re.FindAllString(content, -1)

	for _, v := range matches {
		stripped, err := StripAccentsLower(v)
		if err != nil {
			continue
		}

		tokens = append(tokens, ContentToken{Original: v, Normalized: stripped})
	}

	return tokens
}

func DecreaseNegScopes(negScopes []int) []int {
	decreased := []int{}

	for _, n := range negScopes {
		if n-1 > 0 {
			decreased = append(decreased, n-1)
		}
	}

	return decreased
}

func ExtractMessages(messages *[]FeedMessage, utcNow *time.Time, timeWindowMinutes *int) ([]FeedMessage, []FeedMessage, error) {
	validMessages, err := FilterFuture(*messages, *utcNow)
	if err != nil {
		return nil, nil, err
	}

	windowMessages := []FeedMessage{}

	for _, m := range validMessages {
		if IsMsgWithinWindow(m, *utcNow, *timeWindowMinutes) {
			windowMessages = append(windowMessages, m)
		}
	}

	return validMessages, windowMessages, nil
}

func IsMbrasEmployee(userId string) bool {
	return strings.Contains(strings.ToLower(userId), "mbras")
}
