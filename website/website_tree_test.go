package website

import "testing"

func TestInsertPage(t *testing.T) {
	var p Page
	p.URL = "/some/wonderful/dir/thispage"
	err := SiteContent.addPageToTree(&p)
	if err != nil {
		t.FailNow()
	}
	t.Logf("SiteContent.Tree: %#v\n", SiteContent.Tree)
}
