package services

import (
	"context"

	"check-in/api/internal/dtos"
	"check-in/api/internal/models"
	"check-in/api/internal/repositories"

	"github.com/xdoubleu/essentia/pkg/database"
	errortools "github.com/xdoubleu/essentia/pkg/errors"
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
	user *models.User,
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
		switch err {
		case database.ErrResourceNotFound:
			return nil, errortools.NewNotFoundError("user", id, "id")
		default:
			return nil, err
		}
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
	createUserDto *dtos.CreateUserDto,
	role models.Role,
) (*models.User, error) {
	if v := createUserDto.Validate(); !v.Valid() {
		return nil, errortools.ErrFailedValidation
	}

	passwordHash, err := models.HashPassword(createUserDto.Password)
	if err != nil {
		return nil, err
	}

	user, err := service.users.Create(ctx, createUserDto.Username, passwordHash, role)
	if err != nil {
		switch err {
		case database.ErrResourceConflict:
			return nil, errortools.NewConflictError("user", createUserDto.Username, "username")
		default:
			return nil, err
		}
	}

	return user, nil
}

func (service UserService) Update(
	ctx context.Context,
	id string,
	updateUserDto *dtos.UpdateUserDto,
	role models.Role,
) (*models.User, error) {
	if v := updateUserDto.Validate(); !v.Valid() {
		return nil, errortools.ErrFailedValidation
	}

	user, err := service.GetByID(ctx, id, role)
	if err != nil {
		switch err {
		case database.ErrResourceNotFound:
			return nil, errortools.NewNotFoundError("user", id, "id")
		default:
			return nil, err
		}
	}

	user, err = service.users.Update(ctx, *user, updateUserDto, role)
	if err != nil {
		switch err {
		case database.ErrResourceConflict:
			return nil, errortools.NewConflictError("user", *updateUserDto.Username, "username")
		default:
			return nil, err
		}
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
		switch err {
		case database.ErrResourceNotFound:
			return nil, errortools.NewNotFoundError("user", id, "id")
		default:
			return nil, err
		}
	}

	err = service.users.Delete(ctx, id, role)
	if err != nil {
		return nil, err
	}

	return user, nil
}
