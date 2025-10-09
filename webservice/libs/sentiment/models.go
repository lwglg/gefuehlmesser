package sentiment

import "time"

// FeedMessage
//
// @description    Estrutura de uma típica mensagem,
// @description    associada a um feed de usuário
type FeedMessage struct {
	ID              string   `json:"id" validate:"required,id_field"`                             // O ID da mensagem, e.g. msg_123
	Content         string   `json:"content" validate:"required,max=280"`                         // O conteúdo textual da mensagem
	Timestamp       string   `json:"timestamp" validate:"required,datetime=YYYY-MM-DDTHH:MM:SSZ"` // Data/hora local, em formato UTC (GMT - ZULU Time)
	UserID          string   `json:"user_id" validate:"required,id_field"`                        // ID do usuário ao qual a mensagem está associada, e.g. user_324
	Hashtags        []string `json:"hashtags" validate:"required"`                                // Lista de hashtags associadas às mensagens do feed
	Reactions       int      `json:"reactions" validate:"gte=0"`                                  // Total de reações a mensagem
	Shares          int      `json:"shares" validate:"gte=0"`                                     // Total de cxompartilhamentos da mensagem
	Views           int      `json:"views" validate:"gte=0"`                                      // Total de visualizações da mensagem
	ParsedTimeStamp time.Time
	Sentiment       MessageSentiment
}

// Feed
//
// @description    Estrutura do feed de mensagens,
// @description    associados a um usuário
type Feed struct {
	Messages          []FeedMessage `json:"messages" validate:"required,dive,required"` // A lista de mensagem para aquele feed
	TimeWindowMinutes int           `json:"time_window_minutes" validate:"gte=0"`       // A janela de tempo, em minutos, na qual serão consideradas mensagens para a análise de sentimentos
}

// FeedSentimentDistribution
//
// @description    Contém o percentual de mensagem
// @description    message feed.
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
