package mdast

import (
	"context"
	"fmt"
	"testing"
)

func TestBlockElements(t *testing.T) {
	testCases := []TestCase{
		{"Empty Paragraph", NewNode(NodeParagraph), "\n\n", false},
		{"Paragraph with content", &Node{Type: NodeParagraph, PhrasingChildren: []PhrasingContent{&Node{Type: NodeText, Value: "Hello, world!"}}}, "Hello, world!\n\n", false},
		{"Heading level 1", createHeadingNode(1, "Title"), "# Title\n\n", false},
		{"Heading level 3", createHeadingNode(3, "Subtitle"), "### Subtitle\n\n", false},
		{"Blockquote", createBlockquoteNode("This is a quote."), "> This is a quote.\n\n", false},
		{"Blockquote with multiple lines", createBlockquoteNode("Line 1\nLine 2"), "> Line 1\n> Line 2\n\n", false},
		{"Ordered List", createListNode(true, "First item"), "1. First item\n\n", false},
		{"Unordered List", createListNode(false, "Bullet item"), "- Bullet item\n\n", false},
		{"Table", createTableNode(), "| Cell 1 | Cell 2 |\n| --- | --- |\n\n", false},
		{"Code", createCodeNode("javascript", "console.log('Hello');"), "```javascript\nconsole.log('Hello');\n```\n\n", false},
		{"Code without language", createCodeNode("", "print('Hello')"), "```\nprint('Hello')\n```\n\n", false},
		{"ThematicBreak", NewNode(NodeThematicBreak), "---\n\n", false},
		{"HTML", &Node{Type: NodeHTML, Value: "<div>Test</div>"}, "<div>Test</div>\n\n", false},
		{"Yaml", &Node{Type: NodeYaml, Value: "title: Test\ndate: 2023-05-01"}, "---\ntitle: Test\ndate: 2023-05-01\n---\n\n", false},
		{"Complex Paragraph", createComplexParagraph(), "This is a *complex* paragraph with **nested** elements and a [link](https://example.com).\n\n", false},
		{"Complex List", createComplexList(), "1. First item\n2. Second item with *emphasis*\n   - Subitem\n\n", false},
		{"Definition", createDefinitionNode("id", "https://example.com", "Title"), "[id]: https://example.com \"Title\"\n", false},
		{"FootnoteDefinition", createFootnoteDefinitionNode("1", "Footnote content"), "[^1]: Footnote content\n\n", false},
		{"Nested List", createNestedListNode(), "- Item 1\n  - Subitem 1\n  - Subitem 2\n- Item 2\n\n", false},
		{"Complex Table", createComplexTableNode(), "| Left | Center | Right |\n| --- | --- | --- |\n| 1 | 2 | 3 |\n\n", false},
		{"Complex Blockquote", createComplexBlockquoteNode(), "> # Quoted heading\n>\n> Quoted paragraph\n>\n> - Quoted list item\n\n", false},
		{"Mixed Block Elements", createMixedBlockElements(), "# Heading\n\nParagraph\n\n> Blockquote\n\n- List item\n\n```go\nfmt.Println(\"Hello\")\n```\n\n", false},
		// 添加一个期望错误的测试用例
		{"Invalid Node Type", &Node{Type: "invalid"}, "", true},
	}

	RunTestCases(t, testCases)
}

func ExampleNode_ToMarkdown_heading() {
	node := createHeadingNode(2, "Example Heading")
	result, err := node.ToMarkdown(context.Background())
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Print(result)
	// Output:
	// ## Example Heading
	//
}

func ExampleNode_ToMarkdown_complexList() {
	node := createComplexList()
	result, err := node.ToMarkdown(context.Background())
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Print(result)
	// Output:
	// 1. First item
	// 2. Second item with *emphasis*
	//    - Subitem
	//
}

// 添加更多的 Example 函数...
