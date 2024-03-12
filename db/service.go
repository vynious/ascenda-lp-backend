package db

import (
	"context"
	"database/sql"
	"fmt"
	makerchecker "github.com/vynious/ascenda-lp-backend/types/maker-checker"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"time"
)

type IDBService interface {
	CreateTransaction()
	UpdateTransaction()
	GetMakers()
	GetCheckers()
}

type DBService struct {
	conn    *sql.DB
	timeout time.Duration
}

func SpawnDBService() (*DBService, error) {
	username := os.Getenv("db_user")
	password := os.Getenv("db_pass")
	dbName := os.Getenv("db_name")
	dbHost := os.Getenv("db_host")
	dsn := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, username, dbName, password)
	cc, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to make connection")
	}
	scc, err := cc.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to make connection")
	}
	return &DBService{
		conn: scc,
	}, nil
}

// CreateTransaction creates a maker-checker transaction
func (dbs *DBService) CreateTransaction(ctx context.Context, req *makerchecker.CreateTransactionBody) (*makerchecker.Transaction, error) {
	// todo: add logic
	return &makerchecker.Transaction{}, nil
}

// CloseConn closes connection to db
func (dbs *DBService) CloseConn() error {
	if err := dbs.conn.Close(); err != nil {
		return fmt.Errorf("failed to close connection")
	}
	return nil
}
