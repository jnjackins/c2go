package main

type DeclRefExpr struct {
	Address  string
	Position string
	Type     string
	Lvalue   bool
	For      string
	Address2 string
	Name     string
	Type2    string
	Children []interface{}
}

func parseDeclRefExpr(line string) *DeclRefExpr {
	groups := groupsFromRegex(
		`<(?P<position>.*)>
		 '(?P<type>.*?)'
		.*?
		(?P<lvalue> lvalue)?
		 (?P<for>\w+)
		 (?P<address2>[0-9a-fx]+)
		 '(?P<name>.*?)'
		 '(?P<type2>.*?)'`,
		line,
	)

	return &DeclRefExpr{
		Address:  groups["address"],
		Position: groups["position"],
		Type:     groups["type"],
		Lvalue:   len(groups["lvalue"]) > 0,
		For:      groups["for"],
		Address2: groups["address2"],
		Name:     groups["name"],
		Type2:    groups["type2"],
		Children: []interface{}{},
	}
}

func (n *DeclRefExpr) render() []string {
	name := n.Name

	if name == "argc" {
		name = "len(os.Args)"
		addImport("os")
	} else if name == "argv" {
		name = "os.Args"
		addImport("os")
	}

	return []string{name, n.Type}
}
