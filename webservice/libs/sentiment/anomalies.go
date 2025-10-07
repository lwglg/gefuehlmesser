package sentiment

import (
	"sort"
	"time"
)

func DetectAnomalies(allMsgs *[]FeedMessage) (bool, *string) {
	countMsgs := len(*allMsgs)

	if countMsgs == 0 {
		return false, nil
	}

	// Synchronized posting tolerant: at least 3 messages and all within ±2 seconds
	if countMsgs >= 3 {
		secs := make([]time.Time, countMsgs)
		for i, m := range *allMsgs {
			dt := m.TimeWindow
			secs[i] = dt.Truncate(time.Second)
		}

		minSec := secs[0]
		maxSec := secs[0]

		for _, t := range secs[1:] {
			if t.Before(minSec) {
				minSec = t
			}
			if t.After(maxSec) {
				maxSec = t
			}
		}

		if maxSec.Sub(minSec) <= 2*time.Second {
			result := "synchronized_posting"

			return true, &result
		}
	}

	// Burst: >10 messages from same user in 5 minutes
	byUser := make(map[string][]time.Time)

	for _, m := range *allMsgs {
		userID := m.UserID
		dt := m.TimeWindow
		byUser[userID] = append(byUser[userID], dt)
	}

	for _, tsList := range byUser {
		sort.Slice(tsList, func(i, j int) bool {
			return tsList[i].Before(tsList[j])
		})

		i := 0

		for j := 0; j < len(tsList); j++ {
			for tsList[j].Sub(tsList[i]) > 5*time.Minute {
				i++
			}

			if (j - i + 1) > 10 {
				result := "burst"

				return true, &result
			}
		}
	}

	return false, nil
}
