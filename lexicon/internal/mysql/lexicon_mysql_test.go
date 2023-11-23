package lexicon

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
	"github.com/testcontainers/testcontainers-go/wait"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/libsql/libsql-client-go/libsql"
)

const (
	mysqlImage    = "mysql:8"
	dbName        = "lexicons"
	dbUserName    = "root"
	dbPassword    = "toor"
	testTableName = "lexicon"
)

var randomWordsInsertedInDBOnInit = [...]string{"नमस्ते", "धन्यवाद", "नमस्कार", "सुंदर", "मोक्ष"}

func getMySQLDB(ctx context.Context) (*sql.DB, func()) {
	container, err := mysql.RunContainer(ctx,
		testcontainers.WithImage(mysqlImage),
		mysql.WithDatabase(dbName),
		mysql.WithUsername(dbUserName),
		mysql.WithPassword(dbPassword),
	)

	if err != nil {
		panic(err)
	}

	host, _ := container.Host(ctx)
	port, _ := container.MappedPort(ctx, "3306/tcp")
	p := fmt.Sprint(port.Int())
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUserName, dbPassword, host, p, dbName))

	if err != nil {
		container.Terminate(ctx)
		panic(err)
	}

	// Add initial words to DB
	db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", dbName))
	db.Exec(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s(word VARCHAR(100) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL, PRIMARY KEY (word))", testTableName))
	query := fmt.Sprintf("INSERT INTO %s VALUES (?)", testTableName)

	for _, word := range randomWordsInsertedInDBOnInit {
		if _, err := db.Exec(query, word); err != nil {
			db.Close()
			container.Terminate(ctx)
			panic(err)
		}
	}

	// Clean up the container
	closeFunc := func() {
		if err := db.Close(); err != nil {
			panic(err)
		}

		if err := container.Terminate(ctx); err != nil {
			panic(err)
		}
	}

	return db, closeFunc
}

func getLibSQLDB(ctx context.Context) (*sql.DB, func()) {
	req := testcontainers.ContainerRequest{
		Image:        "ghcr.io/libsql/sqld",
		ExposedPorts: []string{"8080"},
		WaitingFor:   wait.ForLog("Welcome to sqld!"),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		panic(err.Error())
	}

	ip, err := container.Host(ctx)
	if err != nil {
		panic(err.Error())
	}

	mappedPort, err := container.MappedPort(ctx, "8080")
	if err != nil {
		panic(err.Error())
	}

	db, err := sql.Open("libsql", "http://"+ip+":"+mappedPort.Port())
	if err != nil {
		panic(err.Error())
	}

	// Add initial words to DB
	db.Exec(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s(word VARCHAR(100) collate NOCASE, PRIMARY KEY (word))", testTableName))
	query := fmt.Sprintf("INSERT INTO %s VALUES (?)", testTableName)

	for _, word := range randomWordsInsertedInDBOnInit {
		if _, err := db.Exec(query, word); err != nil {
			db.Close()
			container.Terminate(ctx)
			panic(err)
		}
	}

	// Clean up the container
	closeFn := func() {
		if err := db.Close(); err != nil {
			panic(err.Error())
		}

		if err := container.Terminate(ctx); err != nil {
			panic(err.Error())
		}
	}

	return db, closeFn
}

func TestLexiconWithDB_Lookup(t *testing.T) {
	ctx := context.Background()
	ctx, cancelCtx := context.WithTimeout(ctx, 2*time.Minute)
	mysqlDB, closeDB := getMySQLDB(ctx)
	libsqlDB, closeDB2 := getLibSQLDB(ctx)

	defer cancelCtx()
	defer closeDB()
	defer closeDB2()

	type fields struct {
		mySQL  *sql.DB
		libSQL *sql.DB
	}
	type args struct {
		words []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *[]string
		wantErr bool
	}{
		{
			name:    "Given a Lexicon with some words, when Lookup is invoked for the existing word, then return true should be returned",
			fields:  fields{mysqlDB, libsqlDB},
			args:    args{words: []string{"नमस्ते"}},
			want:    &([]string{"नमस्ते"}),
			wantErr: false,
		},
		{
			name:    "Given a Lexicon with some words, when Lookup is invoked for that non existing word, then return value should be false",
			fields:  fields{mysqlDB, libsqlDB},
			args:    args{words: []string{"notexists"}},
			want:    &([]string{}),
			wantErr: false,
		},
		{
			name:    "Given a Lexicon with some words, when Lookup is invoked multiple words some exists and others don't, then return value should be true only for existing words",
			fields:  fields{mysqlDB, libsqlDB},
			args:    args{words: []string{"नमस्कार", "notexists", "सुंदर", "धन्यवाद", "पराक्रम", "पानी"}},
			want:    &([]string{"नमस्कार", "सुंदर", "धन्यवाद"}),
			wantErr: false,
		},
		{
			name:    "Given a Lexicon with some words, when Lookup is invoked for nil words array, then error is expected",
			fields:  fields{mysqlDB, libsqlDB},
			args:    args{words: nil},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Given a Lexicon with some words, when Lookup is invoked for empty words array, then error is expected",
			fields:  fields{mysqlDB, libsqlDB},
			args:    args{words: []string{}},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			test := func(db *sql.DB, dbName string) {
				lxc := &LexiconMySQL{
					db: db,
				}
				got, err := lxc.Lookup(tt.args.words...)
				if (err != nil) != tt.wantErr {
					t.Errorf("[%s] LexiconWithDB.Lookup() error = %v, wantErr %v", dbName, err, tt.wantErr)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("[%s] LexiconWithDB.Lookup() = %v, want %v", dbName, got, tt.want)
				}
			}

			test(tt.fields.mySQL, "mysql")
			test(tt.fields.libSQL, "libsql")
		})
	}
}

