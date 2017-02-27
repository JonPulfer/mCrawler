package queue

import (
	"testing"

	"github.com/JonPulfer/mCrawler/pageloader"
)

func TestAddItemToQueue(t *testing.T) {
	q := NewItems()
	r := &pageloader.Request{URL: "http://golang.org/"}
	q.Add(r)
	if q.Len() != 1 {
		t.Fail()
	}
}

func TestAddTwoItemsToQueue(t *testing.T) {
	q := NewItems()

	r1 := &pageloader.Request{URL: "http://golang.org/"}
	q.Add(r1)

	r2 := &pageloader.Request{URL: "http://google.com/"}
	q.Add(r2)

	if q.Len() != 2 {
		t.Fail()
	}
}

func TestDuplicateItemIgnored(t *testing.T) {
	q := NewItems()

	r1 := &pageloader.Request{URL: "http://golang.org/"}
	q.Add(r1)

	r2 := &pageloader.Request{URL: "http://golang.org/"}
	q.Add(r2)

	if q.Len() != 1 {
		t.Fail()
	}
}

func TestDuplicateIdentifiedWhenNoTrailingSlash(t *testing.T) {
	q := NewItems()

	r1 := &pageloader.Request{URL: "http://golang.org/"}
	q.Add(r1)

	r2 := &pageloader.Request{URL: "http://golang.org"}
	q.Add(r2)

	if q.Len() != 1 {
		t.Fail()
	}
}

func TestDuplicateIdentifiedWithTrailingSlash(t *testing.T) {
	q := NewItems()

	r1 := &pageloader.Request{URL: "http://golang.org"}
	q.Add(r1)

	r2 := &pageloader.Request{URL: "http://golang.org/"}
	q.Add(r2)

	if q.Len() != 1 {
		t.Fail()
	}
}
