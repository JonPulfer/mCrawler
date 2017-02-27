package queue

import (
	"log"
	"strings"
	"sync"

	"github.com/JonPulfer/mCrawler/pageloader"
)

// Items on the queue awaiting processing.
type Items struct {
	sync.RWMutex
	Stack  []*pageloader.Request
	Seen   map[string]bool
	Length int
}

// NewItems returns an initialised item queue.
func NewItems() *Items {
	return &Items{
		Stack: make([]*pageloader.Request, 0, 1),
		Seen:  make(map[string]bool),
	}
}

// Add a new item onto the queue.
func (qi *Items) Add(r *pageloader.Request) {

	if qi.haveSeen(r.URL) {
		return
	}

	qi.Lock()
	log.Println("adding request to the queue")
	qi.Stack = append(qi.Stack, r)
	qi.Seen[r.URL] = true
	qi.Length++
	log.Printf("queue length now: %d\n", qi.Length)
	qi.Unlock()
}

// Next item to process on the queue. This pops the item from the front of the queue.
func (qi *Items) Next() *pageloader.Request {
	qi.Lock()
	log.Println("pulling request from the queue")
	if qi.Length == 0 {
		log.Println("nothing in the queue")
		return nil
	}
	t := qi.Stack[0]
	tCopy := *t
	qi.Stack = qi.Stack[1:]
	qi.Length--
	log.Printf("queue length now: %d\n", qi.Length)
	qi.Unlock()
	return &tCopy
}

// Len of the queue currently.
func (qi *Items) Len() int {
	qi.RLock()
	c := qi.Length
	log.Printf("queue length: %d\n", c)
	qi.RUnlock()
	return c
}

func (qi *Items) haveSeen(s string) bool {
	qi.RLock()
	defer qi.RUnlock()

	if _, ok := qi.Seen[s]; ok {
		return true
	}

	if _, ok := qi.Seen[strings.TrimRight(s, "/")]; ok {
		return true
	}

	if _, ok := qi.Seen[s+"/"]; ok {
		return true
	}

	return false
}
