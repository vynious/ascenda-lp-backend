package db

import (
	"context"

	"github.com/vynious/ascenda-lp-backend/types/points"
)

func (dbs *DBService) GetPoints(ctx context.Context, txnId string, checkerId string, approval bool) (*points.Points, error) {
	// todo: add logic
	return &points.Points{}, nil
}
