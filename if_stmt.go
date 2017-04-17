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
	children := n.Children

	// There is always 4 or 5 children in an IfStmt. For example:
	//
	//     if (i == 0) {
	//         return 0;
	//     } else {
	//         return 1;
	//     }
	//
	// 1. Not sure what this is for. This gets removed.
	// 2. Not sure what this is for.
	// 3. conditional = BinaryOperator: i == 0
	// 4. body = CompoundStmt: { return 0; }
	// 5. elseBody = CompoundStmt: { return 1; }
	//
	// elseBody will be nil if there is no else clause.

	// On linux I have seen only 4 children for an IfStmt with the same
	// definitions above, but missing the first argument. Since we don't
	// know what the first argument is for anyway we will just remove it on
	// Mac if necessary.
	if len(children) == 5 && children[0] != nil {
		panic("non-nil child 0 in ForStmt")
	}
	if len(children) == 5 {
		children = children[1:]
	}

	// From here on there must be 4 children.
	if len(children) != 4 {
		panic(fmt.Sprintf("Expected 4 children in IfStmt, got %#v", children))
	}

	// Maybe we will discover what the nil value is?
	if children[0] != nil {
		panic("non-nil child 0 in ForStmt")
	}

	conditional := renderExpression(children[1])

	// The condition in Go must always be a bool.
	boolCondition := cast(conditional[0], conditional[1], "bool")

	printLine(w, fmt.Sprintf("if %s {", boolCondition), indent)

	render(w, children[2], functionName, indent+1, returnType)

	if children[3] != nil {
		printLine(w, "} else {", indent)
		render(w, children[3], functionName, indent+1, returnType)
	}

	printLine(w, "}", indent)
}
