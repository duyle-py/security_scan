package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Conn *pgxpool.Pool

//Repository contains repo information
type Repository struct {
	Id       int    `json:"id"`
	UserId   string `json:"user_id"`
	RepoName string `json:"repo_name"`
	RepoURL  string `json:"repo_url"`
}

//repositoryFromRequest parse request to Repository Request
func repositoryFromRequest(req *http.Request) *Repository {
	userId := req.URL.Query().Get("user_id")
	repoName := req.URL.Query().Get("repo_name")
	repoUrl := req.URL.Query().Get("repo_url")
	strRepoId := req.URL.Query().Get("id")
	repoId, _ := strconv.Atoi(strRepoId)
	return &Repository{Id: repoId, UserId: userId, RepoURL: repoUrl, RepoName: repoName}
}

//listRepository lists the repositories
func listRepository(conn *pgxpool.Pool, repo *Repository) ([]*Repository, error) {
	const stmt = `
		select id, user_id, repo_name, repo_url
		from repository
		where user_id = $1;
	`
	rows, err := conn.Query(context.Background(), stmt, repo.UserId)
	if err != nil {
		return nil, err
	}
	repositories := make([]*Repository, 0)
	for rows.Next() {
		repo := Repository{}
		err = rows.Scan(&repo.Id, &repo.UserId, &repo.RepoName, &repo.RepoURL)
		if err != nil {
			return nil, err
		}
		repositories = append(repositories, &repo)
	}
	return repositories, nil
}

// postRepository creates a repository
func postRepository(conn *pgxpool.Pool, repo *Repository) error {
	const stmt = `
		insert into repository(user_id, repo_name, repo_url)
		values($1, $2, $3);
	`
	_, err := conn.Exec(context.Background(), stmt, repo.UserId, repo.RepoName, repo.RepoURL)
	return err
}

//putRepository update the repository information
func putRepository(conn *pgxpool.Pool, repo *Repository) error {
	const stmt = `
		insert into repository(id, user_id, repo_name, repo_url)
		values($1, $2, $3, $4) on conflict ON CONSTRAINT repository_pkey
		DO UPDATE SET 
			repo_name = EXCLUDED.repo_name,
			repo_url = EXCLUDED.repo_url;
	`
	_, err := conn.Exec(context.Background(), stmt, repo.Id, repo.UserId, repo.RepoName, repo.RepoURL)
	return err
}

//deleteRepository deletes a repository
func deleteRepository(conn *pgxpool.Pool, repo *Repository) error {
	const stmt = `
		delete from repository where id = $1 and user_id = $2;
	`
	_, err := conn.Exec(context.Background(), stmt, repo.Id, repo.UserId)
	return err
}

//RepositoryHandler is a Handler for Repository Service
func RepositoryHandler(w http.ResponseWriter, r *http.Request) {
	repo := repositoryFromRequest(r)
	var response []byte
	switch r.Method {
	case "GET":
		if repo.UserId == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		results, err := listRepository(Conn, repo)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		response, err = json.Marshal(results)
		if err != nil {
			fmt.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, string(response))
	case "POST":
		if repo.UserId == "" || repo.RepoName == "" || repo.RepoURL == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err := postRepository(Conn, repo)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	case "DELETE":
		if repo.UserId == "" || repo.Id == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err := deleteRepository(Conn, repo)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	case "PUT":
		if repo.UserId == "" || repo.Id == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err := putRepository(Conn, repo)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
