package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"sync"
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
	for i := 0; i < nMessages; i++ {
		conn.Write([]byte("Hello, World\n"))
	}
	conn.CloseWrite()
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
	flag.IntVar(&nClients, "clients", 1000, "Number of clients")
	flag.IntVar(&nMessages, "messages", 10000, "Number of messages per each clients")
	flag.IntVar(&portNo, "port", 8000, "Port number")
	flag.Parse()

	start := make(chan struct{})
	wg := sync.WaitGroup{}

	for i := 0; i < nClients; i++ {
		startClient(&wg, start)
	}
	close(start)
	wg.Wait()
}
