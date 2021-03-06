package main

type ElaboratedType struct {
	Address  string
	Type     string
	Tags     string
	Children []interface{}
}

func parseElaboratedType(line string) *ElaboratedType {
	groups := groupsFromRegex(
		"'(?P<type>.*)' (?P<tags>.+)",
		line,
	)

	return &ElaboratedType{
		Address:  groups["address"],
		Type:     groups["type"],
		Tags:     groups["tags"],
		Children: []interface{}{},
	}
}
