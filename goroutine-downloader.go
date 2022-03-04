package main

import (
    "fmt"
    "io"
    "net/http"
    "os"
    "path/filepath"
)

type repoStats struct {
    name string
    format string
    size int64
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

// We might want to introduce a master/main branch check possibility
func checkRepoFor(repoName string, format *format) *repoStats {
    repoUrl := fmt.Sprintf("https://raw.githubusercontent.com/NovatecConsulting/%v/master/README.%v", repoName, format.fileEnding)
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
        stat := repoStats{repoName, format.name, size}
        return &stat
    } else {
    //    fmt.Println("Unable to find a read me at", repoUrl)
        return nil
    }
}

func main() {
    arguments := os.Args[1:]
    if (len(arguments) == 0) {
        return
    }
    markdown := format { "markdown", "md" }
    asciidoc := format { "asciidoc", "adoc" }
    formats := [2]format { markdown, asciidoc }
    for _, argument := range arguments {
        fmt.Println("Checking repository named " + argument)
        for _, format := range formats {
            stats :=checkRepoFor(argument, &format)
            if(stats != nil) {
                fmt.Println(stats.name, stats.format, stats.size)
            }
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
