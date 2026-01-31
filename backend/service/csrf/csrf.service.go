package csrf

// CSRFService handles CSRF token generation and validation
type CSRFService interface {
	GenerateToken() (string, error)
	ValidateToken(token string) bool
}

// csrfService implements CSRFService
type csrfService struct{}

// NewCSRFService creates a new CSRF service
func NewCSRFService() CSRFService {
	return &csrfService{}
}
