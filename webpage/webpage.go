package webpage

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/JonPulfer/mCrawler/weblink"
	"golang.org/x/net/html"
)

// BaseURL of the site being processed.
var BaseURL *url.URL

// Elements of the webpage extracted during analysis.
type Elements struct {
	Body  *html.Node
	Links []weblink.Resource
}

// NewElements extracts the elements from the page data.
func NewElements(resp *http.Response) *Elements {

	return &Elements{}
}

// ProcessDownloadedPage examines the page accessible through the given reader and extracts the
// elements of interest on the page.
func ProcessDownloadedPage(r io.Reader) (*Elements, error) {

	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}
	elems := &Elements{
		Body: doc,
	}

	log.Println("scanning for links")

	l := scanForLinks(nil, doc)
	if len(l) > 0 {
		clnks, err := cleanDiscoveredLinks(l)
		if err != nil {
			return elems, err
		}
		wLinks := make([]weblink.Resource, 0, 1)
		for _, cl := range clnks {
			log.Printf("found link: %s\n", cl)
			r := weblink.Resource{
				TargetURI: cl,
				Type:      weblink.Hyperlink,
			}
			wLinks = append(wLinks, r)
		}
		elems.Links = wLinks
	}

	return elems, nil
}

func scanForLinks(links []string, n *html.Node) []string {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			links = append(links, a.Val)
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = scanForLinks(links, c)
	}
	return links
}

func cleanDiscoveredLinks(l []string) ([]string, error) {
	cleaned := make([]string, 0, 1)
	urlsFound := make(map[string]bool)
	for _, lk := range l {
		lk = strings.TrimRight(lk, "/")
		cURL, err := url.Parse(lk)
		if err != nil {
			return nil, err
		}
		if cURL.Host == BaseURL.Host || len(cURL.Host) == 0 {
			log.Printf("lk: %s\ncURL: %#v\n", lk, cURL)
			if len(cURL.Path) < 2 || cURL.Path[0] != '/' {
				continue
			}
			if _, ok := urlsFound[cURL.Path]; !ok {
				urlsFound[cURL.Path] = true
			}
		}
	}
	for k := range urlsFound {
		cleaned = append(cleaned, k)
	}

	return cleaned, nil
}
