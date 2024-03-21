package db

import (
	"context"
	"fmt"
	makerchecker "github.com/vynious/ascenda-lp-backend/types"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type IDBService interface {
	CreateTransaction(ctx context.Context, action makerchecker.MakerAction, makerId, description string) (*makerchecker.Transaction, error)
	UpdateTransaction(ctx context.Context, txnId string, checkerId string, approval bool) (*makerchecker.Transaction, error)
	GetCheckers(ctx context.Context, makerId string, role string) ([]string, error)
	GetPoints(ctx context.Context, userId string) ([]string, error)
}

type DBService struct {
	Conn    *gorm.DB
	timeout time.Duration
}

func SpawnDBService() (*DBService, error) {
	dbUser := os.Getenv("dbUser")
	dbName := os.Getenv("dbName")
	dbHost := os.Getenv("dbHost")
	dbPwd := os.Getenv("dbPwd")
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
		dbHost, 5432, dbUser, dbPwd, dbName,
	)
	cc, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to make connection")
	}
	log.Printf("Successfully connected to Database")
	return &DBService{
		Conn: cc,
	}, nil
}

// CloseConn closes connection to db
func (dbs *DBService) CloseConn() error {
	db, _ := dbs.Conn.DB()
	if err := db.Close(); err != nil {
		return fmt.Errorf("failed to close connection")
	}
	return nil
}
