package main

import (
	"fmt"
	"guardrails/src"
	"math/rand"
	"os"
	"os/exec"
	"time"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	if _, err := os.Stat("/tmp/linux"); os.IsNotExist(err) {
		fmt.Println("Downloading linux repo")
		conf := exec.Command("git", "clone", "https://github.com/torvalds/linux")
		conf.Dir = "/tmp"
		conf.Output()
	}

	reqs := make(chan int, 1000)
	resp := make(chan bool, 1000)
	fmt.Println("starting thrasher")

	// 16 concurrent processes
	for i := 0; i < 16; i++ {
		go func() {
			for {
				_ = <-reqs
				src.FindWords("/tmp/linux")
				resp <- true
			}
		}()
	}

	count := 500

	start := time.Now()
	for i := 0; i < count; i++ {
		reqs <- i
	}

	for i := 0; i < count; i++ {
		if <-resp == false {
			fmt.Println("ERROR on", i)
			os.Exit(-1)
		}
	}

	fmt.Println(count, "counts in", time.Since(start))
	fmt.Printf("thats %.2f repo/sec\n", float32(count)/(float32(time.Since(start))/1e9))
}
