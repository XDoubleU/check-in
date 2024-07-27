package services

import (
	"context"

	"check-in/api/internal/dtos"
	"check-in/api/internal/models"
	"check-in/api/internal/repositories"
)

type UserService struct {
	users     repositories.UserRepository
	locations LocationService
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
		return nil, err
	}

	if user.Role != models.DefaultRole {
		return user, nil
	}

	var location *models.Location
	location, err = service.locations.GetByUserID(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	user.Location = location

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
	username string,
	password string,
	role models.Role,
) (*models.User, error) {
	return service.users.Create(ctx, username, password, role)
}

func (service UserService) Update(
	ctx context.Context,
	user *models.User,
	updateUserDto dtos.UpdateUserDto,
	role models.Role,
) error {
	return service.users.Update(ctx, user, updateUserDto, role)
}

func (service UserService) Delete(
	ctx context.Context,
	id string,
	role models.Role,
) error {
	return service.users.Delete(ctx, id, role)
}
