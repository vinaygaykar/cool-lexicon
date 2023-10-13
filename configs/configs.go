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

func GetLexicon(configFileLoc string) (lexicon.Lexicon, error) {
	if len(configFileLoc) == 0 {
		configFileLoc = "cool-lexicon-cfg.json"
	}

	// File exists; read it into `cfg` object
	file, _ := os.Open(configFileLoc)
	defer file.Close()

	decoder := json.NewDecoder(file)
	cfg := Configs{}
	if err := decoder.Decode(&cfg); err != nil {
		return nil, err
	}

	return lexicondb.Open(connectToMySQL(cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database)), nil
}

func connectToMySQL(username, password, host string, port int, database string) *sql.DB {
	url := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", username, password, host, port, database)

	db, err := sql.Open("mysql", url)
	if err != nil {
		panic(err.Error())
	}

	log.Println("DB connected")

	return db
}
