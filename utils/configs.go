// Package configs provides factory methods to configure and get instance of a Lexicon object.
package configs

import (
	"encoding/json"
	"log"
	"os"
)

// A Config holds config values required for this project.
// Values are populated from `config.json` config file as default, unless another file
// is explicitly provided.
// A zero value config is useless and will be reported invalid during validation phase.
type Configs struct {

	// Dbtype mentions type of DB server used.
	Dbtype string `json:"type"`

	// Host address of the database server.
	Host string `json:"host"`

	// Port value of the database server.
	Port int `json:"port"`

	// Database name.
	Database string `json:"database"`

	// Username credentials to use for database login. Not needed if authToken is configured.
	Username string `json:"username"`

	// Password credentials to use for database login. Not needed if authToken is configured.
	Password string `json:"password"`

	// Authentication token. Not needed if username/password is configured.
	AuthToken string `json:"authToken"`
}

func ReadConfigs(filePath string) *Configs {
	cfg := Configs{}

	// read config file into `cfg` object
	file, _ := os.Open(filePath)
	defer file.Close()
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		log.Panic(err.Error())
	}

	// validate
	if len(cfg.Host) == 0 {
		log.Panic("host is invalid")
	}

	return &cfg
}
