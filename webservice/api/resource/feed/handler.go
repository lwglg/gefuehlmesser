package feed

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"

	errorHandler "webservice/api/resource/common/err"
	log "webservice/api/resource/common/log"
	ctxLib "webservice/libs/ctx"
	saLib "webservice/libs/sentiment"
	validationLib "webservice/libs/validator"
)

type API struct {
	logger     *zerolog.Logger
	validator  *validator.Validate
	repository *Repository
}

func New(logger *zerolog.Logger, validator *validator.Validate, analyzer *saLib.SentimentAnalyzer) *API {
	return &API{
		logger:     logger,
		validator:  validator,
		repository: NewRepository(analyzer),
	}
}

// Create godoc
//
//	@summary		Análise de sentimento de feed
//	@description	Realiza a análise determinística de sentimento das mensagens de um feed, dada uma payload de feed válida.
//	@tags			API de análise de sentimentos
//	@accept			json
//	@produce		json
//	@param			body	body	sentiment.Feed	true	"Payload contendo as mensagens de feed e parâmetro de janela temporal."
//	@success		200	{object}	sentiment.FeedSentiment
//	@failure		400	{object}	err.Error
//	@failure		422	{object}	err.Errors
//	@failure		500	{object}	err.Error
//	@router			/api/v1/sentiment/feed [post]
func (api *API) AnalyzeFeed(w http.ResponseWriter, r *http.Request) {
	reqID := ctxLib.RequestID(r.Context())
	feedPayload := &saLib.Feed{}

	err := json.NewDecoder(r.Body).Decode(feedPayload)
	if err != nil {
		api.logger.Error().Str(log.KeyReqID, reqID).Err(err).Msg("")
		errorHandler.BadRequest(w, errorHandler.RespJSONDecodeFailure)
		return
	}

	err = api.validator.Struct(feedPayload)
	if err != nil {
		respBody, err := json.Marshal(validationLib.ToErrResponse(err))
		if err != nil {
			api.logger.Error().Str(log.KeyReqID, reqID).Err(err).Msg("")
			errorHandler.ServerError(w, errorHandler.RespJSONEncodeFailure)
			return
		}

		errorHandler.ValidationErrors(w, respBody)
		return
	}

	analysis, err := api.repository.AnalyzeFeed(feedPayload)
	if err != nil {
		api.logger.Error().Str(log.KeyReqID, reqID).Err(err).Msg("An error occurred while performing the sentiment analysis")
		errorHandler.ServerError(w, errorHandler.RespFeedSentimentFailure)
		return
	}

	err = json.NewEncoder(w).Encode(analysis)
	if err != nil {
		api.logger.Error().Str(log.KeyReqID, reqID).Err(err).Msg("")
		errorHandler.ServerError(w, errorHandler.RespJSONEncodeFailure)
		return
	}
}

// Create godoc
//
//	@summary		Análise de sentimento de mensagem
//	@description	Realiza a análise determinística de sentimento de uma mensagem de feed, dada uma payload de feed válida.
//	@tags			API de análise de sentimentos
//	@accept			json
//	@produce		json
//	@param			body	body	sentiment.FeedMessage	true	"Payload contendo a estrutura de uma mensagen de feed."
//	@success		200	{object}	sentiment.MessageSentiment
//	@failure		400	{object}	err.Error
//	@failure		422	{object}	err.Errors
//	@failure		500	{object}	err.Error
//	@router			/api/v1/sentiment/message [post]
func (api *API) AnalyzeMessage(w http.ResponseWriter, r *http.Request) {
	reqID := ctxLib.RequestID(r.Context())
	feedPayload := &saLib.FeedMessage{}

	err := json.NewDecoder(r.Body).Decode(feedPayload)
	if err != nil {
		api.logger.Error().Str(log.KeyReqID, reqID).Err(err)
		errorHandler.BadRequest(w, errorHandler.RespJSONDecodeFailure)
		return
	}

	err = api.validator.Struct(feedPayload)
	if err != nil {
		respBody, err := json.Marshal(validationLib.ToErrResponse(err))
		if err != nil {
			api.logger.Error().Str(log.KeyReqID, reqID).Err(err)
			errorHandler.ServerError(w, errorHandler.RespJSONEncodeFailure)
			return
		}

		errorHandler.ValidationErrors(w, respBody)
		return
	}

	analysis, err := api.repository.AnalyzeMessage(feedPayload)
	if err != nil {
		api.logger.Error().Str(log.KeyReqID, reqID).Err(err)
		errorHandler.ServerError(w, errorHandler.RespFeedSentimentFailure)
		return
	}

	err = json.NewEncoder(w).Encode(analysis)
	if err != nil {
		api.logger.Error().Str(log.KeyReqID, reqID).Err(err)
		errorHandler.ServerError(w, errorHandler.RespJSONEncodeFailure)
		return
	}
}
