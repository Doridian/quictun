package main

import (
	"context"
	"crypto/tls"
	"flag"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/quic-go/quic-go"
)

var remoteAddr = flag.String("remote-addr", "remote.example.com", "Remote to open tunnel with")
var remotePort = flag.Int("remote-port", 1234, "Port to connect tunnel with")
var boundPort = flag.Int("bind-port", 1235, "Local port to bind in/out to")
var gitVersion = flag.Bool("git-version", false, "Use git rev-parse to send ref to remote")

func main() {
	flag.Parse()

	if *gitVersion {
		gitCmd := exec.Command("git", "rev-parse", "HEAD")
		gitCmd.Stderr = os.Stderr
		stdout, err := gitCmd.StdoutPipe()
		if err != nil {
			panic(err)
		}
		err = gitCmd.Start()
		if err != nil {
			panic(err)
		}
		data, err := io.ReadAll(stdout)
		if err != nil {
			panic(err)
		}
		err = gitCmd.Wait()
		if err != nil {
			panic(err)
		}
		VERSION = strings.Trim(string(data), " \r\n\t")
	}

	var err error
	if *remoteAddr == "" {
		err = runServer()
	} else {
		err = mainRemoteListener()
		if err == nil {
			err = runClient()
		}
	}
	log.Println("main done")
	if err != nil {
		panic(err)
	}
	os.Exit(0)
}

func mainRemoteListener() error {
	sshCmd, err := runRemoteListener()
	if err != nil {
		return err
	}

	go func() {
		goErr := sshCmd.Wait()
		log.Println("mainRemoteListener done")
		if goErr != nil {
			panic(goErr)
		}
		os.Exit(0)
	}()

	return nil
}

func runRemoteListener() (*exec.Cmd, error) {
	sshCmd := &exec.Cmd{}
	sshCmd.Args = []string{"/usr/bin/ssh", *remoteAddr, "go", "run", "github.com/Doridian/quictun@" + VERSION, "-remote-addr", "", "-remote-port", strconv.Itoa(*remotePort), "-bind-port", strconv.Itoa(*boundPort)}
	sshCmd.Path = sshCmd.Args[0]

	stdin, err := sshCmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	stdout, err := sshCmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	sshCmd.Stderr = os.Stderr

	err = sshCmd.Start()
	if err != nil {
		return nil, err
	}

	err = writeConfig(stdin)
	if err != nil {
		return nil, err
	}

	err = readRemoteConfig(stdout)
	if err != nil {
		return nil, err
	}

	return sshCmd, nil
}

func runServer() error {
	err := writeConfig(os.Stdout)
	if err != nil {
		return err
	}

	err = readRemoteConfig(os.Stdin)
	if err != nil {
		return err
	}

	return nil
}

func runClient() error {
	time.Sleep(time.Hour * 24)

	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quictun"},
	}
	conn, err := quic.DialAddr(context.Background(), *remoteAddr, tlsConf, nil)
	if err != nil {
		return err
	}
	defer conn.CloseWithError(0, "")

	stream, err := conn.OpenStreamSync(context.Background())
	if err != nil {
		return err
	}
	defer stream.Close()

	return nil
}
