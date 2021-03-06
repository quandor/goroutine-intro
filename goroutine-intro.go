package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"
)

func printRepoStats(repoNames []string) {
	for _, repoName := range repoNames {
		printRepoStat(repoName)
	}
}

func printRepoStatsWaitGroup(repoNames []string) {
	var wg sync.WaitGroup
	for _, repoName := range repoNames {
		wg.Add(1)
		// new variable necessary. otherwise we just print the stat
		// for the first repo
		repoName := repoName
		go func() {
			defer wg.Done()
			printRepoStat(repoName)
		}()
	}

	wg.Wait()
}

func printRepoStatsChannel(repoNames []string) {
	stats := make(chan readmeStats, 2)

	go func() {
		var wg sync.WaitGroup
		for _, repoName := range repoNames {
			wg.Add(1)
			go func(repoName string) {
				defer wg.Done()
				fmt.Println("Checking repository named", repoName)
				stat := getRepoStats(repoName)
				stats <- stat
			}(repoName)
		}

		wg.Wait()
		close(stats)
	}()

	for stat := range stats {
		fmt.Printf("Quality for %v is %v\n", stat.name, determineQuality(stat))
	}
}

func main() {
	arguments := os.Args[1:]
	if len(arguments) == 0 {
		return
	}

	useWaitGroup := arguments[0] == "--waitgroup"
	useChannel := arguments[0] == "--channel"

	if useWaitGroup {
		printRepoStatsWaitGroup(arguments[1:])
	} else if useChannel {
		printRepoStatsChannel(arguments[1:])
	} else {
		printRepoStats(arguments)
	}
}

func checkRepoFor(repoName string, branchName string, format *format) *readmeStats {
	repoUrl := fmt.Sprintf("https://raw.githubusercontent.com/NovatecConsulting/%v/%v/README.%v", repoName, branchName, format.fileEnding)
	resp, err := http.Get(repoUrl)
	check(err)

	if resp.StatusCode == 200 {
		var size int64
		if resp.ContentLength != -1 {
			size = resp.ContentLength
		} else {
			content := resp.Body
			defer content.Close()
			bytes, err := io.ReadAll(content)
			check(err)
			size = int64(len(bytes))
		}
		stat := readmeStats{repoName, format.name, size, repoUrl}
		return &stat
	} else {
		return nil
	}
}

func getRepoStats(repoName string) readmeStats {
	delayRandomly()
	markdown := format{"markdown", "md"}
	asciidoc := format{"asciidoc", "adoc"}
	formats := [2]format{markdown, asciidoc}
	for _, branchName := range [2]string{"master", "main"} {
		for _, format := range formats {
			stats := checkRepoFor(repoName, branchName, &format)
			if stats != nil {
				return *stats
			}
		}
	}
	return readmeStats{repoName, "", 0, ""}
}

func delayRandomly() {
	rand.Seed(time.Now().UnixNano())
	//randomDelay, _ := time.ParseDuration(fmt.Sprintf("%vms", rand.Intn(600 -10) + 10))
	randomDelay, _ := time.ParseDuration("0ms")
	fmt.Println("We will delay by", randomDelay)
	time.Sleep(randomDelay)
}

func printRepoStat(repoName string) {
	fmt.Println("Checking repository named " + repoName)
	stats := getRepoStats(repoName)
	fmt.Printf("Quality for %v is %v\n", stats.name, determineQuality(stats))
}

func determineQuality(stats readmeStats) string {
	if stats.size > 3000 {
		return "Impressive"
	} else if stats.size > 2000 {
		return "Fantastic"
	} else if stats.size > 1000 {
		return "Excellent"
	} else if stats.size > 500 {
		return "Great"
	} else if stats.size > 0 {
		return "Good"
	} else {
		return "Standard"
	}
}

type readmeStats struct {
	name   string
	format string
	size   int64
	url    string
}

type format struct {
	name       string
	fileEnding string
}

func check(e error) {
	if e != nil {
		os.Exit(1)
	}
}
