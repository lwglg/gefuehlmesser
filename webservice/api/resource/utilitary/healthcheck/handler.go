package healthcheck

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog"

	e "webservice/api/resource/common/err"
	l "webservice/api/resource/common/log"
	ctx "webservice/libs/ctx"
)

type API struct {
	logger *zerolog.Logger
}

func New(logger *zerolog.Logger) *API {
	return &API{
		logger: logger,
	}
}

// Read godoc
//
//	@summary		Healthcheck
//	@description	Realiza uma requisição de healthcheck ao servidor, retornnando uma mensagem simples.
//	@tags			API utilitária
//	@produce		json
//	@success		200 {object}	healthcheck.HealthcheckData
//	@failure		500	{object}	err.Error
//	@router			/utilitary/healthcheck [get]
func (a *API) Read(w http.ResponseWriter, r *http.Request) {
	reqID := ctx.RequestID(r.Context())

	data := HealthcheckData{
		Status:    http.StatusOK,
		RequestId: reqID,
		Message:   "I'm alive",
	}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		a.logger.Error().Str(l.KeyReqID, reqID).Err(err).Msg("Error while trasmitting healthcheck data")
		e.ServerError(w, e.RespHealthCheckFailure)

		return
	}
}
