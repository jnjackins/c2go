package main

import "io"

type ReturnStmt struct {
	Address  string
	Position string
	Children []interface{}
}

func parseReturnStmt(line string) *ReturnStmt {
	groups := groupsFromRegex(
		"<(?P<position>.*)>",
		line,
	)

	return &ReturnStmt{
		Address:  groups["address"],
		Position: groups["position"],
		Children: []interface{}{},
	}
}

func (n *ReturnStmt) renderLine(w io.Writer, functionName string, indent int, returnType string) {
	r := "return"

	if len(n.Children) > 0 && functionName != "main" {
		re := renderExpression(n.Children[0])
		r = "return " + cast(re[0], re[1], "int")
	}

	printLine(w, r, indent)
}
