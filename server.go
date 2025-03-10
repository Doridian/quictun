package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/quic-go/quic-go"
)

func generateServerTLSConfig() (*tls.Config, error) {
	tlsCert, err := generateCert()
	if err != nil {
		return nil, err
	}

	return &tls.Config{
		ClientAuth:           tls.RequireAndVerifyClientCert,
		GetCertificate:       fixedCertGetter[*tls.ClientHelloInfo](tlsCert),
		GetClientCertificate: getRemoteCertificate[*tls.CertificateRequestInfo],
		NextProtos:           []string{"quictun"},
	}, nil
}

func runServer() error {
	tlsConfig, err := generateServerTLSConfig()
	if err != nil {
		return err
	}

	err = configLoop(os.Stdin, os.Stdout)
	if err != nil {
		return err
	}

	listener, err := quic.ListenAddr(fmt.Sprintf(":%d", cfg.QUICPort), tlsConfig, nil)
	if err != nil {
		return err
	}
	defer listener.Close()

	conn, err := listener.Accept(context.Background())
	if err != nil {
		return err
	}

	stream, err := conn.AcceptStream(context.Background())
	if err != nil {
		return err
	}
	quicStream = stream
	defer stream.Close()

	log.Printf("Stream open!")
	for {
		time.Sleep(1 * time.Second)
	}
}
