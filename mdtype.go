package mdast

import "context"

// Position 表示节点在源文件中的位置
type Position struct {
	Start  Point
	End    Point
	Indent []int
}

// Point 表示位置的具体点
type Point struct {
	Line   int
	Column int
	Offset int
}

// NodeType 表示 Markdown AST 节点的类型
type NodeType string

// 定义所有的节点类型
const (
	NodeRoot               NodeType = "root"
	NodeParagraph          NodeType = "paragraph"
	NodeHeading            NodeType = "heading"
	NodeText               NodeType = "text"
	NodeEmphasis           NodeType = "emphasis"
	NodeStrong             NodeType = "strong"
	NodeDelete             NodeType = "delete"
	NodeLink               NodeType = "link"
	NodeImage              NodeType = "image"
	NodeCode               NodeType = "code"
	NodeInlineCode         NodeType = "inlineCode"
	NodeBlockquote         NodeType = "blockquote"
	NodeList               NodeType = "list"
	NodeListItem           NodeType = "listItem"
	NodeTable              NodeType = "table"
	NodeTableRow           NodeType = "tableRow"
	NodeTableCell          NodeType = "tableCell"
	NodeThematicBreak      NodeType = "thematicBreak"
	NodeBreak              NodeType = "break"
	NodeHTML               NodeType = "html"
	NodeDefinition         NodeType = "definition"
	NodeImageReference     NodeType = "imageReference"
	NodeLinkReference      NodeType = "linkReference"
	NodeFootnote           NodeType = "footnote"
	NodeFootnoteReference  NodeType = "footnoteReference"
	NodeFootnoteDefinition NodeType = "footnoteDefinition"
	NodeYaml               NodeType = "yaml"
)

// AlignType 表示表格列的对齐方式
type AlignType string

const (
	AlignNone   AlignType = ""
	AlignLeft   AlignType = "left"
	AlignRight  AlignType = "right"
	AlignCenter AlignType = "center"
)

// ReferenceType 表示引用的类型
type ReferenceType string

const (
	ReferenceShortcut  ReferenceType = "shortcut"
	ReferenceCollapsed ReferenceType = "collapsed"
	ReferenceFull      ReferenceType = "full"
)

// Node 结构体现在包含特定类型的子节点字段
type Node struct {
	Type             NodeType
	Value            string
	FlowChildren     []FlowContent
	PhrasingChildren []PhrasingContent
	ListChildren     []ListContent
	TableChildren    []TableContent
	Data             DataTable
	parent           *Node
}

// IsBlock 检查节点是否为块级元素
func (nt NodeType) IsBlock() bool {
	switch nt {
	case NodeRoot, NodeParagraph, NodeHeading, NodeBlockquote, NodeList, NodeListItem, NodeTable, NodeTableRow, NodeThematicBreak, NodeCode, NodeHTML, NodeYaml, NodeFootnoteDefinition:
		return true
	default:
		return false
	}
}

// IsInline 检查节点是否为内联元素
func (nt NodeType) IsInline() bool {
	switch nt {
	case NodeText, NodeEmphasis, NodeStrong, NodeDelete, NodeLink, NodeImage, NodeInlineCode, NodeBreak, NodeFootnoteReference:
		return true
	default:
		return false
	}
}

// IsListContent 检查节点是否为列表内容
func (nt NodeType) IsListContent() bool {
	return nt == NodeListItem
}

// IsTableContent 检查节点是否为表格内容
func (nt NodeType) IsTableContent() bool {
	return nt == NodeTableRow
}

// IsRowContent 检查节点是否为行内容
func (nt NodeType) IsRowContent() bool {
	return nt == NodeTableCell
}

// Content 接口定义了所有内容节点的基本方法
type Content interface {
	ToMarkdown(ctx context.Context) (string, error)
	GetType() NodeType
}

// FlowContent 接口定义了流式内容节点
type FlowContent interface {
	Content
	IsFlow()
}

// PhrasingContent 接口定义了短语内容节点
type PhrasingContent interface {
	Content
	IsPhrasing()
}

// ListContent 接口定义了列表内容节点
type ListContent interface {
	Content
	IsList()
}

// TableContent 接口定义了表格内容节点
type TableContent interface {
	Content
	IsTable()
}
