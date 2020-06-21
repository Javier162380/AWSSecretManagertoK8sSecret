package env

import "github.com/joho/godotenv"

// GenerateEnvFile method to create an env var from a secret service repository
func GenerateEnvFile(secretdata map[string]string, envpath string) error {
	return godotenv.Write(secretdata, envpath)

}

// LoadEnvFile method to load env file from an specific path
func LoadEnvFile(envpath string) (map[string]string, error) {
	return godotenv.Read(envpath)

}
