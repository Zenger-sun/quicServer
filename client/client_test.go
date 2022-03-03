package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/binary"
	"fmt"
	"github.com/lucas-clemente/quic-go"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/lucas-clemente/quic-go/http3"
)

const (
	URL  = "https://127.0.0.1:8000/"
	ADDR = "127.0.0.1:8000"
	CA   = "../cert/ca.crt"
)

func Test_httpServer(t *testing.T) {
	rootCa, err := os.ReadFile(CA)
	if err != nil {
		panic("failed to read root certificate")
	}

	roots := x509.NewCertPool()
	ok := roots.AppendCertsFromPEM(rootCa)
	if !ok {
		panic("failed to parse root certificate")
	}

	tr := &http3.RoundTripper{
		TLSClientConfig: &tls.Config{
			RootCAs: roots,
		},
	}

	client := &http.Client{Transport: tr}

	_, err = client.Get(URL)
	if err != nil {
		fmt.Errorf("get %s error", URL)
		panic(err)
	}
}

func Test_echoServer(t *testing.T) {
	rootCa, err := os.ReadFile(CA)
	if err != nil {
		panic("failed to read root certificate")
	}

	roots := x509.NewCertPool()
	ok := roots.AppendCertsFromPEM(rootCa)
	if !ok {
		panic("failed to parse root certificate")
	}

	tlsConf := &tls.Config{
		RootCAs:    roots,
		NextProtos: []string{"quic-echo-example"},
	}

	session, err := quic.DialAddr(ADDR, tlsConf, nil)
	if err != nil {
		panic(err)
	}

	stream, err := session.OpenStreamSync(context.Background())
	if err != nil {
		panic(err)
	}

	for i := 0; i < 100; i++ {
		msg := fmt.Sprintf("message%d", i)
		buff := new(bytes.Buffer)
		binary.Write(buff, binary.LittleEndian, int32(len(msg)))
		buff.Write([]byte(msg))

		_, err = stream.Write(buff.Bytes())
		if err != nil {
			panic(err)
		}

		pack := make([]byte, buff.Len() - 4)
		io.ReadFull(stream, pack)
		fmt.Println(string(pack))
	}
}

func Benchmark_HttpServer(b *testing.B) {
	rootCa, err := os.ReadFile(CA)
	if err != nil {
		panic("failed to read root certificate")
	}

	roots := x509.NewCertPool()
	ok := roots.AppendCertsFromPEM(rootCa)
	if !ok {
		panic("failed to parse root certificate")
	}

	tr := &http3.RoundTripper{
		TLSClientConfig: &tls.Config{
			RootCAs: roots,
		},
	}

	client := &http.Client{Transport: tr}

	for i := 0; i < b.N; i++ {
		_, err := client.Get(URL)
		if err != nil {
			fmt.Errorf("get %s error", URL)
			panic(err)
		}
	}
}

func Benchmark_StreamServer(b *testing.B) {
	rootCa, err := os.ReadFile(CA)
	if err != nil {
		panic("failed to read root certificate")
	}

	roots := x509.NewCertPool()
	ok := roots.AppendCertsFromPEM(rootCa)
	if !ok {
		panic("failed to parse root certificate")
	}

	tlsConf := &tls.Config{
		RootCAs:    roots,
		NextProtos: []string{"quic-echo-example"},
	}

	session, err := quic.DialAddr(ADDR, tlsConf, nil)
	if err != nil {
		panic(err)
	}

	stream, err := session.OpenStreamSync(context.Background())
	if err != nil {
		panic(err)
	}

	for i := 0; i < b.N; i++ {
		msg := fmt.Sprintf("message%d", i)
		buff := new(bytes.Buffer)
		binary.Write(buff, binary.LittleEndian, int32(len(msg)))
		buff.Write([]byte(msg))

		_, err = stream.Write(buff.Bytes())
		if err != nil {
			panic(err)
		}

		pack := make([]byte, buff.Len() - 4)
		io.ReadFull(stream, pack)
	}
}