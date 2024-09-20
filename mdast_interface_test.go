package mdast

import (
	"context"
	"errors"
	"testing"
)

type customFlowContent struct {
	markdown string
	err      error
}

func (c customFlowContent) ToMarkdown(context.Context) (string, error) {
	return c.markdown, c.err
}

func (c customFlowContent) GetType() NodeType {
	return NodeParagraph
}

func (c customFlowContent) IsFlow() {}

type customPhrasingContent struct {
	markdown string
	err      error
}

func (c customPhrasingContent) ToMarkdown(context.Context) (string, error) {
	return c.markdown, c.err
}

func (c customPhrasingContent) GetType() NodeType {
	return NodeText
}

func (c customPhrasingContent) IsPhrasing() {}

type customListContent struct {
	markdown string
	err      error
}

func (c customListContent) ToMarkdown(context.Context) (string, error) {
	return c.markdown, c.err
}

func (c customListContent) GetType() NodeType {
	return NodeListItem
}

func (c customListContent) IsList() {}

type customTableCellContent struct {
	markdown string
	err      error
}

func (c customTableCellContent) ToMarkdown(context.Context) (string, error) {
	return c.markdown, c.err
}

func (c customTableCellContent) GetType() NodeType {
	return NodeTableCell
}

func (c customTableCellContent) IsTable() {}

type customTableRowContent struct {
	markdown string
	err      error
}

func (c customTableRowContent) ToMarkdown(context.Context) (string, error) {
	return c.markdown, c.err
}

func (c customTableRowContent) GetType() NodeType {
	return NodeTableRow
}

func (c customTableRowContent) IsTable() {}

func TestCustomContentImplementationsRenderWithoutPanic(t *testing.T) {
	root := NewNode(NodeRoot)
	root.AddFlowChild(customFlowContent{markdown: "custom flow\n\n"})
	assertMarkdown(t, "custom flow content", root, "custom flow\n\n")

	paragraph := NewNode(NodeParagraph)
	paragraph.AddPhrasingChild(customPhrasingContent{markdown: "custom inline"})
	assertMarkdown(t, "custom phrasing content", paragraph, "custom inline\n\n")

	list := NewNode(NodeList)
	list.AddListChild(customListContent{markdown: "- custom item"})
	assertMarkdown(t, "custom list content", list, "- custom item\n\n")

	row := NewNode(NodeTableRow)
	row.AddTableChild(customTableCellContent{markdown: "custom cell"})
	assertMarkdown(t, "custom table cell content", row, "| custom cell |")
}

func TestCustomContentImplementationErrorsPropagate(t *testing.T) {
	wantErr := errors.New("custom content failed")
	testCases := []struct {
		name string
		node *Node
	}{
		{
			name: "flow",
			node: func() *Node {
				root := NewNode(NodeRoot)
				root.AddFlowChild(customFlowContent{err: wantErr})
				return root
			}(),
		},
		{
			name: "phrasing",
			node: func() *Node {
				paragraph := NewNode(NodeParagraph)
				paragraph.AddPhrasingChild(customPhrasingContent{err: wantErr})
				return paragraph
			}(),
		},
		{
			name: "list",
			node: func() *Node {
				list := NewNode(NodeList)
				list.AddListChild(customListContent{err: wantErr})
				return list
			}(),
		},
		{
			name: "table cell",
			node: func() *Node {
				row := NewNode(NodeTableRow)
				row.AddTableChild(customTableCellContent{err: wantErr})
				return row
			}(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := tc.node.ToMarkdown(context.Background())
			if !errors.Is(err, wantErr) {
				t.Errorf("Node.ToMarkdown(%s) error = %v, want %v", tc.name, err, wantErr)
			}
		})
	}
}

func TestCustomTableHeaderRowReportsUnsupportedSeparator(t *testing.T) {
	table := NewNode(NodeTable)
	table.AddTableChild(customTableRowContent{markdown: "| custom header |"})

	assertMarkdownErrorContains(t, "custom table header row", table, "cannot render table separator")
}
