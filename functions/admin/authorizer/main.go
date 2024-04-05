package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/uuid"
	"github.com/vynious/ascenda-lp-backend/db"
	"github.com/vynious/ascenda-lp-backend/types"
	"github.com/vynious/ascenda-lp-backend/util"
	"log"
)

var (
	DBService *db.DBService
	err       error
)

func init() {
	DBService, err = db.SpawnDBService()
	if err != nil {
		log.Fatalf(err.Error())
	}

}

func AuthorizerHandler(ctx context.Context, req events.APIGatewayCustomAuthorizerRequestTypeRequest) (events.APIGatewayCustomAuthorizerResponse, error) {

	/*
		Get the req.Arn, basically `arn:aws:execute-api:us-east-1:123456789012:a1b2c3d4e5/Prod/GET/points`.
		Check what kind of functionality user is accessing based of their RolePermissionList.
		So example if user has CanRead == true for resource == "points_ledger".
		Allow access to the resource.
	*/

	token := req.Headers["Authorization"]
	method := req.HTTPMethod
	route := req.Path[1:]

	roleName, err := util.GetRoleWithCognito(token)
	if err != nil {

	}
	var role types.Role
	role, err = db.RetrieveRoleWithRoleName(ctx, DBService, roleName)
	if err != nil {

	}
	var permissions types.RolePermissionList
	permissions = role.Permissions

	return GeneratePolicy(permissions, uuid.NewString(), route, method, req.MethodArn), nil
}

func GeneratePolicy(permissions []types.RolePermission, principalId, route, method, arn string) events.APIGatewayCustomAuthorizerResponse {
	authResponse := events.APIGatewayCustomAuthorizerResponse{
		PrincipalID: principalId,
	}
	authResponse.PolicyDocument = events.APIGatewayCustomAuthorizerPolicy{
		Version:   "2012-10-17",
		Statement: []events.IAMPolicyStatement{},
	}

	var resource string

	effect := "deny"

	if route == "user" || route == "users" {
		resource = "user_storage"
	} else if route == "points" {
		resource = "points_ledger"
	} else if route == "logs" || route == "log" {
		resource = "logs"
	} else if route == "maker-checker" {
		// no resource but need custom checker

	}

	for _, permission := range permissions {
		if permission.Resource == resource {
			switch method {
			case "GET":
				if permission.CanRead {
					effect = "allow"
				}
			case "PUT":
				if permission.CanUpdate {
					effect = "allow"
				}
			case "DELETE":
				if permission.CanDelete {
					effect = "allow"
				}
			case "POST":
				if permission.CanCreate {
					effect = "allow"
				}
			case "OPTIONS":
				log.Printf("options method going through")
			default:
				log.Printf("unchecked method made")
			}

			statement := generateStatement("execute-api:Invoke", effect, arn)
			authResponse.PolicyDocument.Statement = append(authResponse.PolicyDocument.Statement, statement)
		}
	}
	return authResponse
}

func generateStatement(action, effect, resource string) events.IAMPolicyStatement {
	return events.IAMPolicyStatement{
		Action:   []string{action},
		Effect:   effect,
		Resource: []string{resource},
	}
}

func main() {
	lambda.Start(AuthorizerHandler)
}
