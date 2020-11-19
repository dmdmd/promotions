package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type service struct {
	db *sql.DB
}

type products struct {
	Products []product `json:"products"`
}
type product struct {
	Sku      string `json:"sku"`
	Name     string `json:"name"`
	Category string `json:"category"`
	Price    int    `json:"price"`
}

var productsTableContent = products{
	Products: []product{
		{
			Sku:      "000001",
			Name:     "BV Lean leather ankle boots",
			Category: "boots",
			Price:    89000,
		},
		{
			Sku:      "000002",
			Name:     "BV Lean leather ankle boots",
			Category: "boots",
			Price:    99000,
		},
		{
			Sku:      "000003",
			Name:     "Ashlington leather ankle boots",
			Category: "boots",
			Price:    71000,
		},
		{
			Sku:      "000004",
			Name:     "Naima embellished suede sandals",
			Category: "sandals",
			Price:    79500,
		},
		{
			Sku:      "000005",
			Name:     "Nathane leather sneakers",
			Category: "sneakers",
			Price:    59000,
		},
	},
}

func main() {
	const dbPath = "./database.db"
	var deleteDB *bool = flag.Bool("r", false, "Recreate")
	flag.Parse()

	if *deleteDB {
		os.Remove(dbPath)
	}

	var database *sql.DB = getDB(dbPath)

	var s = &service{
		db: database,
	}

	http.Handle("/products", http.HandlerFunc(s.handler))

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func (s *service) handler(w http.ResponseWriter, r *http.Request) {
	products := queryDB(s.db)

	json.NewEncoder(w).Encode(map[string]interface{}{"products": products})
}

func dbNotExists(filename string) bool {
	_, err := os.Stat(filename)
	return os.IsNotExist(err)
}

func getDB(filename string) *sql.DB {
	if dbNotExists(filename) {
		return createDB(filename)
	}
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		log.Fatal(err)
	}
	// defer db.Close()

	return db
}

func createDB(filename string) *sql.DB {
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		log.Fatal(err)
	}
	// defer db.Close()

	createProductTable(db)

	return db
}

func createProductTable(db *sql.DB) {
	sqlStmt := `
	CREATE TABLE products (
		sku TEXT NOT NULL PRIMARY KEY,
		name TEXT,
		category TEXT,
		price INTEGER
	);
	DELETE FROM products;
	`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare("insert into products(sku, name, category, price) values(?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	for i := 0; i < len(productsTableContent.Products); i++ {
		p := productsTableContent.Products[i]

		fmt.Println("inserting product", p)
		_, err = stmt.Exec(p.Sku, p.Name, p.Category, p.Price)
		if err != nil {
			log.Fatal(err)
		}
	}
	tx.Commit()
}

func queryDB(db *sql.DB) []map[string]interface{} {
	var products []map[string]interface{}

	rows, err := db.Query("select sku, name, category, price as original from products")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var (
		sku, name, category string
		original            int
	)
	for rows.Next() {
		rows.Scan(&sku, &name, &category, &original)
		fmt.Println(strings.Join([]string{sku, name, category, fmt.Sprint(original)}, ","))

		product := map[string]interface{}{
			"sku":      sku,
			"name":     name,
			"category": category,
			"original": original,
		}
		products = append(products, product)
	}

	return products
}
