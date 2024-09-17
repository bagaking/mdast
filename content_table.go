package mdast

import (
	"context"
	"strings"
)

// TableToMarkdown 将表格内容转换为 Markdown
func TableToMarkdown(ctx context.Context, n *Node) (string, error) {
	var result strings.Builder
	alignments, ok := n.Data[NDK_Align].([]AlignType)
	if !ok {
		alignments = []AlignType{}
	}

	for i, row := range n.TableChildren {
		rowContent, err := TableRowToMarkdown(ctx, row.(*Node))
		if err != nil {
			return "", err
		}
		result.WriteString(rowContent + "\n")
		if i == 0 {
			result.WriteString("|")

			for j := range n.TableChildren[0].(*Node).TableChildren {
				align := AlignNone
				if j < len(alignments) {
					align = alignments[j]
				}
				switch align {
				case AlignLeft:
					result.WriteString(" :--- |")
				case AlignRight:
					result.WriteString(" ---: |")
				case AlignCenter:
					result.WriteString(" :---: |")
				default:
					result.WriteString(" --- |")
				}
			}
			result.WriteString("\n")
		}
	}
	result.WriteString("\n")
	return result.String(), nil
}

func TableRowToMarkdown(ctx context.Context, n *Node) (string, error) {
	cells, err := tableChildrenToMarkdownSlice(ctx, n)
	if err != nil {
		return "", err
	}
	return "| " + strings.Join(cells, " | ") + " |", nil
}

func TableCellToMarkdown(ctx context.Context, n *Node) (string, error) {
	return phrasingChildrenToMarkdown(ctx, n)
}

func tableChildrenToMarkdownSlice(ctx context.Context, n *Node) ([]string, error) {
	result := make([]string, len(n.TableChildren))
	for i, child := range n.TableChildren {
		content, err := TableCellToMarkdown(ctx, child.(*Node))
		if err != nil {
			return nil, err
		}
		result[i] = content
	}
	return result, nil
}
