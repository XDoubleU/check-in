package services

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"net/http"
	"time"

	"github.com/xdoubleu/essentia/pkg/database"
	errortools "github.com/xdoubleu/essentia/pkg/errors"
	"github.com/xhit/go-str2duration/v2"

	"check-in/api/internal/dtos"
	"check-in/api/internal/models"
	"check-in/api/internal/repositories"
)

type AuthService struct {
	auth  repositories.AuthRepository
	users UserService
}

func (service AuthService) SignInUser(ctx context.Context, signInDto *dtos.SignInDto) (*models.User, error) {
	if v := signInDto.Validate(); !v.Valid() {
		return nil, errortools.ErrFailedValidation
	}

	user, err := service.users.GetByUsername(ctx, signInDto.Username)
	if err != nil {
		switch err {
		case database.ErrResourceNotFound:
			return nil, errortools.NewUnauthorizedError(errors.New("invalid credentials"))
		default:
			return nil, err
		}
	}

	match, _ := user.CompareHashAndPassword(signInDto.Password)
	if !match {
		return nil, errortools.NewUnauthorizedError(errors.New("invalid credentials"))
	}

	return user, nil
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
	ttl, err := str2duration.ParseDuration(expiry)
	if err != nil {
		return nil, err
	}

	token, err := service.generateToken(userID, ttl, scope)
	if err != nil {
		return nil, err
	}

	err = service.auth.CreateToken(ctx, token)
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
	ctx context.Context,
	scope models.Scope,
	tokenValue string,
) (*http.Cookie, error) {
	err := service.auth.DeleteToken(
		ctx,
		service.hashTokenValue(tokenValue),
	)
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
	return service.auth.DeleteExpiredTokens(ctx)
}

func (service AuthService) GetToken(
	ctx context.Context,
	scope models.Scope,
	tokenValue string,
) (*models.Token, *models.User, error) {
	token, userID, userRole, err := service.auth.GetToken(
		ctx,
		scope,
		service.hashTokenValue(tokenValue),
	)
	if err != nil {
		return nil, nil, err
	}

	user, err := service.users.GetByID(ctx, *userID, *userRole)
	if err != nil {
		return nil, nil, err
	}

	return token, user, nil
}

func (service AuthService) DeleteAllTokensForUser(
	ctx context.Context,
	userID string,
) error {
	return service.auth.DeleteAllTokensForUser(ctx, userID)
}

func (service AuthService) SetTokenAsUsed(
	ctx context.Context,
	tokenValue string,
) error {
	return service.auth.SetTokenAsUsed(ctx, tokenValue)
}

func (service AuthService) generateToken(
	userID string,
	ttl time.Duration,
	scope models.Scope,
) (*models.Token, error) {
	//nolint:exhaustruct //other fields are optional
	token := &models.Token{
		UserID: userID,
		Expiry: time.Now().Add(ttl),
		Scope:  scope,
	}

	randomBytes := make([]byte, 16) //nolint:mnd //no magic number

	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	token.Plaintext = base64.StdEncoding.EncodeToString(randomBytes)
	hash := service.hashTokenValue(token.Plaintext)
	token.Hash = hash[:]

	return token, nil
}

func (service AuthService) hashTokenValue(tokenValue string) [32]byte {
	return sha256.Sum256([]byte(tokenValue))
}
