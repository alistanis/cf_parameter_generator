package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
)

var (
	inFile          string
	outFile         string
	min             bool
	numIndentSpaces int
)

func init() {
	flag.StringVar(&inFile, "f", "", "The file to read from to generate parameters.")
	flag.StringVar(&outFile, "o", "", "Optional: Specify a file name to write out parameters.")
	flag.BoolVar(&min, "min", false, "If given, will generate minified output.")
	flag.IntVar(&numIndentSpaces, "spaces", 2, "The number of spaces used to indent the file if not generating minified output.")

}

func main() {
	flag.Parse()

	if inFile == "" {
		exitError(errors.New("Missing required argument -f"))
	}

	m := make(map[string]interface{})
	data, err := ioutil.ReadFile(inFile)
	if err != nil {
		exitError(err)
	}
	err = json.Unmarshal(data, &m)
	if err != nil {
		exitError(err)
	}

	if _, ok := m["Parameters"]; !ok {
		exitError(errors.New("Parameters not found in file."))
	}

	params := m["Parameters"].(map[string]interface{})

	pl := make(paramList, 0)

	for k := range params {
		p := parameter{ParameterKey: k}
		pl = append(pl, p)
	}

	sort.Sort(pl)
	data, err = marshal(pl)
	if err != nil {
		exitError(err)
	}
	if outFile != "" {
		err = ioutil.WriteFile(outFile, data, 0644)
		if err != nil {
			exitError(err)
		}
	} else {
		fmt.Print(string(data))
	}

}

func marshal(i interface{}) ([]byte, error) {
	if min {
		return json.Marshal(i)
	}
	return json.MarshalIndent(i, "", strings.Repeat(" ", numIndentSpaces))
}

type parameter struct {
	ParameterKey   string
	ParameterValue string
}

type paramList []parameter

func (p paramList) Less(i, j int) bool {
	return p[i].ParameterKey < p[j].ParameterKey
}

func (p paramList) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p paramList) Len() int {
	return len(p)
}


func exitError(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(-1)
}
