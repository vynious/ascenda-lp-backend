package db

import (
	"context"
	"fmt"
	types "github.com/vynious/ascenda-lp-backend/types"
	"gorm.io/gorm/clause"
)

// CreateTransaction creates a maker-checker transaction
func (dbs *DBService) CreateTransaction(ctx context.Context, action types.MakerAction, makerId, description string) (*types.Transaction, error) {

	tx := dbs.Conn.WithContext(ctx)

	txn := &types.Transaction{
		MakerId:     makerId,
		Description: description,
		Action:      action,
	}

	if err := tx.Create(&txn).Error; err != nil {
		return nil, fmt.Errorf(err.Error())
	}

	return txn, nil
}

func (dbs *DBService) GetCheckers(ctx context.Context, role string) ([]string, error) {
	var checkersEmail []string

	// find user's email based on maker checker roles mapping

	return checkersEmail, nil
}

func (dbs *DBService) UpdateTransaction(ctx context.Context, txnId string, checkerId string, approval bool) (*types.Transaction, error) {

	// Start a transaction
	tx := dbs.Conn.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	// Define the update map
	decision := map[string]interface{}{
		"CheckerId": checkerId,
		"Approval":  approval,
	}

	// Update the transaction and locks the current entry
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Model(&types.Transaction{}).Where("TransactionId = ?", txnId).Updates(decision).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Retrieve the updated transaction
	var updatedTransaction types.Transaction

	if err := tx.Where("TransactionId = ?", txnId).First(&updatedTransaction).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &updatedTransaction, nil
}
