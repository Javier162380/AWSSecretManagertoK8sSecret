package cmd

import (
	"fmt"
	"log"
	"os"
	secretmanager "secret-moving/awssecretmanager"
	env "secret-moving/envfile"
	"secret-moving/kubernetes"

	"github.com/spf13/cobra"
)

var (
	// root cmd flags
	namespace        string
	secretRepository string
	region           string
	profile          string
	kubeConfigFile   string

	// k8stoenv , awstoenv flags
	envpath string

	//cobra commands
	awstok8s = &cobra.Command{
		Use:   "awstok8s",
		Short: "Move data from a AWSSecretRepository to Kubernetes Secrets",
		Long:  "Command to upload k8s Secrets from an AWSSecretManager",
		RunE:  fromawstok8s,
	}
	k8stoaws = &cobra.Command{
		Use:   "k8stoaws",
		Short: "Move data from a Kubernetes to AWSSecretManager",
		Long:  "Command to upload data to AWSSecretRepository from k8s secrets",
		RunE:  fromk8stoaws,
	}
	awstoenv = &cobra.Command{
		Use:   "awstoenv",
		Short: "Move secrets data from a AWSSecretRepository to an envfile",
		Long:  "Command to download data from an AWSSecretRepository and store it in an envfile",
		RunE:  repositoryEnvWrapper,
	}
	k8stoenv = &cobra.Command{
		Use:   "k8stoenv",
		Short: "Move secrets data from Kubernetes to an envfile",
		Long:  "Command to donwload a secret from a K8s cluster and store it in an envfile",
		RunE:  repositoryEnvWrapper,
	}
	envtok8s = &cobra.Command{
		Use:   "envtok8s",
		Short: "Move secrets data from an envfile to K8s secrets",
		Long:  "Command to generate k8s secrets from envfiles",
		RunE:  repositoryEnvWrapper,
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

	// Adding rootCMDs persistent flags
	rootCMD.PersistentFlags().StringVar(&namespace, "namespace", "",
		"Kubernetes namespace where the secret is going to be created")
	rootCMD.PersistentFlags().StringVar(&secretRepository, "secretrepository", "",
		"Target AWSSecretRepository to create in k8s")
	rootCMD.PersistentFlags().StringVar(&region, "region", "",
		"AWS Region where the secret manager is located")
	rootCMD.PersistentFlags().StringVar(&profile, "profile", "default",
		"AWS Profile for authenticate ssm request")
	rootCMD.PersistentFlags().StringVar(&kubeConfigFile, "kubeconfig",
		fmt.Sprintf("%s/%s", os.Getenv("HOME"), ".kube/config"),
		"KubeConfig file path to auth with Kubernetes")

	// Adding k8stoenv flags.
	k8stoenv.PersistentFlags().StringVar(&envpath, "envpath", "",
		"Env variable path where you are going to create your envfile")

	// Adding awstoenv flags.
	awstoenv.PersistentFlags().StringVar(&envpath, "envpath", "",
		"Env variable path where you are going to create your envfile")

	// Adding envtok8s flags.
	envtok8s.PersistentFlags().StringVar(&envpath, "envpath", "",
		"Env varialbe path to read from to create the secret")

	// Adding subcommand to root cmd.
	rootCMD.AddCommand(awstok8s)
	rootCMD.AddCommand(awstoenv)
	rootCMD.AddCommand(k8stoaws)
	rootCMD.AddCommand(k8stoenv)
	rootCMD.AddCommand(envtok8s)
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
	errr := kubernetes.UploadSecret(secretdata, namespace, secretRepository, kubeConfigFile)

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

	region, profile = rootCommands["region"], rootCommands["profile"]

	errr := secretmanager.UploadSecret(response, secretRepository, region, profile)

	if errr != nil {
		return errr
	}

	return nil

}

func repositoryEnvWrapper(cmd *cobra.Command, args []string) error {
	rootCommands, err := parserootcmdcommand(cmd)

	if err != nil {
		return err
	}
	// generic flag
	secretRepository := rootCommands["secretrepository"]

	//k8s common flags
	namespace, kubeConfigFile := rootCommands["namespace"], rootCommands["kubeconfig"]

	//aws common flags
	region, profile = rootCommands["region"], rootCommands["profile"]

	envpath, err := parsestringflag(cmd, "envpath")

	if err != nil {
		return err
	}

	command := cmd.Use

	switch command {
	case "k8stoenv":
		secretdata, err := kubernetes.DonwloadSecret(namespace, secretRepository, kubeConfigFile)

		if err != nil {
			return err
		}

		errr := env.GenerateEnvFile(secretdata, envpath)

		if errr != nil {
			return errr
		}

		return nil

	case "awstoenv":
		secretdata, err := secretmanager.DownloadSecret(secretRepository, region, profile)

		if err != nil {
			return err
		}

		errr := env.GenerateEnvFile(secretdata, envpath)

		if errr != nil {
			return errr
		}

		return nil

	case "envtok8s":
		secretdata, err := env.LoadEnvFile(envpath)

		if err != nil {
			return err
		}

		errr := kubernetes.UploadSecret(secretdata, namespace, secretRepository, kubeConfigFile)

		if errr != nil {
			return errr
		}

		return nil

	default:
		return nil

	}

}
