package db

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/vynious/ascenda-lp-backend/types"
)

func (dbs *DBService) GetPoints(ctx context.Context) ([]types.Points, error) {
	var pointsRecords []types.Points
	res := dbs.Conn.Find(&pointsRecords)
	if res.Error != nil {
		return nil, fmt.Errorf("database error %s", res.Error)
	}

	return pointsRecords, nil
}

func (dbs *DBService) GetPointsByID(ctx context.Context, accId string) ([]types.Points, error) {

	var pointsRecords []types.Points
	res := dbs.Conn.Where("id = ?", accId).First(&pointsRecords)
	if res.Error != nil {
		return nil, fmt.Errorf("database error %s", res.Error)
	}

	if res.RowsAffected == 0 {
		return nil, fmt.Errorf("account %s does not exist", accId)
	}

	return pointsRecords, nil
}

func (dbs *DBService) GetPointsAccountsByUser(ctx context.Context, userId string) ([]types.Points, error) {

	log.Printf("Test %s", userId)
	var pointsRecords []types.Points
	res := dbs.Conn.Where("user_id = ?", userId).Find(&pointsRecords)
	if res.Error != nil {
		return nil, fmt.Errorf("database error %s", res.Error)
	}

	if res.RowsAffected == 0 {
		return nil, fmt.Errorf("user %s does not exist", userId)
	}

	return pointsRecords, nil
}

func (dbs *DBService) UpdatePoints(ctx context.Context, accId string, newBalance int32) (*types.Points, error) {

	var pointsRecords []types.Points
	pointsRecords, err := dbs.GetPointsByID(ctx, accId)
	if err != nil {
		return nil, err
	}

	res := dbs.Conn.Model(pointsRecords).Update("balance", newBalance).First(&pointsRecords)
	if res.RowsAffected == 0 {
		return nil, res.Error
	}

	return &pointsRecords[0], nil
}

func (dbs *DBService) CreatePointsAccount(ctx context.Context, userId string, newBalance int32) (*types.Points, error) {

	pointsRecord := types.Points{
		ID:      uuid.NewString(),
		UserID:  userId,
		Balance: newBalance,
	}
	res := dbs.Conn.Model(types.Points{}).Create(&pointsRecord)
	if res.RowsAffected == 0 {
		return nil, fmt.Errorf("database error %s", res.Error)
	}

	return &pointsRecord, nil
}

func (dbs *DBService) DeletePointsAccountByUser(ctx context.Context, userId string) (bool, error) {

	res := dbs.Conn.Where("user_id = ?", &userId).Delete(&types.Points{})
	if res.RowsAffected == 0 {
		return false, fmt.Errorf("database error %s", res.Error)
	}

	return true, nil
}

func (dbs *DBService) DeletePointsAccountByID(ctx context.Context, accId string) (bool, error) {
	res := dbs.Conn.Delete(&types.Points{}, &accId)
	if res.RowsAffected == 0 {
		return false, fmt.Errorf("database error %s", res.Error)
	}

	return true, nil
}
