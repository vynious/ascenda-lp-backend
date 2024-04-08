package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/vynious/ascenda-lp-backend/types"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type IDBService interface {
	CreateTransaction(ctx context.Context, action types.MakerAction, makerId, description string) (*types.Transaction, error)
	UpdateTransaction(ctx context.Context, txnId string, checkerId string, approval bool) (*types.Transaction, error)
	GetCheckers(ctx context.Context, makerId string, role string) ([]string, error)
	GetPoints(ctx context.Context, userId string) ([]string, error)
}

type DB struct {
	Conn    *gorm.DB
	timeout time.Duration
}

type DBService struct {
	ConnMap map[string]*DB
	timeout time.Duration
}

func SpawnDBService() (*DBService, error) {
	dbService := &DBService{
		ConnMap: make(map[string]*DB),
	}

	// Connect to banks and populate the ConnMap
	bankNames := []string{"ascenda", "bankb"} // Add more banks as needed
	for _, bank := range bankNames {
		conn, err := connectToDB(bank)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to %s: %v", bank, err)
		}
		dbService.ConnMap[bank] = &DB{Conn: conn}
	}

	log.Println("Successfully connected to all databases")

	return dbService, nil
}

func connectToDB(bank string) (*gorm.DB, error) {
	dbUser := os.Getenv("dbUser")
	// dbName := os.Getenv("dbName")
	dbHost := os.Getenv("dbHost")
	dbPwd := os.Getenv("dbPwd")
	dbName := bank
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
		dbHost, 5432, dbUser, dbPwd, dbName,
	)
	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to make connection")
	}
	log.Printf("Successfully connected to Database")

	return conn, nil
}

// CloseConn closes connection to db
func (DBService *DBService) CloseConnections() error {
	for _, conn := range DBService.ConnMap {
		db, _ := conn.Conn.DB()
		if err := db.Close(); err != nil {
			return fmt.Errorf("failed to close database connection: %v", err)
		}
	}
	return nil
}
