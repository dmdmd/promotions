package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
)

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

type dbManager interface {
	openDB()
	createDB()
	dbNotExists() bool
	queryDB(filterByCategory string, filterByPriceLessThan string) []product
}

type sqliteManager struct {
	dbPath   string
	db       *sql.DB
	dbIsOpen bool
}

func (s *sqliteManager) openDB() {
	if s.dbNotExists() {
		s.createDB()
		s.dbIsOpen = true
		return
	}

	db, err := sql.Open("sqlite3", s.dbPath)

	if err != nil {
		log.Fatal(err)
	}

	s.db = db
}

func (s *sqliteManager) dbNotExists() bool {
	_, err := os.Stat(s.dbPath)
	return os.IsNotExist(err)
}

func (s *sqliteManager) createDB() {
	db, err := sql.Open("sqlite3", s.dbPath)

	if err != nil {
		log.Fatal(err)
	}

	s.createProductTable(db)

	s.db = db
}

func (s *sqliteManager) createProductTable(db *sql.DB) {
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

		log.Println("inserting product", p)
		_, err = stmt.Exec(p.Sku, p.Name, p.Category, p.Price)
		if err != nil {
			log.Fatal(err)
		}
	}
	tx.Commit()
}

func (s *sqliteManager) queryDB(filterByCategory string, filterByPriceLessThan string) []product {
	var (
		products            []product
		queryString         = "SELECT sku, name, category, price FROM products"
		sku, name, category string
		price               int
	)

	s.openDB()

	defer s.db.Close()

	if len(filterByCategory) > 0 {
		queryString += fmt.Sprintf(" WHERE category = \"%v\"", filterByCategory)
	}

	if len(filterByPriceLessThan) > 0 && len(filterByCategory) > 0 {
		queryString += " AND price < " + filterByPriceLessThan
	} else if len(filterByPriceLessThan) > 0 {
		queryString += " WHERE price <= " + filterByPriceLessThan
	}

	log.Println(queryString)

	rows, err := s.db.Query(queryString)

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		rows.Scan(&sku, &name, &category, &price)
		log.Println(strings.Join([]string{sku, name, category, fmt.Sprint(price)}, ","))

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
