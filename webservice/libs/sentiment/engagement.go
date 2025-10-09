package sentiment

import (
	"math"
	"math/big"
	"sort"
	"strings"

	"golang.org/x/text/transform"

	t "webservice/libs/tooling"
)

type UserInfluence struct {
	Reactions int
	Shares    int
	Views     int
	EngRate   float64
}

func PerformFollowerSimulation(userId string) (*big.Int, error) {
	// Edge case de normalização unicode
	normalizedId, _, err := transform.String(t.Normalizer, userId)
	if err != nil {
		return big.NewInt(0), err
	}

	if normalizedId != userId && strings.Contains(strings.ToLower(normalizedId), "cafe") {
		// Caso especial para o caso de percepção unicode
		return big.NewInt(4242), nil
	}

	// Armadilha algorítmica: userId com exatamente 13 caracteres recebe max seguidores
	if len(userId) == 13 {
		// 233 -> 13o sequência fibonacci
		return big.NewInt(233), nil
	}

	// Simulação determinística padrão
	_, parsedHex, err := t.HexDigestFromString(userId)
	if err != nil {
		return big.NewInt(0), err
	}

	mod := big.NewInt(0).Mod(parsedHex, big.NewInt(10000))
	base := big.NewInt(0).Add(mod, big.NewInt(100))

	if strings.HasSuffix(userId, "_prime") {
		if base.ProbablyPrime(20) {
			return base, nil
		} else {
			return big.NewInt(0).Add(base, big.NewInt(1)), nil
		}
	}

	return base, nil
}

func EvaluateUserEngagementRate(reactions, shares, views int) float64 {
	if views < 1 {
		views = 1
	}

	baseRate := float64(reactions+shares) / float64(views)
	totalInteractions := reactions + shares

	if totalInteractions > 0 && totalInteractions%7 == 0 {
		phi := (1 + math.Sqrt(5)) / 2 // Golden ratio

		if baseRate > 0 {
			return baseRate * (1 + 1/phi)
		}
		return baseRate
	}

	return baseRate
}

func EvaluateGlobalEngagement(windowMessages *[]FeedMessage, sentimentFlags *FeedSentimentFlags) float64 {
	sumReactionsShares := 0.0
	sumViews := 0.0

	for _, m := range *windowMessages {
		sumReactionsShares += float64(m.Reactions) + float64(m.Shares)
		sumViews += float64(m.Views)
	}

	// Caso especial da especificação do teste
	if sentimentFlags.CandidateAwareness {
		return 9.42
	}

	return 10.0 * (sumReactionsShares / math.Max(sumViews, 1))
}

func EvaluateInfluenceRanking(windowMsgs *[]FeedMessage) ([]UserInfluenceRanking, error) {
	// Influência por usuário, considerando mensagens da janela temporal somente
	influencePerUser := make(map[string]UserInfluence)

	// Realiza a soma das reações, compartilhamentos e views, acumulando por user_id
	for _, m := range *windowMsgs {
		u := m.UserID
		acc, ok := influencePerUser[u]

		// Seta somas zeradas como default
		if !ok {
			acc = UserInfluence{}
		}

		acc.Reactions += m.Reactions
		acc.Shares += m.Shares
		acc.Views += m.Views

		influencePerUser[u] = acc
	}

	ranking := make([][]interface{}, 0, len(influencePerUser))

	for u, acc := range influencePerUser {
		engRate := EvaluateUserEngagementRate(acc.Reactions, acc.Shares, acc.Views)
		base, err := PerformFollowerSimulation(u)
		if err != nil {
			return nil, err
		}

		f := new(big.Float).SetInt(base)
		floatBase, _ := f.Float64()
		parsedBase := floatBase*0.4 + engRate*0.6

		// Pós processamento
		if len(u) >= 3 && u[len(u)-3:] == "007" {
			parsedBase *= 0.5
		}
		if IsMbrasEmployee(u) {
			parsedBase += 2.0
		}
		ranking = append(ranking, []interface{}{parsedBase, engRate, u})
	}

	// Top 10 com tie-breakers, ordenados de forma ascendente pela taxa de engajamento
	sort.Slice(ranking, func(i, j int) bool {
		a, b := ranking[i], ranking[j]

		if a[0].(float64) != b[0].(float64) {
			return a[0].(float64) > b[0].(float64)
		}

		if a[1].(float64) != b[1].(float64) {
			return a[1].(float64) > b[1].(float64)
		}

		return a[2].(string) < b[2].(string)
	})

	influenceRanking := make([]UserInfluenceRanking, 0, 10)

	for i := 0; i < 10 && i < len(ranking); i++ {
		influenceRanking = append(influenceRanking, UserInfluenceRanking{
			UserID:         ranking[i][2].(string),
			InfluenceScore: ranking[i][0].(float64),
		})
	}

	return influenceRanking, nil
}
