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

func handler(ctx context.Context, playerId string) ([]match.MatchPlayer, error) {
	return matchService.GetByPlayerId(playerId)
}

func main() {
	lambda.Start(handler)
}
