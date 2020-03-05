package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/treezor-bank/go-xml/xmltree"
	"github.com/treezor-bank/go-xml/xsd"
)

var (
	targetNS = flag.String("ns", "", "Namespace of schema to print")
)

func main() {
	log.SetFlags(0)
	flag.Parse()

	if flag.NArg() < 1 {
		log.Fatalf("Usage: %s [-ns xmlns] file.xsd ...", os.Args[0])
	}

	docs := make([][]byte, 0, flag.NArg())

	for i, filename := range flag.Args() {
		data, err := ioutil.ReadFile(filename)
		if err != nil {
			log.Fatal(err)
		}
		docs[i] = data
	}

	filterSchema := make(map[string]struct{})
	for _, doc := range xsd.StandardSchema {
		root, err := xmltree.Parse(doc)
		if err != nil {
			// should never happen
			log.Fatal(err)
		}
		filterSchema[root.Attr("", "targetNamespace")] = struct{}{}
	}

	norm, err := xsd.Normalize(docs...)
	if err != nil {
		log.Fatal(err)
	}

	selected := make([]*xmltree.Element, 0, len(norm))
	for _, root := range norm {
		tns := root.Attr("", "targetNamespace")
		if _, ok := filterSchema[tns]; ok {
			continue
		}
		if *targetNS != "" && *targetNS == tns {
			selected = append(selected, root)
		}
	}

	for _, root := range selected {
		fmt.Printf("%s\n", xmltree.MarshalIndent(root, "", "  "))
	}
}
