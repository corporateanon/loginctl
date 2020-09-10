package main

import (
	"log"
	"time"

	"github.com/corporateanon/loginctl"
)

func main() {
	lctl, err := loginctl.NewFromRegularUsers()
	if err != nil {
		log.Panicln(err)
	}

	log.Println("----USERS LIST----")
	log.Println(lctl.GetUsersList(true))

	log.Println("----MONITORING USER ACTIVITY----")
	for {
		log.Println(lctl.GetSessionInfo())
		time.Sleep(time.Second * 2)
	}
}
