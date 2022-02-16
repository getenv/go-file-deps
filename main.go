package main

import (
	"flag"
	"fmt"
	"sort"
	"strings"

	"golang.org/x/tools/go/packages"
)

// usage: go-file-deps [-r rule] $package
// print the list of files required to build $package

func main() {
	rule := flag.String("r", "",
		"Generate a Makefile rule with dependencies")

	flag.Parse()

	c := packages.Config{
		Mode: packages.NeedDeps | packages.NeedImports | packages.NeedModule |
			packages.NeedFiles | packages.NeedCompiledGoFiles,
	}

	pkgs, err := packages.Load(&c, flag.Arg(0))
	if err != nil {
		panic(err)
	}

	seen := map[*packages.Package]struct{}{}
	walk(seen, pkgs[0])

	fset := map[string]struct{}{}
	for p := range seen {
		if p.Module == nil {
			continue
		}

		for _, f := range p.CompiledGoFiles {
			fset[f] = struct{}{}
		}

		for _, f := range p.GoFiles {
			fset[f] = struct{}{}
		}

		for _, f := range p.OtherFiles {
			fset[f] = struct{}{}
		}
	}

	var files []string

	for f := range fset {
		files = append(files, f)
	}

	sort.Strings(files)

	switch {
	case len(*rule) > 0:
		fmt.Printf("%s: %s\n", *rule, strings.Join(files, " "))

	default:
		for _, f := range files {
			fmt.Println(f)
		}
	}
}

func walk(seen map[*packages.Package]struct{}, p *packages.Package) {
	if _, ok := seen[p]; ok {
		return
	}

	seen[p] = struct{}{}
	for _, pp := range p.Imports {
		walk(seen, pp)
	}
}
