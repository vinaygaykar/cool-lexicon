// Package lexicon provides implementation for the `pkg/lexicon.Lexicon` interface.
// It uses MySQL to store state or words of the Lexicon.
package lexicon

import (
	"database/sql"
	"errors"
	"fmt"

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

func (lxc *LexiconWithDB) CheckIfExists(word string) (bool, error) {
	exists := false
	query := fmt.Sprintf("SELECT EXISTS (SELECT l.word FROM %s l WHERE l.word LIKE ?)", TABLE_NAME)

	row := lxc.db.QueryRow(query, word)
	err := row.Scan(&exists)
	return exists, err
}

func (lxc *LexiconWithDB) GetAllStartingWith(toSearch string) ([]string, error) {
	return lxc.searchSubString(toSearch + "%")
}

func (lxc *LexiconWithDB) GetAllEndingWith(toSearch string) ([]string, error) {
	return lxc.searchSubString("%" + toSearch)
}

func (lxc *LexiconWithDB) searchSubString(toSearch string) ([]string, error) {
	words := make([]string, 0)
	query := fmt.Sprintf("SELECT l.word FROM %s l WHERE l.word LIKE ?", TABLE_NAME)

	res, err := lxc.db.Query(query, toSearch)
	if err != nil {
		return []string {}, err
	}
	defer res.Close()

	for res.Next() {
		var word string
		if err = res.Scan(&word); err != nil {
			return []string {}, err
		}

		words = append(words, word)
	}

	return words, nil
}

func (lxc *LexiconWithDB) AddAll(words []string) error {
	if len(words) == 0 {
		return errors.New("list of words to add is empty")
	}

	query := fmt.Sprintf("INSERT INTO %s VALUES ", TABLE_NAME)
	vals := []interface{}{}
	for _, w := range words {
		query += "(?), "
		vals = append(vals, w)
	}
	// trim the last comma (,)
	query = query[0 : len(query)-2]

	if stmt, err := lxc.db.Prepare(query); err != nil {
		return err
	} else {
		defer stmt.Close()
		if _, err = stmt.Exec(vals...); err != nil {
			return err
		}
	}

	return nil
}

func (lxc *LexiconWithDB) Close() {
	defer lxc.db.Close()
}
