package main

import "io"

type CompoundStmt struct {
	Address  string
	Position string
	Children []interface{}
}

func parseCompoundStmt(line string) *CompoundStmt {
	groups := groupsFromRegex(
		"<(?P<position>.*)>",
		line,
	)

	return &CompoundStmt{
		Address:  groups["address"],
		Position: groups["position"],
		Children: []interface{}{},
	}
}

func (n *CompoundStmt) renderLine(w io.Writer, functionName string, indent int, returnType string) {
	for _, c := range n.Children {
		render(w, c, functionName, indent, returnType)
	}
}
