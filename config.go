package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
)

var VERSION = "main"

var cfg Config
var remoteCfg Config

type Config struct {
	Certificate []byte `json:"certificate"`

	QUICAddr string `json:"quic_addr"`
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

func SplitAddr(addr string) (string, int, error) {
	pSplit := strings.LastIndex(addr, ":")
	if pSplit == -1 {
		return "", 0, fmt.Errorf("invalid address: %s", addr)
	}
	p, err := strconv.Atoi(addr[pSplit+1:])
	if err != nil {
		return "", 0, err
	}
	return addr[:pSplit], p, nil
}
