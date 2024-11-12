package repositories

import (
	"context"
	"crypto/sha256"

	"github.com/XDoubleU/essentia/pkg/database/postgres"

	"check-in/api/internal/models"
	"check-in/api/internal/shared"
)

type AuthRepository struct {
	db            postgres.DB
	getTimeNowUTC shared.UTCNowTimeProvider
}

func (repo AuthRepository) CreateToken(ctx context.Context, token *models.Token) error {
	query := `
		INSERT INTO tokens (hash, user_id, expiry, scope)
		VALUES ($1, $2, $3, $4)
	`

	_, err := repo.db.Exec(
		ctx,
		query,
		token.Hash,
		token.UserID,
		token.Expiry,
		token.Scope,
	)

	return err
}

func (repo AuthRepository) DeleteToken(ctx context.Context, tokenHash [32]byte) error {
	query := `
		DELETE FROM tokens
		WHERE hash = $1
	`

	_, err := repo.db.Exec(ctx, query, tokenHash[:])
	return err
}

func (repo AuthRepository) DeleteExpiredTokens(ctx context.Context) error {
	query := `
		DELETE FROM tokens
		WHERE expiry < $1
	`

	_, err := repo.db.Exec(ctx, query, repo.getTimeNowUTC())
	return err
}

func (repo AuthRepository) GetToken(
	ctx context.Context,
	scope models.Scope,
	tokenHash [32]byte,
) (*models.Token, *string, *models.Role, error) {
	query := `
		SELECT tokens.used, users.id, users.role
		FROM users
		INNER JOIN tokens
		ON tokens.user_id = users.id
		WHERE tokens.hash = $1
		AND tokens.scope = $2
		AND tokens.expiry > $3
	`

	args := []any{tokenHash[:], scope, repo.getTimeNowUTC()}

	var token models.Token
	var userID string
	var userRole models.Role

	err := repo.db.QueryRow(ctx, query, args...).
		Scan(&token.Used, &userID, &userRole)

	if err != nil {
		return nil, nil, nil, postgres.PgxErrorToHTTPError(err)
	}

	return &token, &userID, &userRole, nil
}

func (repo AuthRepository) DeleteAllTokensForUser(
	ctx context.Context,
	userID string,
) error {
	query := `
		DELETE FROM tokens
		WHERE user_id = $1
	`

	_, err := repo.db.Exec(ctx, query, userID)
	return err
}

func (repo AuthRepository) SetTokenAsUsed(
	ctx context.Context,
	tokenValue string,
) error {
	tokenHash := sha256.Sum256([]byte(tokenValue))

	query := `
		UPDATE tokens
		SET used = true
		WHERE hash = $1
	`

	_, err := repo.db.Exec(ctx, query, tokenHash[:])

	return err
}
