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

type Configs struct {
	Dbtype   string `json:"type"`
	Host     string `json:"host"`
	Username string `json:"username"`
	Password string `json:"password"`
	Port     int    `json:"port"`
	Database string `json:"database"`
}

func GetLexicon(configFileLoc string, checkSetup bool) lexicon.Lexicon {
	if len(configFileLoc) == 0 {
		configFileLoc = "cool-lexicon-cfg.json"
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
