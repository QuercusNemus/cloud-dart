package main

import (
	"cloud-dart/player"
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
)

var playerService *player.Service

func init() {
	playerService = player.NewService("Player", "eu-north-1")
}

func handler(ctx context.Context, player player.Player) (player.Player, error) {
	fmt.Println("CreatePlayer invoked!")
	return playerService.Create(player)
}

func main() {
	lambda.Start(handler)
}
