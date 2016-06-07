package textextract

import ("strings"
	"golang.org/x/net/html"
	"bytes"
	"fmt"
	"regexp"
	"errors"
)

var MinScore = 5

func isInAnchor(n *html.Node) bool {
	if n.Parent == nil {
		return false
	}
	if n.Parent.Data == "a" {
		return true
	}
	return isInAnchor(n.Parent)
}

func normaliseText(t string) string {
	r, _ := regexp.Compile("<[^>]*>|\\n|\\t| +")	
	r2, _ := regexp.Compile("^ +| +$")
	return r2.ReplaceAllString(r.ReplaceAllString(
		r.ReplaceAllString(t, " "),
		" "),"")
}

func filter(doc *html.Node, minScore int) *html.Node {
	type NodePair struct {
		Parent *html.Node
		Child *html.Node
	}
	toDelete := []NodePair{}
	var f func(n *html.Node, score int) int
	f = func(n *html.Node, score int) int {
		if n.Type == html.TextNode {
			count := len(strings.Split(normaliseText(n.Data), " "))
			switch {
			case n.Parent.Data == "script":
			case n.Parent.Data == "style":
			case n.Parent.Data == "link":
			case isInAnchor(n):
				score -= 1 + count^2
			default:
				score += count
			}
			return score
		}

		ownScore := score
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			score += f(c, ownScore)
		}

		if score <= minScore && n.Data != "a"  {
			toDelete = append(toDelete, NodePair{n.Parent,n})
		}
		return score
	}
	f(doc,0)

	for _, x := range toDelete {
		if x.Parent != nil {
			x.Parent.RemoveChild(x.Child)
		}
	}
	return doc
}

func ExtractFromHtml(htmlUTF8Str string) (string,error) {
	doc, err := html.Parse(strings.NewReader(htmlUTF8Str))
	if err != nil {
		return "", errors.New("Could not parse HTML string.")
	}
	doc = filter(doc, MinScore)
	var f func(n *html.Node)
	var buffer bytes.Buffer
	f = func(n *html.Node) {
		d := normaliseText(n.Data)
		if n.Type == html.TextNode && d != "" && d!= " " {
			switch n.Parent.Data {
			case "title":
			default:
				buffer.WriteString(fmt.Sprintf("\n%s", d))
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return buffer.String(), nil
}
