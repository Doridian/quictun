package main

import (
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/quic-go/quic-go"
)

var quicStream quic.Stream

func runLocalEndpoint() error {
	if strings.HasPrefix(*localTunAddr, ":") {
		return runLocalListener()
	}
	return runLocalDialer()
}

func runLocalListener() error {
	port := strings.TrimPrefix(*localTunAddr, ":")
	portInt, err := strconv.Atoi(port)
	if err != nil {
		return err
	}

	listener, err := net.ListenTCP("tcp", &net.TCPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: portInt,
	})
	if err != nil {
		return err
	}

	go func() {
		goErr := runOneConnection(listener)
		log.Println("runLocalListener done")
		if goErr != nil {
			log.Fatalln(goErr)
		}
		os.Exit(0)
	}()

	return nil
}

func runLocalDialer() error {
	conn, err := net.Dial("tcp", *localTunAddr)
	if err != nil {
		return err
	}

	go func() {
		goErr := handleEndpointConn(conn)
		log.Println("runLocalDialer done")
		if goErr != nil {
			log.Fatalln(goErr)
		}
		os.Exit(0)
	}()

	return nil
}

func runOneConnection(listener net.Listener) error {
	conn, err := listener.Accept()
	if err != nil {
		return err
	}
	return handleEndpointConn(conn)
}

func handleEndpointConn(conn net.Conn) error {
	defer conn.Close()

	for quicStream == nil {
		time.Sleep(10 * time.Millisecond)
	}

	var errChan = make(chan error, 2)
	defer close(errChan)
	go func() {
		_, goErr := io.Copy(quicStream, conn)
		errChan <- goErr
	}()
	go func() {
		_, err := io.Copy(conn, quicStream)
		errChan <- err
	}()

	for err := <-errChan; err != nil; err = <-errChan {
		return err
	}
	return nil
}