func TestLexiconWithDB_GetAllWordsStartingWith(t *testing.T) {
	ctx := context.Background()
	ctx, cancelCtx := context.WithTimeout(ctx, 2*time.Minute)
	mysqlDB, closeDB := getMySQLDB(ctx)
	libsqlDB, closeDB2 := getLibSQLDB(ctx)

	defer cancelCtx()
	defer closeDB()
	defer closeDB2()

	type fields struct {
		mySQL  *sql.DB
		libSQL *sql.DB
	}

	type args struct {
		substrings []string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *map[string][]string
		wantErr bool
	}{
		{
			name:   "Given a Lexicon with some words, when SearchStartsWith is invoked for existing word, then return all the words starting with the substring sorted lexicographically",
			fields: fields{mysqlDB, libsqlDB},
			args:   args{substrings: []string{"न"}},
			want: &(map[string][]string{
				"न": {"नमस्कार", "नमस्ते"},
			}),
			wantErr: false,
		},
		{
			name:    "Given a Lexicon with some words, when SearchStartsWith is invoked for non-existing word, then return no response for the substring",
			fields:  fields{mysqlDB, libsqlDB},
			args:    args{substrings: []string{"क्र"}},
			want:    &map[string][]string{},
			wantErr: false,
		},
		{
			name:   "Given a Lexicon with some words, when SearchStartsWith is invoked for mix of existing & non existing word, then return all the words starting with the existing words mapped to correct key sorted lexicographically while non existing words have no entry",
			fields: fields{mysqlDB, libsqlDB},
			args:   args{substrings: []string{"नम", "somethingelse", "नमस्", "धन्य"}},
			want: &(map[string][]string{
				"नम":   {"नमस्कार", "नमस्ते"},
				"नमस्": {"नमस्कार", "नमस्ते"},
				"धन्य": {"धन्यवाद"},
			}),
			wantErr: false,
		},
		{
			name:    "Given a Lexicon with some words, when SearchStartsWith is invoked for nil words array, then error is expected",
			fields:  fields{mysqlDB, libsqlDB},
			args:    args{substrings: nil},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Given a Lexicon with some words, when SearchStartsWith is invoked for empty words array, then error is expected",
			fields:  fields{mysqlDB, libsqlDB},
			args:    args{substrings: []string{}},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			test := func(db *sql.DB, dbName string) {
				lxc := &LexiconMySQL{
					db: db,
				}
				got, err := lxc.GetAllWordsStartingWith(tt.args.substrings...)
				if (err != nil) != tt.wantErr {
					t.Errorf("LexiconWithDB.GetAllWordsStartingWith() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("LexiconWithDB.GetAllWordsStartingWith() = %v, want %v", got, tt.want)
				}
			}

			test(tt.fields.mySQL, "mysql")
			test(tt.fields.libSQL, "libsql")
		})
	}
}

