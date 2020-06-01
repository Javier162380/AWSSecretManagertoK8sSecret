package secretmanager

import (
	"encoding/json"
	"fmt"

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

// SecretParser request a secretfile from AWSSSM repository and returns a map[string][string]
func SecretParser(secretrepository string, region string, profile string) map[string]string {

	secretClient := createsession(region, profile)

	secretResponse, err := secretClient.GetSecretValue(&secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretrepository),
	})

	parseOutput := make(map[string]interface{})

	if err != nil {
		log.Fatalf("Unable to retrieve secret, information %s", err)
	}

	err = json.Unmarshal([]byte(*secretResponse.SecretString), &parseOutput)

	if err != nil {
		log.Fatalf("Unable to unmarshal secret response %s", err)
	}

	parseResult := make(map[string]string)

	for key, value := range parseOutput {
		parseResult[key] = fmt.Sprintf("%s", value)

	}
	return parseResult

}
