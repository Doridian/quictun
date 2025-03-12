package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
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

func runServer() error {
	tlsConfig, err := generateServerTLSConfig()
	if err != nil {
		return err
	}

	listener, err := quic.ListenAddr(fmt.Sprintf(":%d", cfg.QUICPort), tlsConfig, nil)
	if err != nil {
		return err
	}
	cfg.QUICPort = listener.Addr().(*net.UDPAddr).Port

	err = configLoop(os.Stdin, os.Stdout)
	if err != nil {
		_ = listener.Close()
		return err
	}

	conn, err := listener.Accept(context.Background())
	_ = listener.Close()
	if err != nil {
		return err
	}
	defer conn.CloseWithError(quic.ApplicationErrorCode(0), "")

	stream, err := conn.AcceptStream(context.Background())
	if err != nil {
		return err
	}

	happyLoop(stream)
	return nil
}
