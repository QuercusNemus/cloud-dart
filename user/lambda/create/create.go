package main

import (
	"cloud-dart/user"
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
)

var u *user.Service

func init() {
	u = user.NewService("User", "eu-north-1")
}

func HandleCreation(ctx context.Context, user user.User) (user.User, error) {
	fmt.Println("Create use invoked.")
	return u.Create(user)
}

func main() {
	lambda.Start(HandleCreation)
}
