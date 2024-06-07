package services

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"time"

	"github.com/XDoubleU/essentia/pkg/database/postgres"
	str2duration "github.com/xhit/go-str2duration/v2"

	"check-in/api/internal/models"
)

type AuthService struct {
	db        postgres.DB
	locations LocationService
}

func (service AuthService) GetCookieName(scope models.Scope) string {
	switch {
	case scope == models.AccessScope:
		return "accessToken"
	case scope == models.RefreshScope:
		return "refreshToken"
	default:
		panic("invalid scope")
	}
}

func (service AuthService) CreateCookie(
	ctx context.Context,
	scope models.Scope,
	userID string,
	expiry string,
	secure bool,
) (*http.Cookie, error) {
	ttl, _ := str2duration.ParseDuration(expiry)
	token, err := service.newToken(ctx, userID, ttl, scope)
	if err != nil {
		return nil, err
	}

	name := service.GetCookieName(scope)

	cookie := http.Cookie{
		Name:     name,
		Value:    token.Plaintext,
		Expires:  token.Expiry,
		SameSite: http.SameSiteStrictMode,
		HttpOnly: true,
		Secure:   secure,
		Path:     "/",
	}

	return &cookie, nil
}

func (service AuthService) DeleteCookie(
	scope models.Scope,
	value string,
) (*http.Cookie, error) {
	err := service.deleteToken(context.Background(), value)
	if err != nil {
		return nil, err
	}

	name := service.GetCookieName(scope)

	return &http.Cookie{
		Name:     name,
		Value:    "",
		MaxAge:   -1,
		SameSite: http.SameSiteStrictMode,
		HttpOnly: true,
		Path:     "/",
	}, nil
}

func (service AuthService) DeleteExpiredTokens(ctx context.Context) error {
	query := `
		DELETE FROM tokens
		WHERE expiry < $1
	`

	_, err := service.db.Exec(ctx, query, time.Now())
	return err
}

func (service AuthService) GetToken(
	ctx context.Context,
	scope models.Scope,
	tokenValue string,
) (*models.Token, *models.User, error) {
	tokenHash := sha256.Sum256([]byte(tokenValue))

	query := `
		SELECT tokens.used, users.id, users.username, users.role, users.password_hash
		FROM users
		INNER JOIN tokens
		ON tokens.user_id = users.id
		WHERE tokens.hash = $1
		AND tokens.scope = $2
		AND tokens.expiry > $3
	`

	args := []any{tokenHash[:], scope, time.Now()}

	var token models.Token
	var user models.User

	err := service.db.QueryRow(ctx, query, args...).
		Scan(&token.Used, &user.ID, &user.Username, &user.Role, &user.PasswordHash)

	if err != nil {
		return nil, nil, handleError(err)
	}

	if user.Role == models.DefaultRole {
		var location *models.Location
		location, err = service.locations.GetByUserID(ctx, user.ID)
		if err != nil {
			return nil, nil, handleError(err)
		}

		user.Location = location
	}

	return &token, &user, nil
}

func (service AuthService) DeleteAllTokensForUser(
	ctx context.Context,
	userID string,
) error {
	query := `
		DELETE FROM tokens
		WHERE user_id = $1
	`

	_, err := service.db.Exec(ctx, query, userID)
	return err
}

func (service AuthService) SetTokenAsUsed(
	ctx context.Context,
	tokenValue string,
) error {
	tokenHash := sha256.Sum256([]byte(tokenValue))

	query := `
		UPDATE tokens
		SET used = true
		WHERE hash = $1
	`

	_, err := service.db.Exec(ctx, query, tokenHash[:])

	return err
}

func (service AuthService) deleteToken(ctx context.Context, value string) error {
	hash := sha256.Sum256([]byte(value))

	query := `
		DELETE FROM tokens
		WHERE hash = $1
	`

	_, err := service.db.Exec(ctx, query, hash[:])
	return err
}

func (service AuthService) newToken(
	ctx context.Context,
	userID string,
	ttl time.Duration,
	scope models.Scope,
) (*models.Token, error) {
	token, err := service.generateToken(userID, ttl, scope)
	if err != nil {
		return nil, err
	}

	err = service.createToken(ctx, token)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (service AuthService) generateToken(
	userID string,
	ttl time.Duration,
	scope models.Scope,
) (*models.Token, error) {
	token := &models.Token{
		UserID: userID,
		Expiry: time.Now().Add(ttl),
		Scope:  scope,
	}

	randomBytes := make([]byte, 16) //nolint:gomnd //no magic number

	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	token.Plaintext = base64.StdEncoding.EncodeToString(randomBytes)
	hash := sha256.Sum256([]byte(token.Plaintext))
	token.Hash = hash[:]

	return token, nil
}

func (service AuthService) createToken(ctx context.Context, token *models.Token) error {
	query := `
		INSERT INTO tokens (hash, user_id, expiry, scope)
		VALUES ($1, $2, $3, $4)
	`

	_, err := service.db.Exec(
		ctx,
		query,
		token.Hash,
		token.UserID,
		token.Expiry,
		token.Scope,
	)

	return err
}
