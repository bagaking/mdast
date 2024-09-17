package mdast

import (
	"context"
	"fmt"
	"strings"
)

// ListToMarkdown 将列表内容转换为 Markdown
func ListToMarkdown(ctx context.Context, n *Node) (string, error) {
	var result strings.Builder
	ordered, ok := n.Data.GetBool(NDK_Ordered)
	if !ok {
		return "", fmt.Errorf("missing required 'ordered' property for list")
	}
	// spread is optional, thus we don't need to check it
	spread, _ := n.Data.GetBool(NDK_Spread)

	for i, child := range n.ListChildren {
		if child.GetType() != NodeListItem {
			return "", fmt.Errorf("unexpected node type in list: %s", child.GetType())
		}
		itemContent, err := listItemToMarkdown(ctx, child.(*Node), i+1, ordered, spread)
		if err != nil {
			return "", fmt.Errorf("error processing list item: %w", err)
		}
		result.WriteString(itemContent)
		if i < len(n.ListChildren)-1 {
			if spread {
				result.WriteString("\n\n")
			} else {
				result.WriteString("\n")
			}
		}
	}
	return result.String() + "\n", nil
}

func listItemToMarkdown(ctx context.Context, n *Node, index int, ordered bool, spread bool) (string, error) {
	indent := "   "
	var prefix string
	if ordered {
		prefix = fmt.Sprintf("%d. ", index)
	} else {
		prefix = "- "
	}

	var result strings.Builder
	result.WriteString(prefix)

	for i, child := range n.FlowChildren {
		if i > 0 {
			if spread {
				result.WriteString("\n\n" + strings.Repeat(" ", len(prefix)))
			} else {
				result.WriteString("\n" + strings.Repeat(" ", len(prefix)))
			}
		}
		childContent, err := FlowToMarkdown(ctx, child.(*Node))
		if err != nil {
			return "", fmt.Errorf("error processing list item child: %w", err)
		}
		if child.GetType() == NodeList {
			childLines := strings.Split(strings.TrimRight(childContent, "\n"), "\n")
			for j, line := range childLines {
				if j > 0 {
					result.WriteString("\n" + indent + strings.Repeat(" ", len(prefix)))
				}
				result.WriteString(line)
			}
		} else {
			result.WriteString(strings.TrimSpace(childContent))
		}
	}

	return result.String(), nil
}
