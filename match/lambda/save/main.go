package main

import (
	"cloud-dart/match"
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"os"
)

var matchService *match.Service

func init() {
	region, ok := os.LookupEnv("DYNAMODB_AWS_REGION")
	if !ok {
		panic("DYNAMODB_AWS_REGION not set")
	}

	tableName, ok := os.LookupEnv("DYNAMODB_TABLE")
	if !ok {
		panic("DYNAMODB_TABLE not set")
	}
	matchService = match.NewService(tableName, region)
}

type MatchInput struct {
	Arguments struct {
		Input struct {
			Match   match.Match `json:"Match"`
			Players []string    `json:"Players"`
		} `json:"input"`
	} `json:"arguments"`
}

func handler(ctx context.Context, request MatchInput) (match.Match, error) {
	return matchService.Save(request.Arguments.Input.Match, request.Arguments.Input.Players)
}

func main() {
	lambda.Start(handler)
}
