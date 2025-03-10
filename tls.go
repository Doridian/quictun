package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"log"
	"math/big"
	"time"
)

func generateCert() (tls.Certificate, error) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return tls.Certificate{}, err
	}
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(time.Hour),
		DNSNames:     []string{"quictun"},
	}
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
	cfg.Certificate = certDER

	return tls.X509KeyPair(certPEM, keyPEM)
}

func fixedCertGetter[T interface{}](cert tls.Certificate) func(T) (*tls.Certificate, error) {
	return func(T) (*tls.Certificate, error) {
		return &cert, nil
	}
}

func verifyRemoteCert(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
	remoteExpected, err := x509.ParseCertificate(remoteCfg.Certificate)
	if err != nil {
		log.Printf("Failed to parse remote certificate: %v", err)
		return err
	}

	for _, rawCert := range rawCerts {
		cert, err := x509.ParseCertificate(rawCert)
		if err != nil {
			return err
		}
		if cert.Equal(remoteExpected) {
			log.Printf("Peer cert found!")
			return nil
		}
	}

	return errors.New("no matching certificate found")
}

func CommonTLSConfig() *tls.Config {
	return &tls.Config{
		ClientAuth:            tls.RequireAnyClientCert,
		InsecureSkipVerify:    true,
		VerifyPeerCertificate: verifyRemoteCert,
		NextProtos:            []string{"quictun"},
		ServerName:            "quictun",
	}
}
