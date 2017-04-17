package main

type CStyleCastExpr struct {
	Address  string
	Position string
	Type     string
	Kind     string
	Children []interface{}
}

func parseCStyleCastExpr(line string) *CStyleCastExpr {
	groups := groupsFromRegex(
		"<(?P<position>.*)> '(?P<type>.*?)' <(?P<kind>.*)>",
		line,
	)

	return &CStyleCastExpr{
		Address:  groups["address"],
		Position: groups["position"],
		Type:     groups["type"],
		Kind:     groups["kind"],
		Children: []interface{}{},
	}
}

func (n *CStyleCastExpr) render() []string {
	children := n.Children
	return renderExpression(children[0])
}
