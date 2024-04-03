package db

import (
	"context"
	"encoding/json"
	"fmt"

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

func (dbs *DBService) GetTransactions(ctx context.Context) (*[]types.Transaction, error) {
	var transactions []types.Transaction

	tx := dbs.Conn.WithContext(ctx)
	if result := tx.Find(&transactions); result.Error != nil {
		return nil, fmt.Errorf("failed to get all transactions: %v", result.Error.Error())
	}
	return &transactions, nil
}

// GetTransactionsByMakerIdByStatus Gets the transactions by maker_id and the status
func (dbs *DBService) GetTransactionsByMakerIdByStatus(ctx context.Context, makerId string, status string) (*[]types.Transaction, error) {
	var transactions []types.Transaction

	tx := dbs.Conn.WithContext(ctx)
	if result := tx.
		Where("maker_id = ?", makerId).
		Where("status = ?", status).
		Find(&transactions); result.Error != nil {
		return nil, fmt.Errorf("failed to get all transactions by maker_id: %v", result.Error.Error())
	}
	return &transactions, nil
}

//func (dbs *DBService) GetPendingTransactionsByApprovalChain(ctx context.Context, checkerRole string) (*[]types.Transaction, error) {
//	/*
//			get all pending transaction
//			from each pending transaction get the checkerid
//			from the checkerid get the checkerrole
//			based off the checkerrole, check with mapping if mapping checkerrole: [makerrole1, makerrole2]
//			if makerrole is inside the map
//			return the transactions
//
//		- select * from
//
//
//
//
//
//		- available checker roles
//		select * from makercheckermap
//		where makerrole = <makerrole>
//
//
//		- pending transactions
//		select checkerid from transactions
//		where status = <status>
//
//	*/
//	var transactions []types.Transaction
//
//	tx := dbs.Conn.WithContext(ctx)
//	result := tx.
//		Where("status = pending")
//
//	return nil, nil
//}

// GetCompletedTransactionsByCheckerId This function assumes that all transactions with a value checker_id has been completed.
func (dbs *DBService) GetCompletedTransactionsByCheckerId(ctx context.Context, checkerId string) (*[]types.Transaction, error) {
	var transactions []types.Transaction

	tx := dbs.Conn.WithContext(ctx)

	if result := tx.Where("checker_id = ?", checkerId).Find(&transactions); result.Error != nil {
		return nil, fmt.Errorf("failed to get completed transactions by checkerid: %v", result.Error.Error())
	}

	return &transactions, nil
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
		// if you cannot process the transaction, rollback updates for maker-checker transaction
		var makerAction types.MakerAction
		if err := json.Unmarshal(updatedTransaction.Action, &makerAction); err != nil {
			tx.Rollback()
			return nil, err
		}
		if err := dbs.ProcessTransaction(ctx, &makerAction); err != nil {
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

func (dbs *DBService) ProcessTransaction(ctx context.Context, action *types.MakerAction) error {
	switch action.ActionType {
	case "UpdatePoints":
		var updatePointsRequestBody types.UpdatePointsRequestBody
		if err := json.Unmarshal(action.RequestBody, &updatePointsRequestBody); err != nil {
			return err
		}
		if _, err := dbs.UpdatePoints(ctx, updatePointsRequestBody); err != nil {
			return err
		}
	case "UpdateUser":
		var updateUserRequestBody types.UpdateUserRequestBody
		if err := json.Unmarshal(action.RequestBody, &updateUserRequestBody); err != nil {
			return err
		}
		// calls update user.
	case "":

	default:
		return fmt.Errorf("action method does not exist")
	}
	return nil
}
