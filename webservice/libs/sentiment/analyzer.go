package sentiment

import (
	"fmt"
	"strings"
	"time"
)

type SentimentAnalyzerMethods interface {
	BuildFeedSentimentFlags(validMessages *[]FeedMessage, utcNow *time.Time) FeedSentimentFlags
	BuildFeedSentimentDistribution(validMessages *[]FeedMessage, windowMessages *[]FeedMessage, utcNow *time.Time) (*FeedSentimentDistribution, error)
	AnalyzeMessage(userId string) (*MessageSentiment, error)
	AnalyzeFeed(feed *Feed) (*FeedSentiment, error)
}

type SentimentAnalyzer struct{}

func New() *SentimentAnalyzer {
	return &SentimentAnalyzer{}
}

func (analyzer *SentimentAnalyzer) AnalyzeMessage(feedMessage FeedMessage) (*MessageSentiment, error) {
	if CheckCandidateAwareneness(feedMessage.Content) {
		return &MessageSentiment{0.0, "meta"}, nil
	}

	tokens := TokenizeContent(feedMessage.Content)
	totalWords := max(len(tokens), 1)

	label := ""
	postivesSum := 0.0
	negationsSum := 0.0
	nextMultiplier := 10.0
	negationScopes := []int{}

	for _, item := range tokens {
		if len(Filter(Intensifiers, func(v string) bool { return v == item.Normalized })) > 0 {
			nextMultiplier = 1.5
			negationScopes = DecreaseNegScopes(negationScopes)

			continue
		}

		if len(Filter(Negations, func(v string) bool { return v == item.Normalized })) > 0 {
			negationScopes = append(negationScopes, 3)

			continue
		}

		polarity := 0

		if len(Filter(PositiveWords, func(v string) bool { return v == item.Normalized })) > 0 {
			polarity = 1
		} else {
			polarity = -1
		}

		if polarity != 0 {
			value := 1.8 * nextMultiplier
			nextMultiplier = 1.0

			// Aplica paridade de negações acumuladas. Consome todas as negações ativas
			if len(negationScopes)%2 == 1 {
				polarity *= -1
			}

			negationScopes = nil

			// MBRAS - positivos em dobro, após intensificação/negação
			if IsMbrasEmployee(feedMessage.UserID) && polarity > 0 {
				value *= 2.0
			}

			if polarity > 0 {
				postivesSum += value
			} else {
				negationsSum += value
			}
		} else if negationScopes != nil {
			negationScopes = DecreaseNegScopes(negationScopes)
		}
	}

	score := (postivesSum - negationsSum) / float64(totalWords)

	if score > 0.1 {
		label = "positive"
	} else if score < 0.1 {
		label = "negative"
	} else {
		label = "neutral"
	}

	return &MessageSentiment{score, label}, nil
}

func (analyzer *SentimentAnalyzer) BuildFeedSentimentFlags(validMessages *[]FeedMessage) FeedSentimentFlags {
	mBrasEmployee := len(Filter(*validMessages, func(v FeedMessage) bool { return IsMbrasEmployee(v.UserID) })) > 0

	specialPattern := len(Filter(*validMessages, func(v FeedMessage) bool {
		return len(v.Content) == 42 && strings.Contains(strings.ToLower(v.Content), "mbras")
	})) > 0

	candidateAwareness := len(Filter(*validMessages, func(v FeedMessage) bool { return CheckCandidateAwareneness(v.Content) })) > 0

	return FeedSentimentFlags{
		CandidateAwareness: candidateAwareness,
		MbrasEmployee:      mBrasEmployee,
		SpecialPattern:     specialPattern,
	}
}

func (analyzer *SentimentAnalyzer) BuildFeedSentimentDistribution(validMessages *[]FeedMessage, windowMessages *[]FeedMessage) (*FeedSentimentDistribution, error) {
	counts := make(map[string]float64)
	includeForDist := 0

	for _, m := range *validMessages {
		sentiment, err := analyzer.AnalyzeMessage(m)
		if err != nil {
			return nil, err
		}

		m.Sentiment = *sentiment

		if m.Sentiment.Label == "meta" {
			// Apenas mensagens dentro da janela temporal contam para a distribuição
			if len(Filter(*windowMessages, func(v FeedMessage) bool { return m.UserID == v.UserID && m.ID == v.ID })) > 0 {
				counts[m.Sentiment.Label] += 1
				includeForDist += 1
			}
		}
	}

	distribution := FeedSentimentDistribution{}

	// Distribuição de sentimentos em valorees percentuais
	if includeForDist == 0 {
		distribution = FeedSentimentDistribution{
			Positive: 0.0,
			Negative: 0.0,
			Neutral:  0.0,
		}
	} else {
		distribution = FeedSentimentDistribution{
			Positive: 100.0 * counts["positive"] / float64(includeForDist),
			Negative: 100.0 * counts["negative"] / float64(includeForDist),
			Neutral:  100.0 * counts["neutral"] / float64(includeForDist),
		}
	}

	return &distribution, nil
}

func (analyzer *SentimentAnalyzer) AnalyzeFeed(feed *Feed) (*FeedSentiment, error) {
	utcNow := time.Now().UTC()

	fmt.Println(feed.TimeWindowMinutes)
	fmt.Println(utcNow.String())

	validMessages, windowMessages, err := ExtractMessages(&feed.Messages, &utcNow, &feed.TimeWindowMinutes)
	if err != nil {
		fmt.Println("Erro ao extrair mensagens válidas e de janela temporal!")
		return nil, err
	}

	fmt.Printf("LEN windowMessages %d \n", len(windowMessages))
	fmt.Printf("LEN validMessages %d \n", len(validMessages))

	sentimentDistribution, err := analyzer.BuildFeedSentimentDistribution(&validMessages, &windowMessages)
	if err != nil {
		fmt.Println("Erro ao gerar distribuição de sentimentos!")
		return nil, err
	}

	sentimentFlags := analyzer.BuildFeedSentimentFlags(&validMessages)
	trendingTopics := ExtractTrendingTopics(&windowMessages, &utcNow)
	engagementScore := EvaluateGlobalEngagement(&windowMessages, &sentimentFlags)
	influenceRanking, err := EvaluateInfluenceRanking(&windowMessages)
	if err != nil {
		fmt.Println("Erro ao calcular o ranking de influência por usuário!")
		return nil, err
	}

	anomalyFlag, anomalyType := DetectAnomalies(&feed.Messages)
	anomaly := AnomalyType[bool, string]{
		Flag: anomalyFlag,
		Type: *anomalyType,
	}

	analysis := FeedSentiment{
		SentimentDistribution: *sentimentDistribution,
		EngagementScore:       engagementScore,
		TrendingTopics:        trendingTopics,
		InfluenceRanking:      influenceRanking,
		AnomalyDetected:       anomalyFlag,
		AnomalyType:           anomaly,
		Flags:                 sentimentFlags,
	}

	return &analysis, nil
}
