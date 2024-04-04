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

const (
	authorized = false
)

func init() {
	DBService, err = db.SpawnDBService()
	if err != nil {
		log.Fatalf(err.Error())
	}

}

func AuthorizerHandler(ctx context.Context, req events.APIGatewayCustomAuthorizerRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
	token := req.AuthorizationToken

	roleName, err := util.GetRoleWithCognito(token)
	if err != nil {

	}
	var role types.Role
	role, err = db.RetrieveRoleWithRoleName(ctx, DBService, roleName)
	if err != nil {

	}
	var permissions types.RolePermissionList
	permissions = role.Permissions

	return generatePolicy(permissions, uuid.NewString(), req.MethodArn), nil
}

func generatePolicy(permissions []types.RolePermission, principalId, resource string) events.APIGatewayCustomAuthorizerResponse {
	authResponse := events.APIGatewayCustomAuthorizerResponse{PrincipalID: principalId}

	authResponse.PolicyDocument = events.APIGatewayCustomAuthorizerPolicy{
		Version:   "2012-10-17",
		Statement: []events.IAMPolicyStatement{},
	}

	for _, permission := range permissions {
		// For each permission, check if the operation is allowed and add the corresponding policy statement
		if permission.Resource == resource {
			if permission.CanRead {
				statement := generateStatement("execute-api:Invoke", "Allow", resource)
				authResponse.PolicyDocument.Statement = append(authResponse.PolicyDocument.Statement, statement)
			}
			if permission.CanCreate {
				statement := generateStatement("execute-api:Invoke", "Allow", resource)
				authResponse.PolicyDocument.Statement = append(authResponse.PolicyDocument.Statement, statement)
			}
			if permission.CanUpdate {
				statement := generateStatement("execute-api:Invoke", "Allow", resource)
				authResponse.PolicyDocument.Statement = append(authResponse.PolicyDocument.Statement, statement)
			}
			if permission.CanDelete {
				statement := generateStatement("execute-api:Invoke", "Allow", resource)
				authResponse.PolicyDocument.Statement = append(authResponse.PolicyDocument.Statement, statement)
			}
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
