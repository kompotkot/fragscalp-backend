/*
Streamglass data extractors from Binance.
*/
package streamglass

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/gorilla/websocket"
	// streamglass "streamglass/pkg"
)

func BinanceDepth(chOut chan<- []Row, pair Pair, depth int, updateSpeed int, extensionLen int) {
	/*
		Fetch depth of asks and bids from Exchange with WebSockets.

		Docs:
		https://binance-docs.github.io/apidocs/spot/en/#partial-book-depth-streams
	*/
	wsClient, _, wsErr := websocket.DefaultDialer.Dial(fmt.Sprintf("wss://stream.binance.com:9443/ws/%s@depth%d@%dms", pair.Pair, depth, updateSpeed), nil)
	if wsErr != nil {
		log.Printf("[Depth] Connection WebSocket error: %d", wsErr)
	}

	for {
		// Start fetching messages from WebSocket connection
		_, message, readErr := wsClient.ReadMessage()
		if readErr != nil {
			log.Printf("[Depth] Read stream message error: %d", readErr)
			return
		}

		// Convert array of bytes to raw AsksBids structure
		asksBidsRaw := AsksBidsBinance{}
		jsonParseErr := json.Unmarshal([]byte(message), &asksBidsRaw)
		if jsonParseErr != nil {
			fmt.Println(jsonParseErr.Error())
		}

		// Convert raw AsksBidsRaw to temporary structure of not extended AsksBids
		asksTemp := make([]RowTemp, 0, depth)
		bidsTemp := make([]RowTemp, 0, depth)
		// Reverse asks in opposite order
		for i := len(asksBidsRaw.Asks) - 1; i >= 0; i-- {
			priceTemp, _ := strconv.ParseFloat(asksBidsRaw.Asks[i][0], 12)
			quantityTemp, _ := strconv.ParseFloat(asksBidsRaw.Asks[i][1], 12)
			asksTemp = append(asksTemp, RowTemp{Position: "ask", Price: int(priceTemp * float64(pair.ToDecimal)), Quantity: quantityTemp})
		}
		for _, bid := range asksBidsRaw.Bids {
			priceTemp, _ := strconv.ParseFloat(bid[0], 12)
			quantityTemp, _ := strconv.ParseFloat(bid[1], 12)
			bidsTemp = append(bidsTemp, RowTemp{Position: "bid", Price: int(priceTemp * float64(pair.ToDecimal)), Quantity: quantityTemp})
		}

		// Extend AsksBids, concat them together and extend in Rows structure
		topAsk := asksTemp[0] // Highest AskBid in Glass stack
		botAsk := asksTemp[len(asksTemp)-1]
		topBid := bidsTemp[0]
		botBid := bidsTemp[len(bidsTemp)-1] // Lowest AskBid in Glass stack

		// Fill the holes id Depth
		// TODO: Optimize this loop - remove inner loops
		rows := make([]Row, 0, topAsk.Price-botBid.Price+1) // Future fullfilled list of AskBids with empty values
	WS_DEPTH_LOOP:
		for p := topAsk.Price; p >= botBid.Price; p-- {
			if topAsk.Price >= p && p >= botAsk.Price {
				var tempRow Row
				// tempRow.Position = "ask"
				tempRow.Price = float64(p) / float64(pair.ToDecimal)
				if tempRow.Price == float64(int64(tempRow.Price)) {
					tempRow.Position = "ask_hard"
				} else {
					tempRow.Position = "ask"
				}
				tempRow.Quantity = 0
				for i, ask := range asksTemp {
					if p == ask.Price {
						tempRow.Quantity = asksTemp[i].Quantity
					}
				}
				rows = append(rows, tempRow)
				continue WS_DEPTH_LOOP
			} else if topBid.Price >= p && p >= botBid.Price {
				var tempRow Row
				tempRow.Price = float64(p) / float64(pair.ToDecimal)
				if tempRow.Price == float64(int64(tempRow.Price)) {
					tempRow.Position = "bid_hard"
				} else {
					tempRow.Position = "bid"
				}
				tempRow.Quantity = 0
				for i, bid := range bidsTemp {
					if p == bid.Price {
						tempRow.Quantity = bidsTemp[i].Quantity
					}
				}
				rows = append(rows, tempRow)
				continue WS_DEPTH_LOOP
			} else {
				rows = append(rows, Row{Position: "mid", Price: float64(p) / float64(pair.ToDecimal), Quantity: 0})
			}
		}

		chOut <- rows
	}
	close(chOut)
}

func BinanceTrade(chOut chan<- Trade, pair Pair) {
	/*
		Fetch trades from Exchange with WebSockets.
		If IsBuyerMarketMaker == True -> "green"
		If IsBuyerMarketMaker == False -> "red"

		Docs:
		https://binance-docs.github.io/apidocs/spot/en/#trade-streams
	*/
	wsClient, _, wsErr := websocket.DefaultDialer.Dial(fmt.Sprintf("wss://stream.binance.com:9443/ws/%s@trade", pair.Pair), nil)
	if wsErr != nil {
		log.Printf("[Trade] Connection WebSocket error: %d", wsErr)
	}

	for {
		trade := Trade{}

		// Start fetching messages from WebSocket connection
		_, message, readErr := wsClient.ReadMessage()
		if readErr != nil {
			log.Printf("[Trade] Read stream message error: %d", readErr)
			return
		}

		// Convert array of bytes to raw Trade structure
		tradeRaw := TradeBinance{}
		jsonParseErr := json.Unmarshal([]byte(message), &tradeRaw)
		if jsonParseErr != nil {
			fmt.Println(jsonParseErr.Error())
		}

		// Convert raw TradeBinanceRaw to final result of Trade
		if tradeRaw.EventType == "trade" {
			priceTemp, _ := strconv.ParseFloat(tradeRaw.Price, 12)
			quantityTemp, _ := strconv.ParseFloat(tradeRaw.Quantity, 12)

			if quantityTemp > pair.TradeMinValue {
				trade = Trade{
					Price:              priceTemp,
					Quantity:           quantityTemp,
					IsBuyerMarketMaker: tradeRaw.IsBuyerMarketMaker,
					TradeTime:          tradeRaw.TradeTime,
				}

				chOut <- trade
			}
		}
	}
	close(chOut)
}
