package fetcher

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)

// Fetcher container for keeping track of TotalJobs, ExecutedJobs, NextJob. Contains Mutex for synchronization
type Fetcher struct {
	TotalJobs      int
	NextJob        int
	SuccessfulJobs int
	ExecutedJobs   int
	StatusCode     int
	Read           []byte
	Urls           []string
	Mutex          *sync.Mutex
	Total          int
	Fetching       bool
}

// Init copies urls into Fetcher's Urls, assigns urls length to TotalJobs and initializes Fetcher strucrure variables
// returns nothing
func (f *Fetcher) Init(urls []string) {
	f.TotalJobs = len(urls)
	f.NextJob = 0
	f.SuccessfulJobs = 0
	f.ExecutedJobs = 0
	f.Total = 0
	f.Urls = urls
	f.Mutex = &sync.Mutex{}
	f.Fetching = false
}

// getNextUrl synchronously updates Fetcher's NextJob unit and
// returns url of next job if all jobs are not processed else returns "nil"
func (f *Fetcher) getNextUrl() string {
	f.Mutex.Lock()
	if f.NextJob >= f.TotalJobs {
		f.Mutex.Unlock()
		return "nil"
	}
	fmt.Println("Fetching from => ", f.Urls[f.NextJob])
	url := f.Urls[f.NextJob]
	f.NextJob += 1
	f.Mutex.Unlock()
	return url
}

// StartFetching creates channel for fetching web documents and updates Total & SuccessfulJobs counter on execution of successful job
// returns error if StartFetching is alreay running
func (fetch *Fetcher) StartFetching() error {

	fetch.Mutex.Lock()
	if fetch.Fetching {
		fetch.Mutex.Unlock()
		return errors.New("Already fetching")
	}
	fetch.Fetching = true
	fetch.Mutex.Unlock()

	listener1 := make(chan int)
	listener2 := make(chan int)
	listener3 := make(chan int)

	go fetchPage(fetch.getNextUrl(), listener1)
	go fetchPage(fetch.getNextUrl(), listener2)
	go fetchPage(fetch.getNextUrl(), listener3)

	for {
		if fetch.ExecutedJobs == fetch.TotalJobs {
			break
		}
		select {
		case rec1 := <-listener1:
			fetch.ExecutedJobs += 1
			if rec1 != -1 {
				fetch.Total += rec1
				fetch.SuccessfulJobs += 1
			}
			go fetchPage(fetch.getNextUrl(), listener1)

		case rec2 := <-listener2:
			fetch.ExecutedJobs += 1
			if rec2 != -1 {
				fetch.Total += rec2
				fetch.SuccessfulJobs += 1
			}
			go fetchPage(fetch.getNextUrl(), listener2)

		case rec3 := <-listener3:
			fetch.ExecutedJobs += 1
			if rec3 != -1 {
				fetch.Total += rec3
				fetch.SuccessfulJobs += 1
			}
			go fetchPage(fetch.getNextUrl(), listener3)
		}
	}

	return nil
}

// fetchPage fetches web page and passes size of contents fetched to a channel on successful fetch
// on error -1 as a value is passed to channel
func fetchPage(url string, listener chan int) {
	if url == "nil" {
		return
	}
	response, error := http.Get(url)
	if error != nil {
		fmt.Println("Failed fetching =>", url)
		listener <- -1
		return
	}
	contents, err := ioutil.ReadAll(response.Body)
	if err == nil {
		fmt.Println("Done fetching =>", url, "Size=>", len(contents))
	}
	listener <- len(contents)
}
