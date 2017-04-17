package main

import (
	"fmt"
	"io"
	"strings"
)

type RecordDecl struct {
	Address    string
	Position   string
	Prev       string
	Position2  string
	Kind       string
	Name       string
	Definition bool
	Children   []interface{}
}

func parseRecordDecl(line string) *RecordDecl {
	groups := groupsFromRegex(
		`(?P<prev>prev 0x[0-9a-f]+ )?
		<(?P<position>.*)>
		 (?P<position2>[^ ]+ )?
		(?P<kind>struct|union)
		(?P<name>.*)`,
		line,
	)

	definition := false
	name := strings.TrimSpace(groups["name"])
	if name == "definition" {
		name = ""
		definition = true
	}
	if strings.HasSuffix(name, " definition") {
		name = name[0 : len(name)-11]
		definition = true
	}

	return &RecordDecl{
		Address:    groups["address"],
		Position:   groups["position"],
		Prev:       groups["prev"],
		Position2:  strings.TrimSpace(groups["position2"]),
		Kind:       groups["kind"],
		Name:       name,
		Definition: definition,
		Children:   []interface{}{},
	}
}

func (n *RecordDecl) renderLine(w io.Writer, functionName string, indent int, returnType string) {
	name := strings.TrimSpace(n.Name)
	if name == "" || typeIsAlreadyDefined(name) {
		return
	}

	typeIsNowDefined(name)

	if n.Kind == "union" {
		return
	}

	printLine(w, fmt.Sprintf("type %s %s {", name, n.Kind), indent)
	if len(n.Children) > 0 {
		for _, c := range n.Children {
			render(w, c, functionName, indent+1, "")
		}
	}

	printLine(w, "}\n", indent)
}
