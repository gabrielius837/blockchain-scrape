package blockchainscrape

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type Database interface {
	Exec(query string, args ...any) (sql.Result, error)
}

const (
	INSERT_ADDRESS = "INSERT OR IGNORE INTO address (id, block, flag) VALUES(X'%X', %d, %d);"
	INSERT_TX      = "INSERT OR IGNORE INTO tx (id, \"from\", \"to\") VALUES(X'%X', X'%X', X'%X');"
	UNKNOWN        = 0
	WALLET         = 1
	CREATE         = 2
)

func GetAddressInsert(address []byte, number uint64, flag int) string {
	return fmt.Sprintf(INSERT_ADDRESS, address, number, flag)
}

func GetTxInsert(tx Transaction) string {
	return fmt.Sprintf(INSERT_TX, tx.Hash, tx.From, tx.To)
}

func initDb(dbName string, schemaFile string) (*sql.DB, error) {
	_, fsErr := os.Stat(dbName)
	log.Printf("Connecting to %s...\n", dbName)
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		return nil, err
	}
	log.Println("Connected successfully!")

	if os.IsNotExist(fsErr) {
		log.Printf("Initializing schema with %s...\n", schemaFile)
		bytes, err := os.ReadFile(schemaFile)
		if err != nil {
			return nil, err
		}

		script := string(bytes)

		_, err = db.Exec(script)
		if err != nil {
			return nil, err
		}

		log.Println("Initialized successfully!")
	}

	return db, nil
}

func ReadBlockNumber(db *sql.DB) (uint64, error) {
	var number uint64
	err := db.QueryRow("select number from block;").Scan(&number)
	if err != nil {
		return 0, err
	}

	return number, nil
}

func WriteBlockNumber(db Database, number uint64) error {
	result, err := db.Exec("UPDATE block set number = ? where number + 1 = ?", number, number)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		return fmt.Errorf("new block number %d expected to affect a single row instead %d rows where affected", number, rowsAffected)
	}

	return nil
}

func PersistBlock(db *sql.DB, block Block) error {
	transaction, err := db.Begin()
	if err != nil {
		return err
	}

	// While the majority of miner addresses are likely to be individual wallet addresses,
	// it is possible for a miner to use a smart contract address as their miner address.
	// In fact, some mining pools use a smart contract to distribute rewards among their members,
	// and in such cases, the miner address listed in the block would be the smart contract address.
	command := GetAddressInsert(block.Miner[:], uint64(block.Number), UNKNOWN)
	_, err = transaction.Exec(command)
	if err != nil {
		transaction.Rollback()
		return err
	}

	for _, tx := range block.Transactions {
		command = GetAddressInsert(tx.From[:], uint64(block.Number), WALLET)
		_, err = transaction.Exec(command)
		if err != nil {
			transaction.Rollback()
			return err
		}

		isNull := tx.To.IsNull()

		if isNull {
			contract, err := tx.ComputeContractAddres()
			if err != nil {
				transaction.Rollback()
				return err
			}

			command = GetAddressInsert(contract[:], uint64(block.Number), CREATE)
			_, err = transaction.Exec(command)
			if err != nil {
				transaction.Rollback()
				return err
			}
		} else {

			command = GetAddressInsert(tx.To[:], uint64(block.Number), WALLET)
			_, err = transaction.Exec(command)
			if err != nil {
				transaction.Rollback()
				return err
			}
		}

		command = GetTxInsert(tx)
		_, err = transaction.Exec(command)
		if err != nil {
			transaction.Rollback()
			return err
		}
	}

	err = WriteBlockNumber(transaction, uint64(block.Number))
	if err != nil {
		transaction.Rollback()
		return err
	}

	err = transaction.Commit()
	if err != nil {
		return err
	}

	fmt.Printf("transaction for block %d was committed successfully\n", uint64(block.Number))

	return nil
}
