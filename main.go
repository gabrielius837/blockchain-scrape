package main

import (
	"fmt"
	"os"

	scraper "github.com/gabrielius837/blockchain-scraper"
)

func main() {
	config := scraper.GetConfig()

	db := config.InitDb()
	defer db.Close()

	number, err := scraper.ReadBlockNumber(db)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	i := uint64(1)
	go config.CatchSignal()
	for ; config.Run; i++ {
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
