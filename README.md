# Scan Repo
- Goal: < 500 lines

- 410 lines
```  410 total
  146 ./src/repository.go
  133 ./src/scan.go
   93 ./src/lib.go
   29 ./cmd/main.go
    9 ./src/lib_test.go
```

- benchmark

  - `time go run tools/thrasher.go`
  - Use Linux repo ~ 4GB to test grep performance
```
starting thrasher
500 counts in 1m19.904484041s
thats 6.26 repo/sec
go run tools/thrasher.go  731.90s user 504.58s system 1543% cpu 1:20.12 total

```