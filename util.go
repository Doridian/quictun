package main

import (
	"log"
	"time"

	"github.com/quic-go/quic-go"
)

func happyLoop(stream quic.Stream) {
	defer stream.Close()
	quicStream = stream

	sendReady()
	log.Printf("Stream open!")
	for {
		time.Sleep(1 * time.Second)
	}
}
