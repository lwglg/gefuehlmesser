package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"

	httpSwagger "github.com/swaggo/http-swagger"

	"webservice/api/resource/feed"
	"webservice/api/resource/utilitary/healthcheck"
	"webservice/api/router/middleware"
	"webservice/api/router/middleware/requestlog"

	docs "webservice/docs"
	sa "webservice/libs/sentiment"
)

func New(l *zerolog.Logger, v *validator.Validate, a *sa.SentimentAnalyzer) *chi.Mux {
	r := chi.NewRouter()

	SetSwaggetInfo(docs.SwaggerInfo)

	r.Route("/utilitary", func(r chi.Router) {
		r.Use(middleware.ContentTypeJSON)

		healthCheckAPI := healthcheck.New(l)
		r.Method(http.MethodGet, "/healthcheck", requestlog.NewHandler(healthCheckAPI.Read, l))
	})

	r.Route("/api/v1", func(r chi.Router) {
		r.Use(middleware.RequestID)
		r.Use(middleware.ContentTypeJSON)

		r.Method(http.MethodGet, "/swagger/*", requestlog.NewHandler(
			httpSwagger.Handler(
				httpSwagger.URL("http://localhost:8080/api/v1/swagger/doc.json"),
			), l),
		)

		feedAPI := feed.New(l, v, a)
		r.Method(http.MethodPost, "/feeds/analyze", requestlog.NewHandler(feedAPI.AnalyzeFeed, l))
	})

	return r
}
