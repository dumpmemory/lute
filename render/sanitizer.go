// Lute - 一款结构化的 Markdown 引擎，支持 Go 和 JavaScript
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package render

import (
	"bytes"
	"regexp"
	"strings"

	"github.com/88250/lute/html"
	"github.com/88250/lute/util"
	"github.com/microcosm-cc/bluemonday"
)

func SanitizeLinkDestBytes(src []byte) []byte {
	return []byte(SanitizeLinkDest(string(src)))
}

func SanitizeLinkDest(src string) string {
	ret := strings.ReplaceAll(src, "\"", "__@QUOTE@__")
	ret = strings.ReplaceAll(ret, " ", "__@SPACE@__")
	ret = strings.ReplaceAll(ret, "#", "__@HASH@__")
	ret = strings.ReplaceAll(ret, "&", "__@AMP@__")
	ret = strings.ReplaceAll(ret, "|", "__@PIPE@__")
	ret = "<a href=\"" + ret + "\"></a>"
	sanitizer := newSanitizer()
	ret = sanitizer.Sanitize(ret)
	ret = util.TagAttrVal(ret, "href")
	ret = strings.ReplaceAll(ret, "__@QUOTE@__", "\"")
	ret = strings.ReplaceAll(ret, "__@SPACE@__", " ")
	ret = strings.ReplaceAll(ret, "__@HASH@__", "#")
	ret = strings.ReplaceAll(ret, "__@AMP@__", "&")
	ret = strings.ReplaceAll(ret, "__@PIPE@__", "|")
	ret = strings.TrimSpace(ret)
	if strings.HasPrefix(ret, "javascript:") {
		return ""
	}
	return ret
}

func Sanitize(str string) string {
	return string(sanitize([]byte(str)))
}

func sanitize(tokens []byte) []byte {
	ret := newSanitizer().SanitizeBytes(tokens)
	node, err := html.Parse(bytes.NewBuffer(ret))
	if nil != err {
		return ret
	}
	node = node.FirstChild // html 节点
	if nil != node.FirstChild && nil != node.FirstChild.FirstChild {
		if "meta" == node.FirstChild.FirstChild.Data {
			return ret
		}
	}

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode {
			for i := len(n.Attr) - 1; i >= 0; i-- {
				attr := n.Attr[i]
				if attr.Key == "href" || attr.Key == "src" {
					val := strings.ToLower(strings.TrimSpace(attr.Val))
					if strings.HasPrefix(val, "javascript:") || strings.HasPrefix(val, "data:") {
						n.Attr = append(n.Attr[:i], n.Attr[i+1:]...)
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(node)

	ret = util.DomHTML(node.LastChild)
	ret = []byte(strings.TrimSuffix(strings.TrimPrefix(string(ret), "<body>"), "</body>"))
	return ret
}

func newSanitizer() *bluemonday.Policy {
	ret := bluemonday.NewPolicy()
	ret.AllowStandardAttributes()
	ret.AllowDataAttributes()
	ret.AllowStandardURLs()
	ret.AllowImages()
	ret.AllowLists()
	ret.AllowStyling()
	ret.AllowTables()

	ret.AllowAttrs("href", "target").OnElements("a")
	ret.AllowAttrs("align").OnElements("p", "div")
	ret.AllowAttrs("src", "scrolling", "border", "frameborder", "framespacing", "allowfullscreen", "data-subtype", "updated", "style").OnElements("iframe")
	ret.AllowAttrs("content").OnElements("meta")
	ret.AllowAttrs("loading", "title", "style").OnElements("img")
	ret.AllowAttrs("controls", "autoplay", "loop", "muted", "src", "style").OnElements("video", "audio")
	ret.AllowAttrs("type", "allowscriptaccess").OnElements("embed")
	ret.AllowAttrs("open").OnElements("details")

	ret.AllowURLSchemesMatching(regexp.MustCompile("(?i)^[a-z][a-z0-9+.-]*$"))

	ret.RequireParseableURLs(true)
	ret.RequireNoFollowOnLinks(false)

	ret.AllowElements("details", "summary")
	return ret
}
