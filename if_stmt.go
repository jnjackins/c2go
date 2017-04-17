package main

import (
	"fmt"
	"io"
)

type IfStmt struct {
	Address  string
	Position string
	Children []interface{}
}

func parseIfStmt(line string) *IfStmt {
	groups := groupsFromRegex(
		"<(?P<position>.*)>",
		line,
	)

	return &IfStmt{
		Address:  groups["address"],
		Position: groups["position"],
		Children: []interface{}{},
	}
}

func (n *IfStmt) renderLine(w io.Writer, functionName string, indent int, returnType string) {
	// TODO: The first two children of an IfStmt appear to always be null.
	// Are there any cases where they are used?
	children := n.Children[2:]

	e := renderExpression(children[0])
	printLine(w, fmt.Sprintf("if %s {", cast(e[0], e[1], "bool")), indent)

	render(w, children[1], functionName, indent+1, returnType)

	if len(children) > 2 {
		printLine(w, "} else {", indent)
		render(w, children[2], functionName, indent+1, returnType)
	}

	printLine(w, "}", indent)
}
