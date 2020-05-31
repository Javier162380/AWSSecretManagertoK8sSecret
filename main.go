package main

import (
	"fmt"
	secretmanager "secret-moving/awssecretmanager"
	"secret-moving/cmd"
)

func main() {
	cmd.Execute()

	secret := secretmanager.SecretParser(cmd.SecretRepository, cmd.Region, cmd.Profile)
	fmt.Printf("%s", secret)
}
