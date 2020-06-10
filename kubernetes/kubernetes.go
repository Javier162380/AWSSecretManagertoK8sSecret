package kubernetes

import (
	"encoding/base64"
	"log"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func kubernetesauth(kubeconfig string) (*kubernetes.Clientset, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil

}

// DonwloadSecret download a secret file from k8s secret and return a map
func DonwloadSecret(namespace string, secretrepository string,
	kubeconfig string) (map[string]string, error) {

	// creates the clientset
	clientset, err := kubernetesauth(kubeconfig)

	if err != nil {
		return nil, err
	}

	secretResponse, err := clientset.CoreV1().Secrets(namespace).Get(secretrepository, metav1.GetOptions{})

	if err != nil {
		return nil, err
	}

	rawsecret := secretResponse.Data
	processecret := make(map[string]string)

	for key, value := range rawsecret {
		base64string := base64.StdEncoding.EncodeToString(value)
		stringdecoded, err := base64.StdEncoding.DecodeString(base64string)
		if err != nil {
			return nil, err
		}
		processecret[key] = string(stringdecoded)

	}
	return processecret, nil

}

// UploadSecret create a secret k8s from a string aws secret
func UploadSecret(secretdata map[string]string, namespace string,
	secretrepository string, kubeconfig string) error {

	// creates the clientset
	clientset, err := kubernetesauth(kubeconfig)

	if err != nil {
		return err
	}

	// get all the secrets in a given namespace
	secretsMetadata, err := clientset.CoreV1().Secrets(namespace).List(metav1.ListOptions{})

	if err != nil {
		log.Fatalf("Unable to retrieve secret cluster information %s", err)
	}
	secretsList := secretsMetadata.Items

	for _, secret := range secretsList {
		if secret.ObjectMeta.Name == secretrepository {
			_, err := clientset.CoreV1().Secrets(namespace).Update(&v1.Secret{
				StringData: secretdata,
				ObjectMeta: metav1.ObjectMeta{
					Name: secretrepository,
				},
			})

			if err != nil {
				return err
			}

			return nil

		}

	}

	_, errr := clientset.CoreV1().Secrets(namespace).Create(&v1.Secret{
		StringData: secretdata,
		ObjectMeta: metav1.ObjectMeta{
			Name: secretrepository,
		},
	})

	if errr != nil {
		return errr
	}

	return nil
}
