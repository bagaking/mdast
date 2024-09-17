package mdast

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCase 定义了一个通用的测试用例结构
type TestCase struct {
	Name      string
	Node      *Node
	Expected  string
	ExpectErr bool // 新增字段，表示是否期望错误
}

// RunTestCases 运行一组测试用例
func RunTestCases(t *testing.T, testCases []TestCase) {
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			result, err := tc.Node.ToMarkdown(context.Background())
			if tc.ExpectErr {
				assert.Error(t, err, "Expected an error but got none")
			} else {
				assert.NoError(t, err, "Unexpected error")
				assert.Equal(t, tc.Expected, result, "Markdown conversion should match")
			}
		})
	}
}

// 辅助函数
func createLinkNode(text, url string) *Node {
	node := NewNode(NodeLink)
	node.SetData(NDK_URL, url)
	node.AddPhrasingChild(&Node{Type: NodeText, Value: text})
	return node
}

func createImageNode(alt, url string) *Node {
	node := NewNode(NodeImage)
	node.SetData(NDK_Alt, alt)
	node.SetData(NDK_URL, url)
	return node
}

func createHeadingNode(depth int, text string) *Node {
	node := NewNode(NodeHeading)
	node.SetData(NDK_Depth, depth)
	node.AddPhrasingChild(&Node{Type: NodeText, Value: text})
	return node
}

func createBlockquoteNode(text string) *Node {
	node := NewNode(NodeBlockquote)
	child := NewNode(NodeParagraph)
	child.AddPhrasingChild(&Node{Type: NodeText, Value: text})
	node.AddFlowChild(child)
	return node
}

func createListNode(ordered bool, text string) *Node {
	node := NewNode(NodeList)
	node.SetData(NDK_Ordered, ordered)
	listItem := NewNode(NodeListItem)
	listItem.AddFlowChild(&Node{Type: NodeParagraph, PhrasingChildren: []PhrasingContent{&Node{Type: NodeText, Value: text}}})
	node.AddListChild(listItem)
	return node
}

func createTableNode() *Node {
	node := NewNode(NodeTable)
	row := NewNode(NodeTableRow)
	row.AddTableChild(&Node{Type: NodeTableCell, PhrasingChildren: []PhrasingContent{&Node{Type: NodeText, Value: "Cell 1"}}})
	row.AddTableChild(&Node{Type: NodeTableCell, PhrasingChildren: []PhrasingContent{&Node{Type: NodeText, Value: "Cell 2"}}})
	node.AddTableChild(row)
	return node
}

func createCodeNode(lang, code string) *Node {
	node := NewNode(NodeCode)
	node.Value = code
	node.SetData(NDK_Lang, lang)
	return node
}

func createDefinitionNode(identifier, url, title string) *Node {
	node := NewNode(NodeDefinition)
	node.SetData(NDK_Identifier, identifier)
	node.SetData(NDK_URL, url)
	node.SetData(NDK_Title, title)
	return node
}

func createImageReferenceNode(identifier, alt string, referenceType ReferenceType) *Node {
	node := NewNode(NodeImageReference)
	node.SetData(NDK_Identifier, identifier)
	node.SetData(NDK_Alt, alt)
	node.SetData(NDK_ReferenceType, referenceType)
	return node
}

func createLinkReferenceNode(identifier, text string, referenceType ReferenceType) *Node {
	node := NewNode(NodeLinkReference)
	node.SetData(NDK_Identifier, identifier)
	node.SetData(NDK_ReferenceType, referenceType)
	node.AddPhrasingChild(&Node{Type: NodeText, Value: text})
	return node
}

func createFootnoteNode(content string) *Node {
	node := NewNode(NodeFootnote)
	node.AddPhrasingChild(&Node{Type: NodeText, Value: content})
	return node
}

func createFootnoteReferenceNode(identifier string) *Node {
	node := NewNode(NodeFootnoteReference)
	node.SetData(NDK_Identifier, identifier)
	node.SetData(NDK_Label, identifier)
	return node
}

