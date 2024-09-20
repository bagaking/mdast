package mdast

import (
	"context"
	"fmt"
	"testing"
)

func TestReferenceElements(t *testing.T) {
	testCases := []TestCase{
		{"Definition", createDefinitionNode("example", "https://example.com", "Example Title"), "[example]: https://example.com \"Example Title\"\n", false},
		{"Definition without title", createDefinitionNode("example", "https://example.com", ""), "[example]: https://example.com\n", false},
		{"ImageReference full", createImageReferenceNode("example", "Alt Text", "full"), "![Alt Text][example]", false},
		{"ImageReference collapsed", createImageReferenceNode("example", "Alt Text", "collapsed"), "![Alt Text][]", false},
		{"ImageReference shortcut", createImageReferenceNode("example", "Alt Text", "shortcut"), "![Alt Text]", false},
		{"LinkReference full", createLinkReferenceNode("example", "Link Text", "full"), "[Link Text][example]", false},
		{"LinkReference collapsed", createLinkReferenceNode("example", "Link Text", "collapsed"), "[Link Text][]", false},
		{"LinkReference shortcut", createLinkReferenceNode("example", "Link Text", "shortcut"), "[Link Text]", false},
		{"Footnote", createFootnoteNode("Footnote content"), "[^Footnote content]", false},
		{"FootnoteReference", createFootnoteReferenceNode("1"), "[^1]", false},
		{"FootnoteDefinition", createFootnoteDefinitionNode("1", "Footnote content"), "[^1]: Footnote content\n\n", false},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			result, err := tc.Node.ToMarkdown(context.Background())
			if tc.ExpectErr {
				if err == nil {
					t.Fatalf("Node.ToMarkdown(%s) error = nil, want error", tc.Name)
				}
			} else {
				if err != nil {
					t.Fatalf("Node.ToMarkdown(%s) error = %v, want nil", tc.Name, err)
				}
				if result != tc.Expected {
					t.Errorf("Node.ToMarkdown(%s) = %q, want %q", tc.Name, result, tc.Expected)
				}
			}
		})
	}
}

func TestDefinitionDestinationAndTitleEscaping(t *testing.T) {
	testCases := []TestCase{
		{
			Name:     "Definition URL with space",
			Node:     createDefinitionNode("example", "https://example.com/a b", ""),
			Expected: "[example]: <https://example.com/a b>\n",
		},
		{
			Name:     "Definition URL with closing paren",
			Node:     createDefinitionNode("example", "https://example.com/a)b", ""),
			Expected: "[example]: <https://example.com/a)b>\n",
		},
		{
			Name:     "Definition title with quote and backslash",
			Node:     createDefinitionNode("example", "https://example.com", `A "quoted" \ title`),
			Expected: "[example]: https://example.com \"A \\\"quoted\\\" \\\\ title\"\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			result, err := tc.Node.ToMarkdown(context.Background())
			if err != nil {
				t.Fatalf("Node.ToMarkdown(%s) error = %v, want nil", tc.Name, err)
			}
			if result != tc.Expected {
				t.Errorf("Node.ToMarkdown(%s) = %q, want %q", tc.Name, result, tc.Expected)
			}
		})
	}
}

func TestReferenceTypes(t *testing.T) {
	testCases := []struct {
		name          string
		referenceType ReferenceType
		expected      string
	}{
		{"Shortcut", "shortcut", "[Link]"},
		{"Collapsed", "collapsed", "[Link][]"},
		{"Full", "full", "[Link][id]"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			node := createLinkReferenceNode("id", "Link", tc.referenceType)
			result, err := node.ToMarkdown(context.Background())
			if err != nil {
				t.Fatalf("Node.ToMarkdown(%s) error = %v, want nil", tc.name, err)
			}
			if result != tc.expected {
				t.Errorf("Node.ToMarkdown(%s) = %q, want %q", tc.name, result, tc.expected)
			}
		})
	}
}

func ExampleNode_ToMarkdown_definition() {
	node := createDefinitionNode("example", "https://example.com", "Example Title")
	markdown, err := node.ToMarkdown(context.Background())
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Print(markdown)
	// Output: [example]: https://example.com "Example Title"
}

func ExampleNode_ToMarkdown_imageReference() {
	node := createImageReferenceNode("example", "Alt Text", "full")
	markdown, err := node.ToMarkdown(context.Background())
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Print(markdown)
	// Output: ![Alt Text][example]
}
