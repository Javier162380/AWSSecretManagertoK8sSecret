package secretk8s

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// CreateSecret create a secret k8s from a string aws secret
func CreateSecret(secretdata map[string]string, namespace string, secretrepository string) {

	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	clientset.CoreV1().Secrets(namespace).Update(&v1.Secret{
		StringData: secretdata,
		ObjectMeta: metav1.ObjectMeta{
			Name: "test",
		},
	})
}
