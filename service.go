package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type service struct {
	db         dbManager
	payloadFmt formatter
}

func newService(recreateDB bool, dbPath string) *service {
	var myDBManager = sqliteManager{
		dbPath: dbPath,
	}

	var myPayloadFormatter = payloadFormatter{}

	if recreateDB {
		os.Remove(dbPath)
		myDBManager.createDB()
	}

	return &service{
		db:         &myDBManager,
		payloadFmt: &myPayloadFormatter,
	}
}

func (s *service) start() {
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

	dbContent := s.db.queryDB(filterByCategory, filterByPriceLessThan)

	finalPayload := s.payloadFmt.formatResponse(dbContent)

	json.NewEncoder(w).Encode(finalPayload)
}