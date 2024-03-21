package db

import (
	"context"
	"github.com/vynious/ascenda-lp-backend/types"
)

func (dbs *DBService) GetPoints(ctx context.Context) (*types.Points, error) {
	// todo: add logic
	return &types.Points{}, nil
}

func (dbs *DBService) GetPointsByUser(ctx context.Context, userId string) (*types.Points, error) {
	// todo: add logic
	return &types.Points{}, nil
}

func (dbs *DBService) UpdatePoints(ctx context.Context, userId string) (*types.Points, error) {
	// todo: add logic
	return &types.Points{}, nil
}
