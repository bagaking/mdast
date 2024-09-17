package mdast

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComplexStructure(t *testing.T) {
	root := NewNode(NodeRoot)
	heading := createHeadingNode(2, "Title")
	root.AddFlowChild(heading)

	para := NewNode(NodeParagraph)
	para.AddPhrasingChild(&Node{Type: NodeText, Value: "This is a "})
	para.AddPhrasingChild(&Node{Type: NodeStrong, PhrasingChildren: []PhrasingContent{&Node{Type: NodeText, Value: "test"}}})
	para.AddPhrasingChild(&Node{Type: NodeText, Value: " paragraph."})
	root.AddFlowChild(para)

	expected := "## Title\n\nThis is a **test** paragraph.\n\n"
	result, err := root.ToMarkdown(context.Background())
	assert.NoError(t, err, "Unexpected error")
	assert.Equal(t, expected, result, "Complex structure conversion should match")
}

func TestTableWithAlignment(t *testing.T) {
	node := createTableNodeWithAlignment()
	expected := "| Left | Center | Right |\n| :--- | :---: | ---: |\n| 1 | 2 | 3 |\n\n"
	result, err := node.ToMarkdown(context.Background())
	assert.NoError(t, err, "Unexpected error")
	assert.Equal(t, expected, result, "Table with alignment should convert correctly")
}

func TestDeepNestedStructure(t *testing.T) {
	root := NewNode(NodeRoot)
	heading := createHeadingNode(1, "First level")
	root.AddFlowChild(heading)

	secondLevel := NewNode(NodeBlockquote)
	secondLevel.AddFlowChild(createHeadingNode(2, "Second level"))
	root.AddFlowChild(secondLevel)

	thirdLevel := NewNode(NodeList)
	thirdLevel.SetData(NDK_Ordered, true)
	thirdLevel.AddListChild(createListNode(true, "Third level item"))
	secondLevel.AddFlowChild(thirdLevel)

	expected := "# First level\n\n> ## Second level\n> \n> 1. Third level item\n\n"
	result, err := root.ToMarkdown(context.Background())
	assert.NoError(t, err, "Unexpected error")
	assert.Equal(t, expected, result, "Deep nested structure should convert correctly")
}

func TestComplexDocument(t *testing.T) {
	doc := createComplexDocument()
	expected := "# Complex Document\n\n## Introduction\n\nThis is a *complex* document with various elements:\n\n1. Lists\n2. Tables\n3. Code blocks\n\n### Lists\n\n- Unordered item 1\n- Unordered item 2\n  1. Nested ordered item\n  2. Another nested item\n\n### Table\n\n| Header 1 | Header 2 | Header 3 |\n| :--- | :---: | ---: |\n| Left | Center | Right |\n| Data | Data | Data |\n\n### Code\n\n```go\nfunc main() {\n    fmt.Println(\"Hello, world!\")\n}\n```\n\n> This is a blockquote with a [link](https://example.com).\n\nFootnote reference[^1]\n\n[^1]: This is a footnote.\n\n"
	result, err := doc.ToMarkdown(context.Background())
	assert.NoError(t, err, "Unexpected error")
	assert.Equal(t, expected, result, "Complex document should convert correctly")
}

func ExampleNode_ToMarkdown_complexStructure() {
	root := NewNode(NodeRoot)
	root.AddFlowChild(createHeadingNode(1, "Complex Example"))
	para := NewNode(NodeParagraph)
	para.AddPhrasingChild(&Node{Type: NodeText, Value: "This is a "})
	para.AddPhrasingChild(&Node{Type: NodeEmphasis, PhrasingChildren: []PhrasingContent{&Node{Type: NodeText, Value: "complex"}}})
	para.AddPhrasingChild(&Node{Type: NodeText, Value: " example with "})
	para.AddPhrasingChild(&Node{Type: NodeStrong, PhrasingChildren: []PhrasingContent{&Node{Type: NodeText, Value: "nested"}}})
	para.AddPhrasingChild(&Node{Type: NodeText, Value: " elements."})
	root.AddFlowChild(para)

	// 添加调试信息
	fmt.Fprintf(os.Stderr, "Root node type: %s\n", root.Type)
	fmt.Fprintf(os.Stderr, "Number of children: %d\n", len(root.FlowChildren))
	for i, child := range root.FlowChildren {
		fmt.Fprintf(os.Stderr, "Child %d type: %s\n", i, child.GetType())
	}

	markdown, err := root.ToMarkdown(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}
	fmt.Print(markdown)
	// Output:
	// # Complex Example
	//
	// This is a *complex* example with **nested** elements.
	//
}

