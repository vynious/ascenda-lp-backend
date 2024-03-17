package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	makerchecker "github.com/vynious/ascenda-lp-backend/types/maker-checker"
	"github.com/vynious/ascenda-lp-backend/types/points"
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
	conn    *sql.DB
	timeout time.Duration
}

func SpawnDBService() (*DBService, error) {
	secretName := "itsa-g1t2/rds"
	region := "ap-southeast-1"

	config, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		log.Fatal(err)
	}

	// Create Secrets Manager client
	svc := secretsmanager.NewFromConfig(config)

	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
	}

	result, err := svc.GetSecretValue(context.TODO(), input)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Decrypts secret using the associated KMS key.
	secretString := *result.SecretString

	username := os.Getenv("db_user")
	// password := os.Getenv("db_pass")
	dbName := os.Getenv("db_name")
	dbHost := os.Getenv("db_host")
	dsn := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, username, dbName, secretString)
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
func (dbs *DBService) CreateTransaction(ctx context.Context, action makerchecker.MakerAction, makerId, description string) (*makerchecker.Transaction, error) {
	// todo: add logic
	return &makerchecker.Transaction{}, nil
}

func (dbs *DBService) GetCheckers(ctx context.Context, makerId string, role string) ([]string, error) {
	var checkersEmail []string
	// todo: add logic
	return checkersEmail, nil
}

func (dbs *DBService) UpdateTransaction(ctx context.Context, txnId string, checkerId string, approval bool) (*makerchecker.Transaction, error) {
	// todo: add logic
	return &makerchecker.Transaction{}, nil
}

func (dbs *DBService) GetPoints(ctx context.Context) (*points.Points, error) {
	// todo: add logic
	return &points.Points{}, nil
}

// CloseConn closes connection to db
func (dbs *DBService) CloseConn() error {
	if err := dbs.conn.Close(); err != nil {
		return fmt.Errorf("failed to close connection")
	}
	return nil
}
