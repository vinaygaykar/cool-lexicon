// Package lexicon provides implementation for the `pkg/lexicon.Lexicon` interface.
// It uses MySQL to store state or words of the Lexicon.
package lexicon

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var (
	errNilOrEmptyWords = errors.New("list of words is nil or empty")
)

const TABLE_NAME = "lexicon"

func Open(db *sql.DB) *LexiconWithDB {
	if db == nil {
		log.Panicln("database value is nil")
	}

	return &LexiconWithDB{db: db}
}

type LexiconWithDB struct {
	db *sql.DB
}

func (lxc *LexiconWithDB) Lookup(words ...string) (*[]string, error) {
	if len(words) == 0 {
		return nil, errNilOrEmptyWords
	}

	query := fmt.Sprintf("SELECT EXISTS (SELECT l.word FROM %s l WHERE l.word LIKE ?)", TABLE_NAME)
	exists := make([]string, 0)

	for _, word := range words {
		exist := false
		row := lxc.db.QueryRow(query, word)
		if err := row.Scan(&exist); err == nil && exist {
			exists = append(exists, word)
		}
	}

	return &exists, nil
}

func (lxc *LexiconWithDB) GetAllWordsStartingWith(substrings ...string) (*map[string][]string, error) {
	if len(substrings) == 0 {
		return nil, errNilOrEmptyWords
	}

	result := make(map[string][]string, 0)

	for _, substring := range substrings {
		words, err := lxc.searchSubString(substring + "%")
		if err == nil && len(words) != 0 {
			result[substring] = words
		}
	}

	return &result, nil
}

func (lxc *LexiconWithDB) GetAllWordsEndingWith(substrings ...string) (*map[string][]string, error) {
	if len(substrings) == 0 {
		return nil, errNilOrEmptyWords
	}

	result := make(map[string][]string, 0)

	for _, substring := range substrings {
		words, err := lxc.searchSubString("%" + substring)
		if err == nil && len(words) != 0 {
			result[substring] = words
		}
	}

	return &result, nil
}

func (lxc *LexiconWithDB) searchSubString(toSearch string) ([]string, error) {
	words := make([]string, 0)
	query := fmt.Sprintf("SELECT l.word FROM %s l WHERE l.word LIKE ?", TABLE_NAME)

	res, err := lxc.db.Query(query, toSearch)
	if err != nil {
		return []string{}, err
	}
	defer res.Close()

	for res.Next() {
		var word string
		if err = res.Scan(&word); err != nil {
			return []string{}, err
		}

		words = append(words, word)
	}

	return words, nil
}

func (lxc *LexiconWithDB) Add(words ...string) error {
	if len(words) == 0 {
		return errNilOrEmptyWords
	}

	query := fmt.Sprintf("INSERT INTO %s VALUES ", TABLE_NAME)
	vals := []interface{}{}
	for _, w := range words {
		query += "(?), "
		vals = append(vals, w)
	}
	// trim the last comma (,)
	query = query[0 : len(query)-2]

	if stmt, err := lxc.db.Prepare(query); err == nil {
		defer stmt.Close()
		if _, err = stmt.Exec(vals...); err != nil {
			return err
		}
	} else {
		return err
	}

	return nil
}

func (lxc *LexiconWithDB) Close() {
	defer lxc.db.Close()
}
