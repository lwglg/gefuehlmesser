package healthcheck

type HealthcheckData struct {
	Status    int    `json:"status_code"`
	RequestId string `json:"request_id"`
	Message   string `json:"message"`
}
