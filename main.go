package main

import (
	secretmanager "AWSSecretManagertoK8sSecret/awssecretmanager"
	"AWSSecretManagertoK8sSecret/cmd"
	"AWSSecretManagertoK8sSecret/kubernetes"
	"log"
)

func main() {
	cmd.Execute()
	secretdata := secretmanager.SecretParser(cmd.SecretRepository, cmd.Region, cmd.Profile)
	_, err := kubernetes.CreateSecret(secretdata, cmd.Namespace, cmd.SecretRepository)

	if err != nil {
		log.Fatalf("Unable to create kubernetes secret %s", err)
	}
}
