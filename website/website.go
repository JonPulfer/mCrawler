package website

import (
	"fmt"
	"log"
	"net/url"
	"strings"
	"sync"

	"github.com/JonPulfer/mCrawler/webpage"
)

// SiteContent discovered during analysis.
var SiteContent *Content

func init() {
	var sc Content
	p := make(map[string]Page)
	sc.Pages = p
	sc.Tree = NewNode()
	sc.Tree.Name = "/"
	SiteContent = &sc
}

// Content found on the website.
type Content struct {
	sync.RWMutex
	Pages Pages
	Tree  *Node
}

// Node of tree.
type Node struct {
	Left        *Node
	Name        string
	Directories []*Node
	Pages       []*Page
}

// NewNode returns an initialised node ready for use.
func NewNode() *Node {
	var n Node
	n.Pages = make([]*Page, 0, 1)
	n.Directories = make([]*Node, 0, 1)
	return &n
}

// String method to print the node details.
func (n *Node) String() string {
	output := fmt.Sprintf("tree:\n\t%s\n", n.Name)
	for _, ps := range n.Pages {
		output = output + fmt.Sprintf("\t\t%s\n", ps)
	}
	for _, ds := range n.Directories {
		output = output + fmt.Sprintf("\t%s\n", ds)
	}

	return output
}

// IsNode checks whether this is the node identified by the given string.
func (n *Node) IsNode(s string) bool {
	return n.Name == s
}

// GetDirectory returns the node representing the child directory with the provided name or creates
// the directory in the node and returns that.
func (n *Node) GetDirectory(s string) *Node {
	if len(s) == 0 {
		return n
	}
	for _, v := range n.Directories {
		if v.IsNode(s) {
			return v
		}
	}
	return n.InsertDirectory(s)
}

// InsertDirectory into the node.
func (n *Node) InsertDirectory(s string) *Node {
	nd := NewNode()
	nd.Left = n
	nd.Name = s
	n.Directories = append(n.Directories, nd)
	log.Printf("Inserting directory %s into %#v\n", s, n)
	return nd
}

// InsertPage into the node.
func (n *Node) InsertPage(p *Page) {
	n.Pages = append(n.Pages, p)
	log.Printf("Inserting page %s into node: %#v\n", p.URL, n)
}

// GetParent returns the parent node.
func (n *Node) GetParent() *Node {
	return n.Left
}

func (c *Content) addPageToTree(p *Page) error {
	pURL, err := url.Parse(p.URL)
	if err != nil {
		return err
	}

	parts := strings.Split(pURL.Path, "/")
	SiteContent.Lock()
	t := SiteContent.Tree
	for i := 0; i < len(parts)-1; i++ {
		t = t.GetDirectory(parts[i])
	}
	t.InsertPage(p)
	SiteContent.Unlock()
	return nil
}

func (c *Content) insertPageIntoContent(p Page) error {
	SiteContent.Lock()
	SiteContent.Pages[p.URL] = p
	log.Printf("page %s inserted into SiteContent\n", p.URL)
	SiteContent.Unlock()
	return nil
}

// Analysis holds the data extracted during the analysis of the website.
type Analysis struct {
	SiteDepth int
	PageCount int
	RootPage  *Page
}

// Page found on the site
type Page struct {
	URL       string
	LinkCount int
	Elements  *webpage.Elements
}

func (p *Page) String() string {
	urlParts := strings.Split(p.URL, "/")
	return urlParts[len(urlParts)-1]
}

// Pages in the site.
type Pages map[string]Page

// InsertPage into the site content page list.
func InsertPage(p Page) error {
	SiteContent.RLock()
	_, ok := SiteContent.Pages[p.URL]
	SiteContent.RUnlock()

	if !ok {
		if err := SiteContent.insertPageIntoContent(p); err != nil {
			return err
		}
		if err := SiteContent.addPageToTree(&p); err != nil {
			return err
		}
	}

	return nil
}

// PageVisited indicates whether the page has already been visited.
func PageVisited(p Page) bool {
	SiteContent.RLock()
	_, ok := SiteContent.Pages[p.URL]
	SiteContent.RUnlock()
	return ok
}

// ListTree currently in the site content.
func ListTree() {
	SiteContent.RLock()
	log.Printf("%s\n", SiteContent.Tree)
	SiteContent.RUnlock()
}
