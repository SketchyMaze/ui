package ui

import (
	"fmt"
	"strings"
)

// WidgetTree returns a string representing the tree of widgets starting
// at a given widget.
func WidgetTree(root Widget) []string {
	var crawl func(int, Widget) []string

	crawl = func(depth int, node Widget) []string {
		var (
			prefix    = strings.Repeat("  ", depth)
			size      = node.Size()
			width     = size.W
			height    = size.H
			fixedSize = node.FixedSize()

			lines = []string{
				fmt.Sprintf("%s%s   P:%s   S:%dx%d (fixedSize: %+v)",
					prefix, node.ID(), node.Point(), width, height, fixedSize,
				),
			}
		)

		for _, child := range node.Children() {
			lines = append(lines, crawl(depth+1, child)...)
		}

		return lines
	}

	return crawl(0, root)
}
