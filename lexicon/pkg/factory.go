package lexicon

import (
	"fmt"
	"log"

	"github.com/vinaygaykar/cool-lexicon/lexicon/internal/sql"
	"github.com/vinaygaykar/cool-lexicon/utils"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/libsql/libsql-client-go/libsql"
	_ "github.com/mattn/go-sqlite3"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// VerifyDB verifies connection to the provided database and performs required migrations.
// Works for MySQL & libSQL.
// If DB connection or migration fails then the function will panic.
func VerifyDB(cfg *configs.Configs) {
	dbUrl, _ := getDBUrlAndDriver(cfg)
	var m *migrate.Migrate
	var err error
	if cfg.Dbtype == "mysql" {
		m, err = migrate.New("file://db/migrations/mysql", dbUrl)
	} else if cfg.Dbtype == "libsql" || cfg.Dbtype == "turso" {
		if db, err := sql.Open("libsql", dbUrl); err != nil {
			log.Panicf("[migrations] [%s] : could not connect to server : %s\n", cfg.Dbtype, err.Error())
		} else {
			if driver, err := sqlite3.WithInstance(db, &sqlite3.Config{}); err != nil {
				log.Panicf("[migrations] [%s] : could not create driver : %s\n", cfg.Dbtype, err.Error())
			} else {
				if m, err = migrate.NewWithDatabaseInstance("file://db/migrations/libsql", "sqlite3", driver); err != nil {
					log.Panicf("[migrations] [%s] : %s\n", cfg.Dbtype, err.Error())
				}
			}
		}
	} else {
		log.Panicln("invalid db type provided in the configs")
	}

	if err != nil {
		log.Panicf("[migrations] [%s] : %s\n", cfg.Dbtype, err.Error())
	} else {
		m.Up()
	}
}

// GetInstance returns an instance of Lexicon object configured as per the configs.
// If configs are nil or invalid then this function will panic. 
// If internal system connection fails then the function will panic.
func GetInstance(cfg *configs.Configs) *lexicon.LexiconSQL {
	if cfg == nil {
		log.Panic("config is nil")
	}

	dbUrl, driver := getDBUrlAndDriver(cfg)
	db, err := sql.Open(driver, dbUrl)
	if err != nil {
		log.Panicln(err.Error())
	}

	log.Printf("connected to %s @ %s:%d\n", cfg.Dbtype, cfg.Host, cfg.Port)
	return lexicon.Open(db, driver)
}

func getDBUrlAndDriver(cfg *configs.Configs) (dbUrl, driver string) {
	if cfg.Dbtype == "libsql" {
		driver = cfg.Dbtype
		dbUrl = fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	} else if cfg.Dbtype == "turso" {
		driver = "libsql"
		dbUrl = fmt.Sprintf("%s?authToken=%s", cfg.Host, cfg.AuthToken)
	} else if cfg.Dbtype == "mysql" {
		driver = cfg.Dbtype
		dbUrl = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
	} else {
		log.Panicln("invalid db type provided in the configs")
	}

	return dbUrl, driver
}
