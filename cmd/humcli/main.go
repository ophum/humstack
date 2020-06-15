package main

import (
	"log"

	"github.com/ophum/humstack/pkg/humcli/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		log.Println(err)
	}
}
