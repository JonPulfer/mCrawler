package pageloader

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/JonPulfer/mCrawler/webpage"
	"github.com/JonPulfer/mCrawler/website"
)

// Request for the Loader to enable a link to be processed.
type Request struct {
	URL string
}

// Worker for loading webpage URLs received on the channel.
func Worker(q <-chan *Request, pq chan<- *Request, wg *sync.WaitGroup) {
	defer wg.Done()

	for lr := range q {
		log.Printf("processing %s\n", lr.URL)
		if err := downloadPage(lr, pq); err != nil {
			errMsg := fmt.Sprintf("Loader: error downloading %s: %s\n", lr.URL, err.Error())
			log.Println(errMsg)
		}
	}
}

func downloadPage(lr *Request, pq chan<- *Request) error {
	log.Printf("downloading page %s\n", lr.URL)

	if webpage.BaseURL == nil {
		bURL, err := url.Parse(lr.URL)
		if err != nil {
			return err
		}
		webpage.BaseURL = bURL
	} else {
		lr.URL = strings.TrimRight(webpage.BaseURL.String(), "/") + lr.URL
	}
	if website.PageVisited(website.Page{URL: lr.URL}) {
		log.Printf("already visited %s\n", lr.URL)
		return nil
	}

	cl := http.Client{}
	req, err := http.NewRequest(http.MethodGet, lr.URL, nil)
	if err != nil {
		return err
	}

	resp, err := cl.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return processResponseStatus(resp)
	}

	elems, err := webpage.ProcessDownloadedPage(resp.Body)
	if err != nil {
		return err
	}

	page := website.Page{
		URL:      lr.URL,
		Elements: elems,
	}
	if err := website.InsertPage(page); err != nil {
		return err
	}

	for _, lnk := range elems.Links {
		req := &Request{
			URL: lnk.TargetURI,
		}
		newP := website.Page{URL: lnk.TargetURI}
		if website.PageVisited(newP) {
			continue
		}

		if !website.PageVisited(website.Page{URL: lnk.TargetURI}) {
			log.Printf("adding %s to queue", lnk.TargetURI)
			pq <- req
			log.Println("added to queue")
		}
	}

	return nil
}

func processResponseStatus(resp *http.Response) error {

	switch resp.StatusCode {
	case http.StatusNotFound:
		return fmt.Errorf("page not found")
	case http.StatusServiceUnavailable:
		return fmt.Errorf("server reported unavailable, could retry")
	}

	return fmt.Errorf("unhandled status code received: %d", resp.StatusCode)
}
