package mdast

import (
	"context"
	"strings"
	"testing"
)

// TestCase defines a markdown rendering test case.
type TestCase struct {
	Name      string
	Node      *Node
	Expected  string
	ExpectErr bool
}

// RunTestCases runs markdown rendering test cases.
func RunTestCases(t *testing.T, testCases []TestCase) {
	t.Helper()

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			result, err := tc.Node.ToMarkdown(context.Background())
			if tc.ExpectErr {
				if err == nil {
					t.Fatalf("Node.ToMarkdown(%s) error = nil, want error", tc.Name)
				}
				return
			}
			if err != nil {
				t.Fatalf("Node.ToMarkdown(%s) error = %v, want nil", tc.Name, err)
			}
			if result != tc.Expected {
				t.Errorf("Node.ToMarkdown(%s) = %q, want %q", tc.Name, result, tc.Expected)
			}
		})
	}
}

func assertMarkdown(t *testing.T, name string, node *Node, want string) {
	t.Helper()

	got, err := node.ToMarkdown(context.Background())
	if err != nil {
		t.Fatalf("Node.ToMarkdown(%s) error = %v, want nil", name, err)
	}
	if got != want {
		t.Errorf("Node.ToMarkdown(%s) = %q, want %q", name, got, want)
	}
}

func assertMarkdownErrorContains(t *testing.T, name string, node *Node, want string) {
	t.Helper()

	_, err := node.ToMarkdown(context.Background())
	if err == nil {
		t.Fatalf("Node.ToMarkdown(%s) error = nil, want error containing %q", name, want)
	}
	if !strings.Contains(err.Error(), want) {
		t.Errorf("Node.ToMarkdown(%s) error = %q, want containing %q", name, err.Error(), want)
	}
}

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

func createNestedListNode() *Node {
	list := NewNode(NodeList)
	list.SetData(NDK_Ordered, false)

	item1 := NewNode(NodeListItem)
	para1 := NewNode(NodeParagraph)
	para1.AddPhrasingChild(&Node{Type: NodeText, Value: "Item 1"})
	item1.AddFlowChild(para1)

	subList := NewNode(NodeList)
	subList.SetData(NDK_Ordered, false)
	subItem1 := NewNode(NodeListItem)
	paraSub1 := NewNode(NodeParagraph)
	paraSub1.AddPhrasingChild(&Node{Type: NodeText, Value: "Subitem 1"})
	subItem1.AddFlowChild(paraSub1)
	subList.AddListChild(subItem1)

	subItem2 := NewNode(NodeListItem)
	paraSub2 := NewNode(NodeParagraph)
	paraSub2.AddPhrasingChild(&Node{Type: NodeText, Value: "Subitem 2"})
	subItem2.AddFlowChild(paraSub2)
	subList.AddListChild(subItem2)

	item1.AddFlowChild(subList)
	list.AddListChild(item1)

	item2 := NewNode(NodeListItem)
	para2 := NewNode(NodeParagraph)
	para2.AddPhrasingChild(&Node{Type: NodeText, Value: "Item 2"})
	item2.AddFlowChild(para2)
	list.AddListChild(item2)

	return list
}

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

func createComplexBlockquoteNode() *Node {
	blockquote := NewNode(NodeBlockquote)
	blockquote.AddFlowChild(createHeadingNode(1, "Quoted heading"))
	blockquote.AddFlowChild(createParagraphNode("Quoted paragraph"))
	list := createListNode(false, "Quoted list item")
	blockquote.AddFlowChild(list)
	return blockquote
}

func createMixedBlockElements() *Node {
	root := NewNode(NodeRoot)
	root.AddFlowChild(createHeadingNode(1, "Heading"))
	root.AddFlowChild(createParagraphNode("Paragraph"))
	root.AddFlowChild(createBlockquoteNode("Blockquote"))
	root.AddFlowChild(createListNode(false, "List item"))
	root.AddFlowChild(createCodeNode("go", "fmt.Println(\"Hello\")"))
	return root
}

func createParagraphNode(content string) *Node {
	paragraph := NewNode(NodeParagraph)
	paragraph.AddPhrasingChild(&Node{Type: NodeText, Value: content})
	return paragraph
}
