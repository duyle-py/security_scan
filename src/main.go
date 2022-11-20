package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	conn, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	//Pass connection directly
	Conn = conn
	http.HandleFunc("/repository", RepositoryHandler)
	http.HandleFunc("/scan", ScanHandler)
	http.ListenAndServe(fmt.Sprintf(":3000"), nil)
}
