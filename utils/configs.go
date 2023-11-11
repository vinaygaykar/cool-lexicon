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

	// Dbtype mentions type of DB server used. Optional as currently only MySQL is supported.
	Dbtype string `json:"type"`

	// Host address of the database
	Host string `json:"host"`

	// Username credentials to use for database login
	Username string `json:"username"`

	// Password credentials to use for database login
	Password string `json:"password"`

	// Port value of the database connection
	Port int `json:"port"`

	// Database to connect to
	Database string `json:"database"`
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
	if len(cfg.Username) == 0 {
		log.Panic("mysql username is invalid")
	}

	if len(cfg.Password) == 0 {
		log.Panic("mysql password is invalid")
	}

	if len(cfg.Host) == 0 {
		log.Panic("mysql host is invalid")
	}

	if len(cfg.Database) == 0 {
		log.Panic("mysql database is invalid")
	}

	return &cfg
}
