package main

import (
	"fmt"
	"io"
)

type WhileStmt struct {
	Address  string
	Position string
	Children []interface{}
}

func parseWhileStmt(line string) *WhileStmt {
	groups := groupsFromRegex(
		"<(?P<position>.*)>",
		line,
	)

	return &WhileStmt{
		Address:  groups["address"],
		Position: groups["position"],
		Children: []interface{}{},
	}
}

func (n *WhileStmt) renderLine(w io.Writer, functionName string, indent int, returnType string) {
	// TODO: The first child of a WhileStmt appears to always be null.
	// Are there any cases where it is used?
	children := n.Children[1:]

	e := renderExpression(children[0])
	printLine(w, fmt.Sprintf("for %s {", cast(e[0], e[1], "bool")), indent)

	// FIXME: Does this do anything?
	render(w, children[1], functionName, indent+1, returnType)

	printLine(w, "}", indent)
}
