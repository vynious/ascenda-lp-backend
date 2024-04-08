package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/joho/godotenv"
	"github.com/vynious/ascenda-lp-backend/db"
	aws_helpers "github.com/vynious/ascenda-lp-backend/functions/users/aws-helpers"
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

	clearDatabase(DB)

	// TODO: Add all models to be migrated here
	models := []interface{}{&types.Transaction{}, &types.Points{}, &types.User{}, &types.Role{}, &types.RolePermission{}, types.ApprovalChainMap{}}
	if err := DB.Conn.AutoMigrate(models...); err != nil {
		log.Fatalf("Failed to auto-migrate models")
	}
	log.Print("Successfully auto-migrated models")

	defer DB.CloseConn()

	seedRolesAndPermissions(DB)
	seedApprovalChainMap(DB)
	seedFile("users", DB)
	seedFile("points", DB)
	seedCustomUsers(DB)

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
	cognitoClient := aws_helpers.InitCognitoClient()

	var usersRecords []types.User
	for i, record := range records {
		if i == 0 {
			continue
		}

		var rolePtr *string = nil

		if record[4] != "" {
			role := record[4]
			rolePtr = &role
		}

		data := types.User{
			Id:        record[0],
			Email:     record[1],
			FirstName: record[2],
			LastName:  record[3],
			// if no role specified, customer role (no admin access)
			RoleName: rolePtr,
		}

		cognitoUserAttributes := []*cognitoidentityprovider.AttributeType{
			{
				Name:  aws.String("email"),
				Value: aws.String(data.Email),
			},
			{
				Name:  aws.String("email_verified"),
				Value: aws.String("true"),
			},
			{
				Name:  aws.String("custom:userID"),
				Value: aws.String(data.Id),
			},
		}

		if rolePtr != nil {
			cognitoUserAttributes = append(cognitoUserAttributes, &cognitoidentityprovider.AttributeType{
				Name:  aws.String("custom:role"),
				Value: aws.String(*rolePtr),
			})
		}

		cognitoInput := &cognitoidentityprovider.AdminCreateUserInput{
			ForceAliasCreation: aws.Bool(true),
			DesiredDeliveryMediums: []*string{
				aws.String("EMAIL"),
			},
			UserAttributes: cognitoUserAttributes,
			UserPoolId:     aws.String(os.Getenv("COGNITO_USER_POOL_ID")),
			Username:       aws.String(data.Email),
			MessageAction:  aws.String("SUPPRESS"),
		}
		log.Println("creating user in user pool")
		_, err := cognitoClient.AdminCreateUser(cognitoInput)
		if err != nil {
			log.Panicln(err)
		}
		log.Println("User created in user pool")

		usersRecords = append(usersRecords, data)
	}

	res := DB.Conn.CreateInBatches(usersRecords, batchsize)
	if res.Error != nil {
		log.Fatalf("Database error %s", res.Error)
	}
}

func clearDatabase(DB *db.DBService) {
	// Specify the order of deletion based on foreign key dependencies
	models := []interface{}{&types.RolePermission{}, &types.Transaction{}, &types.Points{}, types.ApprovalChainMap{}, &types.User{}, &types.Role{}, types.MakerAction{}}
	for _, model := range models {
		if result := DB.Conn.Unscoped().Where("1 = 1").Delete(model); result.Error != nil {
			log.Fatalf("Failed to clear table for model %v: %v", model, result.Error)
		}
	}
	log.Println("Successfully cleared the database")
}

func seedRolesAndPermissions(DB *db.DBService) {
	// Owner, Manager, Engineer, Product Manager
	var roles types.RoleList = types.RoleList{
		types.Role{
			RoleName: "owner",
			Permissions: types.RolePermissionList{
				types.RolePermission{
					Resource:  "user_storage",
					CanCreate: true,
					CanRead:   true,
					CanUpdate: true,
					CanDelete: true,
				},
				types.RolePermission{
					Resource:  "points_ledger",
					CanRead:   true,
					CanUpdate: true,
				},
				types.RolePermission{
					Resource: "logs",
					CanRead:  true,
				},
			},
		},
		types.Role{
			RoleName: "manager",
			Permissions: types.RolePermissionList{
				types.RolePermission{
					Resource:  "user_storage",
					CanCreate: true,
					CanRead:   true,
					CanUpdate: true,
				},
				types.RolePermission{
					Resource:  "points_ledger",
					CanRead:   true,
					CanUpdate: true,
				},
				types.RolePermission{
					Resource: "logs",
					CanRead:  true,
				},
			},
		},
		types.Role{
			RoleName: "engineer",
			Permissions: types.RolePermissionList{
				types.RolePermission{
					Resource: "user_storage",
					CanRead:  true,
				},
				types.RolePermission{
					Resource: "points_ledger",
					CanRead:  true,
				},
				types.RolePermission{
					Resource: "logs",
					CanRead:  true,
				},
			},
		},
		types.Role{
			RoleName: "product_manager",
			Permissions: types.RolePermissionList{
				types.RolePermission{
					Resource: "user_storage",
					CanRead:  true,
				},
				types.RolePermission{
					Resource: "points_ledger",
					CanRead:  true,
				},
			},
		},
	}
	for _, role := range roles {
		res := DB.Conn.Create(&role)
		if res.Error != nil {
			log.Fatalf("Error creating roles/permissions: %v", res.Error)
		}
	}
	log.Printf("Successful roles and perms seed")
}

