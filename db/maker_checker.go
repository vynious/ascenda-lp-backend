package db

import (
	"context"
	"fmt"

	"github.com/vynious/ascenda-lp-backend/types"
)

// CreateTransaction creates a maker-checker transaction
func (dbs *DBService) CreateTransaction(ctx context.Context, action types.MakerAction, makerId, description string) (*types.Transaction, error) {
	// todo: add logic
	tx := dbs.Conn.WithContext(ctx)

	txn := &types.Transaction{
		MakerId:     makerId,
		Description: description,
		Action:      action,
	}

	if err := tx.Create(&txn).Error; err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	return &types.Transaction{}, nil
}

func (dbs *DBService) GetCheckers(ctx context.Context, makerId string, role string) ([]string, error) {
	var checkersEmail []string
	// todo: add logic
	return checkersEmail, nil
}

func (dbs *DBService) UpdateTransaction(ctx context.Context, txnId string, checkerId string, approval bool) (*types.Transaction, error) {
	// todo: add logic
	return &types.Transaction{}, nil
}
