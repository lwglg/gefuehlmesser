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

func (r *Repository) AnalyzeFeed(feed *sa.Feed) (*sa.FeedSentimentAnalysis, error) {
	analysis, err := r.analyzer.AnalyzeFeed(feed)
	if err != nil {
		return nil, err
	}

	return analysis, nil
}

func (r *Repository) AnalyzeMessage(feedMessage *sa.FeedMessage) (*sa.MessageSentimentAnalysis, error) {
	sentiment, err := r.analyzer.AnalyzeMessage(*feedMessage)
	if err != nil {
		return nil, err
	}

	return sentiment, nil
}
