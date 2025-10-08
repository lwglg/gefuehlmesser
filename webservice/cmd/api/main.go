package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"webservice/api/router"
	"webservice/config"
	"webservice/libs/logger"
	"webservice/libs/sentiment"
	"webservice/libs/validator"
)

//	@title			Gefuehlmesser API
//	@version		1.0.0
//	@description	RESTful API providing sentimnet analysis for feeds' messages

//	@contact.name	Guilherme Gonçalves
//	@contact.url	https://github.com/lwglg

//	@license.name	MIT License
//	@license.url	https://github.com/lwglg/gefuehlmesser/blob/main/LICENSE

// @servers.url	localhost:8080/v1
func main() {
	c := config.New()
	l := logger.New(c.Server.Debug)
	v := validator.New()
	a := sentiment.New()
	r := router.New(l, v, a, c)

	s := &http.Server{
		Addr:         fmt.Sprintf(":%d", c.Server.Port),
		Handler:      r,
		ReadTimeout:  c.Server.TimeoutRead,
		WriteTimeout: c.Server.TimeoutWrite,
		IdleTimeout:  c.Server.TimeoutIdle,
	}

	closed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		l.Info().Msgf("Shutting down server %v", s.Addr)

		ctx, cancel := context.WithTimeout(context.Background(), c.Server.TimeoutIdle)
		defer cancel()

		if err := s.Shutdown(ctx); err != nil {
			l.Error().Err(err).Msg("Server shutdown failure")
		}

		close(closed)
	}()

	l.Info().Msgf("Starting server %v", s.Addr)
	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		l.Fatal().Err(err).Msg("Server startup failure")
	}

	<-closed
	l.Info().Msgf("Server shutdown successfully")
}
