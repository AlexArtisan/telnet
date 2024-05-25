package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
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
		for {
			buf := make([]byte, 1e6)

			n, err := conn.Read(buf)

			if err != nil {
				if err.Error() == "EOF" {
					os.Exit(0)
				}

				log.Fatalln("Read error:", err)
			}

			if n > 0 {
				fmt.Println(string(buf[:n-1]))
			}
		}
	}()

	go func() {
		for {
			reader := bufio.NewReader(os.Stdin)

			str, err := reader.ReadString('\n')
			if err != nil {
				log.Fatalln("Stdin read error:", err)
			}

			_, err = conn.Write([]byte(str))
			if err != nil {
				log.Fatalln("Send error:", err)
			}
		}
	}()

	wg.Wait()

}
