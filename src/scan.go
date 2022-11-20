package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

//SecurityScan contains the secutity scan results
type SecurityScan struct {
	Id          int         `json:"id"`
	Status      string      `json:"status"`
	UserId      string      `json:"user_id"`
	RepoURL     string      `json:"repo_url"`
	RepoName    string      `json:"repo_name"`
	Findings    []*Findings `json:"findings"`
	rawFindings []byte
	QueuedAt    *time.Time `json:"queued_at"`
	ScaningAt   *time.Time `json:"scanning_at"`
	FinishedAt  *time.Time `json:"finished_at"`
}

//listScan lists the security scan results
func listScan(conn *pgxpool.Pool, scan *SecurityScan) ([]*SecurityScan, error) {
	const stmt = `
		select id, status, user_id, repo_url, repo_name, findings, queued_at, scanning_at, finished_at
		from security_scan
		where user_id = $1;
	`
	rows, err := conn.Query(context.Background(), stmt, scan.UserId)
	if err != nil {
		return nil, err
	}
	results := make([]*SecurityScan, 0)
	for rows.Next() {
		scan := SecurityScan{}
		err = rows.Scan(&scan.Id, &scan.Status, &scan.UserId, &scan.RepoURL, &scan.RepoName, &scan.rawFindings, &scan.QueuedAt, &scan.ScaningAt, &scan.FinishedAt)
		json.Unmarshal(scan.rawFindings, &scan.Findings)
		if err != nil {
			return nil, err
		}
		results = append(results, &scan)
	}
	return results, nil
}

//postScan triggers to scan a repository
func postScan(conn *pgxpool.Pool, repo *Repository) (*SecurityScan, error) {
	const stmt = `
		insert into security_scan(status, user_id, repo_url, repo_name, queued_at, scanning_at)
		select 'Queued', r.user_id, r.repo_url, r.repo_name, now(), now()
		from repository r
		where r.id = $1 and r.user_id = $2
		RETURNING id, repo_url
	`
	row := conn.QueryRow(context.Background(), stmt, repo.Id, repo.UserId)
	var scan SecurityScan
	err := row.Scan(&scan.Id, &scan.RepoURL)
	return &scan, err
}

func putScan(conn *pgxpool.Pool, scan *SecurityScan) error {
	const stmt = `
		update security_scan
		set status = $2,
			scanning_at = $3,
			finished_at = $4,
			findings = $5
		where id = $1;
	`
	_, err := conn.Exec(context.Background(), stmt, scan.Id, scan.Status, scan.ScaningAt, scan.FinishedAt, scan.rawFindings)
	return err
}

//ScanHandler is a Handler for Scan Service
func ScanHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		arg_repo_id := r.URL.Query().Get("repo_id")
		repo_id, _ := strconv.Atoi(arg_repo_id)
		user_id := r.URL.Query().Get("user_id")
		scan, err := postScan(Conn, &Repository{Id: repo_id, UserId: user_id})
		if err == pgx.ErrNoRows {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		go func() {
			//Start in progress
			scan.Status = "In Progress"
			now := time.Now()
			scan.ScaningAt = &now
			putScan(Conn, scan)
			findingResults, err := ScanRepository(scan.RepoURL)
			now = time.Now()
			scan.FinishedAt = &now
			//If no error during scanning
			if err == nil {
				binResult, _ := json.Marshal(findingResults)
				scan.rawFindings = binResult
				scan.Status = "Success"
			} else { // if error
				scan.Status = "Failure"
			}
			putScan(Conn, scan)
		}()
		w.WriteHeader(http.StatusCreated)
	case "GET":
		user_id := r.URL.Query().Get("user_id")
		results, err := listScan(Conn, &SecurityScan{UserId: user_id})
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		response, err := json.Marshal(results)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, string(response))
	}
}
