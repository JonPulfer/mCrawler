package main

import (
	"flag"
	"log"
	"sync"
	"time"

	"github.com/JonPulfer/mCrawler/pageloader"
	"github.com/JonPulfer/mCrawler/queue"
	"github.com/JonPulfer/mCrawler/website"
)

var bURL string
var wCount int

func init() {
	flag.StringVar(&bURL, "bURL", "http://golang.org/", "URL of the site to crawl")
	flag.IntVar(&wCount, "w", 5, "Number of workers to download pages")
	flag.Parse()
}

func main() {
	workers := wCount
	baseURL := bURL

	linkQueue := make(chan *pageloader.Request, workers)
	pageQueue := make(chan *pageloader.Request, workers)

	workQueue := queue.NewItems()

	var loaderWG sync.WaitGroup
	var reloadWG sync.WaitGroup

	for i := 0; i < workers; i++ {
		loaderWG.Add(1)
		go pageloader.Worker(linkQueue, pageQueue, &loaderWG)
	}
	reloadWG.Add(1)
	go QueueDiscoveredPage(pageQueue, linkQueue, workQueue, &reloadWG)

	baseRequest := &pageloader.Request{
		URL: baseURL,
	}
	pageQueue <- baseRequest

	// Once the reloading worker has identified that the queue is empty, we know we can close
	// the linkQueue.
	reloadWG.Wait()
	log.Println("closing linkQueue down")
	close(linkQueue)

	loaderWG.Wait()

	log.Printf("finished crawling %s\n", baseURL)
	website.ListTree()
}

// QueueDiscoveredPage receives load requests generated from the items found by the page loaders and
// queues them for processing.
func QueueDiscoveredPage(pageq <-chan *pageloader.Request,
	loadq chan<- *pageloader.Request, workQueue *queue.Items, reloadWG *sync.WaitGroup) {
	defer reloadWG.Done()

	var qCheck int
	time.Sleep(1 * time.Second)

	for {
		select {
		case x := <-pageq:
			if website.PageVisited(website.Page{URL: x.URL}) {
				continue
			}
			workQueue.Add(x)
		default:
			if workQueue.Len() > 0 {
				log.Println("item on queue to process")
				qCheck = 0
				for n := workQueue.Next(); n != nil; n = workQueue.Next() {
					if website.PageVisited(website.Page{URL: n.URL}) {
						continue
					}
					loadq <- n
					break
				}
			} else {
				qCheck++
				if qCheck > wCount {
					log.Printf("checked queue %d times and found it empty, closing down", qCheck)
					return
				}
			}
			time.Sleep(300 * time.Millisecond)
		}
	}
}
