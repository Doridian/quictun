package main

import (
	"encoding/json"
	"io"
	"os"
)

var VERSION = "main"

var cfg = &Config{}
var remoteCfg = &Config{}

type Config struct {
	Certificate string `json:"certificate"`
	Key         string `json:"key"`
}

func readRemoteConfig(r io.ReadCloser) error {
	defer r.Close()

	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, remoteCfg)
}

func writeConfig(w io.WriteCloser) error {
	defer w.Close()

	marshaler := json.NewEncoder(w)
	return marshaler.Encode(cfg)
}
