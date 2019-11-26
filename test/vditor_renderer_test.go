// Lute - A structured markdown engine.
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under the Mulan PSL v1.
// You can use this software according to the terms and conditions of the Mulan PSL v1.
// You may obtain a copy of Mulan PSL v1 at:
//     http://license.coscl.org.cn/MulanPSL
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR
// PURPOSE.
// See the Mulan PSL v1 for more details.

// +build javascript

package test

import (
	"testing"

	"github.com/b3log/lute"
)

var vditorDOM2MdTests = []parseTest{

	{"44", "<p>f<i>o</i>o<wbr></p>", "f*o*o\n"},
	{"43", "<ul data-tight=\"true\"><li data-marker=\"*\">foo<br></li><ul><li data-marker=\"*\">b<wbr></li></ul></ul>", "* foo<br />\n  * b\n"},
	{"42", "<div class=\"vditor-panel vditor-panel--none\" contenteditable=\"false\" style=\"display: block; top: 5px; left: 567px;\"><input class=\"vditor-input\" placeholder=\"row\" style=\"width: 42px; text-align: center;\"> x <input class=\"vditor-input\" placeholder=\"column\" style=\"width: 42px; text-align: center;\"></div>", "\n"},
	{"41", "<pre><code class=\"language-go\"><wbr></code></pre>", "```go\n\n```\n"},
	{"40", "<p>f<span data-marker=\"*\">o</span>ob<wbr></p>", "foob\n"},
	{"39", "<p><b>foo<wbr></b></p>", "**foo**\n"},
	{"38", "<p>```java</p><p><wbr><br></p>", "```java\n\n<br />\n"},
	{"37", "<ul data-tight=\"true\"><li data-marker=\"*\">foo<wbr></li><li data-marker=\"*\"></li><li data-marker=\"*\"><br></li></ul>", "* foo\n*\n* <br />\n"},
	{"36", "<ul data-tight=\"true\"><li data-marker=\"*\">1<em data-marker=\"*\">2</em></li><li data-marker=\"*\"><em data-marker=\"*\"><wbr><br></em></li></ul>", "* 1*2*\n* *<br />*\n"},
	{"35", "<ul data-tight=\"true\"><li data-marker=\"*\"><wbr><br></li></ul>", "* <br />\n"},
	{"34", "<p>中<wbr>文</p>", "中文\n"},
	{"33", "<ol data-tight=\"true\"><li data-marker=\"1.\">foo</li></ul>", "1. foo\n"},
	{"32", "<ul data-tight=\"true\"><li data-marker=\"*\">foo<wbr></li></ul>", "* foo\n"},
	{"31", "<ul><li data-marker=\"*\">foo<ul><li data-marker=\"*\">bar</li></ul></li></ul>", "* foo\n  * bar\n"},
	{"30", "<ul><li data-marker=\"*\">foo</li><li data-marker=\"*\"><ul><li data-marker=\"*\"><br /></li></ul></li></ul>", "* foo\n* * <br />\n"},
	{"29", "<p><s>del</s></p>", "~~del~~\n"},
	{"29", "<p>[]()</p>", "[]()\n"},
	{"28", ":octocat:", ":octocat:\n"},
	{"27", "<table><thead><tr><th>abc</th><th>def</th></tr></thead></table>\n", "|abc|def|\n|---|---|\n"},
	{"26", "<p><del data-marker=\"~~\">Hi</del> Hello, world!</p>", "~~Hi~~ Hello, world!\n"},
	{"25", "<p><del data-marker=\"~\">Hi</del> Hello, world!</p>", "~Hi~ Hello, world!\n"},
	{"24", "<ul><li data-marker=\"*\" class=\"vditor-task\"><input checked=\"\" disabled=\"\" type=\"checkbox\" /> foo<wbr></li></ul>", "* [X] foo\n"},
	{"23", "<ul><li data-marker=\"*\" class=\"vditor-task\"><input disabled=\"\" type=\"checkbox\" /> foo<wbr></li></ul>", "* [ ] foo\n"},
	{"22", "><wbr>", ">\n"},
	{"21", "<p>> foo<wbr></p>", "> foo\n"},
	{"20", "<p>foo</p><p><wbr><br></p>", "foo\n\n<br />\n"},
	{"19", "<ul><li data-marker=\"*\">foo</li></ul><div><wbr><br></div>", "* foo\n\n<br />\n"},
	{"18", "<p><em data-marker=\"*\">foo<wbr></em></p>", "*foo*\n"},
	{"17", "foo bar", "foo bar\n"},
	{"16", "<p><em><strong>foo</strong></em></p>", "***foo***\n"},
	{"15", "<p><strong data-marker=\"__\">foo</strong></p>", "__foo__\n"},
	{"14", "<p><strong data-marker=\"**\">foo</strong></p>", "**foo**\n"},
	{"13", "<h2>foo</h2><p>para<em>em</em></p>", "## foo\n\npara*em*\n"},
	{"12", "<a href=\"/bar\" title=\"baz\">foo</a>", "[foo](/bar \"baz\")\n"},
	{"11", "<img src=\"/bar\" alt=\"foo\" />", "![foo](/bar)\n"},
	{"10", "<img src=\"/bar\" />", "![](/bar)\n"},
	{"9", "<a href=\"/bar\">foo</a>", "[foo](/bar)\n"},
	{"8", "foo<br />bar", "foo<br />bar\n"},
	{"7", "<code>foo</code>", "`foo`\n"},
	{"6", "<pre><code>foo</code></pre>", "```\nfoo\n```\n"},
	{"5", "<ul><li data-marker=\"*\">foo</li></ul>", "* foo\n"},
	{"4", "<blockquote>foo</blockquote>", "> foo\n"},
	{"3", "<h2>foo</h2>", "## foo\n"},
	{"2", "<p><strong><em>foo</em></strong></p>", "***foo***\n"},
	{"1", "<p><strong>foo</strong></p>", "**foo**\n"},
	{"0", "<p>foo</p>", "foo\n"},
}

