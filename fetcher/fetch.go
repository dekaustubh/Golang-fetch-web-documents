package fetcher

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)

// Fetcher container for keeping track of TotalJobs, executedJobs, nextJob. Contains mutex for synchronization
type Fetcher struct {
	TotalJobs      int
	nextJob        int
	SuccessfulJobs int
	executedJobs   int
	StatusCode     int
	Read           []byte
	urls           []string
	mutex          *sync.Mutex
	Total          int
	fetching       bool
}

// Init copies urls into Fetcher's urls, assigns urls length to totalJobs and initializes Fetcher strucrure variables
// returns nothing
func (f *Fetcher) Init(urls []string) {
	f.TotalJobs = len(urls)
	f.nextJob = 0
	f.SuccessfulJobs = 0
	f.executedJobs = 0
	f.Total = 0
	f.urls = urls
	f.mutex = &sync.Mutex{}
	f.fetching = false
}

// getNextUrl synchronously updates Fetcher's nextJob unit and
// returns url of next job if all jobs are not processed else returns "nil"
func (f *Fetcher) getNextUrl() string {
	f.mutex.Lock()
	if f.nextJob >= f.TotalJobs {
		f.mutex.Unlock()
		return "nil"
	}
	fmt.Println("fetching from => ", f.urls[f.nextJob])
	url := f.urls[f.nextJob]
	f.nextJob += 1
	f.mutex.Unlock()
	return url
}

// StartFetching creates channel for fetching web documents and updates Total & SuccessfulJobs counter on execution of successful job
// returns error if StartFetching is alreay running
func (fetch *Fetcher) StartFetching() error {

	if fetch.mutex == nil {
		return errors.New("Not initialized")
	}

	fetch.mutex.Lock()
	if fetch.fetching {
		fetch.mutex.Unlock()
		return errors.New("Already fetching")
	}
	fetch.fetching = true
	fetch.mutex.Unlock()

	listener1 := make(chan int)
	listener2 := make(chan int)
	listener3 := make(chan int)

	go fetchPage(fetch.getNextUrl(), listener1)
	go fetchPage(fetch.getNextUrl(), listener2)
	go fetchPage(fetch.getNextUrl(), listener3)

	for {
		if fetch.executedJobs == fetch.TotalJobs {
			break
		}
		select {
		case rec1 := <-listener1:
			fetch.executedJobs += 1
			if rec1 != -1 {
				fetch.Total += rec1
				fetch.SuccessfulJobs += 1
			}
			go fetchPage(fetch.getNextUrl(), listener1)

		case rec2 := <-listener2:
			fetch.executedJobs += 1
			if rec2 != -1 {
				fetch.Total += rec2
				fetch.SuccessfulJobs += 1
			}
			go fetchPage(fetch.getNextUrl(), listener2)

		case rec3 := <-listener3:
			fetch.executedJobs += 1
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
		fmt.Println(" response content length ", response.ContentLength)
		fmt.Println("Done fetching =>", url, "Size=>", len(contents))
	}
	listener <- len(contents)
}
