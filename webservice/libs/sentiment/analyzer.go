package sentiment

type SentimentAnalysisMethods interface {
	AnalyzeFeed(feed *Feed) *SentimentAnalysis
}

type SentimentAnalyzer struct{}

func New() *SentimentAnalyzer {
	return &SentimentAnalyzer{}
}

func (analyzer *SentimentAnalyzer) AnalyzeFeed(feed *Feed) *SentimentAnalysis {
	sentimentDistribution := SentimentDistribution{
		Positive: 2.1,
		Negative: 2.2,
		Neutral:  2.3,
	}

	analysisFlags := SentimentAnalysisFlags{
		MbrasEmployee:      false,
		SpecialPattern:     false,
		CandidateAwareness: true,
	}

	analysis := SentimentAnalysis{
		SentimentDistribution: sentimentDistribution,
		EngagementScore:       23.2,
		TrendingTopics:        []string{"foo", "bar", "foobar"},
		AnomalyDetected:       true,
		AnomalyType: AnomalyType[bool, string]{
			Flag: true,
			Type: "a-anomaly",
		},
		Flags: analysisFlags,
	}

	return &analysis
}
