package main

import (
	"fetch_source/fetcher"
	"fmt"
)

const (
	URLCOUNT = 20
)

// fillUrls initializes urls
func fillUrls(urls []string) {
	urls[0] = "http://www.tonycuffe.com/mp3/tail%20toddle.mp3"
	urls[1] = "https://facebook.com"
	urls[2] = "https://google.com"
	urls[3] = "https://slack.com"
	urls[4] = "https://facebook.com"
	urls[5] = "https://google.com"
	urls[6] = "https://ssasdfdsfds.com" // this should fail
	urls[7] = "https://facebook.com"
	urls[8] = "https://google.com"
	urls[9] = "https://google.com"
	urls[10] = "https://facebook.com"
	urls[11] = "https://github.com"
	urls[12] = "https://groups.google.com"
	urls[13] = "https://facebook.com"
	urls[14] = "https://google.com"
	urls[15] = "https://gobyexample.com"
	urls[16] = "https://gmail.com"
	urls[17] = "https://stackoverflow.com"
	urls[18] = "https://slack.com"
	urls[19] = "https://golang.org"
}

// main starts execution for fetching urls, calculates average and prints them on successful execution
func main() {
	urls := make([]string, URLCOUNT)
	fillUrls(urls)
	var fetch = fetcher.Fetcher{}
	fetch.Init(urls)
	if err := fetch.StartFetching(); err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("")
	fmt.Println(" success in fetching ", fetch.SuccessfulJobs, " documents ")
	fmt.Println(" failed to fetch ", (fetch.TotalJobs - fetch.SuccessfulJobs), " documents ")
	fmt.Println(" average byte size of ", fetch.SuccessfulJobs, " fetched documents : ", (fetch.Total / fetch.SuccessfulJobs))
}
