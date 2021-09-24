# streamglass-backend

Backend API Go client to process raw data from different sources.
Clean, organize compress it and send to web client [fragscalp-frontend](https://github.com/kompotkot/fragscalp-frontend).

## Startup

- Build package

```bash
go build -o build/server main.go
```

- Start the server

```bash
./build/server -host 0.0.0.0 -port 7881
```

- With test argument server doesn't open WS connection and print output to terminal

```bash
go run main.go -test
```
