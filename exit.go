package main

import (
	"log"
	"os"
	"os/exec"
	"os/signal"
)

var bgCommands = []*exec.Cmd{}

func init() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	go func() {
		for range sigChan {
			closeProgram()
		}
	}()
}

func addBGCommand(cmd *exec.Cmd) {
	bgCommands = append(bgCommands, cmd)
}

func killAllBGCommands() {
	log.Println("Killing all background commands")

	for _, cmd := range bgCommands {
		log.Printf("Killing command: %v (res %v)", cmd,
			cmd.Process.Kill())
	}
}

func closeProgram() {
	killAllBGCommands()

	log.Println("Closing program")
	os.Exit(0)
}

func fatalProgram(val interface{}) {
	killAllBGCommands()

	log.Fatalln(val)
	os.Exit(1)
}
