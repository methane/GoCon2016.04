package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"sync"
	"time"
)

var (
	nClients  int
	nMessages int
	portNo    int
)

func reader(wg *sync.WaitGroup, conn *net.TCPConn) {
	io.Copy(ioutil.Discard, conn)
	wg.Done()
}

func sender(start <-chan struct{}, conn *net.TCPConn) {
	<-start
	defer conn.CloseWrite()
	for i := 0; i < nMessages; i++ {
		_, err := conn.Write([]byte("Hello, World\n"))
		if err != nil {
			panic(err)
		}
	}
}

func startClient(wg *sync.WaitGroup, start <-chan struct{}) {
	conn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", portNo))
	if err != nil {
		panic(err)
	}
	tconn := conn.(*net.TCPConn)
	go reader(wg, tconn)
	go sender(start, tconn)
}

func main() {
	flag.IntVar(&nClients, "clients", 100, "Number of clients")
	flag.IntVar(&nMessages, "messages", 1000, "Number of messages per each clients")
	flag.IntVar(&portNo, "port", 8000, "Port number")
	flag.Parse()

	start := make(chan struct{})
	wg := sync.WaitGroup{}
	wg.Add(nClients)

	for i := 0; i < nClients; i++ {
		startClient(&wg, start)
	}
	t0 := time.Now()
	close(start)
	wg.Wait()
	fmt.Println(time.Now().Sub(t0))
}
