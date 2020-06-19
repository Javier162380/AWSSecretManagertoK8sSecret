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

// UploadSecret upload a secret file to AWSSecretManager repository from a map string k8s string.
func UploadSecret(secretdata map[string]string, secretrepository string, region string, profile string) error {

	rawdata, err := json.Marshal(&secretdata)
	if err != nil {
		return err
	}
	awssecretdata := aws.String(string(rawdata))
	secretid := aws.String(secretrepository)

	secretClient := createsession(region, profile)
	secretsResponse, err := secretClient.ListSecrets(&secretsmanager.ListSecretsInput{})

	if err != nil {
		return err
	}

	secretsList := secretsResponse.SecretList

	for _, secret := range secretsList {

		if secret.Name == aws.String(secretrepository) {

			_, err := secretClient.UpdateSecret(&secretsmanager.UpdateSecretInput{
				SecretId:     secretid,
				SecretString: awssecretdata,
			})

			if err != nil {
				return err
			}

			return nil

		}

	}

	_, errr := secretClient.CreateSecret(&secretsmanager.CreateSecretInput{
		Name:         aws.String(secretrepository),
		SecretString: awssecretdata,
	})

	if errr != nil {
		return errr
	}

	return nil
}
