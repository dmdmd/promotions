package main

import (
	"flag"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	const dbPath = "./database.db"
	var deleteDB *bool = flag.Bool("r", false, "Recreate")
	flag.Parse()

	var myDBManager = sqliteManager{
		dbPath: dbPath,
	}

	var myPayloadFormatter = payloadFormatter{}

	if *deleteDB {
		os.Remove(dbPath)
		myDBManager.createDB()
	}

	var s = &service{
		db:         &myDBManager,
		payloadFmt: &myPayloadFormatter,
	}

	s.start()
}