func createFootnoteDefinitionNode(identifier, content string) *Node {
	node := NewNode(NodeFootnoteDefinition)
	node.SetData(NDK_Identifier, identifier)
	node.SetData(NDK_Label, identifier)
	node.AddFlowChild(&Node{Type: NodeParagraph, PhrasingChildren: []PhrasingContent{&Node{Type: NodeText, Value: content}}})
	return node
}

func createTableNodeWithAlignment() *Node {
	node := NewNode(NodeTable)
	node.SetData(NDK_Align, []AlignType{AlignLeft, AlignCenter, AlignRight})
	header := NewNode(NodeTableRow)
	header.AddTableChild(&Node{Type: NodeTableCell, PhrasingChildren: []PhrasingContent{&Node{Type: NodeText, Value: "Left"}}})
	header.AddTableChild(&Node{Type: NodeTableCell, PhrasingChildren: []PhrasingContent{&Node{Type: NodeText, Value: "Center"}}})
	header.AddTableChild(&Node{Type: NodeTableCell, PhrasingChildren: []PhrasingContent{&Node{Type: NodeText, Value: "Right"}}})
	node.AddTableChild(header)
	row := NewNode(NodeTableRow)
	row.AddTableChild(&Node{Type: NodeTableCell, PhrasingChildren: []PhrasingContent{&Node{Type: NodeText, Value: "1"}}})
	row.AddTableChild(&Node{Type: NodeTableCell, PhrasingChildren: []PhrasingContent{&Node{Type: NodeText, Value: "2"}}})
	row.AddTableChild(&Node{Type: NodeTableCell, PhrasingChildren: []PhrasingContent{&Node{Type: NodeText, Value: "3"}}})
	node.AddTableChild(row)
	return node
}

// 添加新的辅助函数
func createComplexParagraph() *Node {
	para := NewNode(NodeParagraph)
	para.AddPhrasingChild(&Node{Type: NodeText, Value: "This is a "})
	para.AddPhrasingChild(&Node{Type: NodeEmphasis, PhrasingChildren: []PhrasingContent{&Node{Type: NodeText, Value: "complex"}}})
	para.AddPhrasingChild(&Node{Type: NodeText, Value: " paragraph with "})
	para.AddPhrasingChild(&Node{Type: NodeStrong, PhrasingChildren: []PhrasingContent{&Node{Type: NodeText, Value: "nested"}}})
	para.AddPhrasingChild(&Node{Type: NodeText, Value: " elements and a "})
	para.AddPhrasingChild(createLinkNode("link", "https://example.com"))
	para.AddPhrasingChild(&Node{Type: NodeText, Value: "."})
	return para
}

func createComplexList() *Node {
	list := NewNode(NodeList)
	list.SetData(NDK_Ordered, true)

	item1 := NewNode(NodeListItem)
	item1.AddFlowChild(&Node{Type: NodeParagraph, PhrasingChildren: []PhrasingContent{&Node{Type: NodeText, Value: "First item"}}})

	item2 := NewNode(NodeListItem)
	item2Para := NewNode(NodeParagraph)
	item2Para.AddPhrasingChild(&Node{Type: NodeText, Value: "Second item with "})
	item2Para.AddPhrasingChild(&Node{Type: NodeEmphasis, PhrasingChildren: []PhrasingContent{&Node{Type: NodeText, Value: "emphasis"}}})
	item2.AddFlowChild(item2Para)

	sublist := NewNode(NodeList)
	sublist.SetData(NDK_Ordered, false)
	subItem := NewNode(NodeListItem)
	subItem.AddFlowChild(&Node{Type: NodeParagraph, PhrasingChildren: []PhrasingContent{&Node{Type: NodeText, Value: "Subitem"}}})
	sublist.AddListChild(subItem)
	item2.AddFlowChild(sublist)

	list.AddListChild(item1)
	list.AddListChild(item2)

	return list
}

