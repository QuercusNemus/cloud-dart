package main

import (
	"cloud-dart/user"
	"fmt"
)

var service *user.Service

func init() {
	service = user.NewService("User", "eu-north-1")
}

func main() {
	users, err := service.GetAll()
	if err != nil {
		return
	}

	for _, user := range users {
		fmt.Printf("%+v\n", user)
	}
}
