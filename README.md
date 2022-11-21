# Security Scan
- This project scan small security issues related to public_key and private_key.
- I think this should be simple so the total of core codes in under 500 lines.
- `main` is out file is really small ~13MB.
- I only use `github.com/jackc/pgx/v5` lib to work with Postgres Database

## Project Structure
- The project structure is really simple
  - src: contains all core source codes
    - main.go: defines how to run the services
    - lib.go & lib_test.go: contains how to download and scan keys in the repositories
    - repository.go: contains all MVC of Repository Service. Repository Service is defined at "/repository" API
    - scan.go contains all MVC of Scan Service. Scan Service is defined at "/scan" API
  - tools: contains deployment, stress-test and integration test tools
  - migrations: contains the migration instructions

## Database Schema
- repository table
```
  id SERIAL primary key = unique repository id 
  user_id text not null = the unique user id
  repo_url text not null = the repository url (github or gitlab)
  repo_name text not null = the repository name (defined by user)
```
- security_scan table
```
    id  serial primary key = the unique scan id
    status text not null = the scan status, value is ( Queued" | "In Progress" | "Success" | "Failure")
    user_id text not null = the unique user id 
    repo_url text not null = the repository url (github or gitlab)
    repo_name text not null = the repository name (defined by user)
    findings bytea = the securiity scan results
    queued_at timestamptz = the time which the scan is queued
    scanning_at timestamptz = the time which the scan is scanning
    finished_at timestamptz = = the time which the scan is finished
```

## API
### Repository Service
  - GET /repository?user_id=<user_id>
    - 200, return this payload 
    ```[
      {
        "id": <internal_repo_id>,
        "user_id": <user_id>,
        "repo_url": <repo_url>,
        "repo_name": <repo_name>
      }
    ]
    ```
  - POST /repository?user_id=<user_id>&repo_name=<repo_name>&repo_url=<repo_url>
    - 201 = written
  - PUT /repository?user_id=<user_id>&repo_id=<repo_id>&repo_name=<repo_name>&repo_url=<repo_url>
    - 200 = updated
  - DELETE /repository?user_id=<user_id>&repo_id=<repo_id>
    - 204 = deleted
  
### Scan Service
- GET /scan?user_id=<user_id>
    - 200, return this payload 
    ```
    [
      {
        "id": <internal_scan_id>,
        "status": <status>,
        "user_id": <user_id>,
        "repo_url": <repo_url>,
        "repo_name": <repo_name>,
        "findings": [{
          "type": "sast",
          "ruleId": "G402",
          "location": {
            "path": "connectors/apigateway.go",
            "positions": {
              "begin": {
                "line": 60
              }
            }
          },
          "metadata": {
            "description": "TLS InsecureSkipVerify set true.",
            "severity": "HIGH"
          }],
        "queued_at": <queued_at>,
        "scanning_at": <scanning_at>,
        "finished_at": <finished_at>,
      }
    ]

    ```
  - POST /scan?user_id=<user_id>&repo_id=<repo_id>
    - 201 = requested to scan

## Usage
```
# post a repository
curl -v -L -X POST "localhost:3000/repository?user_id=1&repo_url=https://github.com/duyle-py/security_scan&repo_name=sec"

# get with user_id = 1
curl -v -L "localhost:3000/repository?user_id=1"

# result is [{"id":29,"user_id":"1","repo_name":"sec","repo_url":"https://github.com/duyle-py/security_scan"}]

# put with new name
curl -v -L -X PUT "localhost:3000/repository?user_id=1&id=29&repo_name=new_sec&repo_url=https://github.com/duyle-py/security_scan"

# get with user_id = 1 
curl -v -L "localhost:3000/repository?user_id=1" 
# result is [{"id":29,"user_id":"1","repo_name":"new_sec","repo_url":"https://github.com/duyle-py/security_scan"}] 

# post scan
curl -v -L -X POST "localhost:3000/scan?user_id=1&repo_id=29" 

# get scan
curl -v -L  "localhost:3000/scan?user_id=1" 
# result is [{"id":3,"status":"Success","user_id":"1","repo_url":"https://github.com/duyle-py/security_scan","repo_name":"new_sec","findings":[{"type":"sast","ruleId":"G402","location":{"path":"src/lib_test.go","position":{"begin":{"line":"20"}}},"metadata":{"description":"private_key exists here","severity":"HIGH"}},{"type":"sast","ruleId":"G402","location":{"path":"src/lib.go","position":{"begin":{"line":"58"}}},"metadata":{"description":"private_key exists here","severity":"HIGH"}},{"type":"sast","ruleId":"G402","location":{"path":"src/lib_test.go","position":{"begin":{"line":"20"}}},"metadata":{"description":"public_key exists here","severity":"HIGH"}},{"type":"sast","ruleId":"G402","location":{"path":"src/lib.go","position":{"begin":{"line":"58"}}},"metadata":{"description":"public_key exists here","severity":"HIGH"}}],"queued_at":"2022-11-20T10:38:28.167228Z","scanning_at":"2022-11-20T10:38:29.336548Z","finished_at":"2022-11-20T10:38:29.336548Z"}]   

```

## Infra design
- My project is designed for non-blocking security scan for repositories, so we can scale multiple machines.
- Because we have two main services here: Repository and Scan Services. 
  - Repository Service is used for CRUD Repository only.
  - Scan Service contains two Methods are POST which scan repo and GET which list scan results.
    - POST method is really heavy, we can split it into a service. Scale it into multiple machine.
- My codebase is really simple now, so split it into a microservice is really simple.

## Local Starting
```
docker compose up --build
```

## Intergration testing
```
docker compose run test
```

## Unit testing
```
go test src/*.go -v
```

## Stress test 
  - `time go run tools/thrasher.go`
  - I used linux repository ~ 4GB to test grep performance, this repository is already downloaded in my computer. My specification is i7-11800H @ 2.30GHz × 16 cores and 64GB RAMs.
```
starting thrasher
500 counts in 1m19.904484041s
thats 6.26 repo/sec
go run tools/thrasher.go  731.90s user 504.58s system 1543% cpu 1:20.12 total

```

## Extra features
- Fast downloading repository. Now i'm using `git clone`, it will fetch all codes into a machine. We can use `git checkout` to fetch some parts we need. But need to design how to store and redirect repository in multiple machine or a machine. 
- Grep searching. Check performance of `ripgrep` lib and replace grep `https://healeycodes.com/beating-grep-with-go` 
- 