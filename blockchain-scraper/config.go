package blockchainscrape

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/joho/godotenv"
)

type Config struct {
	ApiKey   string
	Database string
	Run      bool
}

func readConfig() (*Config, error) {
	errors := []string{}
	config := &Config{Run: true}
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

func (config *Config) CatchSignal() {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT)
	<-signalChannel
	log.Println("gracefully shutting down")
	config.Run = false
}

func GetConfig() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	config, err := readConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return config
}

func (config *Config) InitDb() *sql.DB {
	db, err := initDb(config.Database, "init.sql")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return db
}
