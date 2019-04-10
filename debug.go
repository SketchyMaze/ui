package ui

import "strings"

// WidgetTree returns a string representing the tree of widgets starting
// at a given widget.
func WidgetTree(root Widget) []string {
	var crawl func(int, Widget) []string
	crawl = func(depth int, node Widget) []string {
		var (
			prefix = strings.Repeat("  ", depth)
			lines  = []string{prefix + node.ID()}
		)

		for _, child := range node.Children() {
			lines = append(lines, crawl(depth+1, child)...)
		}

		return lines
	}

	return crawl(0, root)
}
