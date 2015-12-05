package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const version = "0.1"

var indexFile, outputFile string
var verboseMode, hasOutputSetting bool
var totalSource uint
var spPattern = regexp.MustCompile(`\s+`)
var outputPattern = regexp.MustCompile(`#output\s+\S+`)
var includePattern = regexp.MustCompile(`#include\s+\S+`)

func init() {
	var helpOnly bool
	// defines flags
	flag.StringVar(&indexFile, "index", "", "path of the index file, defaults to stdin")
	flag.StringVar(&outputFile, "output", "", "path of the output file, defaults to stdout")
	flag.BoolVar(&verboseMode, "verbose", false, "show more info while running")
	flag.BoolVar(&helpOnly, "help", false, "show this help")

	flag.Parse()

	if helpOnly {
		fmt.Printf("Usage of Text Builder v%v:\n", version)
		flag.PrintDefaults()
		os.Exit(0)
	}

	log.SetOutput(os.Stdout)
	log.SetFlags(0)
}

func main() {
	startTime := time.Now()
	info("Text Builder version", version)

	lines := parseFile(indexFile)

	var writer io.Writer
	if len(outputFile) == 0 {
		writer = os.Stdout
	} else {
		file, err := os.Create(outputFile)
		check(err)
		defer file.Close()
		writer = file
	}

	fmt.Fprint(writer, strings.Join(lines, "\n"))
	info("Built", totalSource, "files in", time.Now().Sub(startTime).Seconds(), "s.")
}

func parseFile(path string) []string {
	totalSource++

	var lines []string
	var reader io.Reader
	remoteMode := isRemoteURL(path)

	if len(path) == 0 {
		reader = os.Stdin
	} else if !remoteMode {
		info("reading file:", path)
		file, err := os.Open(path)
		check(err)
		defer file.Close()
		reader = file
	} else {
		info("loading remote:", path)
		response, err := http.Get(path)
		check(err)
		defer response.Body.Close()

		if response.StatusCode < 200 || response.StatusCode >= 400 {
			log.Fatalln("failed with status", response.Status)
		}

		reader = response.Body
	}

	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		line := scanner.Text()

		if !hasOutputSetting {
			hasOutputSetting = true
			if outputPattern.MatchString(line) {
				out := strings.TrimSpace(spPattern.Split(line, 2)[1])
				if len(out) > 0 {
					out = resolvePath(path, out)
					info("Specified output file is", out)
					outputFile = out
				}
				continue
			}
		}

		if includePattern.MatchString(line) {
			nextFilename := strings.TrimSpace(spPattern.Split(line, 2)[1])
			if len(nextFilename) > 0 {
				nextFilename = resolvePath(path, nextFilename)
				info("including content", nextFilename)
				lines = append(lines, parseFile(nextFilename)...)
			}
		} else {
			lines = append(lines, line)
		}
	}

	err := scanner.Err()
	check(err)

	return lines
}

func check(e error) {
	if e != nil {
		log.Fatalln(e)
	}
}

func isRemoteURL(url string) bool {
	return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
}

func resolvePath(basePath string, relativePath string) string {
	if isRemoteURL(basePath) || isRemoteURL(relativePath) || filepath.IsAbs(relativePath) {
		return relativePath
	}
	return filepath.Join(filepath.Dir(basePath), relativePath)
}

func info(v ...interface{}) {
	if !verboseMode {
		return
	}
	log.Println(v...)
}
