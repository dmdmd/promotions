package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type service struct {
	db *sql.DB
}

func main() {
	const dbPath = "./foo.db"
	var deleteDB *bool = flag.Bool("r", false, "Recreate")
	flag.Parse()

	if *deleteDB {
		os.Remove(dbPath)
	}

	var database *sql.DB = getDB(dbPath)

	var s = &service{
		db: database,
	}

	queryDB(database)

	http.Handle("/products", http.HandlerFunc(s.handler))

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func (s *service) handler(w http.ResponseWriter, r *http.Request) {

}

func dbNotExists(filename string) bool {
	_, err := os.Stat(filename)
	return os.IsNotExist(err)
}

func getDB(filename string) *sql.DB {
	if dbNotExists(filename) {
		return createDB()
	}
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		log.Fatal(err)
	}
	// defer db.Close()

	return db
}

func createDB() *sql.DB {
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		log.Fatal(err)
	}
	// defer db.Close()

	sqlStmt := `
	CREATE TABLE products (
		id INTEGER NOT NULL PRIMARY KEY,
		sku TEXT,
		name TEXT,
		category TEXT,
		price INTEGER
	);
	DELETE FROM products;
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return db
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare("insert into products(id, sku, name, category, price) values(?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	categories := []string{
		"boots",
		"sandals",
		"sneakers",
	}

	for i := 0; i < 9; i++ {
		cat := categories[i%3]
		_, err = stmt.Exec(i, "00000"+fmt.Sprint(i), "leather "+cat, cat, rand.Intn(100)*100)
		if err != nil {
			log.Fatal(err)
		}
	}
	tx.Commit()

	return db
}

func queryDB(db *sql.DB) {
	rows, err := db.Query("select id, sku, name, category, price as original from products")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var (
		id, sku, name, category string
		original                int
	)
	for rows.Next() {
		rows.Scan(&id, &sku, &name, &category, &original)
		fmt.Println(strings.Join([]string{sku, name, category, fmt.Sprint(original)}, ","))
	}
}