func TestLexiconWithDB_GetAllWordsEndingWith(t *testing.T) {
	ctx := context.Background()
	ctx, cancelCtx := context.WithTimeout(ctx, 2*time.Minute)
	mysqlDB, closeDB := getMySQLDB(ctx)
	libsqlDB, closeDB2 := getLibSQLDB(ctx)

	defer cancelCtx()
	defer closeDB()
	defer closeDB2()

	type fields struct {
		mySQL  *sql.DB
		libSQL *sql.DB
	}
	type args struct {
		substrings []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *map[string][]string
		wantErr bool
	}{
		{
			name:   "Given a Lexicon with some words, when SearchEndsWith is invoked for existing word, then return all the words ending with the substring sorted lexicographically",
			fields: fields{mysqlDB, libsqlDB},
			args:   args{substrings: []string{"र"}},
			want: &(map[string][]string{
				"र": {"नमस्कार", "सुंदर"},
			}),
			wantErr: false,
		},
		{
			name:    "Given a Lexicon with some words, when SearchEndsWith is invoked for non-existing word, then return no response for the substring",
			fields:  fields{mysqlDB, libsqlDB},
			args:    args{substrings: []string{"क्र"}},
			want:    &map[string][]string{},
			wantErr: false,
		},
		{
			name:   "Given a Lexicon with some words, when SearchEndsWith is invoked for mix of existing & non existing word, then return all the words ending with the existing words mapped to correct key sorted lexicographically while non existing words have no entry",
			fields: fields{mysqlDB, libsqlDB},
			args:   args{substrings: []string{"र", "somethingelse", "क्ष", "वाद"}},
			want: &(map[string][]string{
				"र":   {"नमस्कार", "सुंदर"},
				"क्ष": {"मोक्ष"},
				"वाद": {"धन्यवाद"},
			}),
			wantErr: false,
		},
		{
			name:    "Given a Lexicon with some words, when SearchEndsWith is invoked for nil words array, then error is expected",
			fields:  fields{mysqlDB, libsqlDB},
			args:    args{substrings: nil},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Given a Lexicon with some words, when SearchEndsWith is invoked for empty words array, then error is expected",
			fields:  fields{mysqlDB, libsqlDB},
			args:    args{substrings: []string{}},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			test := func(db *sql.DB, dbName string) {
				lxc := &LexiconMySQL{
					db: db,
				}
				got, err := lxc.GetAllWordsEndingWith(tt.args.substrings...)
				if (err != nil) != tt.wantErr {
					t.Errorf("LexiconWithDB.GetAllWordsEndingWith() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("LexiconWithDB.GetAllWordsEndingWith() = %v, want %v", got, tt.want)
				}
			}

			test(tt.fields.mySQL, "mysql")
			test(tt.fields.libSQL, "libsql")
		})
	}
}

func TestLexiconWithDB_Add(t *testing.T) {
	ctx := context.Background()
	ctx, cancelCtx := context.WithTimeout(ctx, 2*time.Minute)
	mysqlDB, closeDB := getMySQLDB(ctx)
	libsqlDB, closeDB2 := getLibSQLDB(ctx)

	defer cancelCtx()
	defer closeDB()
	defer closeDB2()

	type fields struct {
		mySQL  *sql.DB
		libSQL *sql.DB
	}
	type args struct {
		words []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "Given a Lexicon with some words, when Add is invoked with empty array, then error is expected",
			fields:  fields{mysqlDB, libsqlDB},
			args:    args{words: []string{}},
			wantErr: true,
		},
		{
			name:    "Given a Lexicon with some words, when Add is invoked on a new non existent word, then no error is expected",
			fields:  fields{mysqlDB, libsqlDB},
			args:    args{words: []string{"देव"}},
			wantErr: false,
		},
		{
			name:    "Given a Lexicon with some words, when Lookup is invoked multiple words some exists and others don't, then return value should be true only for existing words",
			fields:  fields{mysqlDB, libsqlDB},
			args:    args{words: []string{"नमस्कार"}},
			wantErr: true,
		},
		{
			name:    "Given a Lexicon with some words, when Add is invoked for nil words array, then error is expected",
			fields:  fields{mysqlDB, libsqlDB},
			args:    args{words: nil},
			wantErr: true,
		},
		{
			name:    "Given a Lexicon with some words, when Add is invoked for empty words array, then error is expected",
			fields:  fields{mysqlDB, libsqlDB},
			args:    args{words: []string{}},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			test := func(db *sql.DB, dbName string) {
				lxc := &LexiconMySQL{
					db: db,
				}
				if err := lxc.Add(tt.args.words...); (err != nil) != tt.wantErr {
					t.Errorf("LexiconWithDB.Add() error = %v, wantErr %v", err, tt.wantErr)
				}
			}

			test(tt.fields.mySQL, "mysql")
			test(tt.fields.libSQL, "libsql")
		})
	}
}
