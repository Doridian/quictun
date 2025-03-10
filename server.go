package main

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"

	"github.com/quic-go/quic-go"
)

func generateTLSConfig() (*tls.Config, error) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		return nil, err
	}

	keyBytes, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		return nil, err
	}

	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "ECDSA PRIVATE KEY", Bytes: keyBytes})

	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	cfg.Certificate = string(certPEM)

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		panic(err)
	}

	return &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		NextProtos:   []string{"quictun"},
	}, nil
}

func runServer() error {
	tlsConfig, err := generateTLSConfig()
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
		panic(err)
	}
	defer stream.Close()

	return err
}
