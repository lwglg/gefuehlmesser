package sentiment

import (
	"regexp"
	"strings"
	"time"

	t "webservice/libs/tooling"
)

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

func IsMsgWithinWindow(m FeedMessage, utcNow time.Time, minutes int) bool {
	dt := m.TimeWindow

	return dt.After(utcNow.Add(-time.Duration(minutes)*time.Minute)) && (dt.Equal(utcNow) || dt.Before(utcNow))
}

func CheckCandidateAwareneness(content string, referenceStr ...string) bool {
	reference := "teste tecnico mbras"

	if len(referenceStr) > 0 {
		reference = referenceStr[0]
	}

	strippedNorm, _ := t.StripAccentsLower(content)
	strippedRef, _ := t.StripAccentsLower(reference)

	return strippedNorm == strippedRef
}

func TokenizeContent(content string) []ContentToken {
	tokens := [](ContentToken){}

	re := regexp.MustCompile(TokenizationRegexString)
	matches := re.FindAllString(content, -1)

	for _, v := range matches {
		stripped, err := t.StripAccentsLower(v)
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
