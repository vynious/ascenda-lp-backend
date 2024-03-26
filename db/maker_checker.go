package db

import (
	"context"
	"fmt"
	"github.com/vynious/ascenda-lp-backend/types"
	"gorm.io/gorm/clause"
)

// CreateTransaction creates a maker-checker transaction
func (dbs *DBService) CreateTransaction(ctx context.Context, action types.MakerAction, makerId string) (*types.Transaction, error) {

	tx := dbs.Conn.WithContext(ctx)

	txn := &types.Transaction{
		MakerId: makerId,
		Action:  action,
	}

	if err := tx.Create(&txn).Error; err != nil {
		return nil, err
	}

	return txn, nil
}

func (dbs *DBService) GetTransaction(ctx context.Context, txnId string) (*types.Transaction, error) {
	var transaction types.Transaction

	tx := dbs.Conn.WithContext(ctx)

	if err := tx.Where("TransactionId = ?", txnId).First(&transaction).Error; err != nil {
		return nil, err
	}

	return &transaction, nil
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

	if approval == true {
		go func() {
			if err := dbs.ProcessTransaction(&updatedTransaction.Action); err != nil {
			}
		}()

	}

	return &updatedTransaction, nil
}

func (dbs *DBService) GetCheckers(ctx context.Context, makerRole string) ([]string, error) {
	var checkersEmail []string // need to convert to aws.String() ??

	// mapping
	roleMap := map[string][]string{
		"product_owner": {"owner"},
		"engineer":      {"manager", "owner"},
	}

	checkerRole := roleMap[makerRole]

	// find user's email based on maker checker roles mapping
	tx := dbs.Conn.WithContext(ctx).Begin()

	// SELECT email FROM users WHERE role_id IN (1, 2, 3);
	if err := tx.Model(&types.User{}).Where("Role IN ?", checkerRole).Pluck("Email", checkersEmail).Error; err != nil {
		return nil, err
	}
	return checkersEmail, nil
}

func (dbs *DBService) ProcessTransaction(action *types.MakerAction) error {

	switch action.ActionType {
	case "UpdatePoints":

		requestBody, ok := action.RequestBody.(types.UpdatePointsRequestBody)
		if !ok {
			return fmt.Errorf("request body for update points does not match")
		}

	case "UpdateUser":

	}
	return nil
}
