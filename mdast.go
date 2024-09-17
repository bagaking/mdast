package mdast

import (
	"context"
	"fmt"
	"strings"
)

// NewNode 创建一个新的 Node
func NewNode(nodeType NodeType) *Node {
	return &Node{
		Type:             nodeType,
		FlowChildren:     []FlowContent{},
		PhrasingChildren: []PhrasingContent{},
		ListChildren:     []ListContent{},
		TableChildren:    []TableContent{},
		Data:             make(DataTable),
	}
}

// IsFlow 实现 FlowContent 接口
func (n *Node) IsFlow() {}

// IsPhrasing 实现 PhrasingContent 接口
func (n *Node) IsPhrasing() {}

// IsList 实现 ListContent 接口
func (n *Node) IsList() {}

// IsTable 实现 TableContent 接口
func (n *Node) IsTable() {}

var (
	_ FlowContent = &Node{}
)

// AddFlowChild 添加一个流式子节点
func (n *Node) AddFlowChild(child FlowContent) {
	n.FlowChildren = append(n.FlowChildren, child)
	if childNode, ok := child.(*Node); ok {
		childNode.parent = n
	}
}

// AddPhrasingChild 添加一个短语子节点
func (n *Node) AddPhrasingChild(child PhrasingContent) {
	n.PhrasingChildren = append(n.PhrasingChildren, child)
	if childNode, ok := child.(*Node); ok {
		childNode.parent = n
	}
}

// AddListChild 添加一个列表子节点
func (n *Node) AddListChild(child ListContent) {
	n.ListChildren = append(n.ListChildren, child)
	if childNode, ok := child.(*Node); ok {
		childNode.parent = n
	}
}

// AddTableChild 添加一个表格子节点
func (n *Node) AddTableChild(child TableContent) {
	n.TableChildren = append(n.TableChildren, child)
	if childNode, ok := child.(*Node); ok {
		childNode.parent = n
	}
}

// SetData 设置节点数据
func (n *Node) SetData(key DataKey, value any) {
	n.Data[key] = value
}

// GetData 获取节点数据
func (n *Node) GetData(key DataKey) (any, bool) {
	value, exists := n.Data[key]
	return value, exists
}

// GetType 实现 Content 接口
func (n *Node) GetType() NodeType {
	return n.Type
}

// ToMarkdown 将节点转换为 Markdown 文本
func (n *Node) ToMarkdown(ctx context.Context) (string, error) {
	switch n.Type {
	case NodeRoot:
		return flowChildrenToMarkdown(ctx, n)
	case NodeParagraph, NodeHeading, NodeBlockquote, NodeCode, NodeThematicBreak,
		NodeHTML, NodeYaml, NodeDefinition, NodeFootnoteDefinition:
		return FlowToMarkdown(ctx, n)
	case NodeList:
		return ListToMarkdown(ctx, n)
	case NodeTable:
		return TableToMarkdown(ctx, n)
	case NodeTableRow:
		return TableRowToMarkdown(ctx, n)
	case NodeTableCell:
		return TableCellToMarkdown(ctx, n)
	case NodeText, NodeEmphasis, NodeStrong, NodeDelete, NodeLink,
		NodeImage, NodeInlineCode, NodeBreak,
		NodeLinkReference, NodeImageReference,
		NodeFootnoteReference:
		return InlineToMarkdown(ctx, n)
	default:
		return "", fmt.Errorf("unknown node type: %s", n.Type)
	}
}

// phrasingChildrenToMarkdown 将子节点转换为 Markdown 文本
func phrasingChildrenToMarkdown(ctx context.Context, n *Node) (string, error) {
	var result strings.Builder
	for _, child := range n.PhrasingChildren {
		childContent, err := InlineToMarkdown(ctx, child.(*Node))
		if err != nil {
			return "", err
		}
		result.WriteString(childContent)
	}
	return result.String(), nil
}

// flowChildrenToMarkdown 将流式子节点转换为 Markdown 文本
func flowChildrenToMarkdown(ctx context.Context, n *Node) (string, error) {
	var result strings.Builder
	for _, child := range n.FlowChildren {
		childContent, err := FlowToMarkdown(ctx, child.(*Node))
		if err != nil {
			return "", err
		}
		result.WriteString(childContent)
	}
	return result.String(), nil
}