func TestVditorDOM2Md(t *testing.T) {
	luteEngine := lute.New()

	for _, test := range vditorDOM2MdTests {
		md, err := luteEngine.VditorDOM2Md(test.from)
		if nil != err {
			t.Fatalf("test case [%s] unexpected: %s", test.name, err)
		}

		if test.to != md {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal html\n\t%q", test.name, test.to, md, test.from)
		}
	}
}

var vditorRendererTests = []*parseTest{

	{"27", "<p>foo</p>\n<div class=\"vditor-wysiwyg__block\" data-type=\"html\"><textarea class=\"vditor-reset\"><audio controls=\"controls\" src=\"http://localhost:8080/upload/file/2019/11/1440573175609-96444c00.mp3\"></audio></textarea></div>\n<p>bar</p>", "<p>foo</p>\n<div class=\"vditor-wysiwyg__block\" data-type=\"html\"><textarea class=\"vditor-reset\">&lt;audio controls=&#34;controls&#34; src=&#34;http://localhost:8080/upload/file/2019/11/1440573175609-96444c00.mp3&#34;&gt;&lt;/audio&gt;</textarea></div>\n<p>bar</p>"},
	{"27", "<p><wbr></p>", "<p>\n<wbr></p>"},
	{"26", "<p>![alt](src \"title\")</p>", "<p><img src=\"src\" alt=\"alt\" title=\"title\" /></p>"},
	{"25", "<pre><code class=\"language-java\"><wbr>\n</code></pre>", "<div class=\"vditor-wysiwyg__block\" data-type=\"pre\"><pre><code><wbr>\n</code></pre></div>"},
	{"24", "<ul data-tight=\"true\"><li data-marker=\"*\"><wbr><br></li></ul>", "<ul data-tight=\"true\"><li data-marker=\"*\"><wbr><br /></li></ul>"},
	{"23", "<ol><li data-marker=\"1.\">foo</li></ol>", "<ol data-tight=\"true\"><li data-marker=\"1.\">foo</li></ol>"},
	{"22", "<ul><li data-marker=\"*\">foo</li><li data-marker=\"*\"><ul><li data-marker=\"*\"><wbr><br /></li></ul></li></ul>", "<ul data-tight=\"true\"><li data-marker=\"*\">foo</li><li data-marker=\"*\"><ul data-tight=\"true\"><li data-marker=\"*\"><wbr><br /></li></ul></li></ul>"},
	{"21", "<p>[foo](/bar \"baz\")</p>", "<p><a href=\"/bar\" title=\"baz\">foo</a></p>"},
	{"20", "<p>[foo](/bar)</p>", "<p><a href=\"/bar\">foo</a></p>"},
	{"19", "<p>[foo]()</p>", "<p>[foo]()</p>"},
	{"18", "<p>[](/bar)</p>", "<p>[](/bar)</p>"},
	{"17", "<p>[]()</p>", "<p>[]()</p>"},
	{"16", "<p>[](</p>", "<p>[](</p>"},
	{"15", "<p><img alt=\"octocat\" class=\"emoji\" src=\"https://cdn.jsdelivr.net/npm/vditor/dist/images/emoji/octocat.png\" title=\"octocat\" /></p>", "<p><img alt=\"octocat\" class=\"emoji\" src=\"https://cdn.jsdelivr.net/npm/vditor/dist/images/emoji/octocat.png\" title=\"octocat\" /></p>"},
	{"14", ":octocat:", "<p><img alt=\"octocat\" class=\"emoji\" src=\"https://cdn.jsdelivr.net/npm/vditor/dist/images/emoji/octocat.png\" title=\"octocat\" /></p>"},
	{"13", "<div class=\"vditor-block\"><table><thead><tr><th>abc</th><th>def</th></tr></thead></table></div>", "<div class=\"vditor-block\"><table><thead><tr><th>abc</th><th>def</th></tr></thead></table></div>\n"},
	{"12", "<p><s data-marker=\"~~\">Hi</s> Hello, world!</p>", "<p><s data-marker=\"~~\">Hi</s> Hello, world!</p>"},
	{"11", "<p><del data-marker=\"~\">Hi</del> Hello, world!</p>", "<p><s data-marker=\"~\">Hi</s> Hello, world!</p>"},
	{"10", "<ul><li data-marker=\"*\" class=\"vditor-task\"><input checked=\"\" type=\"checkbox\" /> foo<wbr></li></ul>", "<ul data-tight=\"true\"><li data-marker=\"*\" class=\"vditor-task\"><input checked=\"\" type=\"checkbox\" /> foo<wbr></li></ul>"},
	{"9", "<ul><li data-marker=\"*\" class=\"vditor-task\"><input type=\"checkbox\" /> foo<wbr></li></ul>", "<ul data-tight=\"true\"><li data-marker=\"*\" class=\"vditor-task\"><input type=\"checkbox\" /> foo<wbr></li></ul>"},
	{"8", "> <wbr>", "<p>> <wbr></p>"},
	{"7", "><wbr>", "<p>><wbr></p>"},
	{"6", "<p>> foo<wbr></p>", "<blockquote><p>foo<wbr></p></blockquote>"},
	{"5", "<p>foo</p><p><wbr><br></p>", "<p>foo</p><p>\n<wbr><br /></p>"},
	{"4", "<ul data-tight=\"true\"><li data-marker=\"*\">foo<wbr></li></ul>", "<ul data-tight=\"true\"><li data-marker=\"*\">foo<wbr></li></ul>"},
	{"3", "<p><em data-marker=\"*\">foo<wbr></em></p>", "<p><em data-marker=\"*\">foo<wbr></em></p>"},
	{"2", "<p>foo<wbr></p>", "<p>foo<wbr></p>"},
	{"1", "<p><strong data-marker=\"**\">foo</strong></p>", "<p><strong data-marker=\"**\">foo</strong></p>"},
	{"0", "<p>foo</p>", "<p>foo</p>"},
}

func TestVditorRenderer(t *testing.T) {
	luteEngine := lute.New()

	for _, test := range vditorRendererTests {
		html, err := luteEngine.SpinVditorDOM(test.from)
		if nil != err {
			t.Fatalf("unexpected: %s", err)
		}

		if test.to != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal html\n\t%q", test.name, test.to, html, test.from)
		}
	}
}
