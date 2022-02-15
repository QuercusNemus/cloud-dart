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
	m := match.Match{
		MatchId: "49d34bb7-862a-4e46-5b3d-3664049d7778",
		SortKey: "MATCH",
	}
	get, err := matchService.Delete(m)
	if err != nil {
		panic(err)
		return
	}

	fmt.Println(get)
}
