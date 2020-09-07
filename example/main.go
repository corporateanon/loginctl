package main

import (
	"log"
	"time"

	"github.com/corporateanon/loginctl"
)

func main() {
	// lctl := loginctl.New([]string{"root", "user", "dev"})
	lctl, err := loginctl.NewFromRegularUsers()
	if err != nil {
		log.Panicln(err)
	}

	for {
		log.Println(lctl.GetSessionInfo())
		time.Sleep(time.Second * 2)
	}

}
