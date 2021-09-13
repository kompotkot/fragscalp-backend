package main

import (
	"flag"
	"fmt"

	"streamglass/pkg"
)

func main() {
	/*
		Entrypoint.
	*/
	var fHost, fPort string
	flag.StringVar(&fHost, "host", "0.0.0.0", "Server host")
	flag.StringVar(&fPort, "port", "7881", "Server port")

	var fTest bool
	flag.BoolVar(&fTest, "test", false, "Test extractors without opening WebSocket connection")
	flag.Parse()

	if fTest {
		chTest := make(chan []byte, 1)
		pair := streamglass.BINANCE_ETHBTC
		depth := 20
		updateSpeed := 100
		go streamglass.Extractor(chTest, pair, depth, updateSpeed)
		for {
			select {
			case val := <-chTest:
				fmt.Println(string(val))
			}
		}
	} else {
		streamglass.Server(fHost, fPort)
	}
}
