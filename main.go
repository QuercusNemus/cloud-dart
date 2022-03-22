package main

import (
	"cloud-dart/match"
	"cloud-dart/player"
)

var userService *player.Service
var matchService *match.Service

func init() {
	userService = player.NewService("Player", "eu-north-1")
	matchService = match.NewService("Matches", "eu-north-1")
}

func main() {
}
