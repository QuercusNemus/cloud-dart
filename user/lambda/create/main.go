package main

import (
	"cloud-dart/user"
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"os"
)

var u *user.Service

func init() {
	region, ok := os.LookupEnv("DYNAMODB_AWS_REGION")
	if !ok {
		panic("DYNAMODB_AWS_REGION not set")
	}

	tableName, ok := os.LookupEnv("DYNAMODB_TABLE")
	if !ok {
		panic("DYNAMODB_TABLE not set")
	}

	u = user.NewService(tableName, region)
}

func HandleCreation(ctx context.Context, user user.User) (user.User, error) {
	fmt.Println("Create use invoked.")
	return u.Create(user)
}

func main() {
	lambda.Start(HandleCreation)
}
