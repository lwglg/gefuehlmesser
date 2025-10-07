package feed

import (
	sa "webservice/libs/sentiment"
)

type Repository struct {
	analyzer *sa.SentimentAnalyzer
}

func NewRepository(analyzer *sa.SentimentAnalyzer) *Repository {
	return &Repository{
		analyzer: analyzer,
	}
}

func (r *Repository) AnalyzeFeed(feed *sa.Feed) (*sa.SentimentAnalysis, error) {
	analysis := r.analyzer.AnalyzeFeed(feed)

	return analysis, nil
}
