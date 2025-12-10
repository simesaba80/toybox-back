package usecase

type TokenProvider interface {
	GenerateToken(userID string) (string, error)
	RegenerateToken(refreshToken string) (string, error)
}
