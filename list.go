// Lute - A structured markdown engine.
// Copyright (C) 2019-present, b3log.org
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package lute

import (
	"fmt"
)

type ListType int

const (
	ListTypeBullet  = 0
	ListTypeOrdered = 1
)

type List struct {
	*Node
	int
	items
	t *Tree

	ListType ListType
	Start    int
	Tight    bool

	IndentSpaces int // #4 Indentation https://spec.commonmark.org/0.29/#list-items
	Marker       string
	WNSpaces     int // W + N https://spec.commonmark.org/0.29/#list-items
}

func (n *List) String() string {
	return fmt.Sprintf("%s", n.Children())
}

func (n *List) HTML() string {
	content := html(n.Children())

	if NodeListItem == n.Parent.NodeType {
		return fmt.Sprintf("\n<ul>\n%s</ul>\n", content)
	}

	return fmt.Sprintf("<ul>\n%s</ul>\n", content)
}

func newList(indentSpaces int, marker string, wnSpaces int, t *Tree, token *item) (ret *Node, list *List) {
	ret = &Node{NodeType: NodeList, Parent: t.context.CurNode}
	list = &List{
		ret, token.pos, items{}, t,
		ListTypeBullet,
		1,
		false,
		indentSpaces,
		marker,
		wnSpaces,
	}
	t.context.CurNode = ret

	return
}

func (t *Tree) parseList(line items) (ret *Node) {
	spaces, tabs, tokens, firstNonWhitespace := t.nonWhitespace(line)
	marker := firstNonWhitespace
	indentSpaces := spaces + tabs*4
	line = line[len(tokens):]
	spaces, tabs, _, firstNonWhitespace = t.nonWhitespace(line)
	w := len(marker.val)
	n := spaces + tabs*4
	wnSpaces := w + n
	oldContextIndentSpaces := t.context.IndentSpaces
	t.context.IndentSpaces = indentSpaces + wnSpaces
	if 4 <= n { // rule 2 in https://spec.commonmark.org/0.29/#list-items
		line = indentOffset(line, w+1, t)
		t.context.IndentSpaces = 2
	} else {
		line = indentOffset(line, indentSpaces+wnSpaces, t)
	}
	ret, list := newList(indentSpaces, marker.val, wnSpaces, t, marker)
	tight := false

	defer func() { t.context.IndentSpaces = oldContextIndentSpaces }()

	for {
		n, i := t.parseListItem(line)
		if nil == n {
			break
		}
		list.Append(n)

		if i.Tight {
			tight = true
		}

		line = t.nextLine()
		if line.isEOF() {
			break
		}

		t.skipWhitespaces(line)
		if marker.val != line[0].val {
			// TODO: 考虑有序列表序号递增
			t.backupLine(line)

			break
		} else {
			line = line[len(marker.val):]
		}
	}

	list.Tight = tight

	return
}

// https://spec.commonmark.org/0.29/#lists
func (t *Tree) isList(line items) bool {
	if 2 > len(line) { // at least marker and newline
		return false
	}

	line = t.skipWhitespaces(line)
	firstNonWhitespace := line[0]
	if "*" != firstNonWhitespace.val && "-" != firstNonWhitespace.val && "+" != firstNonWhitespace.val {
		// TODO: 有序列表判断
		return false
	}

	if itemSpace != line[1].typ && itemTab != line[1].typ && itemNewline != line[1].typ {
		return false
	}

	return true
}
