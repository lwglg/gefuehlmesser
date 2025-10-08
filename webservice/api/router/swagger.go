package router

import (
	"fmt"
	"strings"

	swag "github.com/swaggo/swag"

	cfg "webservice/config"
	t "webservice/libs/tooling"
)

func sanitizeHttpSchemes(c *cfg.Conf) []string {
	delimiter := ","
	supportedSchemes := []string{"http", "https"}
	sanitizedStr := strings.ToLower(strings.TrimSpace(c.SwaggerUI.HttpSchemes))
	rawList := strings.Split(sanitizedStr, delimiter)

	schemes := []string{}

	for _, item := range rawList {
		sanitized := strings.TrimSpace(item)

		// Check if each parsed item is a valid HTTP protocol
		if t.HasAny(supportedSchemes, func(v string) bool { return v == sanitized }) {
			schemes = append(schemes, strings.TrimSpace(item))
		}
	}

	return schemes
}

// Programmatically set Swagger UI information
func SetSwaggetInfo(si *swag.Spec, c *cfg.Conf) {
	host := fmt.Sprint(c.SwaggerUI.Host, ":", c.Server.Port)

	si.Title = c.SwaggerUI.Title
	si.Description = c.SwaggerUI.Description
	si.Version = c.SwaggerUI.Version
	si.Host = host
	si.Schemes = sanitizeHttpSchemes(c)
}
