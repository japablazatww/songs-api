package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/japablazatww/song-searcher/cmd/api/internal/core/domain"
)

func (r *PostgresRepository) CreateApp(app *domain.AppAuth) error {
	query := `
        INSERT INTO app_auth (app_name, client_id, client_secret, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id`

	return r.DB.QueryRow(
		query,
		app.AppName,
		app.ClientID,
		app.ClientSecret,
		app.CreatedAt,
		app.UpdatedAt,
	).Scan(&app.ID)
}

func (r *PostgresRepository) GetAppByClientCredentials(clientID, clientSecret string) (*domain.AppAuth, error) {
	fmt.Println("ME EJECUTO TAMBIEN")
	app := &domain.AppAuth{}
	query := `
        SELECT id, app_name, client_id, client_secret, created_at, updated_at
        FROM app_auth
        WHERE client_id = $1 AND client_secret = $2`

	err := r.DB.QueryRow(query, clientID, clientSecret).Scan(
		&app.ID,
		&app.AppName,
		&app.ClientID,
		&app.ClientSecret,
		&app.CreatedAt,
		&app.UpdatedAt,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return app, nil
}

func (r *PostgresRepository) UpdateAppToken(clientID string) error {
	_, err := r.DB.Exec("UPDATE app_auth SET updated_at = $1 WHERE client_id = $2", time.Now(), clientID)
	return err
}

func (r *PostgresRepository) GetAppByClientID(clientID string) (*domain.AppAuth, error) {
	var app domain.AppAuth
	err := r.DB.QueryRow("SELECT id, app_name, client_id, client_secret, created_at, updated_at FROM app_auth WHERE client_id = $1", clientID).Scan(
		&app.ID,
		&app.AppName,
		&app.ClientID,
		&app.ClientSecret,
		&app.CreatedAt,
		&app.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("error retrieving app by client ID: %w", err)
	}

	return &app, nil
}

func (r *PostgresRepository) ValidateToken(tokenString string) error {

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		fmt.Println("JWT_SECRET not set in .env file en ValidateToken")
		return errors.New("JWT_SECRET not set in .env file")
	}

	token, err := jwt.ParseWithClaims(tokenString, &domain.Claims{}, func(token *jwt.Token) (interface{}, error) {

		// se verifica el tipo de signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// se retorna la clave secreta
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return err
	}

	claims, ok := token.Claims.(*domain.Claims)
	if !ok {
		return errors.New("invalid token claims")
	}

	// Verificar si el clientID existe en la base de datos
	app, err := r.GetAppByClientID(claims.ClientID)
	if err != nil || app == nil {
		return errors.New("invalid client ID")
	}

	return nil
}
