package kubernetes

import (
	"log"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// CreateSecret create a secret k8s from a string aws secret
func CreateSecret(secretdata map[string]string, namespace string,
	secretrepository string, kubeconfig string) (*v1.Secret, error) {

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatalf("Uanble to generate k8s config %s", config)
	}

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Unable to authenticate in kubernetes %s", err)
	}

	// get all the secrets in a given namespace
	secretsMetadata, err := clientset.CoreV1().Secrets(namespace).List(metav1.ListOptions{})
	secretsList := secretsMetadata.Items

	for _, secret := range secretsList {
		if secret.ObjectMeta.Name == secretrepository {
			return clientset.CoreV1().Secrets(namespace).Update(&v1.Secret{
				StringData: secretdata,
				ObjectMeta: metav1.ObjectMeta{
					Name: secretrepository,
				},
			})

		}

	}

	return clientset.CoreV1().Secrets(namespace).Create(&v1.Secret{
		StringData: secretdata,
		ObjectMeta: metav1.ObjectMeta{
			Name: secretrepository,
		},
	})
}
