package main

import (
	"cloud-dart/match"
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"os"
)

const (
	ContentTypeCreate = "CREATE"
	ContentTypeGet    = "GET"
	ContentTypeDelete = "DELETE"
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
	ContentType string      `json:"ContentType"`
	Match       match.Match `json:"Match"`
	Players     []string    `json:"Players"`
}

type Input struct {
	Input MatchInput `json:"input"`
}

func handler(ctx context.Context, input Input) (match.Match, error) {
	var err error

	switch input.ContentType {
	case ContentTypeCreate:
		return matchService.Save(input.Match, input.Players)
	case ContentTypeGet:
		return matchService.Get(input.Match)
	case ContentTypeDelete:
		return matchService.Delete(input.Match)
	default:
		err = fmt.Errorf("unable to recognise content-type: %v", input.Input.ContentType)
	}

	return match.Match{}, err
}

func main() {
	lambda.Start(handler)
}
