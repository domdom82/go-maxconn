package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net"
	"os"
	"time"
)

var address = flag.String("address", "", "the address to connect to e.g. www.myhost.com:443")
var useTLS = flag.Bool("tls", true, "whether to use TLS to connect or not")
var maxConn = flag.Int("connections", 100, "maximum number of concurrent connections to open")
var wait = flag.Duration("wait", 5*time.Minute, "time to wait before tearing down connections again")

func main() {

	flag.Parse()

	if *address == "" {
		flag.Usage()
		os.Exit(1)
	}

	if *useTLS {
		connectTLS(*address, *maxConn, *wait)
	} else {
		connectTCP(*address, *maxConn, *wait)
	}

}

func connectTCP(addr string, maxConn int, wait time.Duration) {
	var connections []net.Conn

	d := &net.Dialer{
		KeepAlive: 10 * time.Second,
	}

	fmt.Println("Opening TCP connections...")
	for i := 1; i <= maxConn; i++ {
		conn, err := d.Dial("tcp", addr)
		if err != nil {
			fmt.Printf("%v (%d)\n", err, i)
		}
		if conn != nil {
			fmt.Printf("%s -> %s (%d)\n", conn.LocalAddr().String(), conn.RemoteAddr().String(), i)
		}

		connections = append(connections, conn)
	}

	fmt.Printf("\nWaiting for %s...\n\n", wait)
	time.Sleep(wait)
	fmt.Println("Closing TCP connections...")
	for _, connection := range connections {
		err := connection.Close()

		if err != nil {
			fmt.Println(err)
		}
	}

	fmt.Println("Done.")
}

func connectTLS(addr string, maxConn int, wait time.Duration) {
	var connections []*tls.Conn

	d := &net.Dialer{
		KeepAlive: 10 * time.Second,
	}

	fmt.Println("Opening TLS connections...")
	for i := 1; i <= maxConn; i++ {
		conn, err := tls.DialWithDialer(d, "tcp", addr, nil)
		if err != nil {
			fmt.Printf("%v (%d)\n", err, i)
		}
		if conn != nil {
			fmt.Printf("%s -> %s (%d)\n", conn.LocalAddr().String(), conn.RemoteAddr().String(), i)
		}

		connections = append(connections, conn)
	}

	fmt.Printf("\nWaiting for %s...\n\n", wait)
	time.Sleep(wait)
	fmt.Println("Closing TLS connections...")
	for _, connection := range connections {
		err := connection.Close()
		if err != nil {
			fmt.Println(err)
		}
	}

	fmt.Println("Done.")
}
