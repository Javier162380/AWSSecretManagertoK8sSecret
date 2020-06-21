package main

import (
	"log"
	"secret-moving/cmd"
)

func main() {

	err := cmd.Execute()

	if err != nil {
		log.Fatalf("Unable to execute command error %s", err)
	}

}
