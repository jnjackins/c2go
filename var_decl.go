package main

import (
	"fmt"
	"io"
	"strings"
)

type VarDecl struct {
	Address   string
	Position  string
	Position2 string
	Name      string
	Type      string
	Type2     string
	IsExtern  bool
	IsUsed    bool
	IsCInit   bool
	Children  []interface{}
}

func parseVarDecl(line string) *VarDecl {
	groups := groupsFromRegex(
		`<(?P<position>.*)>(?P<position2> .+:\d+)?
		(?P<used> used)?
		(?P<name> \w+)?
		 '(?P<type>.+?)'
		(?P<type2>:'.*?')?
		(?P<extern> extern)?
		(?P<cinit> cinit)?`,
		line,
	)

	type2 := groups["type2"]
	if type2 != "" {
		type2 = type2[2 : len(type2)-1]
	}

	return &VarDecl{
		Address:   groups["address"],
		Position:  groups["position"],
		Position2: strings.TrimSpace(groups["position2"]),
		Name:      strings.TrimSpace(groups["name"]),
		Type:      groups["type"],
		Type2:     type2,
		IsExtern:  len(groups["extern"]) > 0,
		IsUsed:    len(groups["used"]) > 0,
		IsCInit:   len(groups["cinit"]) > 0,
		Children:  []interface{}{},
	}
}

func (n *VarDecl) render() []string {
	theType := resolveType(n.Type)
	name := n.Name

	// Go does not allow the name of a variable to be called "type".
	// For the moment I will rename this to avoid the error.
	if name == "type" {
		name = "type_"
	}

	suffix := ""
	if len(n.Children) > 0 {
		children := n.Children
		suffix = fmt.Sprintf(" = %s", renderExpression(children[0])[0])
	}

	if suffix == " = (0)" {
		suffix = " = nil"
	}

	return []string{fmt.Sprintf("var %s %s%s", name, theType, suffix), "unknown3"}
}

func (n *VarDecl) renderLine(w io.Writer, functionName string, indent int, returnType string) {
}
