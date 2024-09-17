package mdast

import (
	"context"
	"fmt"
	"testing"
)

// 辅助函数
func createLinkNodeWithTitle(text, url, title string) *Node {
	node := createLinkNode(text, url)
	node.SetData(NDK_Title, title)
	return node
}

func createImageNodeWithTitle(alt, url, title string) *Node {
	node := createImageNode(alt, url)
	node.SetData(NDK_Title, title)
	return node
}

func createNestedInlineNode() *Node {
	node := NewNode(NodeParagraph) // 使用 NodeParagraph
	node.AddPhrasingChild(&Node{Type: NodeText, Value: "This is "})
	emphasis := &Node{Type: NodeEmphasis}
	emphasis.AddPhrasingChild(&Node{Type: NodeText, Value: "emphasized and "})
	strong := &Node{Type: NodeStrong}
	strong.AddPhrasingChild(&Node{Type: NodeText, Value: "strong"})
	emphasis.AddPhrasingChild(strong)
	emphasis.AddPhrasingChild(&Node{Type: NodeText, Value: " text"})
	node.AddPhrasingChild(emphasis)
	node.AddPhrasingChild(&Node{Type: NodeText, Value: " with "})
	node.AddPhrasingChild(&Node{Type: NodeInlineCode, Value: "code"})
	return node
}

// 测试用例

func TestInlineElements(t *testing.T) {
	testCases := []TestCase{
		{"Text", &Node{Type: NodeText, Value: "Hello"}, "Hello", false},
		{"Text with special characters", &Node{Type: NodeText, Value: "Hello * _ ` [ ]"}, "Hello * _ ` [ ]", false},
		{"Emphasis", &Node{Type: NodeEmphasis, PhrasingChildren: []PhrasingContent{&Node{Type: NodeText, Value: "em"}}}, "*em*", false},
		{"Strong", &Node{Type: NodeStrong, PhrasingChildren: []PhrasingContent{&Node{Type: NodeText, Value: "strong"}}}, "**strong**", false},
		{"Delete", &Node{Type: NodeDelete, PhrasingChildren: []PhrasingContent{&Node{Type: NodeText, Value: "deleted"}}}, "~~deleted~~", false},
		{"InlineCode", &Node{Type: NodeInlineCode, Value: "code"}, "`code`", false},
		{"InlineCode with backticks", &Node{Type: NodeInlineCode, Value: "code with ` backticks"}, "`code with ` backticks`", false},
		{"Link", createLinkNode("Example", "https://example.com"), "[Example](https://example.com)", false},
		{"Link with title", createLinkNodeWithTitle("Example", "https://example.com", "Title"), "[Example](https://example.com \"Title\")", false},
		{"Image", createImageNode("Alt text", "https://example.com/image.png"), "![Alt text](https://example.com/image.png)", false},
		{"Image with title", createImageNodeWithTitle("Alt text", "https://example.com/image.png", "Title"), "![Alt text](https://example.com/image.png \"Title\")", false},
		{"Break", NewNode(NodeBreak), "\n", false},
		{"Nested Inline", createNestedInlineNode(), "This is *emphasized and **strong** text* with `code`", false},
	}

	RunTestCases(t, testCases)
}

func ExampleNode_ToMarkdown_emphasis() {
	node := &Node{Type: NodeEmphasis, PhrasingChildren: []PhrasingContent{&Node{Type: NodeText, Value: "emphasized text"}}}
	result, err := node.ToMarkdown(context.Background())
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Print(result)
	// Output: *emphasized text*
}

func ExampleNode_ToMarkdown_nestedInline() {
	node := createNestedInlineNode()
	result, err := node.ToMarkdown(context.Background())
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Print(result)
	// Output: This is *emphasized and **strong** text* with `code`
}

func TestAdditionalInlineElements(t *testing.T) {
	testCases := []TestCase{
		{"HTML", &Node{Type: NodeHTML, Value: "<span>inline html</span>"}, "<span>inline html</span>", false},
		{"LinkReference", createLinkReferenceNode("ref", "Link", "full"), "[Link][ref]", false},
		{"ImageReference", createImageReferenceNode("ref", "Alt text", "full"), "![Alt text][ref]", false},
		{"FootnoteReference", createFootnoteReferenceNode("1"), "[^1]", false},
	}

	RunTestCases(t, testCases)
}