// 辅助函数
func createComplexDocument() *Node {
	doc := NewNode(NodeRoot)
	doc.AddFlowChild(createHeadingNode(1, "Complex Document"))
	doc.AddFlowChild(createHeadingNode(2, "Introduction"))

	intro := NewNode(NodeParagraph)
	intro.AddPhrasingChild(&Node{Type: NodeText, Value: "This is a "})
	intro.AddPhrasingChild(&Node{Type: NodeEmphasis, PhrasingChildren: []PhrasingContent{&Node{Type: NodeText, Value: "complex"}}})
	intro.AddPhrasingChild(&Node{Type: NodeText, Value: " document with various elements:"})
	doc.AddFlowChild(intro)

	list := NewNode(NodeList)
	list.SetData("ordered", true)
	list.AddListChild(&Node{Type: NodeListItem, PhrasingChildren: []PhrasingContent{&Node{Type: NodeText, Value: "Lists"}}})
	list.AddListChild(&Node{Type: NodeListItem, PhrasingChildren: []PhrasingContent{&Node{Type: NodeText, Value: "Tables"}}})
	list.AddListChild(&Node{Type: NodeListItem, PhrasingChildren: []PhrasingContent{&Node{Type: NodeText, Value: "Code blocks"}}})
	doc.AddFlowChild(list)

	doc.AddFlowChild(createHeadingNode(3, "Lists"))
	doc.AddFlowChild(createComplexList())

	doc.AddFlowChild(createHeadingNode(3, "Table"))
	doc.AddFlowChild(createTableNodeWithAlignment())

	doc.AddFlowChild(createHeadingNode(3, "Code"))
	doc.AddFlowChild(createCodeNode("go", "func main() {\n    fmt.Println(\"Hello, world!\")\n}"))

	quote := createBlockquoteNode("This is a blockquote with a ")
	quote.AddFlowChild(createLinkNode("link", "https://example.com"))
	doc.AddFlowChild(quote)

	doc.AddFlowChild(&Node{
		Type: NodeParagraph,
		PhrasingChildren: []PhrasingContent{
			&Node{Type: NodeText, Value: "Footnote reference"},
			&Node{Type: NodeFootnoteReference, Data: map[DataKey]any{"identifier": "1"}},
		},
	})

	doc.AddFlowChild(createFootnoteDefinitionNode("1", "This is a footnote."))

	return doc
}

func TestComplexListStructure(t *testing.T) {
	root := NewNode(NodeRoot)
	list := NewNode(NodeList)
	list.SetData(NDK_Ordered, true)

	item1 := NewNode(NodeListItem)
	item1.AddFlowChild(createParagraphNode("First level ordered item"))

	nestedList := NewNode(NodeList)
	nestedList.SetData(NDK_Ordered, false)
	nestedItem := NewNode(NodeListItem)
	nestedItem.AddFlowChild(createParagraphNode("Second level unordered item"))
	nestedList.AddListChild(nestedItem)

	item1.AddFlowChild(nestedList)
	list.AddListChild(item1)

	item2 := NewNode(NodeListItem)
	item2.AddFlowChild(createParagraphNode("Another first level item"))
	list.AddListChild(item2)

	root.AddFlowChild(list)

	expected := "1. First level ordered item\n\n   - Second level unordered item\n\n2. Another first level item\n\n"
	result, err := root.ToMarkdown(context.Background())
	assert.NoError(t, err, "Unexpected error")
	assert.Equal(t, expected, result, "Complex list structure should convert correctly")
}

func TestListWithMissingOrderedProperty(t *testing.T) {
	root := NewNode(NodeRoot)
	list := NewNode(NodeList)
	// 故意不设置 NDK_Ordered 属性
	item := NewNode(NodeListItem)
	item.AddFlowChild(createParagraphNode("Test item"))
	list.AddListChild(item)
	root.AddFlowChild(list)

	_, err := root.ToMarkdown(context.Background())
	assert.Error(t, err, "Expected an error due to missing 'ordered' property")
	assert.Contains(t, err.Error(), "missing required 'ordered' property for list")
}

func TestListWithInvalidChild(t *testing.T) {
	root := NewNode(NodeRoot)
	list := NewNode(NodeList)
	list.SetData(NDK_Ordered, true)
	// 故意添加一个非 ListItem 的子节点
	list.AddListChild(createParagraphNode("Invalid child"))
	root.AddFlowChild(list)

	_, err := root.ToMarkdown(context.Background())
	assert.Error(t, err, "Expected an error due to invalid child node")
	assert.Contains(t, err.Error(), "unexpected node type in list")
}