func seedApprovalChainMap(DB *db.DBService) {
	var approvalChainMaps = []struct {
		MakerRoleName   string
		CheckerRoleName string
	}{
		{"product_manager", "owner"},
		{"engineer", "manager"},
		{"engineer", "owner"},
	}

	for _, acm := range approvalChainMaps {
		var makerRole, checkerRole types.Role

		// Find MakerRole and CheckerRole based on RoleName
		if err := DB.Conn.Where("role_name = ?", acm.MakerRoleName).First(&makerRole).Error; err != nil {
			log.Fatalf("Maker role not found: %s", acm.MakerRoleName)
		}

		if err := DB.Conn.Where("role_name = ?", acm.CheckerRoleName).First(&checkerRole).Error; err != nil {
			log.Fatalf("Checker role not found: %s", acm.CheckerRoleName)
		}

		newACM := types.ApprovalChainMap{
			MakerRoleID:   makerRole.Id,
			CheckerRoleID: checkerRole.Id,
		}

		// Create ApprovalChainMap entry
		res := DB.Conn.Create(&newACM)
		if res.Error != nil {
			log.Fatalf("Error creating approval chain map: %v", res.Error)
		}
	}
}

// SeedCustomUser creates a user with a specified role
func seedCustomUsers(DB *db.DBService) {
	cognitoClient := aws_helpers.InitCognitoClient()
	// Define users
	users := []types.User{
		{
			Id:        "123-456-789",
			Email:     "shawn.thiah.2022@scis.smu.edu.sg",
			FirstName: "shawn",
			LastName:  "thiah",
			RoleID:    getRoleID(DB, "product_manager"),
			RoleName:  aws.String("product_manager"),
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
		},
		{
			Id:        "234-567-890",
			Email:     "jingjie.lim.2022@scis.smu.edu.sg",
			FirstName: "jingjie",
			LastName:  "lim",
			RoleID:    getRoleID(DB, "owner"),
			RoleName:  aws.String("owner"),
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
		},
	}

	for _, user := range users {
		cognitoUserAttributes := []*cognitoidentityprovider.AttributeType{
			{
				Name:  aws.String("email"),
				Value: aws.String(user.Email),
			},
			{
				Name:  aws.String("email_verified"),
				Value: aws.String("true"),
			},
			{
				Name:  aws.String("custom:userID"),
				Value: aws.String(user.Id),
			},
			{
				Name:  aws.String("custom:role"),
				Value: aws.String(*user.RoleName),
			},
		}

		cognitoInput := &cognitoidentityprovider.AdminCreateUserInput{
			ForceAliasCreation: aws.Bool(true),
			DesiredDeliveryMediums: []*string{
				aws.String("EMAIL"),
			},
			UserAttributes: cognitoUserAttributes,
			UserPoolId:     aws.String(os.Getenv("COGNITO_USER_POOL_ID")),
			Username:       aws.String(user.Email),
			MessageAction:  aws.String("SUPPRESS"),
		}
		log.Println("creating user in user pool")

		_, err := cognitoClient.AdminCreateUser(cognitoInput)
		if err != nil {
			log.Panicln(err)
		}
		log.Println("User created in user pool")

		log.Printf("adding users to rds")
		if err := DB.Conn.Create(&user).Error; err != nil {
			log.Fatalf("Failed to create user: %v", err)
		}
	}
}

func getRoleID(DB *db.DBService, roleName string) *uint {
	var role types.Role
	if err := DB.Conn.Where("role_name = ?", roleName).First(&role).Error; err != nil {
		log.Fatalf("Role not found: %v", err)
		return nil
	}
	return &role.Id
}
