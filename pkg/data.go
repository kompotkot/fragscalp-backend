/*
Streamglass data types and interfaces.
*/
package streamglass

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// Pairs
// TODO: Make as const?
var BINANCE_BNBUSDT = Pair{Pair: "bnbusdt", PrecisionStep: 0.1, ToDecimal: 10, TradeMinValue: 0.01}
var BINANCE_ETHUSDT = Pair{Pair: "ethusdt", PrecisionStep: 0.01, ToDecimal: 100, TradeMinValue: 0.01}
var BINANCE_ETHBTC = Pair{Pair: "ethbtc", PrecisionStep: 0, ToDecimal: 1000000, TradeMinValue: 0.0001}
var BINANCE_SOLUSDT = Pair{Pair: "solusdt", PrecisionStep: 0.01, ToDecimal: 100, TradeMinValue: 0.01}

// Symbol description
type Pair struct {
	Pair          string  `json:"pair"`
	PrecisionStep float64 `json:"precision"`
	ToDecimal     int     `json:"toDecimal"`
	TradeMinValue float64 `json:"tradeMinValue"`
}

// Row final representation with price and quantity
// in Glass fullfilled with data for Frontend
type Row struct {
	Position string `json:"position"` // Row position in Glass (ask, mid, bid)

	Price    float64 `json:"price"`
	Quantity float64 `json:"quantity"`

	// PriceTextColor          string `json:"priceTextColor"`          // Text color for price in Glass
	// QuantityTextColor       string `json:"quantityTextColor"`       // Text color for quantity in Glass
	// BackgroundPriceColor    string `json:"backgroundPriceColor"`    // Background row color for price in Glass
	// BackgroundQuantityColor string `json:"backgroundQuantityColor"` // Background row color for quantity in Glass
}

// Temporary Row representation
type RowTemp struct {
	Position string  `json:"position"`
	Price    int     `json:"price"`
	Quantity float64 `json:"quantity"`
}

// List of Rows fetched from Binance
type AsksBidsBinance struct {
	LastUpdateId int        `json:"lastUpdateId"`
	Asks         [][]string `json:"asks"`
	Bids         [][]string `json:"bids"`
}

// Trade final representation fullfilled with data for Frontend
type Trade struct {
	Price     float64 `json:"price"`
	Quantity  float64 `json:"quantity"`
	TradeTime uint64  `json:"tradeTime"`

	// True = ask was hitted and mark green, if false = bid was hitted and mark trade as red
	IsBuyerMarketMaker bool `json:"isBuyerMarketMaker"`

	// TradeTextColor       string `json:"tradeTextColor"`       // Text color for price in Trade screen
	// TradeBackgroundColor string `json:"tradeBackgroundColor"` // Background for trade bubble in Trade screen
}

// Trade fetched from Binance
type TradeBinance struct {
	EventType          string `json:"e"`
	EventTime          uint64 `json:"E"`
	Symbol             string `json:"s"`
	TradeId            uint64 `json:"t"`
	Price              string `json:"p"`
	Quantity           string `json:"q"`
	BuyerOrderId       uint64 `json:"b"`
	SellerOrderId      uint64 `json:"a"`
	TradeTime          uint64 `json:"T"`
	IsBuyerMarketMaker bool   `json:"m"`
	Ignore             bool   `json:"M"`
}

// Transaction final representation from txpool
type Txpool struct {
	Hash      string  `json:"hash"`
	Quantity  float64 `json:"quantity"`
	MinedTime string  `json:"minedTime"` // If transaction was mined, it sets timestamp, if not - null

	// QuantityTextColor       string `json:"quantityTextColor"`       // Text color for quantity in Glass
	// BackgroundQuantityColor string `json:"backgroundQuantityColor"` // Background row color for quantity in Glass
}

// Txpool raw transaction from Geth
type TxGeth struct {
	Hash  common.Hash  `json:"hash"`
	Value *hexutil.Big `json:"value"` // Quantity
}

type WSResponse struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}
