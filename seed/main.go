package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/vynious/ascenda-lp-backend/db"
	"github.com/vynious/ascenda-lp-backend/types"
)

var (
	DB        *db.DBService
	batchsize = 100
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to load .env")
	}

	var DB, err = db.SpawnDBService()
	if err != nil {
		log.Fatalf("Error spawn DB service...")
	}

	// TODO: Add all models to be migrated here
	models := []interface{}{&types.MakerAction{}, &types.Points{}, &types.User{}}
	if err := DB.Conn.AutoMigrate(models...); err != nil {
		log.Fatalf("Failed to auto-migrate models")
	}
	log.Print("Successfully auto-migrated models")

	defer DB.CloseConn()

	seedFile("users", DB)
	seedFile("points", DB)
}

func seedFile(filename string, DB *db.DBService) {
	file, err := os.Open(fmt.Sprintf("./seed/data/%s.csv", filename))
	if err != nil {
		log.Fatalf("Error opening %s.csv: %v", filename, err)
	}
	defer file.Close()

	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		log.Fatalf("Error reading %s.csv: %v", filename, err)
	}

	switch filename {
	case "points":
		seedPoints(records, DB)
	case "users":
		seedUsers(records, DB)
	default:
		log.Fatalf("Unsupported file type: %s", filename)
	}
}

func seedPoints(records [][]string, DB *db.DBService) {
	var pointsRecords []types.Points
	for i, record := range records {
		if i == 0 {
			continue
		}
		balance, _ := strconv.Atoi(record[2]) // convert to int
		data := types.Points{
			ID:      record[0],
			UserID:  record[1],
			Balance: int32(balance),
		}
		pointsRecords = append(pointsRecords, data)
	}

	res := DB.Conn.CreateInBatches(pointsRecords, batchsize)
	if res.Error != nil {
		log.Fatalf("Database error %s", res.Error)
	}
}

func seedUsers(records [][]string, DB *db.DBService) {
	var usersRecords []types.User
	for i, record := range records {
		if i == 0 {
			continue
		}
		data := types.User{
			ID:        record[0],
			Email:     record[1],
			FirstName: record[2],
			LastName:  record[3],
			Role:      record[4],
		}
		usersRecords = append(usersRecords, data)
	}

	res := DB.Conn.CreateInBatches(usersRecords, batchsize)
	if res.Error != nil {
		log.Fatalf("Database error %s", res.Error)
	}
}
