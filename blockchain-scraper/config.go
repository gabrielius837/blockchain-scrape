package blockchainscrape

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	ApiKey   string
	Database string
}

func ReadConfig() (*Config, error) {
	errors := []string{}
	config := &Config{}
	var key, value string
	var exists bool
	errMsg := "%s is missing"

	key = "API_KEY"
	value, exists = os.LookupEnv(key)
	if !exists || len(value) == 0 {
		errors = append(errors, fmt.Sprintf(errMsg, key))
	} else {
		config.ApiKey = value
	}

	key = "DATABASE"
	value, exists = os.LookupEnv(key)
	if !exists || len(value) == 0 {
		errors = append(errors, fmt.Sprintf(errMsg, key))
	} else if !(strings.HasSuffix(value, ".db")) {
		errors = append(errors, fmt.Sprintf("'%s' is missing suffix '.db'", value))
	} else {
		config.Database = value
	}

	if len(errors) > 0 {
		error := strings.Join(errors, "\n")
		return nil, fmt.Errorf(error)
	}

	return config, nil
}
