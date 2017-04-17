package main

type BreakStmt struct {
	Address  string
	Position string
	Children []interface{}
}

func parseBreakStmt(line string) *BreakStmt {
	groups := groupsFromRegex(
		"<(?P<position>.*)>",
		line,
	)

	return &BreakStmt{
		Address:  groups["address"],
		Position: groups["position"],
		Children: []interface{}{},
	}
}

func (n *BreakStmt) render() []string {
	return []string{"break", ""}
}
