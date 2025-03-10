package main

import (
	"log"
	"os"
	"os/exec"
)

var bgCommands = []*exec.Cmd{}

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
