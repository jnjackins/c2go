package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
)

var (
	printAst = flag.Bool("print-ast", false, "Print AST before translated Go code.")
)

func readAST(data []byte) []string {
	uncolored := regexp.MustCompile(`\x1b\[[\d;]+m`).ReplaceAll(data, []byte{})
	return strings.Split(string(uncolored), "\n")
}

func convertLinesToNodes(lines []string) []interface{} {
	nodes := []interface{}{}
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		// It is tempting to discard null AST nodes, but these may
		// have semantic importance: for example, they represent omitted
		// for-loop conditions, as in for(;;).
		line = strings.Replace(line, "<<<NULL>>>", "NullStmt", 1)

		indentAndType := regexp.MustCompile("^([|\\- `]*)(\\w+)").FindStringSubmatch(line)
		if len(indentAndType) == 0 {
			panic(fmt.Sprintf("Cannot understand line '%s'", line))
		}

		offset := len(indentAndType[1])
		node := Parse(line[offset:])

		indentLevel := len(indentAndType[1]) / 2
		nodes = append(nodes, []interface{}{indentLevel, node})
	}

	return nodes
}

// buildTree convert an array of nodes, each prefixed with a depth into a tree.
func buildTree(nodes []interface{}, depth int) []interface{} {
	if len(nodes) == 0 {
		return []interface{}{}
	}

	// Split the list into sections, treat each section as a a tree with its own root.
	sections := [][]interface{}{}
	for _, node := range nodes {
		if node.([]interface{})[0] == depth {
			sections = append(sections, []interface{}{node})
		} else {
			sections[len(sections)-1] = append(sections[len(sections)-1], node)
		}
	}

	results := []interface{}{}
	for _, section := range sections {
		slice := []interface{}{}
		for _, n := range section {
			if n.([]interface{})[0].(int) > depth {
				slice = append(slice, n)
			}
		}

		children := buildTree(slice, depth+1)
		result := section[0].([]interface{})[1]

		if len(children) > 0 {
			c := reflect.ValueOf(result).Elem().FieldByName("Children")
			slice := reflect.AppendSlice(c, reflect.ValueOf(children))
			c.Set(slice)
		}

		results = append(results, result)
	}

	return results
}

func ToJSON(tree []interface{}) []map[string]interface{} {
	r := make([]map[string]interface{}, len(tree))

	for j, n := range tree {
		rn := reflect.ValueOf(n).Elem()
		r[j] = make(map[string]interface{})
		r[j]["node"] = rn.Type().Name()

		for i := 0; i < rn.NumField(); i++ {
			name := strings.ToLower(rn.Type().Field(i).Name)
			value := rn.Field(i).Interface()

			if name == "children" {
				v := value.([]interface{})

				if len(v) == 0 {
					continue
				}

				value = ToJSON(v)
			}

			r[j][name] = value
		}
	}

	return r
}

func translate(path string) string {

	// Preprocess

	tmp, err := ioutil.TempDir("", "c2go")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tmp)

	pp, err := os.Create(filepath.Join(tmp, "pp.c"))
	if err != nil {
		log.Fatal(err)
	}
	defer pp.Close()

	args := strings.Fields(os.Getenv("CFLAGS"))
	args = append(args, "-E", path)
	ppCmd := exec.Command("clang", args...)
	ppCmd.Stdout = pp

	var errBuf bytes.Buffer
	ppCmd.Stderr = &errBuf

	if err := ppCmd.Run(); err != nil {
		io.Copy(os.Stderr, &errBuf)
		log.Fatal(err)
	}
	pp.Close()

	// Generate AST

	args = strings.Fields(os.Getenv("CFLAGS"))
	args = append(args, "-Xclang", "-ast-dump", "-fsyntax-only", pp.Name())
	astCmd := exec.Command("clang", args...)
	var astBuf bytes.Buffer
	astCmd.Stdout = &astBuf
	errBuf.Reset()
	astCmd.Stderr = &errBuf

	if err := astCmd.Run(); err != nil {
		io.Copy(os.Stderr, &errBuf)
		log.Fatal(err)
	}

	lines := readAST(astBuf.Bytes())
	if *printAst {
		for _, l := range lines {
			fmt.Println(l)
		}
		fmt.Println()
	}
	nodes := convertLinesToNodes(lines)
	tree := buildTree(nodes, 0)

	var goBuf bytes.Buffer

	// Generate Go code from AST

	fmt.Fprint(&goBuf, "package main\n\nimport (\n")
	for _, importName := range Imports {
		fmt.Fprintf(&goBuf, "\t\"%s\"\n", importName)
	}
	fmt.Fprintf(&goBuf, ")\n\n")

	render(&goBuf, tree[0], "", 0, "")

	return goBuf.String()
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("c2go: ")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s <file.c>\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	fmt.Println(translate(flag.Arg(0)))
}
