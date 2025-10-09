package sentiment

import (
	"math"
	"sort"
	"strings"
	"time"
)

func ExtractTrendingTopics(windowMsgs *[]FeedMessage, utcNow *time.Time) []string {
	weights := make(map[string]float64)
	counts := make(map[string]int)
	sentimentWeights := make(map[string]float64)

	for _, m := range *windowMsgs {
		sentimentMultiplier := 1.0

		switch m.Sentiment.Label {
		case "positive":
			sentimentMultiplier = 1.2
		case "negative":
			sentimentMultiplier = 0.8
		}

		for _, h := range m.Hashtags {
			tag := strings.ToLower(h)

			dt := m.ParsedTimeStamp
			deltaMin := (*utcNow).Sub(dt).Minutes()

			if deltaMin < 0 {
				deltaMin = 0
			}

			baseWeight := 1.0 + (1.0 / math.Max(deltaMin, 0.01))

			if len(tag) > 8 {
				lengthFactor := math.Log10(float64(len(tag))) / math.Log10(8)
				baseWeight *= lengthFactor
			}

			weight := baseWeight * sentimentMultiplier
			weights[tag] += weight
			counts[tag]++
			sentimentWeights[tag] += sentimentMultiplier
		}
	}

	type kv struct {
		key string
		val float64
	}

	items := make([]kv, 0, len(weights))

	for k, v := range weights {
		items = append(items, kv{k, v})
	}

	sort.Slice(items, func(i, j int) bool {
		ki, kj := items[i].key, items[j].key
		vi, vj := items[i].val, items[j].val

		if vi != vj {
			return vi > vj
		}

		if counts[ki] != counts[kj] {
			return counts[ki] > counts[kj]
		}

		if sentimentWeights[ki] != sentimentWeights[kj] {
			return sentimentWeights[ki] > sentimentWeights[kj]
		}

		return ki < kj
	})

	top := min(len(items), 5)
	result := make([]string, top)

	for i := range top {
		result[i] = items[i].key
	}

	return result
}
