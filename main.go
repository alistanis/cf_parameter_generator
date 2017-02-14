package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
)

var (
	inFile  string
	outFile string
)

func init() {
	flag.StringVar(&inFile, "f", "", "The file to read from to generate parameters")
	flag.StringVar(&outFile, "o", "", "Optional: Specify a file name to write out parameters")
}

func main() {
	flag.Parse()

	if inFile == "" {
		fmt.Fprintln(os.Stderr, "Missing required argument -f")
	}

	m := make(map[string]interface{})
	data, err := ioutil.ReadFile(inFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}
	err = json.Unmarshal(data, &m)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}

	if _, ok := m["Parameters"]; !ok {
		fmt.Fprintln(os.Stderr, "Parameters not found in file")
		os.Exit(-1)
	}

	params := m["Parameters"].(map[string]interface{})

	keys := make([]string, 0)

	for k := range params {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	s := "[\n"

	for i, k := range keys {
		s += fmt.Sprintf(`	{
		"ParameterKey": "%s",
		"ParameterValue": ""
	}`, k)
		if i != len(keys)-1 {
			s += `,
`
		}
	}

	s += "\n]"

	if outFile != "" {
		err = ioutil.WriteFile(outFile, []byte(data), 0644)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	} else {
		fmt.Print(s)
	}

}
