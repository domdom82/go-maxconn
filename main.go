package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"math"
	"net"
	"os"
	"time"
)

var address = flag.String("address", "", "the address to connect to e.g. www.myhost.com:443")
var useTLS = flag.Bool("tls", true, "whether to use TLS to connect or not")
var useWS = flag.Bool("ws", false, "whether to do a WebSocket upgrade or not")
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

	if *useWS {
		connectWs(*address, *maxConn, *wait, *useTLS)
	} else {
		connect(*address, *maxConn, *wait, *useTLS)
	}

}

func connectWs(addr string, maxConn int, wait time.Duration, useTLS bool) {
	var connections []*websocket.Conn

	c := &tls.Config{
		InsecureSkipVerify: true,
	}

	d := websocket.Dialer{
		TLSClientConfig: c,
	}

	scheme := ""

	if useTLS {
		scheme = "wss://"
	} else {
		scheme = "ws://"
	}

	url := fmt.Sprintf("%s%s", scheme, addr)

	fmt.Println("Opening connections...")
	tTotalStart := time.Now()
	for i := 1; i <= maxConn; i++ {
		go func(i int) {
			tStart := time.Now()
			var conn *websocket.Conn
			var err error
			conn, _, err = d.Dial(url, nil)
			tEnd := time.Now()
			dur := tEnd.Sub(tStart)
			totalDur := tEnd.Sub(tTotalStart)
			rate := float64(i) / float64(totalDur/time.Second)
			percent := (float64(i) / float64(maxConn)) * float64(100)
			step := int(math.Floor(float64(maxConn) * 0.05)) //log every 5%
			if err != nil {
				fmt.Printf("%v (%d / %d %.1f%%, took %s) (rate: %.1f/s, time: %s)\n", err, i, maxConn, percent, dur, rate, totalDur)
			} else if conn != nil {
				if maxConn > 100 && i%step == 0 || maxConn <= 100 {
					fmt.Printf("%s -> %s (%d / %d %.1f%%, took %s) (rate: %.1f/s, time: %s)\n", conn.LocalAddr().String(), conn.RemoteAddr().String(), i, maxConn, percent, dur, rate, totalDur)
				}
				connections = append(connections, conn)
			}
			if dur < timePerConn {
				time.Sleep(timePerConn - dur)
			}
		}(i)
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
		percent := (float64(i) / float64(maxConn)) * float64(100)
		step := int(math.Floor(float64(maxConn) * 0.05)) //log every 5%
		if err != nil {
			fmt.Printf("%v (%d / %d %.1f%%, took %s) (rate: %.1f/s, time: %s)\n", err, i, maxConn, percent, dur, rate, totalDur)
		} else if conn != nil {
			if maxConn > 100 && i%step == 0 || maxConn <= 100 {
				fmt.Printf("%s -> %s (%d / %d %.1f%%, took %s) (rate: %.1f/s, time: %s)\n", conn.LocalAddr().String(), conn.RemoteAddr().String(), i, maxConn, percent, dur, rate, totalDur)
			}
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
