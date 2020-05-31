package secretmanager

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/prometheus/common/log"
)

func createsession(region string, profile string) *secretsmanager.SecretsManager {
	session := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region: aws.String(region),
		},
		Profile:           profile,
		SharedConfigState: session.SharedConfigEnable,
	},
	))
	return secretsmanager.New(session)

}

// SecretParser a secret file from AWSSSM repository
func SecretParser(secretrepository string, region string, profile string) string {

	secretClient := createsession(region, profile)

	secretResponse, err := secretClient.GetSecretValue(&secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretrepository),
	})

	if err != nil {
		log.Fatalf("Unable to retrieve secret, information %s", err)
	}

	return *secretResponse.SecretString

}
