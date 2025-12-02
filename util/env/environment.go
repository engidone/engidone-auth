package env

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func getSecret(secretPath string) string {
	file, err := os.Open(secretPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}

func Get(name string) string {
	value := strings.TrimSpace(os.Getenv(name))
	if value != "" {
		return value
	}

	return getSecret(fmt.Sprintf("/run/secrets/%s", strings.ToLower(name)))
}
