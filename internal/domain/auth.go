package domain

type AuthService interface {
	GenerateToken(userID string) (string, error)
	VerifyToken(token string) (bool, error)
}
