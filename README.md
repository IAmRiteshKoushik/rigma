# WebSocket Orderbook (Simulated)

This project is an introduction to web-sockets using Go lang without using 
any framework. It uses the standard packages and simulates an orderbook.

### Steps to run:
1. Repo setup
```go
git clone https://github.com/IAmRiteshKoushik/go.sock-server
cd go.sock-server
go mod tidy
```

2. Open two browsers side by side -> inspect element -> go to `console`
3. Type the following in both console:
```js
let socket = new WebSocket("ws://localhost:3000/orderbook")
```
```js
socket.onmessage = (event) => { console.log("Order received: ", event.data); }
```

You should see a stream of messages being sent to both instances of your browser 
console at a one second interval.
