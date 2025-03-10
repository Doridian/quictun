package main

import (
	"log"
	"os"
)

func closeProgram() {
	log.Println("Closing program")
	os.Exit(0)
}

func fatalProgram(val interface{}) {
	log.Fatalln(val)
	os.Exit(1)
}
