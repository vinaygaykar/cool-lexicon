package lexicon

import (
	"fmt"
	"log"

	"github.com/vinaygaykar/cool-lexicon/lexicon/internal/mysql"
	"github.com/vinaygaykar/cool-lexicon/utils"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/libsql/libsql-client-go/libsql"
)

// GetInstance returns an instance of Lexicon object configured using properties as described in configFileLoc.
// If configs are nil or invalid then this function will panic. If internal systme connection fails function will panic.
// If shouldPerformSetupCheck is true then system checks are performed to make sure everything is setup as expected.
// If system is setup is incorrectly then it will "try" to correct the setup or end up panicking. This field
// is useful during troubleshooting and should only be set once during first run, on later runs if this value
// is set it won't cause any harm but might slow down the operations.
func GetInstance(shouldPerformSetupCheck bool, cfg *configs.Configs) *lexicon.LexiconMySQL {
	if cfg == nil {
		log.Panic("config is nil")
	}

	var driver string
	var url string
	if cfg.Dbtype == "libsql" {
		driver = cfg.Dbtype
		url = fmt.Sprintf("http://%s:%d", cfg.Host, cfg.Port)
		if shouldPerformSetupCheck {
			performSetupChecksLibSQL(url)
		}
	} else if cfg.Dbtype == "mysql" {
		driver = cfg.Dbtype
		url = fmt.Sprintf("%s:%s@tcp(%s:%d)", cfg.Username, cfg.Password, cfg.Host, cfg.Port)
		if shouldPerformSetupCheck {
			performSetupChecksMySQL(url, cfg.Database)
		}
		url += "/" + cfg.Database
	} else {
		log.Panicln("invalid db type provided in the configs")
	}

	db, err := sql.Open(driver, url)
	if err != nil {
		log.Panicln(err.Error())
	}

	log.Printf("Connected to %s @ %s:%d\n", cfg.Dbtype, cfg.Host, cfg.Port)
	return lexicon.Open(db, driver)
}

func performSetupChecksMySQL(url, database string) {
	log.Printf("Connecting to mysql")
	db, err := sql.Open("mysql", url)
	if err != nil {
		log.Panicln(err.Error())
	}
	defer db.Close()

	log.Printf("Creating database %s if it does not already exists\n", database)
	if _, err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", database)); err != nil {
		log.Panicln(err.Error())
	}
	db.Close()

	log.Printf("Connected to mysql")
	db, err = sql.Open("mysql", url+"/"+database)
	if err != nil {
		log.Panicln(err.Error())
	}
	defer db.Close()

	log.Println("Creating table lexicon if it does not already exists")
	if _, err = db.Exec("CREATE TABLE IF NOT EXISTS lexicon(word VARCHAR(100) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL, PRIMARY KEY (word))"); err != nil {
		log.Panicln(err.Error())
	}

	log.Println("All checks completed")
}

func performSetupChecksLibSQL(url string) {
	log.Printf("Connecting to libsql")
	db, err := sql.Open("libsql", url)
	if err != nil {
		log.Panicln(err.Error())
	}
	defer db.Close()

	log.Println("Creating table lexicon if it does not already exists")
	if _, err = db.Exec("CREATE TABLE IF NOT EXISTS lexicon(word VARCHAR(100) COLLATE NOCASE, PRIMARY KEY (word))"); err != nil {
		log.Panicln(err.Error())
	}

	log.Println("All checks completed")
}
