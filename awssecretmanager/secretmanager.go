package secretmanager

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
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

// DownloadSecret request a secretfile from AWSSSM repository and returns a map[string][string]
func DownloadSecret(secretrepository string, region string, profile string) (map[string]string, error) {

	secretClient := createsession(region, profile)

	secretResponse, err := secretClient.GetSecretValue(&secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretrepository),
	})

	if err != nil {
		return nil, err
	}

	parseOutput := make(map[string]interface{})
	err = json.Unmarshal([]byte(*secretResponse.SecretString), &parseOutput)

	if err != nil {
		return nil, err
	}

	parseResult := make(map[string]string)

	for key, value := range parseOutput {
		parseResult[key] = fmt.Sprintf("%s", value)

	}

	if len(parseResult) == 0 {
		return nil, errors.New(`secret repository is empty.
			secret repository ougth to contain secret to generate k8s secret`)

	}

	return parseResult, nil

}
