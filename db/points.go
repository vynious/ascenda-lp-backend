package db

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/vynious/ascenda-lp-backend/types"
	"github.com/vynious/ascenda-lp-backend/util"
)

func (dbs *DB) GetPoints(ctx context.Context) ([]types.Points, error) {
	var pointsRecords []types.Points
	res := dbs.Conn.Find(&pointsRecords)
	if res.Error != nil {
		// Check if userLocation is part of the context
		userLocation, locationOk := ctx.Value("userLocation").(string)
		if locationOk {
			logEntry := types.Log{
				Type:   "Points",
				Action: "Failed to query points",
				// UserId:       ctx.Value("userId").(string),
				UserLocation: userLocation,
			}
			if err := util.CreateLogEntry(logEntry); err != nil {
				log.Printf("Error creating log entry: %v", err)
			}
		}
		return nil, fmt.Errorf("database error: %s", res.Error)
	}
	// Check if userLocation is part of the context
	userLocation, locationOk := ctx.Value("userLocation").(string)
	if locationOk {
		logEntry := types.Log{
			Type:   "Points",
			Action: "Queried all points",
			// UserId:       ctx.Value("userId").(string),
			UserLocation: userLocation,
		}
		if err := util.CreateLogEntry(logEntry); err != nil {
			log.Printf("Error creating log entry: %v", err)
		}
	}

	return pointsRecords, nil
}

func (dbs *DB) GetPointsAccountById(ctx context.Context, accId string) ([]types.Points, error) {

	var pointsRecords []types.Points
	res := dbs.Conn.Where("id = ?", accId).First(&pointsRecords)
	if res.Error != nil {
		// Check if userLocation is part of the context
		userLocation, locationOk := ctx.Value("userLocation").(string)
		if locationOk {
			logEntry := types.Log{
				Type:   "Points",
				Action: "Failed to connect to db when querying user point account by id",
				// UserId:       ctx.Value("userId").(string),
				UserLocation: userLocation,
			}
			if err := util.CreateLogEntry(logEntry); err != nil {
				log.Printf("Error creating log entry: %v", err)
			}
		}
		return nil, fmt.Errorf("database error %s", res.Error)
	}

	if res.RowsAffected == 0 {
		// Check if userLocation is part of the context
		userLocation, locationOk := ctx.Value("userLocation").(string)
		if locationOk {
			logEntry := types.Log{
				Type:   "Points",
				Action: fmt.Sprintf("Tried to query for a user that does not exist %s", accId),
				// UserId:       ctx.Value("userId").(string),
				UserLocation: userLocation,
			}
			if err := util.CreateLogEntry(logEntry); err != nil {
				log.Printf("Error creating log entry: %v", err)
			}
		}
		return nil, fmt.Errorf("account %s does not exist", accId)
	}
	// Check if userLocation is part of the context
	userLocation, locationOk := ctx.Value("userLocation").(string)
	if locationOk {
		logEntry := types.Log{
			Type:   "Points",
			Action: fmt.Sprintf("Successfully queried points for account id for %s", accId),
			// UserId:       ctx.Value("userId").(string),
			UserLocation: userLocation,
		}
		if err := util.CreateLogEntry(logEntry); err != nil {
			log.Printf("Error creating log entry: %v", err)
		}
	}
	return pointsRecords, nil
}

func (dbs *DB) GetPointsAccountsByUser(ctx context.Context, userId string) ([]types.Points, error) {

	log.Printf("Test %s", userId)
	var pointsRecords []types.Points
	res := dbs.Conn.Where("user_id = ?", userId).Find(&pointsRecords)
	if res.Error != nil {
		// Check if userLocation is part of the context
		userLocation, locationOk := ctx.Value("userLocation").(string)
		if locationOk {
			logEntry := types.Log{
				Type:   "Points",
				Action: fmt.Sprintf("Failed to connect to db when getting points account by user for %s", userId),
				// UserId:       ctx.Value("userId").(string),
				UserLocation: userLocation,
			}
			if err := util.CreateLogEntry(logEntry); err != nil {
				log.Printf("Error creating log entry: %v", err)
			}
		}
		return nil, fmt.Errorf("database error %s", res.Error)
	}

	if res.RowsAffected == 0 {
		// Check if userLocation is part of the context
		userLocation, locationOk := ctx.Value("userLocation").(string)
		if locationOk {
			logEntry := types.Log{
				Type:   "Points",
				Action: fmt.Sprintf("Tried querying an account that does not exist when getting points account by user for %s", userId),
				// UserId:       ctx.Value("userId").(string),
				UserLocation: userLocation,
			}
			if err := util.CreateLogEntry(logEntry); err != nil {
				log.Printf("Error creating log entry: %v", err)
			}
		}
		return nil, fmt.Errorf("user %s does not exist", userId)
	}
	// Check if userLocation is part of the context
	userLocation, locationOk := ctx.Value("userLocation").(string)
	if locationOk {
		logEntry := types.Log{
			Type:   "Points",
			Action: fmt.Sprintf("Queried getting points account by user for %s", userId),
			// UserId:       ctx.Value("userId").(string),
			UserLocation: userLocation,
		}
		if err := util.CreateLogEntry(logEntry); err != nil {
			log.Printf("Error creating log entry: %v", err)
		}
	}
	return pointsRecords, nil
}

