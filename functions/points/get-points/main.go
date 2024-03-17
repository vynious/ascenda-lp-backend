package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/vynious/ascenda-lp-backend/db"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DBService *db.DBService
	RDSClient *rds.Client
	err       error
)

func init() {
	log.Printf("INIT")

	// Initialise global variable DBService tied to Aurora
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		fmt.Println("Couldn't load default configuration. Have you set up your AWS account?")
		fmt.Println(err)
		return
	}
	log.Printf(cfg.AppID)

	dbUser := "postgres"
	password := ""
	dbName := "ascenda"
	dbHost := ""
	// dbPort := 5432
	// region := "ap-southeast-1"

	// authenticationToken, err := auth.BuildAuthToken(
	// 	context.TODO(), dbHost, region, dbUser, cfg.Credentials)
	// if err != nil {
	// 	panic("failed to create authentication token: " + err.Error())
	// }

	// dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
	// 	dbHost, dbPort, dbUser, authenticationToken, dbName,
	// )
	// conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	// log.Printf(conn.Name())

	// TEST GORM

	dsn := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, dbUser, dbName, password)
	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	// conn, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to RDS")
	}
	log.Printf("Connected to DB %s", conn.Name())
	// log.Printf("%s", conn.Stats())
}

func main() {
	// we are simulating a lambda behind an ApiGatewayV2
	lambda.Start(handler)
}

func handler(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return events.APIGatewayV2HTTPResponse{
		StatusCode: 200,
		Body:       "Hello",
	}, nil
}
