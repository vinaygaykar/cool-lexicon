package lexicon

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/vinaygaykar/cool-lexicon/configs"
	"github.com/vinaygaykar/cool-lexicon/lexicon/internal"
)

// GetInstance returns an instance of Lexicon object configured using properties as described in configFileLoc.
// If configs are nil or invalid then this function will panic. If internal systme connection fails function will panic.
// If shouldPerformSetupCheck is true then system checks are performed to make sure everything is setup as expected.
// If system is setup is incorrectly then it will "try" to correct the setup or end up panicking. This field
// is useful during troubleshooting and should only be set once during first run, on later runs if this value
// is set it won't cause any harm but might slow down the operations.
func GetInstance(shouldPerformSetupCheck bool, cfg *configs.Configs) *lexicon.LexiconWithDB {
	if cfg == nil {
		log.Panic("config is nil")
	}

	if shouldPerformSetupCheck {
		performSetupChecks(cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
	}

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database))
	if err != nil {
		log.Panic(err.Error())
	}

	log.Printf("Connected to MySQL @ %s:%d/%s\n", cfg.Host, cfg.Port, cfg.Database)
	return lexicon.Open(db)
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
	if _, err = db.Exec(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s.lexicon(word VARCHAR(100) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL, PRIMARY KEY (word))", database)); err != nil {
		log.Panic(err.Error())
	}

	log.Println("DB checks comlpeted")
	log.Println("All checks completed")
}