func (dbs *DB) UpdatePoints(ctx context.Context, req types.UpdatePointsRequestBody) (*types.Points, error) {

	var pointsRecords []types.Points
	pointsRecords, err := dbs.GetPointsAccountById(ctx, req.ID)
	if err != nil {
		// Check if userLocation is part of the context
		userLocation, locationOk := ctx.Value("userLocation").(string)
		if locationOk {
			logEntry := types.Log{
				Type:   "Points",
				Action: fmt.Sprintf("Failed to update points with %+v", req),
				// UserId:       ctx.Value("userId").(string),
				UserLocation: userLocation,
			}
			if err := util.CreateLogEntry(logEntry); err != nil {
				log.Printf("Error creating log entry: %v", err)
			}
		}
		return nil, err
	}

	res := dbs.Conn.Model(pointsRecords).Update("balance", req.NewBalance).First(&pointsRecords)
	if res.RowsAffected == 0 {
		// Check if userLocation is part of the context
		userLocation, locationOk := ctx.Value("userLocation").(string)
		if locationOk {
			logEntry := types.Log{
				Type:   "Points",
				Action: fmt.Sprintf("No points updated with with %+v", req),
				// UserId:       ctx.Value("userId").(string),
				UserLocation: userLocation,
			}
			if err := util.CreateLogEntry(logEntry); err != nil {
				log.Printf("Error creating log entry: %v", err)
			}
		}
		return nil, res.Error
	}
	// Check if userLocation is part of the context
	userLocation, locationOk := ctx.Value("userLocation").(string)
	if locationOk {
		logEntry := types.Log{
			Type:   "Points",
			Action: fmt.Sprintf("Updated points with %+v", req),
			// UserId:       ctx.Value("userId").(string),
			UserLocation: userLocation,
		}
		if err := util.CreateLogEntry(logEntry); err != nil {
			log.Printf("Error creating log entry: %v", err)
		}
	}

	return &pointsRecords[0], nil
}

func (dbs *DB) CreatePointsAccount(ctx context.Context, req types.CreatePointsAccountRequestBody) (*types.Points, error) {

	pointsRecord := types.Points{
		ID:      uuid.NewString(),
		UserID:  *req.UserID,
		Balance: *req.NewBalance,
	}
	res := dbs.Conn.Model(types.Points{}).Create(&pointsRecord)
	if res.RowsAffected == 0 {
		// Check if userLocation is part of the context
		userLocation, locationOk := ctx.Value("userLocation").(string)
		if locationOk {
			logEntry := types.Log{
				Type:   "Points",
				Action: fmt.Sprintf("Failed to create points account with %+v", req),
				// UserId:       ctx.Value("userId").(string),
				UserLocation: userLocation,
			}
			if err := util.CreateLogEntry(logEntry); err != nil {
				log.Printf("Error creating log entry: %v", err)
			}
		}
		return nil, fmt.Errorf("database error %s", res.Error)
	}
	// Check if userLocation is part of the context
	userLocation, locationOk := ctx.Value("userLocation").(string)
	if locationOk {
		logEntry := types.Log{
			Type:   "Points",
			Action: fmt.Sprintf("Created points account with %+v", req),
			// UserId:       ctx.Value("userId").(string),
			UserLocation: userLocation,
		}
		if err := util.CreateLogEntry(logEntry); err != nil {
			log.Printf("Error creating log entry: %v", err)
		}
	}
	return &pointsRecord, nil
}

func (dbs *DB) DeletePointsAccountByUser(ctx context.Context, userId string) (bool, error) {
	res := dbs.Conn.Where("user_id = ?", &userId).Delete(&types.Points{})
	if res.RowsAffected == 0 {
		// Check if userLocation is part of the context
		userLocation, locationOk := ctx.Value("userLocation").(string)
		if locationOk {
			logEntry := types.Log{
				Type:   "Points",
				Action: fmt.Sprintf("Faled to delete points account for user %s", userId),
				// UserId:       ctx.Value("userId").(string),
				UserLocation: userLocation,
			}
			if err := util.CreateLogEntry(logEntry); err != nil {
				log.Printf("Error creating log entry: %v", err)
			}
		}
		return false, fmt.Errorf("database error %s", res.Error)
	}
	// Check if userLocation is part of the context
	userLocation, locationOk := ctx.Value("userLocation").(string)
	if locationOk {
		logEntry := types.Log{
			Type:   "Points",
			Action: fmt.Sprintf("Deleted points account for user %s", userId),
			// UserId:       ctx.Value("userId").(string),
			UserLocation: userLocation,
		}
		if err := util.CreateLogEntry(logEntry); err != nil {
			log.Printf("Error creating log entry: %v", err)
		}
	}

	return true, nil
}

func (dbs *DB) DeletePointsAccountByID(ctx context.Context, accId string) (bool, error) {
	res := dbs.Conn.Delete(&types.Points{}, &accId)
	if res.RowsAffected == 0 {
		// Check if userLocation is part of the context
		userLocation, locationOk := ctx.Value("userLocation").(string)
		if locationOk {
			logEntry := types.Log{
				Type:   "Points",
				Action: fmt.Sprintf("Faled to delete points account for id %s", accId),
				// UserId:       ctx.Value("userId").(string),
				UserLocation: userLocation,
			}
			if err := util.CreateLogEntry(logEntry); err != nil {
				log.Printf("Error creating log entry: %v", err)
			}
		}
		return false, fmt.Errorf("database error %v", res.Error)
	}
	// Check if userLocation is part of the context
	userLocation, locationOk := ctx.Value("userLocation").(string)
	if locationOk {
		logEntry := types.Log{
			Type:   "Points",
			Action: fmt.Sprintf("Deleted points account for id %s", accId),
			// UserId:       ctx.Value("userId").(string),
			UserLocation: userLocation,
		}
		if err := util.CreateLogEntry(logEntry); err != nil {
			log.Printf("Error creating log entry: %v", err)
		}
	}

	return true, nil
}
