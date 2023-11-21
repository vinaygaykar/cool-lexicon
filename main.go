package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/libsql/libsql-client-go/libsql"
)

var dbUrl = "http://127.0.0.1:8080"

func main() {
	if _, err := sql.Open("libsql", dbUrl); err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db %s: %s", dbUrl, err)
		os.Exit(1)
	}
}