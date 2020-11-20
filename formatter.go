package main

import "log"

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

type formatter interface {
	formatResponse(ps []product) payload
	getPriceAndDiscount(p product) (int, *string)
}

type payloadFormatter struct {}

func (s *payloadFormatter) formatResponse(ps []product) payload {
	var products = make([]productPayload, len(ps))
	for i := 0; i < len(ps); i++ {
		p := ps[i]
		finalPrice, discount := s.getPriceAndDiscount(p)

		products[i] = productPayload{
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

		log.Println("formated payload", products[i], "end")
	}

	return payload{
		Products: products,
	}
}

func (s *payloadFormatter) getPriceAndDiscount(p product) (int, *string) {
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
