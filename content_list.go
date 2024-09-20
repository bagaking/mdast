package mdast

import (
	"context"
	"fmt"
	"strings"
)

// ListToMarkdown 将列表内容转换为 Markdown
func ListToMarkdown(ctx context.Context, n *Node) (string, error) {
	ordered, _ := n.Data.GetBool(NDK_Ordered)
	spread, _ := n.Data.GetBool(NDK_Spread)
	start, ok := n.Data.GetInt(NDK_Start)
	if !ok || start < 1 {
		start = 1
	}

	lines := make([]string, 0, len(n.ListChildren))
	for i, child := range n.ListChildren {
		if child.GetType() != NodeListItem {
			return "", fmt.Errorf("unexpected node type in list: %s", child.GetType())
		}
		itemContent, err := listContentToMarkdown(ctx, child, start+i, ordered)
		if err != nil {
			return "", fmt.Errorf("error processing list item: %w", err)
		}
		lines = append(lines, itemContent)
	}

	separator := "\n"
	if spread {
		separator = "\n\n"
	}
	return strings.Join(lines, separator) + "\n\n", nil
}

func listContentToMarkdown(ctx context.Context, child ListContent, index int, ordered bool) (string, error) {
	if childNode, ok := child.(*Node); ok {
		return listItemToMarkdown(ctx, childNode, index, ordered)
	}
	return child.ToMarkdown(ctx)
}

func listItemToMarkdown(ctx context.Context, n *Node, index int, ordered bool) (string, error) {
	var prefix string
	if ordered {
		prefix = fmt.Sprintf("%d. ", index)
	} else {
		prefix = "- "
	}

	parts := make([]string, 0, len(n.FlowChildren)+1)
	if len(n.PhrasingChildren) > 0 {
		content, err := phrasingChildrenToMarkdown(ctx, n)
		if err != nil {
			return "", err
		}
		parts = append(parts, content)
	}

	for _, child := range n.FlowChildren {
		childContent, err := flowChildToMarkdown(ctx, child)
		if err != nil {
			return "", fmt.Errorf("error processing list item child: %w", err)
		}
		parts = append(parts, strings.TrimSpace(childContent))
	}

	separator := "\n"
	if spread, _ := n.Data.GetBool(NDK_Spread); spread {
		separator = "\n\n"
	}
	content := strings.Join(parts, separator)
	if content == "" {
		return strings.TrimSpace(prefix), nil
	}

	return prefix + indentContinuation(content, len(prefix)), nil
}

func indentContinuation(content string, width int) string {
	lines := strings.Split(content, "\n")
	for i := 1; i < len(lines); i++ {
		if lines[i] != "" {
			lines[i] = strings.Repeat(" ", width) + lines[i]
		}
	}
	return strings.Join(lines, "\n")
}
