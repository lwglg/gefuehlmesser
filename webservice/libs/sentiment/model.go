package sentiment

type AnomalyType[T, U any] struct {
	Flag T `json:"flag"`
	Type U `json:"type"`
}

type FeedMessage struct {
	ID        string   `json:"id" validate:"required,msg_id"`
	Content   string   `json:"content" validate:"required,max=255"`
	Timestamp string   `json:"timestamp" validate:"required,datetime=YYYY-MM-DDTHH:MM:SSZ"`
	UserID    string   `json:"user_id" validate:"required,user_id"`
	Hashtags  []string `json:"hashtags" validate:"required"`
	Reactions int      `json:"reactions" validate:"required,gte=0"`
	Shares    int      `json:"shares" validate:"required,gte=0"`
	Views     int      `json:"views" validate:"required,gte=0"`
}

type Feed struct {
	Messages          []FeedMessage `json:"messages" validate:"required,dive,required"`
	TimeWindowMinutes int           `json:"time_window_minutes" validate:"required,gte=0"`
}

type InfluenceRanking struct {
	UserID         string  `json:"user_id"`
	InfluenceScore float32 `json:"influece_score"`
}

type SentimentDistribution struct {
	Positive float64 `json:"positive"`
	Negative float64 `json:"negative"`
	Neutral  float64 `json:"neutral"`
}

type SentimentAnalysisFlags struct {
	MbrasEmployee      bool `json:"mbras_employee"`
	SpecialPattern     bool `json:"special_pattern"`
	CandidateAwareness bool `json:"candidate_awareness"`
}

type SentimentAnalysis struct {
	SentimentDistribution SentimentDistribution     `json:"sentiment_distribution"`
	EngagementScore       float64                   `json:"engagement_score"`
	TrendingTopics        []string                  `json:"trending_topics"`
	AnomalyDetected       bool                      `json:"anomaly_detected"`
	AnomalyType           AnomalyType[bool, string] `json:"anomaly_type"`
	Flags                 SentimentAnalysisFlags    `json:"flags"`
}
