package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/lucas-clemente/quic-go"
	"github.com/lucas-clemente/quic-go/http3"
	"io"
	"log"
	"net/http"
)

const (
	WEB_DIR = "./web"
	ADDR = "127.0.0.1:8000"
	CERT = "./cert/server.pem"
	PRIV = "./cert/server.key"
)

type loggingWriter struct{ io.Writer }

func main() {
	log.Println("start quic server on", ADDR)

	//httpServer()
	echoServer()
}

func httpServer() {
	http.Handle("/", http.FileServer(http.Dir(WEB_DIR)))
	err := http3.ListenAndServeQUIC(ADDR, CERT, PRIV, nil)
	if err != nil {
		panic(err)
	}
}

func echoServer() {
	const pack_len = 4
	var err error

	certs := make([]tls.Certificate, 1)
	certs[0], err = tls.LoadX509KeyPair(CERT, PRIV)
	if err != nil {
		panic(err)
	}

	config := &tls.Config{
		Certificates: certs,
		NextProtos:   []string{"quic-echo-example"},
	}

	listener, _ := quic.ListenAddr(ADDR, config, nil)
	defer listener.Close()

	for {
		sess, _ := listener.Accept(context.Background())
		stream, _ := sess.AcceptStream(context.Background())

		go func(stream quic.Stream) {
			defer stream.Close()
			b := make([]byte, 32 * 1024)

			for {
				_, err := io.ReadAtLeast(stream, b, pack_len)
				switch err {
				case nil:
				default:
					fmt.Println(err)
					return
				}

				l := uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24
				pack := b[pack_len:pack_len+l]
				fmt.Println(string(pack))

				stream.Write(pack)
			}

		}(stream)
	}
}