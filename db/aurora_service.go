package db

import (
	"context"
	"fmt"
	makerchecker "github.com/vynious/ascenda-lp-backend/types/maker-checker"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"time"
)

type IDBService interface {

	// Maker-checker functions
	CreateTransaction(ctx context.Context, action makerchecker.MakerAction, makerId, description string) (*Transaction, error)
	UpdateTransaction(ctx context.Context, txnId string, checkerId string, approval bool) (*Transaction, error)
	GetCheckers(ctx context.Context, makerId string, role string) ([]string, error)
}

type DBService struct {
	conn    *gorm.DB
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

	// todo:  Add on schema for AuroraDB
	if err := cc.AutoMigrate(&Transaction{}, &User{}, &MakerChecker{}, &Role{}); err != nil {
		return nil, fmt.Errorf("failed to migrate schema")
	}
	return &DBService{
		conn: cc,
	}, nil
}

// CreateTransaction creates a maker-checker transaction
func (dbs *DBService) CreateTransaction(ctx context.Context, action makerchecker.MakerAction, makerId, description string) (*Transaction, error) {

	tx := dbs.conn.WithContext(ctx)

	txn := &Transaction{
		MakerId:     makerId,
		Description: description,
		Action:      action,
	}

	if err := tx.Create(&txn).Error; err != nil {
		return nil, fmt.Errorf(err.Error())
	}

	return txn, nil
}

func (dbs *DBService) GetCheckers(ctx context.Context, makerId string, role string) ([]string, error) {
	var checkersEmail []string

	//tx := dbs.conn.WithContext(ctx)

	// todo: add logic
	return checkersEmail, nil
}

func (dbs *DBService) UpdateTransaction(ctx context.Context, txnId string, checkerId string, approval bool) (*Transaction, error) {

	//tx := dbs.conn.WithContext(ctx)
	return &Transaction{}, nil
}

// CloseConn closes connection to db
func (dbs *DBService) CloseConn() error {
	//if err := dbs.conn.DB().Close(); err != nil {
	//	return fmt.Errorf("failed to close connection")
	//}
	return nil
}
