package security

import (
	"fmt"
	"os"
)

func GetEnv(env string) (string, error) {
	variable := os.Getenv(env)
	if variable == "" {
		return "", fmt.Errorf("%s not found in .env file", env)
	}
	return variable, nil
}
