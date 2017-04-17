package main

import (
	"fmt"
	"io"
	"log"
	"strings"
)

type expressionRenderer interface {
	// TODO: The two arguments returned are the rendered Go and the C type.
	// This should be made into an appropriate type.
	render() []string
}

type lineRenderer interface {
	renderLine(w io.Writer, functionName string, indent int, returnType string)
}

func printLine(w io.Writer, line string, indent int) {
	fmt.Fprintf(w, "%s%s\n", strings.Repeat("\t", indent), line)
}

func renderExpression(node interface{}) []string {
	if node == nil {
		return []string{""}
	}
	n, ok := node.(expressionRenderer)
	if !ok {
		log.Fatalf("not an expressionRenderer: %T", node)
	}

	return n.render()
}

func getFunctionParams(f *FunctionDecl) []*ParmVarDecl {
	r := []*ParmVarDecl{}
	for _, n := range f.Children {
		if v, ok := n.(*ParmVarDecl); ok {
			r = append(r, v)
		}
	}

	return r
}

func getFunctionReturnType(f string) string {
	// The type of the function will be the complete prototype, like:
	//
	//     __inline_isfinitef(float) int
	//
	// will have a type of:
	//
	//     int (float)
	//
	// The arguments will handle themselves, we only care about the
	// return type ('int' in this case)
	return strings.TrimSpace(strings.Split(f, "(")[0])
}

func render(w io.Writer, node interface{}, functionName string, indent int, returnType string) {
	if n, ok := node.(lineRenderer); ok {
		n.renderLine(w, functionName, indent, returnType)
		return
	}
	printLine(w, renderExpression(node)[0], indent)
}
