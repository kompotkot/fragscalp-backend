package streamglass

import (
	"encoding/json"
)

func Extractor(chOut chan []byte, pair Pair, depth int, updateSpeed int) {
	/*
		Execute required extractors and start streaming data
		one by one.

		:chOut: (chan) Chanel for list of bytes
		:pair: (Pair) Pair with their names, precision step, etc
		:depth: (int) Depth of Binance depth glass (5, 10, 20)
		:updateSpeed: (int) Speed of Binance depth glass update (100, 1000)
	*/
	const (
		EXTENSION_LEN   int = 0   // How far extend Glass depth with empty values
	)

	chDepth := make(chan []Row, 1)
	chTrade := make(chan Trade, 1)
	// chTxpool := make(chan []Txpool, 1)

	go BinanceDepth(chDepth, pair, depth, 100, EXTENSION_LEN)
	go BinanceTrade(chTrade, pair)
	// go GethTxpool(chTxpool, 0.1)

	// CH_LOOP:
	for {
		select {
		case val := <-chDepth:
			message, _ := json.Marshal(WSResponse{Type: "depth", Data: val})
			chOut <- message
		case val := <-chTrade:
			message, _ := json.Marshal(WSResponse{Type: "trade", Data: val})
			chOut <- message
		// case val := <-chTxpool:
		// 	message, _ := json.Marshal(WSResponse{Type: "txpool", Data: val})
		// 	chOut <- message
			// default:
			// 	break CH_LOOP
		}
	}
}
