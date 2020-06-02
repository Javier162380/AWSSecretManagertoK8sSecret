package cmd

import (
	secretmanager "AWSSecretManagertoK8sSecret/awssecretmanager"
	"AWSSecretManagertoK8sSecret/kubernetes"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var (
	namespace        string
	secretRepository string
	region           string
	profile          string
	kubeConfigFile   string
	awstok8s         = &cobra.Command{
		Use:   "awstok8s",
		Short: "Move data from a AWSSecretManager to Kubernetes Secrets",
		Long:  "Command line application to upload k8s Secrets from an AWSSecretManager",
		RunE:  fromawstok8s,
	}
	k8stoaws = &cobra.Command{
		Use:   "k8stoaws",
		Short: "Move data from a Kubernetes to AWSSecretManager",
		Long:  "Command line application to upload AWSSecretRepository from k8s secrets",
		RunE:  fromk8stoaws,
	}
	rootCMD = &cobra.Command{
		Use: "secret-moving",
	}
)

// Execute cli command
func Execute() error {
	err := rootCMD.Execute()

	return err
}

func init() {
	rootCMD.PersistentFlags().StringVar(&namespace, "namespace", "",
		"kubernetes namespace where the secret is going to be created")
	rootCMD.PersistentFlags().StringVar(&secretRepository, "secretrepository", "",
		"Target AWSSecretRepository to create in k8s")
	rootCMD.PersistentFlags().StringVar(&region, "region", "",
		"AWS Region where the secret manager is located")
	rootCMD.PersistentFlags().StringVar(&profile, "profile", "default",
		"AWS Profile for authenticate ssm request")
	rootCMD.PersistentFlags().StringVar(&kubeConfigFile, "kubeconfig",
		fmt.Sprintf("%s/%s", os.Getenv("HOME"), ".kube/config"),
		"KubeConfig file path to auth with Kubernetes")
	rootCMD.AddCommand(awstok8s)
	rootCMD.AddCommand(k8stoaws)
}

func parsestringflag(cmd *cobra.Command, flag string) (string, error) {
	commandresult, err := cmd.Flags().GetString(flag)
	if err != nil {
		return "", err
	}

	return commandresult, nil
}

func parserootcmdcommand(cmd *cobra.Command) (map[string]string, error) {
	rootcommandsmap := make(map[string]string)

	cmdcommands := []string{"region", "profile", "secretrepository", "namespace", "kubeconfig"}
	for _, command := range cmdcommands {
		commandresult, err := parsestringflag(cmd, command)
		if err != nil {
			return nil, err
		}
		rootcommandsmap[command] = commandresult

	}

	return rootcommandsmap, nil
}

func fromawstok8s(cmd *cobra.Command, args []string) error {

	rootCommands, err := parserootcmdcommand(cmd)

	if err != nil {
		return err
	}

	secretRepository, region, profile = rootCommands["secretrepository"], rootCommands["region"], rootCommands["profile"]
	secretdata, err := secretmanager.DownloadSecret(secretRepository, region, profile)

	if err != nil {
		return err
	}

	kubeConfigFile = rootCommands["kubeconfig"]
	_, errr := kubernetes.UploadSecret(secretdata, namespace, secretRepository, kubeConfigFile)

	if errr != nil {
		return errr
	}

	log.Printf("Secret repository %s created/updated succesfully in namespace %s",
		secretRepository, namespace)

	return nil
}

func fromk8stoaws(cmd *cobra.Command, args []string) error {
	rootCommands, err := parserootcmdcommand(cmd)

	if err != nil {
		return err
	}
	namespace, secretRepository, kubeConfigFile = rootCommands["namespace"], rootCommands["secretrepository"], rootCommands["kubeconfig"]

	response, err := kubernetes.DonwloadSecret(namespace, secretRepository, kubeConfigFile)

	if err != nil {
		return err
	}

	log.Printf("%s", response)

	return nil

}
