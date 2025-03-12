package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
)

var VERSION = "main"

var cfg Config
var remoteCfg Config

type Config struct {
	Certificate []byte `json:"certificate"`

	QUICPort int `json:"quic_port"`
}

func configLoop(r io.ReadCloser, w io.WriteCloser) error {
	defer w.Close()
	defer r.Close()

	log.Println("Sending local config")

	cfgBuf := &bytes.Buffer{}
	multiW := io.MultiWriter(w, cfgBuf)
	marshaler := json.NewEncoder(multiW)
	err := marshaler.Encode(cfg)
	if err != nil {
		return err
	}
	log.Printf("Local config sent: %s", cfgBuf.String())

	log.Println("Waiting for remote config")

	cfgBuf.Reset()
	multiR := io.TeeReader(r, cfgBuf)
	unmarshaler := json.NewDecoder(multiR)
	unmarshaler.DisallowUnknownFields()
	err = unmarshaler.Decode(&remoteCfg)
	if err != nil {
		return err
	}

	log.Printf("Remote config received: %s", cfgBuf.String())
	return nil
}
