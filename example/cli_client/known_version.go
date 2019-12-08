package main

import (
	"fmt"
	"log"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/the729/go-libra/client"
)

type knownVersionState struct {
	KnownVersion uint64   `toml:"known_version"`
	Subtrees     [][]byte `toml:"subtrees"`
}

func loadKnownVersion(c *client.Client, filepath string) error {
	v := &knownVersionState{}
	_, err := toml.DecodeFile(filepath, v)
	if err != nil {
		return fmt.Errorf("toml decode file error: %v", err)
	}
	if len(v.Subtrees) == 0 {
		return fmt.Errorf("empty subtree")
	}
	err = c.SetKnownVersion(v.KnownVersion, v.Subtrees)
	if err != nil {
		return fmt.Errorf("set accumulator error: %v", err)
	}
	return nil
}

func saveKnownVersion(c *client.Client, filepath string) error {
	v := &knownVersionState{}
	v.KnownVersion, v.Subtrees = c.GetKnownVersion()
	f, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("create file error: %v", err)
	}
	err = toml.NewEncoder(f).Encode(v)
	if err != nil {
		log.Printf("cannot encode toml file: %v", err)
	}
	return nil
}
