package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"math/big"
)

func generateCert() (tls.Certificate, error) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return tls.Certificate{}, err
	}
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		return tls.Certificate{}, err
	}

	keyBytes, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		return tls.Certificate{}, err
	}

	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "ECDSA PRIVATE KEY", Bytes: keyBytes})

	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	cfg.Certificate = string(certPEM)

	return tls.X509KeyPair(certPEM, keyPEM)
}

func fixedCertGetter[T interface{}](cert tls.Certificate) func(T) (*tls.Certificate, error) {
	return func(T) (*tls.Certificate, error) {
		return &cert, nil
	}
}

func getRemoteCertificate[T interface{}](T) (*tls.Certificate, error) {
	cert, _ := pem.Decode([]byte(remoteCfg.Certificate))
	if cert == nil {
		return nil, fmt.Errorf("Failed to decode certificate")
	}
	return &tls.Certificate{Certificate: [][]byte{cert.Bytes}}, nil
}
