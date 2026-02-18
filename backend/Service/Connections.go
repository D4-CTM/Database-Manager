package service

import (
	"encoding/json"
	"log"
	"os"
	"path"
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

func SaveConnections(path string) {
	if len(Cons) == 0 {
		return 
	}
	Cons.close()

	f, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	j := json.NewEncoder(f)
	j.SetIndent("", "\t")

	if err = j.Encode(Cons); err != nil {
		log.Printf("[ERROR] %v",err)
	}

	for k := range Cons {
		delete(Cons, k)
	}
}

func LoadConnections(credsPath string) error {
	bytes, err := os.ReadFile(credsPath)
	if err != nil {
		Cons = make(Connections)
		if os.IsNotExist(err) {
			if internErr := os.MkdirAll(path.Dir(credsPath), 0700); internErr != nil {
				return internErr
			}
		}
		return err
	}

	if err = json.Unmarshal(bytes, &Cons); err != nil {
		Cons = make(Connections)
		return err
	}

	return nil
}
