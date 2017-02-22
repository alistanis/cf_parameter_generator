package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"

	"reflect"

	"gopkg.in/yaml.v2"
)

var (
	inFile                       string
	outFile                      string
	min                          bool
	numIndentSpaces              int
	overwrite                    bool
	removeOldParamsNotInTemplate bool
	oyaml                        bool
	inyaml                       bool
)

func init() {
	flag.StringVar(&inFile, "f", "", "The file to read from to generate parameters.")
	flag.StringVar(&outFile, "o", "", "Optional: Specify a file name to write out parameters.")
	flag.BoolVar(&min, "min", false, "If given, will generate minified output.")
	flag.IntVar(&numIndentSpaces, "spaces", 2, "The number of spaces used to indent the file if not generating minified output.")
	flag.BoolVar(&overwrite, "overwrite", false, "By default, will update an existing parameters file with newly found parameters, but will not overwrite.")
	flag.BoolVar(&removeOldParamsNotInTemplate, "r", false, "Removes old entries from parameters found in old parameters files.")
	flag.BoolVar(&oyaml, "outyaml", false, "Will output in yaml instead of json.")
	flag.BoolVar(&inyaml, "inyaml", false, "Will expect input as yaml instead of json.")
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
	err = unmarshal(data, &m)
	if err != nil {
		exitError(err)
	}

	if _, ok := m["Parameters"]; !ok {
		exitError(errors.New("Parameters not found in file."))
	}
	params := make(map[string]interface{})
	if inyaml {
		recurseMapInterface(m["Parameters"].(map[interface{}]interface{}), params)
	} else {
		params = m["Parameters"].(map[string]interface{})
	}

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
		if !overwrite {
			odata, err := ioutil.ReadFile(outFile)
			if err == nil {
				// original param list
				opl := make(paramList, 0)
				err = unmarshal(odata, &opl)
				if err != nil {
					exitError(err)
				}
				//new values from template
				for _, nv := range pl {
					found := false

					// old values from params file
					for i, ov := range opl {
						if nv.ParameterKey == ov.ParameterKey {
							found = true
						}
						if i == len(opl)-1 && !found {
							opl = append(opl, nv)
						}
					}
				}

				if removeOldParamsNotInTemplate {
					//new values from template
				outer:
					for i := 0; i < len(opl); i++ {
						found := false
						// old values from params file
						for _, nv := range pl {
							if nv.ParameterKey == opl[i].ParameterKey {
								found = true
							}
						}

						if !found {
							fmt.Println("Removing value", opl[i])
							opl = append(opl[:i], opl[i+1:]...)
							i--
							continue outer
						}
					}
				}
				sort.Sort(opl)
				data, err = marshal(opl)
				if err != nil {
					exitError(err)
				}
			}
		}

		err = ioutil.WriteFile(outFile, data, 0644)
		if err != nil {
			exitError(err)
		}
	} else {
		fmt.Print(string(data))
	}

}

func marshal(i interface{}) ([]byte, error) {
	if oyaml {
		return yaml.Marshal(i)
	}
	if min {
		return json.Marshal(i)
	}
	return json.MarshalIndent(i, "", strings.Repeat(" ", numIndentSpaces))
}

func unmarshal(data []byte, i interface{}) error {
	if inyaml {
		return yaml.Unmarshal(data, i)
	}
	return json.Unmarshal(data, i)
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

func recurseMapInterface(m map[interface{}]interface{}, newMap map[string]interface{}) {
	nm := make(map[string]interface{})
	for k, v := range m {
		nk := ""
		switch t := k.(type) {
		case string:
			nk = t
		case int:
			nk = strconv.Itoa(t)
		case float64:
			nk = strconv.FormatFloat(t, 'E', -1, 64)
		case map[interface{}]interface{}:
			recurseMapInterface(t, newMap)
		}

		var nv interface{}
		switch t := v.(type) {
		case map[interface{}]interface{}:
			m := make(map[string]interface{})
			nm[nk] = m
			recurseMapInterface(t, m)
		case []interface{}:
			recurseArray(nk, t, nm)
		default:
			nv = t
		}
		if nv != nil {
			nm[nk] = nv
		}
	}

	for k, v := range nm {
		newMap[k] = v
	}

}

func recurseArray(k string, slc []interface{}, container interface{}) {
	nslc := make([]interface{}, 0)
	for _, i := range slc {
		var v interface{}
		switch i := i.(type) {
		case []interface{}:
			recurseArray(k, i, &nslc)
		case map[interface{}]interface{}:
			m := make(map[string]interface{})
			nslc = append(nslc, m)
			recurseMapInterface(i, m)
		default:
			v = i
		}
		if v != nil {
			nslc = append(nslc, v)
		}
	}

	rv := reflect.ValueOf(container)
	if rv.Kind() == reflect.Ptr {
		rv = reflect.Indirect(rv)
	}
	switch rv.Kind() {
	case reflect.Slice:
		*container.(*[]interface{}) = append(*container.(*[]interface{}), nslc)
	case reflect.Map:
		m := container.(map[string]interface{})
		m[k] = nslc
	}
}
