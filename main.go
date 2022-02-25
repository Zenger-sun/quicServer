package main

import (
	"log"
	"net/http"

	"github.com/lucas-clemente/quic-go/http3"
)

const (
	WEB_DIR = "./web"
	ADDR = "127.0.0.1:8000"
	CERT = "./cert/server.pem"
	PRIV = "./cert/server.key"
)

func main() {
	log.Println("start https server on", ADDR)

	http.Handle("/", http.FileServer(http.Dir(WEB_DIR)))
	err := http3.ListenAndServeQUIC(ADDR, CERT, PRIV, nil)
	if err != nil {
		panic(err)
	}
}