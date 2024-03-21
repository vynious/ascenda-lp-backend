package db

import (
	"context"

	"github.com/vynious/ascenda-lp-backend/types/points"
)

func (dbs *DBService) GetPoints(ctx context.Context) (*points.Points, error) {
	// todo: add logic
	return &points.Points{}, nil
}

func (dbs *DBService) GetPointsByUser(ctx context.Context, userId string) (*points.Points, error) {
	// todo: add logic
	return &points.Points{}, nil
}

func (dbs *DBService) UpdatePoints(ctx context.Context, userId string) (*points.Points, error) {
	// todo: add logic
	return &points.Points{}, nil
}
