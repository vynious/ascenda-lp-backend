package db

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/vynious/ascenda-lp-backend/types"
	"gorm.io/gorm/clause"
)

// CreateTransaction creates a maker-checker transaction
func (dbs *DBService) CreateTransaction(ctx context.Context, action types.MakerAction, makerId string) (*types.Transaction, error) {

	tx := dbs.Conn.WithContext(ctx)

	jsonMsgAction, _ := json.Marshal(action)

	txn := &types.Transaction{
		TransactionId: uuid.NewString(),
		MakerId:       makerId,
		Action:        jsonMsgAction,
	}

	if err := tx.Create(&txn).Error; err != nil {
		return nil, err
	}

	return txn, nil
}

func (dbs *DBService) GetTransaction(ctx context.Context, txnId string) (*types.Transaction, error) {
	var transaction types.Transaction

	tx := dbs.Conn.WithContext(ctx)

	if err := tx.Where("transaction_id = ?", txnId).First(&transaction).Error; err != nil {
		return nil, err
	}

	return &transaction, nil
}

func (dbs *DBService) UpdateTransaction(ctx context.Context, txnId string, checkerId string, approval bool) (*types.Transaction, error) {


	tx := dbs.Conn.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	// Define the update map
	decision := map[string]interface{}{
		"CheckerId": checkerId,
		"Approval":  approval,
		"Status":    "completed",
	}

	// Update the transaction and locks the current entry
	if err := tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Model(&types.Transaction{}).
			Where("transaction_id = ?", txnId).
			Updates(decision).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Retrieve the updated transaction
	var updatedTransaction types.Transaction

	if err := tx.Where("transaction_id = ?", txnId).First(&updatedTransaction).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Commit the transaction

	if approval {
		// if cannot process the transaction, rollback updates for maker-checker transaction
		var makerAction types.MakerAction
		if err := json.Unmarshal(updatedTransaction.Action, &makerAction); err != nil {
			tx.Rollback()
			return nil, err
		}
		if err := dbs.ProcessTransaction(&makerAction); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
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


	if err := tx.
			Model(&types.User{}).
			Where("role IN ?", checkerRole).
			Pluck("Email", &checkersEmail).Error; err != nil {
		return nil, err
	}
	return checkersEmail, nil
}

func (dbs *DBService) ProcessTransaction(action *types.MakerAction) error {

	switch action.ActionType {
	case "UpdatePoints":
		var updatePointsRequestBody types.UpdatePointsRequestBody
		if err := json.Unmarshal(action.RequestBody, &updatePointsRequestBody); err != nil {
			return err
		}
		if _, err := dbs.UpdatePoints(context.Background(), updatePointsRequestBody); err != nil {
			return err
		}
	case "UpdateUser":

	}
	return nil
}
