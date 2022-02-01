package main

import (
	"cloud-dart/match"
	"cloud-dart/player"
	"fmt"
)

var userService *player.Service
var matchService *match.Service

func init() {
	userService = player.NewService("Player", "eu-north-1")
	matchService = match.NewService("Matches", "eu-north-1")
}

func main() {
	users, err := userService.GetAll()
	if err != nil {
		return
	}

	for _, user := range users {
		fmt.Printf("%+v\n", user)
	}

}
