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
var connRate = flag.Int("rate", 0, "Connection rate in connections / second. Zero means no rate limit.")
var wait = flag.Duration("wait", 5*time.Minute, "time to wait before tearing down connections again")

var timePerConn time.Duration

func main() {

	flag.Parse()

	if *address == "" {
		flag.Usage()
		os.Exit(1)
	}

	if *connRate > 0 {
		timePerConn = time.Second / time.Duration(*connRate)
	}

	connect(*address, *maxConn, *wait, *useTLS)

}

func connect(addr string, maxConn int, wait time.Duration, useTLS bool) {
	var connections []net.Conn

	c := &tls.Config{
		InsecureSkipVerify: true,
	}

	d := &net.Dialer{
		KeepAlive: 10 * time.Second,
	}

	fmt.Println("Opening connections...")
	tTotalStart := time.Now()
	for i := 1; i <= maxConn; i++ {
		tStart := time.Now()
		var conn net.Conn
		var err error
		if useTLS {
			conn, err = tls.DialWithDialer(d, "tcp", addr, c)
		} else {
			conn, err = d.Dial("tcp", addr)
		}
		tEnd := time.Now()
		dur := tEnd.Sub(tStart)
		totalDur := tEnd.Sub(tTotalStart)
		rate := float64(i) / float64(totalDur/time.Second)
		if err != nil {
			fmt.Printf("%v (%d, took %s)\n", err, i, dur)
		} else if conn != nil {
			fmt.Printf("%s -> %s (%d, took %s) (rate: %.1f/s, time: %s)\n", conn.LocalAddr().String(), conn.RemoteAddr().String(), i, dur, rate, totalDur)
			connections = append(connections, conn)
		}
		if dur < timePerConn {
			time.Sleep(timePerConn - dur)
		}
	}

	fmt.Printf("\nWaiting for %s...\n\n", wait)
	time.Sleep(wait)
	fmt.Println("Closing connections...")
	for _, connection := range connections {
		err := connection.Close()

		if err != nil {
			fmt.Println(err)
		}
	}

	fmt.Println("Done.")
}
