# AWSSecretManagertoK8sSecret
Simple tool to move secrets from AWSSecretManager to K8sSecrets

# Usage.
``` bash
./secret-moving -h
```

```
Usage:
  secret-moving [command]

Available Commands:
  awstoenv    Move secrets data from a AWSSecretRepository to an envfile
  awstok8s    Move data from a AWSSecretRepository to Kubernetes Secrets
  help        Help about any command
  k8stoaws    Move data from a Kubernetes to AWSSecretManager
  k8stoenv    Move secrets data from Kubernetes to an envfile

Flags:
  -h, --help                      help for secret-moving
      --kubeconfig string         KubeConfig file path to auth with Kubernetes (default "/Users/javier.llorente/.kube/config")
      --namespace string          Kubernetes namespace where the secret is going to be created
      --profile string            AWS Profile for authenticate ssm request (default "default")
      --region string             AWS Region where the secret manager is located
      --secretrepository string   Target AWSSecretRepository to create in k8s

Use "secret-moving [command] --help" for more information about a command.
```

