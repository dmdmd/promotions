package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func getServiceResponse(t *testing.T, target string) []byte {
	const dbPath = "./database.db"
	var deleteDB = true

	var s = newService(deleteDB, dbPath)

	req := httptest.NewRequest("GET", target, nil)
	rr := httptest.NewRecorder()
	s.handler(rr, req)

	resp := rr.Result()

	if status := resp.StatusCode; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		t.Errorf("Error reading products endpoint: %v\n", err)
	}

	return bodyBytes
}

/*
TestHandlerNoParam ...
*/
func TestHandlerNoParam(t *testing.T) {
	var (
		responseBodyJSON payload
		err              error
		bodyBytes        []byte
		wanted           = "{\"products\":[{\"sku\":\"000001\",\"name\":\"BV Lean leather ankle boots\",\"category\":\"boots\",\"price\":{\"original\":89000,\"final\":62299,\"discount\":\"30%\",\"currency\":\"EUR\"}},{\"sku\":\"000002\",\"name\":\"BV Lean leather ankle boots\",\"category\":\"boots\",\"price\":{\"original\":99000,\"final\":69300,\"discount\":\"30%\",\"currency\":\"EUR\"}},{\"sku\":\"000003\",\"name\":\"Ashlington leather ankle boots\",\"category\":\"boots\",\"price\":{\"original\":71000,\"final\":60350,\"discount\":\"15%\",\"currency\":\"EUR\"}},{\"sku\":\"000004\",\"name\":\"Naima embellished suede sandals\",\"category\":\"sandals\",\"price\":{\"original\":79500,\"final\":79500,\"discount\":null,\"currency\":\"EUR\"}},{\"sku\":\"000005\",\"name\":\"Nathane leather sneakers\",\"category\":\"sneakers\",\"price\":{\"original\":59000,\"final\":59000,\"discount\":null,\"currency\":\"EUR\"}}]}\n"
	)

	bodyBytes = getServiceResponse(t, "/products")

	err = json.Unmarshal(bodyBytes, &responseBodyJSON)

	if err != nil {
		t.Errorf("JSON error: %v", err)
	}

	if wanted != string(bodyBytes) {
		t.Errorf("JSON error: got %v, expected %v", string(bodyBytes), wanted)
	}
}

/*
TestHandlerPriceLessThan tests price filter
*/
func TestHandlerPriceLessThan(t *testing.T) {
	var (
		responseBodyJSON payload
		err              error
		bodyBytes        []byte
		maxPrice         = 70000
	)

	bodyBytes = getServiceResponse(t, fmt.Sprintf("/products?priceLessThan=%v", maxPrice))

	err = json.Unmarshal(bodyBytes, &responseBodyJSON)

	if err != nil {
		t.Errorf("JSON error: %v", err)
	}

	for i := 0; i < len(responseBodyJSON.Products); i++ {
		p := responseBodyJSON.Products[i]
		if p.Price.Original > maxPrice {
			t.Errorf("Found a product with incorrect price: got %v, expected less than %v", p.Price.Final, maxPrice)
		}
	}
}

/*
TestHandlerCategory tests category filter
*/
func TestHandlerCategory(t *testing.T) {
	var (
		responseBodyJSON payload
		err              error
		bodyBytes        []byte
		category         = "boots"
	)

	bodyBytes = getServiceResponse(t, fmt.Sprintf("/products?category=%v", category))

	err = json.Unmarshal(bodyBytes, &responseBodyJSON)

	if err != nil {
		t.Errorf("JSON error: %v", err)
	}

	for i := 0; i < len(responseBodyJSON.Products); i++ {
		p := responseBodyJSON.Products[i]
		if p.Category != category {
			t.Errorf("Found a product with incorrect category: got %v, expected %v", p.Category, category)
		}
	}
}

/*
TestHandlerCategoryAndPrice test category and price filter
*/
func TestHandlerCategoryAndPrice(t *testing.T) {
	var (
		responseBodyJSON payload
		err              error
		bodyBytes        []byte
		category         = "boots"
		maxPrice         = 90000
	)

	bodyBytes = getServiceResponse(t, fmt.Sprintf("/products?category=%v&priceLessThan=%v", category, maxPrice))

	err = json.Unmarshal(bodyBytes, &responseBodyJSON)

	if err != nil {
		t.Errorf("JSON error: %v", err)
	}

	for i := 0; i < len(responseBodyJSON.Products); i++ {
		p := responseBodyJSON.Products[i]
		if p.Category != category {
			t.Errorf("Found a product with incorrect category: got %v, expected %v", p.Category, category)
		}
		if p.Price.Original > maxPrice {
			t.Errorf("Found a product with incorrect price: got %v, expected less than %v", p.Price.Final, maxPrice)
		}
	}
}
