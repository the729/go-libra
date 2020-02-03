package main

import (
	"fmt"
	"log"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/the729/go-libra/client"
)

func saveClientState(c *client.Client, filepath string) error {
	v := c.GetState()
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

func newClientFromWaypointOrFile(serverAddr, waypoint, filepath string) (*client.Client, error) {
	if waypoint != "" {
		return client.New(serverAddr, waypoint)
	}

	v := &client.State{}
	_, err := toml.DecodeFile(filepath, v)
	if err != nil {
		return nil, fmt.Errorf("toml decode file error: %v", err)
	}
	return client.NewFromState(serverAddr, v)
}
