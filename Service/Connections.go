package service

import (
	"encoding/json"
	"log"
	"os"
)

type Connections map[string]Credentials

func (c *Connections) close() {
	for n, c := range Cons {
		if err := c.Close(); err != nil {
			log.Printf("[ERROR] %s: %v", n, err)
		}
	}
}

var Cons Connections


func SaveConnections(path string) error {
	if len(Cons) == 0 {
		return nil
	}
	Cons.close()

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	j := json.NewEncoder(f)
	j.SetIndent("", "\t")

	if err = j.Encode(Cons); err != nil {
		return err
	}

	return nil
}

func LoadConnections(path string) error {
	bytes, err := os.ReadFile(path)
	if err != nil {
		Cons = make(Connections)
		return err
	}

	if err = json.Unmarshal(bytes, &Cons); err != nil {
		Cons = make(Connections)
		return err
	}

	return nil
}
