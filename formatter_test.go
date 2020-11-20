package main

import "testing"

func TestGetPriceAndDiscountCategory(t *testing.T) {
	var (
		myFormatter = payloadFormatter{}
		p           = product{
			Sku:      "000001",
			Name:     "BV Lean leather ankle boots",
			Category: "boots",
			Price:    10000,
		}
		price, discount  = myFormatter.getPriceAndDiscount(p)
		expectedPrice    = 7000
		expectedDiscount = "30%"
	)

	if *discount != expectedDiscount {
		t.Errorf("Wrong discount: expected %v, got %v", expectedDiscount, discount)
	}

	if price != expectedPrice {
		t.Errorf("Wrong final price: expected %v, got %v", expectedPrice, price)
	}
}

func TestGetPriceAndDiscountSKU(t *testing.T) {
	var (
		myFormatter = payloadFormatter{}
		p           = product{
			Sku:      "000003",
			Name:     "BV Lean leather ankle boots",
			Category: "sandals",
			Price:    10000,
		}
		price, discount  = myFormatter.getPriceAndDiscount(p)
		expectedPrice    = 8500
		expectedDiscount = "15%"
	)

	if *discount != expectedDiscount {
		t.Errorf("Wrong discount: expected %v, got %v", expectedDiscount, discount)
	}

	if price != expectedPrice {
		t.Errorf("Wrong final price: expected %v, got %v", expectedPrice, price)
	}
}

func TestGetPriceAndDiscountSKUAndCategory(t *testing.T) {
	var (
		myFormatter = payloadFormatter{}
		p           = product{
			Sku:      "000003",
			Name:     "BV Lean leather ankle boots",
			Category: "boots",
			Price:    10000,
		}
		price, discount  = myFormatter.getPriceAndDiscount(p)
		expectedPrice    = 8500
		expectedDiscount = "15%"
	)

	if *discount != expectedDiscount {
		t.Errorf("Wrong discount: expected %v, got %v", expectedDiscount, discount)
	}

	if price != expectedPrice {
		t.Errorf("Wrong final price: expected %v, got %v", expectedPrice, price)
	}
}
