package main

import "io"

type TranslationUnitDecl struct {
	Address  string
	Children []interface{}
}

func parseTranslationUnitDecl(line string) *TranslationUnitDecl {
	groups := groupsFromRegex("", line)

	return &TranslationUnitDecl{
		Address:  groups["address"],
		Children: []interface{}{},
	}
}

func (n *TranslationUnitDecl) renderLine(w io.Writer, functionName string, indent int, returnType string) {
	for _, c := range n.Children {
		render(w, c, functionName, indent, returnType)
	}
}
