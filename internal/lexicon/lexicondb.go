package lexicon

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

const TABLE_NAME = "lexicon"

func Open(db *sql.DB) *LexiconWithDB {
	return &LexiconWithDB{
		db: db,
	}
}

type LexiconWithDB struct {
	db *sql.DB
}

func (lxc *LexiconWithDB) CheckIfExists(word string) bool {
	exists := false
	query := fmt.Sprintf("SELECT EXISTS (SELECT l.word FROM %s l WHERE l.word LIKE ?)", TABLE_NAME)

	row := lxc.db.QueryRow(query, word)
	err := row.Scan(&exists)
	if err != nil {
		panic(err.Error())
	}

	return exists
}

func (lxc *LexiconWithDB) GetAllStartingWith(toSearch string) []string {
	return lxc.searchSubString(toSearch + "%")
}

func (lxc *LexiconWithDB) GetAllEndingWith(toSearch string) []string {
	return lxc.searchSubString("%" + toSearch)
}

func (lxc *LexiconWithDB) searchSubString(toSearch string) []string {
	words := make([]string, 0)
	query := fmt.Sprintf("SELECT l.word FROM %s l WHERE l.word LIKE ?", TABLE_NAME)

	log.Println(query)

	res, err := lxc.db.Query(query, toSearch)
	if err != nil {
		panic(err.Error())
	}

	defer res.Close()

	for res.Next() {
		var word string
		err2 := res.Scan(&word)
		if err2 != nil {
			panic(err2.Error())
		}

		words = append(words, word)
	}

	return words
}

func (lxc *LexiconWithDB) AddAll(words []string) {
	if len(words) == 0 {
		return
	}

	query := fmt.Sprintf("INSERT INTO %s VALUES ", TABLE_NAME)
	vals := []interface{}{}
	for _, w := range words {
		query += "(?), "
		vals = append(vals, w)
	}
	// trim the last comma (,)
	query = query[0 : len(query)-2]

	stmt, err := lxc.db.Prepare(query)
	if err != nil {
		panic(err.Error())
	}

	_, err2 := stmt.Exec(vals...)
	if err2 != nil {
		panic(err2.Error())
	}

	defer stmt.Close()
}

func (lxc *LexiconWithDB) Close() {
	defer lxc.db.Close()
}
