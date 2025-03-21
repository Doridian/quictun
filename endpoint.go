package main

import (
	"io"
	"log"
	"net"
	"strings"
	"time"

	"github.com/quic-go/quic-go"
)

var quicStream quic.Stream

func runLocalEndpoint() error {
	if strings.HasPrefix(*localTunAddr, "@") {
		return runLocalListener()
	}
	return runLocalDialer()
}

func runLocalListener() error {
	addr := strings.TrimPrefix(*localTunAddr, "@")
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	go func() {
		goErr := runOneConnection(listener)
		log.Println("runLocalListener done")
		if goErr != nil {
			fatalProgram(goErr)
		}
		closeProgram()
	}()

	return nil
}

func runLocalDialer() error {
	localTunStr := *localTunAddr

	go func() {
		var conn net.Conn
		var err error

		for {
			conn, err = net.Dial("tcp", localTunStr)
			if err != nil {
				log.Printf("Dial failed, retrying in 1s: %v", err)
				time.Sleep(1 * time.Second)
				continue
			}
			break
		}

		err = handleEndpointConn(conn)
		log.Println("runLocalDialer done")
		if err != nil {
			fatalProgram(err)
		}
		closeProgram()
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

	log.Println("Got conn, copying...")

	var errChan = make(chan error, 2)
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
