package main

import (
	"crypto/tls"
	"io"
	"log"
	"net"
	"os"
	"sync"
)

func dial(network, host, port string) (net.Conn, error) {
	if port == "443" {
		return tls.Dial(network, host+":"+port, &tls.Config{
			ServerName: host,
		})
	}

	return net.Dial(network, host+":"+port)
}

type reader struct {
	r io.Reader
}

func (r *reader) Read(p []byte) (n int, err error) {
	return r.r.Read(p)
}

func newReader(r io.Reader) *reader {
	return &reader{
		r: r,
	}
}

func main() {

	if len(os.Args) < 3 {
		log.Fatalln("Host and port not specified")
	}

	host := os.Args[1]
	port := os.Args[2]

	conn, err := dial("tcp", host, port)
	if err != nil {
		log.Fatalln("Dial error:", err)
	}
	defer conn.Close()

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		if _, err := io.Copy(os.Stdout, conn); err != nil {
			log.Fatalln("Receiving error:", err)
		}
	}()

	go func() {
		defer wg.Done()
		if _, err := io.Copy(conn, newReader(os.Stdin)); err != nil {
			log.Fatalln("Sending error:", err)
		}
	}()

	wg.Wait()

}
