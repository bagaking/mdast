package mdast

import (
	"context"
	"fmt"
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
		rowContent, err := tableRowContentToMarkdown(ctx, row)
		if err != nil {
			return "", err
		}
		result.WriteString(rowContent + "\n")
		if i == 0 {
			header, ok := row.(*Node)
			if !ok {
				return "", fmt.Errorf("cannot render table separator for header node type: %s", row.GetType())
			}
			result.WriteString("|")

			for j := range header.TableChildren {
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
	content, err := phrasingChildrenToMarkdown(ctx, n)
	if err != nil {
		return "", err
	}
	return escapeTableCellContent(content), nil
}

func escapeTableCellContent(content string) string {
	return escapeTableCellPipes(normalizeTableCellLineBreaks(content))
}

func normalizeTableCellLineBreaks(content string) string {
	if !strings.ContainsAny(content, "\r\n") {
		return content
	}
	return strings.NewReplacer(
		"\r\n", " ",
		"\r", " ",
		"\n", " ",
	).Replace(content)
}

func escapeTableCellPipes(content string) string {
	if !strings.Contains(content, "|") {
		return content
	}

	var result strings.Builder
	result.Grow(len(content))
	for i := 0; i < len(content); i++ {
		if content[i] == '|' && !isEscapedMarkdownByte(content, i) {
			result.WriteByte('\\')
		}
		result.WriteByte(content[i])
	}
	return result.String()
}

func isEscapedMarkdownByte(content string, index int) bool {
	backslashes := 0
	for i := index - 1; i >= 0 && content[i] == '\\'; i-- {
		backslashes++
	}
	return backslashes%2 == 1
}

func tableRowContentToMarkdown(ctx context.Context, child TableContent) (string, error) {
	if childNode, ok := child.(*Node); ok {
		return TableRowToMarkdown(ctx, childNode)
	}
	return child.ToMarkdown(ctx)
}

func tableCellContentToMarkdown(ctx context.Context, child TableContent) (string, error) {
	if childNode, ok := child.(*Node); ok {
		return TableCellToMarkdown(ctx, childNode)
	}
	return child.ToMarkdown(ctx)
}

func tableChildrenToMarkdownSlice(ctx context.Context, n *Node) ([]string, error) {
	result := make([]string, len(n.TableChildren))
	for i, child := range n.TableChildren {
		content, err := tableCellContentToMarkdown(ctx, child)
		if err != nil {
			return nil, err
		}
		result[i] = content
	}
	return result, nil
}
