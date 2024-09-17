package mdast

import (
	"context"
	"fmt"
	"strings"
)

// FlowToMarkdown 将流式内容转换为 Markdown
func FlowToMarkdown(ctx context.Context, n *Node) (string, error) {
	switch n.Type {
	case NodeParagraph:
		return paragraphToMarkdown(ctx, n)
	case NodeHeading:
		return headingToMarkdown(ctx, n)
	case NodeBlockquote:
		return blockquoteToMarkdown(ctx, n)
	case NodeCode:
		return codeToMarkdown(ctx, n)
	case NodeThematicBreak:
		return "---\n\n", nil
	case NodeHTML:
		return htmlToMarkdown(ctx, n)
	case NodeYaml:
		return yamlToMarkdown(ctx, n)
	case NodeDefinition:
		return definitionToMarkdown(ctx, n)
	case NodeFootnoteDefinition:
		return footnoteDefinitionToMarkdown(ctx, n)
	case NodeList:
		return ListToMarkdown(ctx, n)
	case NodeFootnote:
		return footnoteToMarkdown(ctx, n)
	default:
		return "", fmt.Errorf("unknown flow node type: %s", n.Type)
	}
}

func paragraphToMarkdown(ctx context.Context, n *Node) (string, error) {
	content, err := phrasingChildrenToMarkdown(ctx, n)
	if err != nil {
		return "", err
	}
	return content + "\n\n", nil
}

func headingToMarkdown(ctx context.Context, n *Node) (string, error) {
	level, _ := n.Data.GetInt(NDK_Depth)
	content, err := phrasingChildrenToMarkdown(ctx, n)
	if err != nil {
		return "", err
	}
	return strings.Repeat("#", level) + " " + content + "\n\n", nil
}

func blockquoteToMarkdown(ctx context.Context, n *Node) (string, error) {
	content, err := flowChildrenToMarkdown(ctx, n)
	if err != nil {
		return "", err
	}
	lines := strings.Split(strings.TrimRight(content, "\n"), "\n")
	for i, line := range lines {
		if line != "" {
			lines[i] = "> " + line
		} else {
			lines[i] = ">"
		}
	}
	return strings.Join(lines, "\n") + "\n\n", nil
}

func codeToMarkdown(ctx context.Context, n *Node) (string, error) {
	lang, _ := n.Data.GetString(NDK_Lang)
	meta, _ := n.Data.GetString(NDK_Meta)
	if meta != "" {
		lang += " " + meta
	}
	return "```" + lang + "\n" + n.Value + "\n```\n\n", nil
}

func htmlToMarkdown(ctx context.Context, n *Node) (string, error) {
	if IsInlineHTML(n) {
		return n.Value, nil // 内联 HTML
	}
	return n.Value + "\n\n", nil // 块级 HTML
}

func yamlToMarkdown(ctx context.Context, n *Node) (string, error) {
	return "---\n" + n.Value + "\n---\n\n", nil
}

func definitionToMarkdown(ctx context.Context, n *Node) (string, error) {
	identifier, ok := n.Data.GetString(NDK_Identifier)
	if !ok {
		return "", fmt.Errorf("missing or invalid identifier for definition")
	}
	url, ok := n.Data.GetString(NDK_URL)
	if !ok {
		return "", fmt.Errorf("missing or invalid URL for definition")
	}
	title, _ := n.Data.GetString(NDK_Title)
	if title != "" {
		return fmt.Sprintf("[%s]: %s \"%s\"\n", identifier, url, title), nil
	}
	return fmt.Sprintf("[%s]: %s\n", identifier, url), nil
}

func footnoteDefinitionToMarkdown(ctx context.Context, n *Node) (string, error) {
	identifier, ok := n.Data.GetString(NDK_Identifier)
	if !ok {
		return "", fmt.Errorf("missing or invalid identifier for footnote definition")
	}
	content, err := flowChildrenToMarkdown(ctx, n)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("[^%s]: %s\n\n", identifier, strings.TrimSpace(content)), nil
}

func footnoteToMarkdown(ctx context.Context, n *Node) (string, error) {
	content, err := phrasingChildrenToMarkdown(ctx, n)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("[^%s]", content), nil
}
