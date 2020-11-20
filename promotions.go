package main

import (
	"flag"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	const dbPath = "./database.db"
	var deleteDB *bool = flag.Bool("r", false, "Recreate")
	flag.Parse()

	var s = newService(*deleteDB, dbPath)

	s.start()
}
