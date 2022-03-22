package main

import (
	"cloud-dart/player"
	"context"
	"github.com/aws/aws-lambda-go/lambda"
)

var playerService *player.Service

func init() {
	playerService = player.NewService("Players", "eu-north-1")
}

type PlayerInput struct {
	Arguments struct {
		Player player.Player `json:"player"`
	} `json:"arguments"`
}

func handler(ctx context.Context, request PlayerInput) (bool, error) {
	err := playerService.Create(request.Arguments.Player)
	if err != nil {
		return false, err
	}
	return true, nil
}

func main() {
	lambda.Start(handler)
}
