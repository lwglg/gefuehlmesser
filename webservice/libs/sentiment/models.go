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
	Hashtags        []string `json:"hashtags" validate:"required,dive,hashtag"`                   // Lista de hashtags associadas às mensagens do feed
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
	Messages          []FeedMessage `json:"messages" validate:"required,dive,required"` // A lista de mensagens para aquele feed
	TimeWindowMinutes int           `json:"time_window_minutes" validate:"gte=0"`       // A janela de tempo, em minutos, na qual serão consideradas mensagens para a análise de sentimentos
}

// FeedSentimentDistribution
//
// @description    Percentuais dos sentimentos detectados
// @description    das mensagens, por categoria
type FeedSentimentDistribution struct {
	Positive float64 `json:"positive"` // Fração de mensagens com sentimento positivo
	Negative float64 `json:"negative"` // Fração de mensagens com sentimento negativo
	Neutral  float64 `json:"neutral"`  // Fração de mensagens com sentimento nêutro
}

// FeedSentimentFlags
//
// @description    Indicadores de casos especiais
// @description    detectados nas mensagens
type FeedSentimentFlags struct {
	MbrasEmployee      bool `json:"mbras_employee"`      // Indica se o "user_id" contém a substring "mbras"
	SpecialPattern     bool `json:"special_pattern"`     // Indica a condição anterior e se também a mensagem possui exatamente 42 caracteres
	CandidateAwareness bool `json:"candidate_awareness"` // Indica se o conteúdo normalizado da mensagem corresponde a "teste tecnico mbras"
}

// ContentToken
//
// @description    Objeto que contem o resultado
// @description    da tokenização de uma string

type ContentToken struct {
	Original   string // String original
	Normalized string // String normalizada
}

// MessageSentiment
//
// @description		Objeto contendo o resultado da
// @description		análise de sentimento para uma mensagem

type MessageSentiment struct {
	Score            float64 `json:"score"`              // Número real, expressando quantitativamente o sentimento da mensagem
	Label            string  `json:"label"`              // Categoria do sentimento ("positive" | "negative" | "neutral"), de acordo com o valor do escore
	ProcessingTimeMs float64 `json:"processing_time_ms"` // Auxiliar. Tempo de processamento da análise em milissegundos
}

// UserInfluenceRanking
//
// @description		Objeto contendo o índice de influência
// @description		no feed para um usuário
type UserInfluenceRanking struct {
	UserID         string  `json:"user_id"`         // Identificação do usuario, e.g. user_123
	InfluenceScore float64 `json:"influence_score"` // Índice de influência
}

// FeedSentiment
//
// @description		Objeto que agrega as métricas determinísticas
// @description		de análise de sentimento para um feed

type FeedSentiment struct {
	SentimentDistribution FeedSentimentDistribution `json:"sentiment_distribution"`
	EngagementScore       float64                   `json:"engagement_score"` // Índice de engajamento
	TrendingTopics        []string                  `json:"trending_topics"`  // Lista unificada de hashtags, citadas nas mensagens
	InfluenceRanking      []UserInfluenceRanking    `json:"influence_ranking"`
	AnomalyDetected       bool                      `json:"anomaly_detected"` // Indica se alguma anomalia no feed foi detectada
	AnomalyType           string                    `json:"anomaly_type"`     // Categoria da anomalia ("synchronized_posting" | "burst")
	Flags                 FeedSentimentFlags        `json:"flags"`
	ProcessingTimeMs      float64                   `json:"processing_time_ms"` // Auxiliar. Tempo de processamento da análise, em milissegundos
}

// FeedSentimentAnalysis
//
// @description		Objeto contendo a análise de sentimento
// @description		para um feed de mensagens

type FeedSentimentAnalysis struct {
	Analysis FeedSentiment `json:"analysis"` // Análise de sentimento para um feed
}

// MessageSentimentAnalysis
//
// @description		Objeto contendo a análise de sentimento
// @description		para uma mensagem

type MessageSentimentAnalysis struct {
	Analysis MessageSentiment `json:"analysis"` // Análise de sentimento para uma mensagem
}
