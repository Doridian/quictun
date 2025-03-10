package main

import (
	"encoding/json"
	"io"
	"log"
)

var VERSION = "main"

var cfg = &Config{}
var remoteCfg = &Config{}

type Config struct {
	Certificate string `json:"certificate"`
	PublicKey   string `json:"public_key"`

	QUICPort  int `json:"quic_port"`
	LocalPort int `json:"local_port"`
}

func readRemoteConfig(r io.ReadCloser) error {
	defer r.Close()

	unmarshaler := json.NewDecoder(r)
	return unmarshaler.Decode(remoteCfg)
}

func writeConfig(w io.WriteCloser) error {
	defer w.Close()

	marshaler := json.NewEncoder(w)
	return marshaler.Encode(cfg)
}

func configLoop(r io.ReadCloser, w io.WriteCloser) error {
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
