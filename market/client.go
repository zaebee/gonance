//Package general provides the binding for Binance Rest APIs market endpoints
package market

import "gonance/client"

type Client struct {
	API client.API
}

//Methods for market endpoints
type MarketInterface interface {
	Prices() (Prices, error)
}

var _ MarketInterface = (*Client)(nil)

type Prices map[string]Price

//Prices returns latest price for all symbols.
func (c *Client) Prices() (Prices, error) {
	priceList := []Price{}
	err := c.API.Request("GET", "/api/v1/ticker/allPrices", nil, &priceList)
	prices := Prices{}
	for _, p := range priceList {
		prices[p.Symbol] = p
	}
	return prices, err
}
