package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

var VERSION = "main"

var cfg = &Config{}
var remoteCfg = &Config{}

type Config struct {
	Certificate string `json:"certificate"`
	Key         string `json:"key"`
}

func readRemoteConfig(r io.Reader) error {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, remoteCfg)
}

func writeConfig(w io.Writer) error {
	marshaler := json.NewEncoder(w)
	return marshaler.Encode(cfg)
}

func configLoop(r io.Reader, w io.Writer) error {
	log.Println("Writing local config")

	err := writeConfig(w)
	if err != nil {
		return err
	}

	log.Println("Waiting for remote config")

	err = readRemoteConfig(r)
	if err != nil {
		return err
	}

	log.Println("Remote config received")
	return nil
}
