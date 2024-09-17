package mdast

import (
	"context"
	"fmt"
	"strings"
)

// InlineToMarkdown 将内联元素转换为 Markdown
func InlineToMarkdown(ctx context.Context, n *Node) (string, error) {
	switch n.Type {
	case NodeText:
		return n.Value, nil
	case NodeEmphasis:
		content, err := phrasingChildrenToMarkdown(ctx, n)
		if err != nil {
			return "", err
		}
		return "*" + content + "*", nil
	case NodeStrong:
		content, err := phrasingChildrenToMarkdown(ctx, n)
		if err != nil {
			return "", err
		}
		return "**" + content + "**", nil
	case NodeDelete:
		content, err := phrasingChildrenToMarkdown(ctx, n)
		if err != nil {
			return "", err
		}
		return "~~" + content + "~~", nil
	case NodeLink:
		return linkToMarkdown(ctx, n)
	case NodeImage:
		return imageToMarkdown(ctx, n)
	case NodeInlineCode:
		return "`" + n.Value + "`", nil
	case NodeBreak:
		return "\n", nil
	case NodeLinkReference:
		return linkReferenceToMarkdown(ctx, n)
	case NodeImageReference:
		return imageReferenceToMarkdown(ctx, n)
	case NodeFootnoteReference:
		identifier, ok := n.Data.GetString(NDK_Identifier)
		if !ok {
			return "", fmt.Errorf("missing or invalid identifier for footnote reference")
		}
		return "[^" + identifier + "]", nil
	case NodeFootnote:
		return footnoteToMarkdown(ctx, n)
	default:
		return "", fmt.Errorf("unknown inline node type: %s", n.Type)
	}
}

func linkToMarkdown(ctx context.Context, n *Node) (string, error) {
	text, err := phrasingChildrenToMarkdown(ctx, n)
	if err != nil {
		return "", err
	}

	url, ok := n.Data.GetString(NDK_URL)
	if !ok {
		return "", fmt.Errorf("missing or invalid URL for link")
	}

	title, _ := n.Data.GetString(NDK_Title)
	// 注意这里我们不检查 ok，因为 title 是可选的

	if title != "" {
		return fmt.Sprintf("[%s](%s \"%s\")", text, url, title), nil
	}
	return fmt.Sprintf("[%s](%s)", text, url), nil
}

func imageToMarkdown(ctx context.Context, n *Node) (string, error) {
	alt, ok := n.Data.GetString(NDK_Alt)
	if !ok {
		return "", fmt.Errorf("missing or invalid alt text for image")
	}

	url, ok := n.Data.GetString(NDK_URL)
	if !ok {
		return "", fmt.Errorf("missing or invalid URL for image")
	}

	title, _ := n.Data.GetString(NDK_Title)
	if title != "" {
		return fmt.Sprintf("![%s](%s \"%s\")", alt, url, title), nil
	}
	return fmt.Sprintf("![%s](%s)", alt, url), nil
}

func linkReferenceToMarkdown(ctx context.Context, n *Node) (string, error) {
	identifier, ok := n.Data.GetString(NDK_Identifier)
	if !ok {
		return "", fmt.Errorf("missing or invalid identifier for link reference")
	}

	referenceType, ok := n.Data.GetReferenceType(NDK_ReferenceType)
	if !ok {
		return "", fmt.Errorf("missing or invalid reference type for link reference")
	}

	label, err := phrasingChildrenToMarkdown(ctx, n)
	if err != nil {
		return "", err
	}

	switch referenceType {
	case ReferenceShortcut:
		return fmt.Sprintf("[%s]", label), nil
	case ReferenceCollapsed:
		return fmt.Sprintf("[%s][]", label), nil
	default: // full
		return fmt.Sprintf("[%s][%s]", label, identifier), nil
	}
}

func imageReferenceToMarkdown(ctx context.Context, n *Node) (string, error) {
	identifier, ok := n.Data.GetString(NDK_Identifier)
	if !ok {
		return "", fmt.Errorf("missing or invalid identifier for image reference")
	}

	alt, ok := n.Data.GetString(NDK_Alt)
	if !ok {
		return "", fmt.Errorf("missing or invalid alt text for image reference")
	}

	referenceType, ok := n.Data.GetReferenceType(NDK_ReferenceType)
	if !ok {
		return "", fmt.Errorf("missing or invalid reference type for image reference")
	}

	switch referenceType {
	case ReferenceShortcut:
		return fmt.Sprintf("![%s]", alt), nil
	case ReferenceCollapsed:
		return fmt.Sprintf("![%s][]", alt), nil
	default: // full
		return fmt.Sprintf("![%s][%s]", alt, identifier), nil
	}
}

// IsInlineHTML 判断给定的 HTML 节点是否为内联元素
func IsInlineHTML(n *Node) bool {
	// 1. 检查父节点类型
	// 如果父节点是段落或其他内联容器，则该 HTML 节点被视为内联元素
	if n.parent != nil && (n.parent.GetType() == NodeParagraph || n.parent.GetType().IsInline()) {
		return true
	}

	// 2. 检查 HTML 内容
	// 常见的块级 HTML 标签
	blockTags := []string{"<div", "<p", "<blockquote", "<pre", "<table", "<ul", "<ol", "<li", "<hr", "<h1", "<h2", "<h3", "<h4", "<h5", "<h6"}
	lowerValue := strings.ToLower(n.Value)
	for _, tag := range blockTags {
		if strings.HasPrefix(lowerValue, tag) {
			return false // 如果是块级标签，则返回 false
		}
	}

	// 3. 检查 HTML 内容是否包含换行符， 如果 HTML 内容中包含换行符，则通常被视为块级元素，返回 false
	return !strings.Contains(n.Value, "\n")
}
