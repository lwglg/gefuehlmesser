package router

import swag "github.com/swaggo/swag"

// Programmatically set swagger info
func SetSwaggetInfo(si *swag.Spec) {
	si.Title = "Gefuehlmesser API"
	si.Description = "API RESTful, escrita em Go, para realizar análises de sentimentos, em tempo real, processando mensagens de feeds para calcular métricas de engajamento usando algorítmos determinísticos."
	si.Version = "1.0.0"
	si.Host = "localhost:8080"
	si.Schemes = []string{"http"}
}
