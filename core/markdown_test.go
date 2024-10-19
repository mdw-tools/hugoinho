package core

import (
	"testing"

	"github.com/mdw-go/testing/should"
)

func TestGoldmarkMarkdownConverterFixture(t *testing.T) {
	should.Run(&GoldmarkMarkdownConverterFixture{T: should.New(t)}, should.Options.UnitTests())
}

type GoldmarkMarkdownConverterFixture struct {
	*should.T

	converter *GoldmarkMarkdownConverter
}

func (this *GoldmarkMarkdownConverterFixture) Setup() {
	this.converter = NewGoldmarkMarkdownConverter()
}

func (this *GoldmarkMarkdownConverterFixture) Test() {
	output, err := this.converter.Convert(MARKDOWN_INPUT)
	this.So(err, should.BeNil)
	this.So(output, should.Equal, EXPECTED_HTML_OUTPUT)

	// Assert correct reuse of internal buffer:
	output2, err2 := this.converter.Convert(MARKDOWN_INPUT)
	this.So(err2, should.BeNil)
	this.So(output2, should.Equal, EXPECTED_HTML_OUTPUT)
}

const MARKDOWN_INPUT = `
# H1

## H2

### H3

#### H4

- a
- b
- c

1. 1
2. 2
3. 3

[link](/target)

` + "`code`" + `

` + "```go\nfenced code\n```" + `

> blockquote

_emphasized_

**bolded**

---

| foo | bar |
| --- | --- |
| baz | bim |

Apple
:   Pomaceous fruit of plants of the genus Malus in the family Rosaceae.

Orange
:   The fruit of an evergreen tree of the genus Citrus.

That's some text with a footnote.[^1]

[^1]: And that's the footnote.

    That's the second paragraph.

`

const EXPECTED_HTML_OUTPUT = `<h1>H1</h1>
<h2>H2</h2>
<h3>H3</h3>
<h4>H4</h4>
<ul>
<li>a</li>
<li>b</li>
<li>c</li>
</ul>
<ol>
<li>1</li>
<li>2</li>
<li>3</li>
</ol>
<p><a href="/target">link</a></p>
<p><code>code</code></p>
<pre><code class="language-go">fenced code
</code></pre>
<blockquote>
<p>blockquote</p>
</blockquote>
<p><em>emphasized</em></p>
<p><strong>bolded</strong></p>
<hr>
<table>
<thead>
<tr>
<th>foo</th>
<th>bar</th>
</tr>
</thead>
<tbody>
<tr>
<td>baz</td>
<td>bim</td>
</tr>
</tbody>
</table>
<dl>
<dt>Apple</dt>
<dd>Pomaceous fruit of plants of the genus Malus in the family Rosaceae.</dd>
<dt>Orange</dt>
<dd>The fruit of an evergreen tree of the genus Citrus.</dd>
</dl>
<p>That's some text with a footnote.<sup id="fnref:1"><a href="#fn:1" class="footnote-ref" role="doc-noteref">1</a></sup></p>
<div class="footnotes" role="doc-endnotes">
<hr>
<ol>
<li id="fn:1">
<p>And that's the footnote.</p>
<p>That's the second paragraph.&#160;<a href="#fnref:1" class="footnote-backref" role="doc-backlink">&#x21a9;&#xfe0e;</a></p>
</li>
</ol>
</div>
`
