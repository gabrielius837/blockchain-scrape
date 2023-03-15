package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	scraper "github.com/gabrielius837/blockchain-scraper"
	"github.com/joho/godotenv"
)

func ExitIfError(err error) {
	if err == nil {
		return
	}
}

func main() {
	config := getConfig()

	db := initDb(config)
	defer db.Close()

	number, err := scraper.ReadBlockNumber(db)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	i := uint64(1)
	run := true
	go catchSignal(&run)
	for ; run; i++ {
		resp, err := scraper.GetBlock(config.ApiKey, i, number+i)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		block := resp.Result
		err = scraper.PersistBlock(db, block)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
	fmt.Println("execution is finished")
}

func catchSignal(run *bool) {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT)
	<-signalChannel
	log.Println("gracefully shutting down")
	*run = false
}

func getConfig() *scraper.Config {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	config, err := scraper.ReadConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return config
}

func initDb(config *scraper.Config) *sql.DB {
	db, err := scraper.InitDb(config.Database, "init.sql")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return db
}
