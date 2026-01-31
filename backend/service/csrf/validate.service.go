package csrf

// ValidateToken validates a CSRF token (basic check for non-empty)
// In production, you'd check against stored tokens in session/Redis
func (s *csrfService) ValidateToken(token string) bool {
	return token != ""
}
