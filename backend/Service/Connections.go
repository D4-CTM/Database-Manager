package service

import (
	"encoding/json"
	"log"
	"os"
	"path"
)

type Connections map[string]Credentials

func (c *Connections) Close() {
	for n, c := range Cons {
		if err := c.Close(); err != nil {
			log.Printf("[ERROR] %s: %v", n, err)
		}
	}
	log.Println("\n[INFO] Connections closed")
}

var (
	Cons Connections
	JsonPath string
)

func SaveConnections() {
	credsPath := path.Join(os.Getenv("CREDS_SUBDIR"), "data.json")
	log.Printf("[INFO] Saving %d credentials at: %s", len(Cons), credsPath)
	if len(Cons) == 0 {
		log.Printf("[WARNING] Credentials map is empty")
	}

	f, err := os.Create(credsPath)
	if err != nil {
		log.Printf("[ERROR] %v", err)
		return 
	}
	j := json.NewEncoder(f)
	j.SetIndent("", "\t")

	if err = j.Encode(Cons); err != nil {
		log.Printf("[ERROR] %v",err)
		return 
	}
	log.Printf("[INFO] Credentials saved")
}

func LoadConnections() error {
	credsPath := path.Join(os.Getenv("CREDS_SUBDIR"), "data.json")
	log.Printf("[INFO] Loading credentials from: %s", credsPath)
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

	log.Printf("[INFO] Credentials loaded: %d", len(Cons))
	return nil
}
