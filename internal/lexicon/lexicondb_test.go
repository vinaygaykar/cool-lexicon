// Package lexicon provides implementation for the `pkg/lexicon.Lexicon` interface.
// It uses MySQL to store state or words of the Lexicon.
package lexicon

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
)

const (
	mysqlImage = "mysql:8"
	dbName = "lexicons"
	dbUserName = "root"
	dbPassword = "toor"
	tableName = "lexicon"
)

func getDB(ctx *context.Context, initialWords []string) (*sql.DB, func()) {
	container, err := mysql.RunContainer(*ctx,
		testcontainers.WithImage(mysqlImage),
		mysql.WithDatabase(dbName),
		mysql.WithUsername(dbUserName),
		mysql.WithPassword(dbPassword),
	)

	if err != nil {
		panic(err)
	}

	host, _ := container.Host(*ctx)
	port, _ := container.MappedPort(*ctx, "3306/tcp")
	p := fmt.Sprint(port.Int())
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUserName, dbPassword, host, p, dbName))

	if err != nil {
		container.Terminate(*ctx)
		panic(err)
	}

	// Add initial words to DB
	db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", dbName))
	db.Exec(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s(word VARCHAR(100) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL)", tableName))
	query := fmt.Sprintf("INSERT INTO %s VALUES (?)", tableName)
	
	for _, word := range initialWords {
		if _, err := db.Exec(query, word); err != nil {
			db.Close()
			container.Terminate(*ctx)
			panic(err)
		}
	}

	// Clean up the container
	closeFunc := func() {
		if err := db.Close(); err != nil {
			panic(err)
		}

		if err := container.Terminate(*ctx); err != nil {
			panic(err)
		}
	}

	return db, closeFunc
}

func TestLexiconWithDB_Lookup(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 2 * time.Minute)
	db, closeFunc := getDB(&ctx, []string{"exists"})

	defer func() {
		defer cancel()
		defer closeFunc()
	}()

	type fields struct {
		db *sql.DB
	}
	type args struct {
		words []string
	}
	
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []bool
		wantErr bool
	}{
		{
			name:    "Given a word exists in db when lookup is called for the word then return true should be returned",
			fields:  fields{db: db},
			args:    args{words: []string{"exists"}},
			want:    []bool{true},
			wantErr: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lxc := &LexiconWithDB{
				db: tt.fields.db,
			}
			
			got, err := lxc.Lookup(tt.args.words...)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("LexiconWithDB.Lookup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LexiconWithDB.Lookup() = %v, want %v", got, tt.want)
			}
		})
	}

}
