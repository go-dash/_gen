package main

import (
	"fmt"
	"os"
	"strings"
	"io/ioutil"
	"regexp"
	"go/build"
	"errors"
)

type config struct {
}

var conf = config{}

func main() {
	args := os.Args[1:]
	for _, arg := range args {
		if isFlag(arg) {
			handleFlag(arg)
		}
	}
	imports, importLocations, err := findAllGodashImports("./")
	if err != nil {
		fmt.Println("ERROR:", err.Error())
		os.Exit(1)
	}
	if len(imports) == 0 {
		fmt.Println("No go-dash imports found in project")
		fmt.Println(" TIP 1: Make sure you're running the generator from project source root")
		fmt.Println(" TIP 2: Use go-dash for \"string\" type by importing \"github.com/go-dash/slice/_string\"")
	}
	for _, imp := range imports {
		err := generateGodashCodeForType(imp[1:], importLocations[imp])
		if err != nil {
			fmt.Println("ERROR:", err.Error())
			os.Exit(1)
		}
		fmt.Println("Generated go-dash", imp)
	}
}

func isFlag(arg string) bool {
	return strings.HasPrefix(arg, "-")
}

func handleFlag(flag string) {
	switch flag {
	case "--version":
		displayVersion()
	case "-v":
		displayVersion()
	case "--help":
		displayUsage()
	case "-h":
		displayUsage()
	}
}

func displayVersion() {
	fmt.Println("_gen 0.0.1 (go-dash code generator)")
	os.Exit(0)
}

func displayUsage() {
	fmt.Println("Usage: _gen [OPTION]")
	fmt.Println("Parse project go files and generate github.com/go-dash implementations required by the project")
	fmt.Println("  -v, --version    Show version info and exit.")
	fmt.Println("  -h, --help       Show this usage text and exit.")
	os.Exit(0)
}

func findAllGodashImports(dir string) ([]string, map[string]string, error) {
	res := []string{}
	location := make(map[string]string)
	files, err := glob("./", ".go")
	if err != nil {
		return nil, nil, err
	}
	if len(files) == 0 {
		return nil, nil, errors.New("no go files in " + dir + ", make sure to run in project root")
	}
	re, err := regexp.Compile(`(?m)\Q"github.com/go-dash/slice/\E(_\w+)\Q"\E\s*(\Q//\E\s*.+)?`)
	if err != nil {
		return nil, nil, err
	}
	for _, file := range files {
		content, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, nil, err
		}
		matches := re.FindAllSubmatch(content, -1)
		for _, match := range matches {
			imp := string(match[1])
			res = append(res, imp)
			if len(match) > 2 {
				location[imp] = trimLeftComment(string(match[2]))
			}
		}
	}
	return unique(res), location, nil
}

func getGodashCodePath() (string, error) {
	pkg, err := build.Import("github.com/go-dash/slice", "", build.FindOnly)
	if err != nil {
		return "", err
	}
	return pkg.Dir, nil
}

func generateGodashCodeForType(imp string, importLocation string) error {
	godashPath, err := getGodashCodePath()
	if err != nil {
		return err
	}
	impPath := godashPath + "/_" + imp
	tpl := "_SimpleType"
	if !isSimpleType(imp) {
		tpl = "_ComplexType"
	}
	tplPath := godashPath + "/templates/" + tpl
	emptyDir(impPath)
	copyDir(tplPath, impPath)
	files, err := glob(impPath, ".go")
	if err != nil {
		return err
	}
	for _, file := range files {
		err := searchAndReplaceInFile(file, tpl, imp)
		if err != nil {
			return err
		}
		if tpl == "_ComplexType" {
			err := searchAndReplaceInFile(file, "_ImportLocation", importLocation)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func isSimpleType(imp string) bool {
	simples := map[string]bool{
		"string":     true,
		"int":        true,
		"uint":       true,
		"bool":       true,
		"uint8":      true,
		"uint16":     true,
		"uint32":     true,
		"uint64":     true,
		"int8":       true,
		"int16":      true,
		"int32":      true,
		"int64":      true,
		"float32":    true,
		"float64":    true,
		"complex64":  true,
		"complex128": true,
		"uintptr":    true,
		"byte":       true,
		"rune":       true,
	}
	_, found := simples[imp]
	return found
}