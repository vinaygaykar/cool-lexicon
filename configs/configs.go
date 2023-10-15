// Package configs provides factory methods to configure and get instance of a Lexicon object.
package configs

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	lexicondb "github.com/vinaygaykar/cool-lexicon/internal/lexicon"
	"github.com/vinaygaykar/cool-lexicon/pkg/lexicon"
)

// A Config holds config values required for this project.
// values are populated from `cool-lexicon-cfg.json` config file as default, unless another file
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

// GetLexicon returns an instance of Lexicon object configured using properties as described in configFileLoc.
// For now instance of LexiconWithDB is provided as an instance of Lexicon, which uses MySQL as data storage solution.
// If configFileLoc is empty or invalid then GetLexicon will panic.
// If checkSetup is true then system checks are performed to make sure everything is setup as expected. 
// If system is setup is incorrectly then it will "try" to correct the setup or end up panicking. This field
// is useful during troubleshooting and should only be set once during first run, on later runs if this value 
// is set it won't cause any harm but might slow down the operations.
func GetLexicon(configFileLoc string, checkSetup bool) lexicon.Lexicon {
	if len(configFileLoc) == 0 {
		log.Panicln("config file not provided")
	}

	cfg := getConfigs(configFileLoc)

	if checkSetup {
		performSetupChecks(cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
	}

	return lexicondb.Open(connectToMySQL(cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database))
}

func getConfigs(configFileLoc string) *Configs {
	cfg := Configs{}

	// read config file into `cfg` object
	file, _ := os.Open(configFileLoc)
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

func performSetupChecks(username, password, host string, port int, database string) {
	log.Printf("Connecting to MySQL @ %s:%d", host, port)
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/", username, password, host, port))
	if err != nil {
		log.Panic(err.Error())
	}
	defer db.Close()

	log.Printf("Creating database %s if it does not already exists\n", database)
	if _, err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", database)); err != nil {
		log.Panic(err.Error())
	}
	db.Close()

	log.Printf("Connected to %s:%d/%s\n", host, port, database)
	db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", username, password, host, port, database))
	if err != nil {
		log.Panic(err.Error())
	}
	defer db.Close()

	log.Println("Creating table lexicon if it does not already exists")
	if _, err = db.Exec(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s.lexicon(word VARCHAR(100) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL)", database)); err != nil {
		log.Panic(err.Error())
	}

	log.Println("DB checks comlpeted")
	log.Println("All checks completed")
}

func connectToMySQL(username, password, host string, port int, database string) *sql.DB {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", username, password, host, port, database))
	if err != nil {
		log.Panic(err.Error())
	}

	log.Printf("Connected to MySQL @ %s:%d/%s\n", host, port, database)
	return db
}
