package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/quic-go/quic-go"
)

func mainRemoteListener() error {
	sshCmd, err := runRemoteListener()
	if err != nil {
		return err
	}

	go func() {
		goErr := sshCmd.Wait()
		log.Println("mainRemoteListener done")
		if goErr != nil {
			log.Fatalln(goErr)
		}
		os.Exit(0)
	}()

	return nil
}

func runRemoteListener() (*exec.Cmd, error) {
	sshCmd := &exec.Cmd{}
	sshCmd.Args = []string{"/usr/bin/ssh", *remoteAddr, "--", "go", "run", "github.com/Doridian/quictun@" + VERSION, "-remote-addr", ":", "-quic-port", strconv.Itoa(*quicPort), "-local-tunnel-addr", *remoteTunAddr}
	sshCmd.Path = sshCmd.Args[0]

	stdin, err := sshCmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	defer stdin.Close()
	stdout, err := sshCmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	defer stdout.Close()
	sshCmd.Stderr = os.Stderr

	log.Println("Starting SSH command")

	err = sshCmd.Start()
	if err != nil {
		return nil, err
	}

	err = configLoop(stdout, stdin)
	if err != nil {
		return nil, err
	}

	return sshCmd, nil
}

func generateClientTLSConfig() (*tls.Config, error) {
	tlsClientCert, err := generateCert()
	if err != nil {
		return nil, err
	}

	tlsConfig := CommonTLSConfig()
	tlsConfig.GetClientCertificate = fixedCertGetter[*tls.CertificateRequestInfo](tlsClientCert)
	return tlsConfig, nil
}

func runClient() error {
	err := mainRemoteListener()
	if err != nil {
		return err
	}

	tlsConf, err := generateClientTLSConfig()
	if err != nil {
		return err
	}

	conn, err := quic.DialAddr(context.Background(), fmt.Sprintf("[%s]:%d", *remoteAddr, remoteCfg.QUICPort), tlsConf, nil)
	if err != nil {
		return err
	}
	defer conn.CloseWithError(0, "")

	stream, err := conn.OpenStreamSync(context.Background())
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
