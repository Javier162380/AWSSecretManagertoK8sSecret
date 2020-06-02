package main

import (
	"AWSSecretManagertoK8sSecret/cmd"
	"log"
)

func main() {

	err := cmd.Execute()

	if err != nil {
		log.Fatalf("Unable to execute command error %s", err)
	}

}
