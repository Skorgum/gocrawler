package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Println("not enough arguments provided.")
		fmt.Println("usage: gocrawler <url> <max_concurrency> <max_pages>")
		os.Exit(1)
	}

	if len(os.Args) > 4 {
		fmt.Println("too many arguments provided.")
		os.Exit(1)
	}

	rawBaseURL := os.Args[1]

	maxConcurrency, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Printf("error parsing concurrency: %v\n", err)
		return
	}
	maxPages, err := strconv.Atoi(os.Args[3])
	if err != nil {
		fmt.Printf("error parsing max pages: %v\n", err)
		return
	}
	cfg, err := configure(rawBaseURL, maxConcurrency, maxPages)
	if err != nil {
		fmt.Printf("error in configuration: %v\n", err)
		return
	}

	fmt.Printf("starting crawl of: %s\n", rawBaseURL)

	cfg.wg.Add(1)
	go cfg.crawlPage(rawBaseURL)
	cfg.wg.Wait()

	if err := writeCSVReport(cfg.pages, "report.csv"); err != nil {
		fmt.Printf("error writing CSV report: %v\n", err)
	} else {
		fmt.Println("CSV report written to report.csv")
	}
}
