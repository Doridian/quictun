package main

import (
	"context"
	"crypto/tls"
	"os"

	"github.com/quic-go/quic-go"
)

func generateServerTLSConfig() (*tls.Config, error) {
	tlsCert, err := generateCert()
	if err != nil {
		return nil, err
	}

	tlsConfig := CommonTLSConfig()
	tlsConfig.GetCertificate = fixedCertGetter[*tls.ClientHelloInfo](tlsCert)
	return tlsConfig, nil
}

func runServer(addr string) error {
	tlsConfig, err := generateServerTLSConfig()
	if err != nil {
		return err
	}

	listener, err := quic.ListenAddr(addr, tlsConfig, nil)
	if err != nil {
		return err
	}
	defer listener.Close()

	cfg.QUICAddr = listener.Addr().String()

	err = configLoop(os.Stdin, os.Stdout)
	if err != nil {
		_ = listener.Close()
		return err
	}

	conn, err := listener.Accept(context.Background())
	if err != nil {
		return err
	}
	defer conn.CloseWithError(quic.ApplicationErrorCode(0), "")

	stream, err := conn.AcceptStream(context.Background())
	if err != nil {
		return err
	}

	return happyLoop(stream)
}
