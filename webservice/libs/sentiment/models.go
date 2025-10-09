package sentiment

import "time"

type FeedMessage struct {
	ID              string   `json:"id" validate:"required,id_field"`
	Content         string   `json:"content" validate:"required,max=280"`
	Timestamp       string   `json:"timestamp" validate:"required,datetime=YYYY-MM-DDTHH:MM:SSZ"`
	UserID          string   `json:"user_id" validate:"required,id_field"`
	Hashtags        []string `json:"hashtags" validate:"required"`
	Reactions       int      `json:"reactions" validate:"gte=0"`
	Shares          int      `json:"shares" validate:"gte=0"`
	Views           int      `json:"views" validate:"gte=0"`
	ParsedTimeStamp time.Time
	Sentiment       MessageSentiment
}

type Feed struct {
	Messages          []FeedMessage `json:"messages" validate:"required,dive,required"`
	TimeWindowMinutes int           `json:"time_window_minutes" validate:"gte=0"`
}

type FeedSentimentDistribution struct {
	Positive float64 `json:"positive"`
	Negative float64 `json:"negative"`
	Neutral  float64 `json:"neutral"`
}

type FeedSentimentFlags struct {
	MbrasEmployee      bool `json:"mbras_employee"`
	SpecialPattern     bool `json:"special_pattern"`
	CandidateAwareness bool `json:"candidate_awareness"`
}

type ContentToken struct {
	Original   string
	Normalized string
}

type MessageSentiment struct {
	Score            float64 `json:"score"`
	Label            string  `json:"label"`
	ProcessingTimeMs float64 `json:"processing_time_ms"`
}

type UserInfluenceRanking struct {
	UserID         string  `json:"user_id"`
	InfluenceScore float64 `json:"influence_score"`
}
type FeedSentiment struct {
	SentimentDistribution FeedSentimentDistribution `json:"sentiment_distribution"`
	EngagementScore       float64                   `json:"engagement_score"`
	TrendingTopics        []string                  `json:"trending_topics"`
	InfluenceRanking      []UserInfluenceRanking    `json:"influence_ranking"`
	AnomalyDetected       bool                      `json:"anomaly_detected"`
	AnomalyType           string                    `json:"anomaly_type"`
	Flags                 FeedSentimentFlags        `json:"flags"`
	ProcessingTimeMs      float64                   `json:"processing_time_ms"`
}

type FeedSentimentAnalysis struct {
	Analysis FeedSentiment `json:"analysis"`
}

type MessageSentimentAnalysis struct {
	Analysis MessageSentiment `json:"analysis"`
}
