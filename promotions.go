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

type payload struct {
	Products []productPayload `json:"products"`
}
type productPayload struct {
	Sku      string              `json:"sku"`
	Name     string              `json:"name"`
	Category string              `json:"category"`
	Price    productPricePayload `json:"price"`
}

type productPricePayload struct {
	Original           int     `json:"original"`
	Final              int     `json:"final"`
	DiscountPercentage *string `json:"discount"`
	Currency           string  `json:"currency"`
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

	http.HandleFunc("/products", s.handler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func (s *service) handler(w http.ResponseWriter, r *http.Request) {
	var filterByCategory, filterByPriceLessThan string
	if len(r.URL.Query()["category"]) > 0 {
		filterByCategory = r.URL.Query()["category"][0]
	}
	if len(r.URL.Query()["priceLessThan"]) > 0 {
		filterByPriceLessThan = r.URL.Query()["priceLessThan"][0]
	}

	dbContent := queryDB(s.db, filterByCategory, filterByPriceLessThan)

	finalPayload := formatResponse(dbContent)

	json.NewEncoder(w).Encode(payload{
		Products: finalPayload,
	})
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

	stmt, err := tx.Prepare("INSERT INTO products(sku, name, category, price) VALUES(?, ?, ?, ?)")
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

func queryDB(db *sql.DB, filterByCategory string, filterByPriceLessThan string) []product {
	var products []product
	var queryString = "SELECT sku, name, category, price FROM products"

	if len(filterByCategory) > 0 {
		queryString += fmt.Sprintf(" WHERE category = \"%v\"", filterByCategory)
	}

	if len(filterByPriceLessThan) > 0 && len(filterByCategory) > 0 {
		queryString += " AND price < " + filterByPriceLessThan
	} else if len(filterByPriceLessThan) > 0 {
		queryString += " WHERE price <= " + filterByPriceLessThan
	}

	fmt.Println(queryString)

	rows, err := db.Query(queryString)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var (
		sku, name, category string
		price               int
	)
	for rows.Next() {
		rows.Scan(&sku, &name, &category, &price)
		fmt.Println(strings.Join([]string{sku, name, category, fmt.Sprint(price)}, ","))

		product := product{
			Sku:      sku,
			Name:     name,
			Category: category,
			Price:    price,
		}
		products = append(products, product)
	}

	return products
}

func formatResponse(ps []product) []productPayload {
	var result = make([]productPayload, len(ps))
	for i := 0; i < len(ps); i++ {
		p := ps[i]
		finalPrice, discount := calculateFinalPriceAndDiscount(p)

		result[i] = productPayload{
			Sku:      p.Sku,
			Name:     p.Name,
			Category: p.Category,
			Price: productPricePayload{
				Original:           p.Price,
				Currency:           "EUR",
				Final:              finalPrice,
				DiscountPercentage: discount,
			},
		}

		fmt.Println("formated payload", result[i], "end")
	}

	return result
}

func calculateFinalPriceAndDiscount(p product) (int, *string) {
	var discount string

	if p.Sku == "000003" {
		discount = "15%"
		return int(float64(p.Price) * 0.7), &discount
	}

	if p.Category == "boots" {
		discount = "30%"
		return int(float64(p.Price) * 0.7), &discount
	}

	return p.Price, nil
}
