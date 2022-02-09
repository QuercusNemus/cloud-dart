package main

import (
	"cloud-dart/match"
	"context"
	"github.com/aws/aws-lambda-go/lambda"
)

var matchService *match.Service

func init() {
	matchService = match.NewService("Matches", "eu-north-1")
}

type Input struct {
	Match   match.Match
	Players []string
}

func handler(ctx context.Context, input Input) (match.Match, error) {
	return matchService.Save(input.Match, input.Players)
}

func main() {
	lambda.Start(handler)
}
