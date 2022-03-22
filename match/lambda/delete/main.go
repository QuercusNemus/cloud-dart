package delete

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
	Match   match.Match `json:"Match"`
	Players []string    `json:"Players"`
}

type Input struct {
	Input MatchInput `json:"input"`
}

func handler(ctx context.Context, input Input) (match.Match, error) {
	return match.Match{}, nil
}

func main() {
	lambda.Start(handler)
}
