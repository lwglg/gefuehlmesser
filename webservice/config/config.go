package config

import (
	"log"
	"time"

	"github.com/joeshaw/envdecode"
)

type Conf struct {
	Server    ConfServer
	SwaggerUI ConfSwaggerUI
}

type ConfSwaggerUI struct {
	Title       string `env:"SWAGGER_UI_API_TITLE"`
	Description string `env:"SWAGGER_UI_API_DESCRIPTION"`
	Version     string `env:"SWAGGER_UI_API_VERSION"`
	Host        string `env:"SWAGGER_UI_HOSTNAME"`
	HttpSchemes string `env:"SWAGGER_UI_HTTP_SCHEMES"` // Comma-separated string, e.g. http,https
}

type ConfServer struct {
	Port         int           `env:"SERVER_PORT,required"`
	TimeoutRead  time.Duration `env:"SERVER_TIMEOUT_READ,required"`
	TimeoutWrite time.Duration `env:"SERVER_TIMEOUT_WRITE,required"`
	TimeoutIdle  time.Duration `env:"SERVER_TIMEOUT_IDLE,required"`
	Debug        bool          `env:"SERVER_DEBUG,required"`
}

func New() *Conf {
	var c Conf
	if err := envdecode.StrictDecode(&c); err != nil {
		log.Fatalf("Failed to decode: %s", err)
	}

	return &c
}
