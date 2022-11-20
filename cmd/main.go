package main

import (
	"context"
	"fmt"
	"guardrails/src"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	DATABASE_URL := os.Getenv("DATABASE_URL")

	conn, err := pgxpool.New(context.Background(), DATABASE_URL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	//Pass connection directly src package
	src.Conn = conn

	http.HandleFunc("/repository", src.RepositoryHandler)
	http.HandleFunc("/scan", src.ScanHandler)

	http.ListenAndServe(fmt.Sprintf(":3000"), nil)
}
