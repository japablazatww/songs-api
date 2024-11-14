package services

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/japablazatww/song-searcher/cmd/api/internal/core/domain"
	"github.com/japablazatww/song-searcher/cmd/api/internal/core/ports"
)

type AuthService struct {
	Repo ports.AuthRepository
}

func NewAuthService(repo ports.AuthRepository) *AuthService {
	return &AuthService{Repo: repo}
}

func generateSecureToken(length int) string {
	b := make([]byte, length)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func (s *AuthService) CreateApp(appName string) (*domain.AppAuth, error) {
	app := &domain.AppAuth{
		AppName:      appName,
		ClientID:     generateSecureToken(16),
		ClientSecret: generateSecureToken(32),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err := s.Repo.CreateApp(app)
	if err != nil {
		fmt.Printf("Error creating app: %v\n", err)
		return nil, err
	}

	return app, nil
}

func (s *AuthService) GenerateToken(clientID, clientSecret string) (string, error) {
	app, err := s.Repo.GetAppByClientCredentials(clientID, clientSecret)
	if err != nil {
		fmt.Println("GetAppByClientCredentials", err)
		return "", err
	}

	if app == nil {
		fmt.Println("app is nil", app)
		return "", fmt.Errorf("invalid credentials")
	}

	// Crea las claims del JWT
	claims := &domain.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "tu_aplicacion",
			Subject:   app.ClientID,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
		ClientID: app.ClientID,
		AppName:  app.AppName,
	}

	// Crea el token JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", fmt.Errorf("JWT_SECRET not set in .env file en GenerateToken")
	}

	// Firma el token con tu clave secreta
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("could not sign token: %w", err)
	}

	return signedToken, nil
}
