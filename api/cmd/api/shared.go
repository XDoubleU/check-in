package main

import (
	"context"
	"math"

	"check-in/api/internal/dtos"
)

type Service[T any] interface {
	GetAllPaginated(ctx context.Context, limit int64, offset int64) ([]*T, error)
	GetTotalCount(ctx context.Context) (*int64, error)
}

func getAllPaginated[T any](
	ctx context.Context,
	service Service[T],
	page int64,
	pageSize int64,
) (*dtos.PaginatedResultDto[T], error) {
	data, err := service.GetAllPaginated(ctx, pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, err
	}

	total, err := service.GetTotalCount(ctx)
	if err != nil {
		return nil, err
	}

	return &dtos.PaginatedResultDto[T]{
		Data: data,
		Pagination: dtos.Pagination{
			Current: page,
			Total:   int64(math.Ceil(float64(*total) / float64(pageSize))),
		},
	}, nil
}
