package main

import (
	"cloud-dart/integration"
	"cloud-dart/match"
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
)

// Input Struct for direct lambda save
type Input struct {
	Arguments struct {
		Match match.Match `json:"match"`
	} `json:"arguments"`
}

type DynamoService struct {
	db integration.DynamoDB
}

func (s *DynamoService) HandleRequest(ctx context.Context, input Input) error {
	err := s.db.Write(ctx, &input.Arguments.Match)
	if err != nil {
		return err
	}
	return nil
}

// main lambda function
func main() {
	defaultConfig, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		panic(err)
	}

	db, err := integration.New(defaultConfig)
	if err != nil {
		panic(err)
	}

	ds := DynamoService{
		db: *db,
	}

	lambda.Start(ds.HandleRequest)
}
