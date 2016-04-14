// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// Chat is a server that lets clients chat with each other.
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime/pprof"
)

type client chan<- string // an outgoing message channel

type cmdEnter struct {
	c client
}

type cmdLeave struct {
	c client
}

type cmdMessage struct {
	m string
}

var (
	cmdQueue = make(chan interface{}, 128)
)

func broadcaster() {
	clients := make(map[client]bool) // all connected clients
	for cmd := range cmdQueue {
		switch cmd := cmd.(type) {
		case cmdEnter:
			clients[cmd.c] = true
		case cmdLeave:
			delete(clients, cmd.c)
			close(cmd.c)
		case cmdMessage:
			for cli := range clients {
				cli <- cmd.m
			}
		}
	}
}

func handleConn(conn net.Conn) {
	ch := make(chan string, 128) // outgoing client messages
	go clientWriter(conn, ch)

	who := conn.RemoteAddr().String()
	ch <- "You are " + who
	cmdQueue <- cmdMessage{who + " has arrived"}
	cmdQueue <- cmdEnter{ch}

	input := bufio.NewScanner(conn)
	for input.Scan() {
		cmdQueue <- cmdMessage{who + ": " + input.Text()}
	}

	cmdQueue <- cmdLeave{ch}
	cmdQueue <- cmdMessage{who + " has left"}
	conn.Close()
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg) // NOTE: ignoring network errors
	}
}

func listen(l net.Listener) {
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	go broadcaster()
	go listen(listener)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	for range c {
		prof := pprof.Lookup("goroutine")
		prof.WriteTo(os.Stderr, 1)
	}
}
