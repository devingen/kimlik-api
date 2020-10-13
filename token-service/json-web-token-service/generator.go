package json_web_token_service

// NewTokenService generates new DatabaseService
func NewTokenService() *JWTService {
	service := &JWTService{}
	service.Init()
	return service
}
