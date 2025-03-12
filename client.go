package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"os/exec"

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
			fatalProgram(goErr)
		}
		closeProgram()
	}()

	return nil
}

func runRemoteListener() (*exec.Cmd, error) {
	addr, port, err := SplitAddr(cfg.QUICAddr)
	if err != nil {
		return nil, err
	}

	sshCmd := &exec.Cmd{}
	sshCmd.Args = []string{"/usr/bin/ssh", addr, "--"}
	if *useBinary == "" {
		sshCmd.Args = append(sshCmd.Args, "go", "run", "github.com/Doridian/quictun@"+VERSION)
	} else {
		sshCmd.Args = append(sshCmd.Args, *useBinary)
	}
	sshCmd.Args = append(sshCmd.Args, "-quic-addr", fmt.Sprintf("@:%d", port), "-local-tunnel-addr", *remoteTunAddr)
	sshCmd.Path = sshCmd.Args[0]
	addBGCommand(sshCmd)

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

	_, remotePort, err := SplitAddr(remoteCfg.QUICAddr)
	if err != nil {
		return nil, err
	}
	localAddr, _, err := SplitAddr(cfg.QUICAddr)
	if err != nil {
		return nil, err
	}
	cfg.QUICAddr = fmt.Sprintf("%s:%d", localAddr, remotePort)

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
	tlsConf, err := generateClientTLSConfig()
	if err != nil {
		return err
	}

	err = mainRemoteListener()
	if err != nil {
		return err
	}

	conn, err := quic.DialAddr(context.Background(), cfg.QUICAddr, tlsConf, nil)
	if err != nil {
		return err
	}
	defer conn.CloseWithError(0, "")

	stream, err := conn.OpenStreamSync(context.Background())
	if err != nil {
		return err
	}

	return happyLoop(stream)
}
