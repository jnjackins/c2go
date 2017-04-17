package main

import "io"

type DeclStmt struct {
	Address  string
	Position string
	Children []interface{}
}

func parseDeclStmt(line string) *DeclStmt {
	groups := groupsFromRegex(
		"<(?P<position>.*)>",
		line,
	)

	return &DeclStmt{
		Address:  groups["address"],
		Position: groups["position"],
		Children: []interface{}{},
	}
}

func (n *DeclStmt) renderLine(w io.Writer, functionName string, indent int, returnType string) {
	for _, child := range n.Children {
		printLine(w, renderExpression(child)[0], indent)
	}
}
