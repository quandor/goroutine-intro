package main

import (
    "fmt"
    "io"
    "net/http"
    "os"
    "path/filepath"
)

type readmeStats struct {
    name string
    format string
    size int64
    url string
}

type format struct {
    name string
    fileEnding string
}


func check(e error) {
    if e != nil {
        panic(e)
    }
}

func writeToFile(file *os.File, content io.Reader) {
    fmt.Println("Writing to file:", file.Name())
    _, err := io.Copy(file, content)
    check(err)
}

func downloadReadMe(repoName string) io.ReadCloser {
    repoUrl := fmt.Sprintf("https://raw.githubusercontent.com/NovatecConsulting/%v/master/README.md", repoName)
    resp, err := http.Get(repoUrl)
    check(err)

    return resp.Body
}

func checkRepoFor(repoName string, branchName string, format *format) *readmeStats {
    repoUrl := fmt.Sprintf("https://raw.githubusercontent.com/NovatecConsulting/%v/%v/README.%v", repoName, branchName, format.fileEnding)
    resp, err := http.Get(repoUrl)
    check(err)

    if(resp.StatusCode == 200) {
        var size int64
        if(resp.ContentLength != -1) {
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

func getRepoStats(repoName string) *readmeStats {

    markdown := format { "markdown", "md" }
    asciidoc := format { "asciidoc", "adoc" }
    formats := [2]format { markdown, asciidoc }
    for _, branchName := range [2]string { "master", "main" } {
        for _, format := range formats {
            stats := checkRepoFor(repoName, branchName, &format)
            if(stats != nil) {
                return stats
            }
       }
   }
   return nil
}

func main() {
    arguments := os.Args[1:]
    if (len(arguments) == 0) {
        return
    }
    for _, argument := range arguments {
        fmt.Println("Checking repository named " + argument)
        stats := getRepoStats(argument)
        if(stats != nil) {
            fmt.Printf("Readme found for %v of type %v with size %v at %v\n", stats.name, stats.format, stats.size, stats.url)
        } else {
            fmt.Println("No reamdme found for " + argument)
        }
    }
}

func main_download() {
    arguments := os.Args[1:]
    if (len(arguments) == 0) {
        return
    }
    dirName, err := os.MkdirTemp("", "goroutine-downloader-")
    check(err)
    fmt.Println("Storing results in ", dirName)
    for _, argument := range arguments {
        absoluteFilePath := filepath.Join(dirName, argument)
        file, err2 := os.Create(absoluteFilePath)
        check(err2)
        defer file.Close()
        content := downloadReadMe(argument)
        defer content.Close()
        writeToFile(file, content)
    }
}
