package blockchainscrape

import (
	"database/sql"
	"io"
	"log"
	"os"
	"testing"
)

// I know, not unit tests but whatever
const (
	SCRIPT = "../init.sql"
)

func TestMain(m *testing.M) {
	log.SetOutput(io.Discard)
	os.Exit(m.Run())
}

func initTestDb(dbName string, t *testing.T) (*sql.DB, uint64) {
	db, err := InitDb(dbName, "../init.sql")
	if err != nil {
		t.Fatal(err)
	}

	number, err := ReadBlockNumber(db)
	if err != nil {
		t.Fatal(err)
	}
	if number != 46146 {
		t.Fatalf("Unexpected block number, got %d, expecting 46146", number)
	}

	return db, number
}

func TestReadBlockNumber_MustReturnNumber_WhenNew(t *testing.T) {
	dbName := "test_readblocknumber.db"
	defer os.Remove(dbName)
	initTestDb(dbName, t)
}

func TestWriteBlockNumber_MustPersist_WhenHigher(t *testing.T) {
	dbName := "test_writeblocknumber_1.db"
	defer os.Remove(dbName)
	db, number := initTestDb(dbName, t)

	increment := uint64(1)
	err := WriteBlockNumber(db, number+increment)

	if err != nil {
		t.Fatal(err)
	}

	newNumber, err := ReadBlockNumber(db)
	if err != nil {
		t.Fatal(err)
	}
	if number+increment != newNumber {
		t.Fatalf("expected for increment to match new number, but %d + %d != %d", number, increment, newNumber)
	}
}

func TestWriteBlockNumber_MustFail_WhenSame(t *testing.T) {
	dbName := "test_writeblocknumber_2.db"
	defer os.Remove(dbName)
	db, number := initTestDb(dbName, t)

	increment := uint64(0)
	err := WriteBlockNumber(db, number+increment)

	if err == nil {
		t.Fatalf("when updating with the same number it must fail, but it did not")
	}
}

func TestWriteBlockNumber_MustFail_WhenLower(t *testing.T) {
	dbName := "test_writeblocknumber_3.db"
	defer os.Remove(dbName)
	db, number := initTestDb(dbName, t)

	increment := uint64(1)
	err := WriteBlockNumber(db, number-increment)

	if err == nil {
		t.Fatalf("when updating with a lower number it must fail, but it did not")
	}
}
