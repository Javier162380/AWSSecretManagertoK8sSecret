package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	Namespace        string
	SecretRepository string
	Region           string
	Profile          string
	KubeConfigFile   string
	rootCMD          = &cobra.Command{
		Use:   "secret-moving",
		Short: "Move data from a AWSSecretManager to Kubernetes Secrets",
		Long:  "Command line application to generate Kubernetes Secrets from an AWSSecretManager",
	}
)

// Execute cli command
func Execute() {
	rootCMD.Execute()
}

func init() {
	rootCMD.PersistentFlags().StringVar(&Namespace, "namespace", "",
		"kubernetes namespace where the secret is going to be created")
	rootCMD.PersistentFlags().StringVar(&SecretRepository, "secretrepository", "",
		"Target AWSSecretRepository to create in k8s")
	rootCMD.PersistentFlags().StringVar(&Region, "region", "",
		"AWS Region where the secret manager is located")
	rootCMD.PersistentFlags().StringVar(&Profile, "profile", "default",
		"AWS Profile for authenticate ssm request")
	rootCMD.PersistentFlags().StringVar(&KubeConfigFile, "kubeconfig",
		fmt.Sprintf("%s/%s", os.Getenv("HOME"), ".kube/config"),
		"KubeConfig file path to auth with Kubernetes")
}
