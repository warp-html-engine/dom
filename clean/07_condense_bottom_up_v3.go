package clean

import (
	"strings"

	"github.com/pbberlin/dom/node"
	"github.com/zew/logx"
	"golang.org/x/net/html"
)

func flattenSubtreeV3(n, nClone *html.Node) {

	// logx.Printf("fsbo\n")
	flattenSubtreeV3Inner(n, nClone, 0)

}

var standard = map[string]bool{

	"title": true,

	"p":   true,
	"div": true,
	"ul":  true,
	"ol":  true,
	"li":  true,
	"h1":  true,
	"h2":  true,

	"em":       true,
	"strong":   true,
	"label":    true,
	"input":    true,
	"textarea": true,

	"form":       true,
	"blockquote": true,
}

func flattenSubtreeV3Inner(n, nClone *html.Node, lvl int) {

	// logx.Printf("fsbi\n")

	for ch := n.FirstChild; ch != nil; ch = ch.NextSibling {

		chClone := node.CloneNode(ch)

		switch {

		case ch.Type == html.ElementNode && standard[ch.Data]:
			nClone.AppendChild(chClone)
			flattenSubtreeV3Inner(ch, chClone, lvl+1)

		case ch.Type == html.ElementNode && ch.Data == "a":
			nClone.AppendChild(chClone)
			flattenSubtreeV3Inner(ch, chClone, lvl+1)

		case ch.Type == html.ElementNode && ch.Data == "img":
			nClone.AppendChild(chClone)

		case ch.Data == "span":
			// logx.Printf(strings.Repeat("  ", lvl) + "span \n")
			for cch := ch.FirstChild; cch != nil; cch = cch.NextSibling {
				// logx.Printf(strings.Repeat("    ", lvl)+"span child %v", cch.Data)
				cchClone := node.CloneNode(cch)
				nClone.AppendChild(cchClone)
				nClone.AppendChild(node.NewTextNode(" "))
				flattenSubtreeV3Inner(cch, cchClone, lvl+1)
			}

		case ch.Type == html.TextNode && ch.Data != "":
			chClone.Data = strings.TrimSpace(chClone.Data)
			chClone.Data += " "
			nClone.AppendChild(chClone)

		default:
			//			nClone.AppendChild(chClone)
			logx.Printf("unhandled %s %s\n", node.NodeTypeStr(ch.Type), ch.Data)

		}

	}

}

func condenseBottomUpV3(n *html.Node, lvl, lvlDo int, unusedTypes map[string]bool) {

	if lvl < lvlDo {

		// Delve deeper until reaching lvlDo
		cs := []*html.Node{}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			cs = append(cs, c)
		}
		for _, c := range cs {
			condenseBottomUpV3(c, lvl+1, lvlDo, unusedTypes)
		}

	} else {

		if n.Type == html.ElementNode {

			nClone := node.CloneNode(n)
			flattenSubtreeV3(n, nClone)

			nParent := n.Parent
			nParent.InsertBefore(nClone, n)
			nParent.RemoveChild(n)

			// 	bx := node.PrintSubtree(nParent)
			// 	fmt.Printf("%s", bx)
		}

	}

}
