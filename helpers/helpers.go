package helpers

import "github.com/joho/godotenv"

// GenerateEnvFile method to create an env var from a secret service repository
func GenerateEnvFile(secretdata map[string]string, envpath string) error {
	return godotenv.Write(secretdata, envpath)

}
