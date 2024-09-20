# mdast

`mdast` is a small Go package for building Markdown AST nodes and rendering
them back to Markdown. It is a renderer/builder helper, not a Markdown parser.

## Status

Early library. The public types are usable for simple Markdown generation, but
edge cases such as escaping, custom node implementations, and full GFM parity
should be validated before depending on the package in production code.

## Install

```bash
go get github.com/bagaking/mdast
```

## Example

```go
package main

import (
	"context"
	"fmt"

	"github.com/bagaking/mdast"
)

func main() {
	root := mdast.NewNode(mdast.NodeRoot)

	heading := mdast.NewNode(mdast.NodeHeading)
	heading.SetData(mdast.NDK_Depth, 1)
	heading.AddPhrasingChild(&mdast.Node{Type: mdast.NodeText, Value: "Release Notes"})
	root.AddFlowChild(heading)

	item := mdast.NewNode(mdast.NodeListItem)
	item.AddFlowChild(&mdast.Node{
		Type: mdast.NodeParagraph,
		PhrasingChildren: []mdast.PhrasingContent{
			&mdast.Node{Type: mdast.NodeText, Value: "Ship the smoke tests"},
		},
	})

	list := mdast.NewNode(mdast.NodeList)
	list.SetData(mdast.NDK_Ordered, true)
	list.SetData(mdast.NDK_Start, 1)
	list.AddListChild(item)
	root.AddFlowChild(list)

	out, err := root.ToMarkdown(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Print(out)
}
```

Output:

```markdown
# Release Notes

1. Ship the smoke tests
```

## Supported Surface

- Flow nodes: paragraphs, headings, blockquotes, code, thematic breaks, HTML,
  YAML, definitions, footnote definitions, lists, and GFM tables.
- Phrasing nodes: text, emphasis, strong, delete, links, images, inline code,
  breaks, link/image references, footnote references, and the package's
  `NodeFootnote` inline helper.
- List rendering follows mdast defaults: missing `ordered` renders unordered
  lists, positive ordered `start` values are honored, missing or `<1` `start`
  values render from `1`, and list/listItem `spread` controls blank line
  separation.

## Validation

```bash
gofmt -w *.go
go test ./...
go vet ./...
go list -f '{{.GoFiles}} {{.Imports}}' .
```

The production package should not import `testing` or assertion libraries.
