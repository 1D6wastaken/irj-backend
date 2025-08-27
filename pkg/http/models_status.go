package http

type ServiceStatus = string

const (
	OK ServiceStatus = "OK"
	KO ServiceStatus = "KO"
)

type Service struct {
	Name     string    `json:"name"`
	Status   string    `json:"status"`
	Details  string    `json:"details,omitempty"`
	Services []Service `json:"services,omitempty"`
}
