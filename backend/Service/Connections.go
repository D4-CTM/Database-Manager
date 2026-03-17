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

type SlaveMasterMap map[string]SlaveMasterPair

var (
	Cons Connections
	SlaveMaster SlaveMasterMap
	JsonPath string
)

func saveToJson(credsPath string, val any) error {
	f, err := os.Create(credsPath)
	if err != nil {
		log.Printf("[ERROR] %v", err)
		return err
	}
	defer f.Close();

	j := json.NewEncoder(f)
	j.SetIndent("", "\t")

	if err = j.Encode(val); err != nil {
		log.Printf("[ERROR] %v",err)
		return err
	}

	return nil
}

func SaveConnections() {
	credsPath := os.Getenv("CREDS_SUBDIR")
	log.Printf("[INFO] Saving %d credentials at: %s", len(Cons), credsPath)
	if len(Cons) == 0 {
		log.Printf("[WARNING] Credentials map is empty")
	}

	if err := saveToJson(path.Join(credsPath, "credentials.json"), Cons); err != nil {
		log.Printf("[ERROR] %v",err)
		return
	}
	if err := saveToJson(path.Join(credsPath, "slave-master.json"), SlaveMaster); err != nil {
		log.Printf("[ERROR] %v",err)
		return
	}
	log.Printf("[INFO] Credentials saved")
}

func readFromJson(credsPath string, val any) error {
	bytes, err := os.ReadFile(credsPath)
	if err != nil {
		if os.IsNotExist(err) {
			if internErr := os.MkdirAll(path.Dir(credsPath), 0700); internErr != nil {
				return internErr
			}
		}
		return err
	}

	if err = json.Unmarshal(bytes, val); err != nil {
		return err
	}

	return nil
}

func LoadConnections() error {
	credsPath := os.Getenv("CREDS_SUBDIR")
	log.Printf("[INFO] Loading credentials from: %s", credsPath)

	var anyCon any = &Cons
	if err := readFromJson(path.Join(credsPath, "credentials.json"), anyCon); err != nil {
		Cons = make(Connections)
		return err;
	}

	var anySMP any = &SlaveMaster
	if err := readFromJson(path.Join(credsPath, "slave-master.json"), anySMP); err != nil {
		SlaveMaster = make(SlaveMasterMap)
		return err;
	}

	log.Printf("[INFO] Data loaded: %d", len(Cons))
	return nil
}
