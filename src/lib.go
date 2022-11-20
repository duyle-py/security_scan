package src

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

func ScanRepository(repository string) ([]*Findings, error) {
	dir, err := ioutil.TempDir("/tmp", "repo")

	if err != nil {
		return nil, err
	}

	conf := exec.Command("git", "clone", repository)
	conf.Dir = dir

	_, err = conf.Output()

	if err != nil {
		return nil, err
	}
	repo_parts := strings.Split(repository, "/")
	repo_name := repo_parts[len(repo_parts)-1]
	res, _ := FindWords(dir + "/" + repo_name)

	defer os.RemoveAll(dir)
	return res, err
}

type Findings struct {
	Type     string   `json:"type"`
	RuleID   string   `json:"ruleId"`
	Location Location `json:"location"`
	Metadata Metadata `json:"metadata"`
}
type Location struct {
	Path     string   `json:"path"`
	Position Position `json:"position"`
}

type Position struct {
	Begin Begin `json:"begin"`
}

type Begin struct {
	Line string `json:"line"`
}
type Metadata struct {
	Description string `json:"description"`
	Severity    string `json:"severity"`
}

func FindWords(dir string) ([]*Findings, error) {
	keyWords := []string{"private_key", "public_key"}
	results := make([]*Findings, 0)

	for _, key := range keyWords {
		conf := exec.Command("grep", "-w", "-r", "-n", key)
		conf.Dir = dir

		out, _ := conf.Output()
		res := strings.Split(string(out), "\n")
		for _, r := range res {
			a := strings.Split(r, ":")
			if len(a) > 2 {
				f :=
					Findings{
						Type:   "sast",
						RuleID: "G402",
						Location: Location{
							Path: a[0],
							Position: Position{
								Begin: Begin{
									Line: a[1],
								},
							},
						},
						Metadata: Metadata{
							Description: fmt.Sprintf("%s exists here", key),
							Severity:    "HIGH",
						},
					}
				results = append(results, &f)
			}
		}
	}
	return results, nil
}
