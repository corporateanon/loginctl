package main

import (
	"log"

	"github.com/corporateanon/loginctl"
)

func main() {
	lctl, err := loginctl.New([]string{"mk", "ar", "te"})
	if err != nil {
		panic(err)
	}
	log.Println(lctl.GetSessionInfo())
}