// createNestedListNode 创建一个嵌套列表节点
func createNestedListNode() *Node {
	list := NewNode(NodeList)
	list.SetData(NDK_Ordered, false)

	item1 := NewNode(NodeListItem)
	// 使用 NodeParagraph 创建段落
	para1 := NewNode(NodeParagraph)
	para1.AddPhrasingChild(&Node{Type: NodeText, Value: "Item 1"})
	item1.AddFlowChild(para1)

	subList := NewNode(NodeList)
	subItem1 := NewNode(NodeListItem)
	// 使用 NodeParagraph 创建段落
	paraSub1 := NewNode(NodeParagraph)
	paraSub1.AddPhrasingChild(&Node{Type: NodeText, Value: "Subitem 1"})
	subItem1.AddFlowChild(paraSub1)
	subList.AddListChild(subItem1)

	subItem2 := NewNode(NodeListItem)
	// 使用 NodeParagraph 创建段落
	paraSub2 := NewNode(NodeParagraph)
	paraSub2.AddPhrasingChild(&Node{Type: NodeText, Value: "Subitem 2"})
	subItem2.AddFlowChild(paraSub2)
	subList.AddListChild(subItem2)

	item1.AddFlowChild(subList)
	list.AddListChild(item1)

	item2 := NewNode(NodeListItem)
	// 使用 NodeParagraph 创建段落
	para2 := NewNode(NodeParagraph)
	para2.AddPhrasingChild(&Node{Type: NodeText, Value: "Item 2"})
	item2.AddFlowChild(para2)
	list.AddListChild(item2)

	return list
}

// createComplexTableNode 创建一个复杂的表格节点
func createComplexTableNode() *Node {
	table := NewNode(NodeTable)
	row1 := NewNode(NodeTableRow)
	row1.AddTableChild(&Node{Type: NodeTableCell, PhrasingChildren: []PhrasingContent{&Node{Type: NodeText, Value: "Left"}}})
	row1.AddTableChild(&Node{Type: NodeTableCell, PhrasingChildren: []PhrasingContent{&Node{Type: NodeText, Value: "Center"}}})
	row1.AddTableChild(&Node{Type: NodeTableCell, PhrasingChildren: []PhrasingContent{&Node{Type: NodeText, Value: "Right"}}})
	table.AddTableChild(row1)

	row2 := NewNode(NodeTableRow)
	row2.AddTableChild(&Node{Type: NodeTableCell, PhrasingChildren: []PhrasingContent{&Node{Type: NodeText, Value: "1"}}})
	row2.AddTableChild(&Node{Type: NodeTableCell, PhrasingChildren: []PhrasingContent{&Node{Type: NodeText, Value: "2"}}})
	row2.AddTableChild(&Node{Type: NodeTableCell, PhrasingChildren: []PhrasingContent{&Node{Type: NodeText, Value: "3"}}})
	table.AddTableChild(row2)

	return table
}

// createComplexBlockquoteNode 创建一个复杂的引用节点
func createComplexBlockquoteNode() *Node {
	blockquote := NewNode(NodeBlockquote)
	blockquote.AddFlowChild(createHeadingNode(1, "Quoted heading"))
	blockquote.AddFlowChild(createParagraphNode("Quoted paragraph"))
	list := createListNode(false, "Quoted list item")
	blockquote.AddFlowChild(list)
	return blockquote
}

// createMixedBlockElements 创建一个混合块元素节点
func createMixedBlockElements() *Node {
	root := NewNode(NodeRoot)
	root.AddFlowChild(createHeadingNode(1, "Heading"))
	root.AddFlowChild(createParagraphNode("Paragraph"))
	root.AddFlowChild(createBlockquoteNode("Blockquote"))
	root.AddFlowChild(createListNode(false, "List item"))
	root.AddFlowChild(createCodeNode("go", "fmt.Println(\"Hello\")"))
	return root
}

// createParagraphNode 创建一个段落节点
func createParagraphNode(content string) *Node {
	paragraph := NewNode(NodeParagraph)
	paragraph.AddPhrasingChild(&Node{Type: NodeText, Value: content})
	return paragraph
}
