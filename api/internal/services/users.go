package services

import (
	"context"
	"errors"

	"github.com/XDoubleU/essentia/pkg/database"
	errortools "github.com/XDoubleU/essentia/pkg/errors"

	"check-in/api/internal/dtos"
	"check-in/api/internal/models"
	"check-in/api/internal/repositories"
)

type UserService struct {
	users repositories.UserRepository
}

func (service UserService) GetTotalCount(ctx context.Context) (*int64, error) {
	return service.users.GetTotalCount(ctx)
}

func (service UserService) GetAll(
	ctx context.Context,
) ([]*models.User, error) {
	return service.users.GetAll(ctx)
}

func (service UserService) GetAllPaginated(
	ctx context.Context,
	_ *models.User,
	limit int64,
	offset int64,
) ([]*models.User, error) {
	return service.users.GetAllPaginated(ctx, limit, offset)
}

func (service UserService) GetByID(
	ctx context.Context,
	id string,
	role models.Role,
) (*models.User, error) {
	user, err := service.users.GetByID(ctx, id, role)
	if err != nil {
		if errors.Is(err, database.ErrResourceNotFound) {
			return nil, errortools.NewNotFoundError("user", id, "id")
		}
		return nil, err
	}

	return user, nil
}

func (service UserService) GetByUsername(
	ctx context.Context,
	username string,
) (*models.User, error) {
	return service.users.GetByUsername(ctx, username)
}

func (service UserService) Create(
	ctx context.Context,
	createUserDto dtos.CreateUserDto,
	role models.Role,
) (*models.User, error) {
	passwordHash, err := models.HashPassword(createUserDto.Password)
	if err != nil {
		return nil, err
	}

	user, err := service.users.Create(ctx, createUserDto.Username, passwordHash, role)
	if err != nil {
		if errors.Is(err, database.ErrResourceConflict) {
			return nil, errortools.NewConflictError(
				"user",
				createUserDto.Username,
				"username",
			)
		}
		return nil, err
	}

	return user, nil
}

func (service UserService) Update(
	ctx context.Context,
	id string,
	updateUserDto dtos.UpdateUserDto,
	role models.Role,
) (*models.User, error) {
	user, err := service.GetByID(ctx, id, role)
	if err != nil {
		if errors.Is(err, database.ErrResourceNotFound) {
			return nil, errortools.NewNotFoundError("user", id, "id")
		}
		return nil, err
	}

	user, err = service.users.Update(ctx, *user, updateUserDto, role)
	if err != nil {
		if errors.Is(err, database.ErrResourceConflict) {
			return nil, errortools.NewConflictError(
				"user",
				*updateUserDto.Username,
				"username",
			)
		}
		return nil, err
	}

	return user, nil
}

func (service UserService) Delete(
	ctx context.Context,
	id string,
	role models.Role,
) (*models.User, error) {
	user, err := service.GetByID(ctx, id, role)
	if err != nil {
		if errors.Is(err, database.ErrResourceNotFound) {
			return nil, errortools.NewNotFoundError("user", id, "id")
		}

		return nil, err
	}

	err = service.users.Delete(ctx, id, role)
	if err != nil {
		return nil, err
	}

	return user, nil
}
