// +build split
// go test -tags=split

package clean

import (
	"bytes"
	"strings"
	"testing"

	"github.com/zew/exceldb/dom/node"
	"github.com/zew/logx"
	"golang.org/x/net/html"
)

var testDocs = make([]string, 2)

func init() {

	// inp
	testDocs[0] = `<!DOCTYPE html><html><head>
		</head><body>

			<p>
				<a href="/some/first/page.html">Links1:
					no img
				</a>
				<a href="/some/first/page.html">Links2:
					<span>text bef</span>
					<img src="/img1src" title="img-title-01" />
					<span>text aft</span>
				</a>
				
			</p>
			<div>
				<div>
					<a href="/some/first/page.html">Links3:
						<span>text2 bef</span>
						<span>
							<span>text3 bef</span>
							<img src="/img1src" title="img-title-02" />
							<span>text3 aft</span>
						</span>
						<span>text2 aft</span>
					</a>
				</div>
			</div>
		</body></html>`

	// want
	testDocs[1] = `<!DOCTYPE html><html><head></head><body>
	<p>
		<a href="/some/first/page.html">Links1:
					no img 
		</a>
		<a href="/some/first/page.html">Links2: 
			<span>text bef 
			</span>
		</a>
		<a href="/img1src" title="img-title-01" cfrom="img">[img] img-title-01 /img1src | 
		</a>
		<a href="/some/first/page.html">
			<span>text aft 
			</span>
		</a>
	</p>
	<div>
		<div>
			<a href="/some/first/page.html">Links3: 
				<span>text2 bef 
				</span>
				<span>
					<span>text3 bef 
					</span>
				</span>
			</a>
			<a href="/img1src" title="img-title-02" cfrom="img">[img] img-title-02 /img1src | 
			</a>
			<a href="/some/first/page.html">
				<span>
					<span>text3 aft 
					</span>
				</span>
				<span>text2 aft 
				</span>
			</a>
		</div>
	</div></body></html>`

}

func Test2(t *testing.T) {

	doc, err := html.Parse(strings.NewReader(testDocs[0]))
	if err != nil {
		logx.Printf("%v", err)
		return
	}
	removeCommentsAndInterTagWhitespace(node.NdX{doc, 0})

	breakoutImgFromLinkTrees(doc)
	img2Link(doc)

	removeCommentsAndInterTagWhitespace(node.NdX{doc, 0})
	reIndent(doc, 0)
	var b bytes.Buffer
	err = html.Render(&b, doc)
	// logx.Printf("%v", err)
	if b.String() != testDocs[1] {
		t.Errorf("output unexpted")
	}

	Bytes2File("outp1_inp.html", []byte(testDocs[0]))
	Dom2File("outp2_got.html", doc)
	Bytes2File("outp3_want.html", []byte(testDocs[1]))

	logx.Printf("end")

}
