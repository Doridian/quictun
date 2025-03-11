package main

import (
	"flag"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

var remoteAddr = flag.String("remote-addr", "remote.example.com", "Remote to open tunnel with")
var quicPort = flag.Int("quic-port", 1234, "Port to connect tunnel with")
var gitVersion = flag.Bool("git-version", false, "Use git rev-parse to send ref to remote")
var useBinary = flag.String("use-binary", "", "Use a specific binary for the remote (instead of \"go run\")")

var localTunAddr = flag.String("local-tunnel-addr", ":1235", ":PORT for listener, otherwise connect IP:PORT")
var remoteTunAddr = flag.String("remote-tunnel-addr", ":1235", ":PORT for listener, otherwise connect IP:PORT")

var readyPidPtr = flag.Int("ready-pid", 0, "PID to send ready signal to")
var readySignalInt = flag.Int("ready-signal", int(syscall.SIGUSR1), "Signal to send to ready-pid")

func sendReady() {
	log.Printf("Ready signal")

	readyPid := *readyPidPtr
	if readyPid == 0 {
		return
	}
	*readyPidPtr = 0
	readySignal := syscall.Signal(*readySignalInt)

	err := syscall.Kill(readyPid, readySignal)
	if err != nil {
		log.Printf("Failed to send %v to %d: %v", readySignal, readyPid, err)
	}
}

func main() {
	flag.Parse()

	if *remoteAddr == ":" {
		log.SetPrefix("server: ")
	} else {
		log.SetPrefix("client: ")
	}

	log.Printf("Embedded version: %s", VERSION)

	if *gitVersion {
		gitCmd := exec.Command("git", "rev-parse", "HEAD")
		gitCmd.Stderr = os.Stderr
		stdout, err := gitCmd.StdoutPipe()
		if err != nil {
			fatalProgram(err)
		}
		err = gitCmd.Start()
		if err != nil {
			fatalProgram(err)
		}
		data, err := io.ReadAll(stdout)
		if err != nil {
			fatalProgram(err)
		}
		err = gitCmd.Wait()
		if err != nil {
			fatalProgram(err)
		}
		VERSION = strings.Trim(string(data), " \r\n\t")
	}

	log.Printf("Used version: %s", VERSION)

	cfg.QUICPort = *quicPort

	err := runLocalEndpoint()
	if err != nil {
		fatalProgram(err)
	}

	if *remoteAddr == ":" {
		err = runServer()
	} else {
		err = runClient()
	}
	log.Println("main done")
	if err != nil {
		fatalProgram(err)
	}
	closeProgram()
}
