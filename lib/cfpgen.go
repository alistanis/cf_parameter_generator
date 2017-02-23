package cfpgen

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"

	"reflect"

	"bytes"

	"gopkg.in/yaml.v2"
)

var (
	ErrMissingInFile     = errors.New("Missing required argument -f")
	ErrMissingParameters = errors.New("Parameters not found in file")
)

// Config represents a config struct holding information about how to format output and where to write said output
type Config struct {
	InFile                       string
	OutFile                      string
	Minimize                     bool
	Indent                       int
	Overwrite                    bool
	RemoveOldParamsNotInTemplate bool
	OutYaml                      bool
	InYaml                       bool
	Verbose                      bool
}

// Generate generates a cloud formation parameters file template and writes either to a file or stdout
func Generate(c *Config) error {
	flag.Parse()

	if c.InFile == "" {
		return ErrMissingInFile
	}

	m := make(map[string]interface{})
	data, err := ioutil.ReadFile(c.InFile)
	if err != nil {
		return err
	}
	err = c.unmarshal(data, &m)
	if err != nil {
		return err
	}

	if _, ok := m["Parameters"]; !ok {
		return ErrMissingParameters
	}
	params := make(map[string]interface{})
	if c.InYaml {
		RecurseMapInterface(m["Parameters"].(map[interface{}]interface{}), params)
	} else {
		params = m["Parameters"].(map[string]interface{})
	}

	pl, err := c.GetParamList(params)
	if err != nil {
		return err
	}
	sort.Sort(pl)

	return c.ProcessPL(pl)
}

// GetParamList takes the given config and generates a paramList ([]parameter)
func (c *Config) GetParamList(params map[string]interface{}) (paramList, error) {
	pl := make(paramList, 0)

	for k := range params {
		m := make(map[string]interface{})
		v := params[k]
		switch val := v.(type) {
		case map[interface{}]interface{}:
			RecurseMapInterface(val, m)
		case map[string]interface{}:
			m = val
		}
		p := parameter{ParameterKey: k}
		for k, v := range m {
			switch k {
			case "Default":
				p.Default = v
			case "AllowedValues":
				v, ok := v.([]interface{})
				if ok {
					p.AllowedValues = v
				} else {
					return nil, errors.New("AllowedValues was not an []interface{}")
				}
			case "Description":
				if s, ok := v.(string); ok {
					p.Description = s
				} else {
					return nil, errors.New("Description was not a string")
				}
			case "Type":
				if s, ok := v.(string); ok {
					p.Type = s
				} else {
					return nil, errors.New("Type was not a string")
				}
			case "AllowedPattern":
				if s, ok := v.(string); ok {
					p.AllowedPattern = s
				} else {
					return nil, errors.New("AllowedPattern was not a string")
				}
			}
		}
		pl = append(pl, p)
	}

	for i := range pl {
		var s string
		s = fmt.Sprintf("Type: %s", pl[i].Type)
		if c.Verbose {
			if pl[i].Default != nil {
				s += fmt.Sprintf(", Default: %v", pl[i].Default)
			}
			if pl[i].AllowedValues != nil {
				s += fmt.Sprintf(", AllowedValues: %v", pl[i].AllowedValues)
			}
			if pl[i].AllowedPattern != "" {
				s += fmt.Sprintf(", AllowedPattern: %s", pl[i].AllowedPattern)
			}
			if pl[i].Description != "" {
				s += fmt.Sprintf(", Description: %s", pl[i].Description)
			}
		}
		pl[i].ParameterValue = s
	}
	return pl, nil
}

// ProcessPL processes a paramList and writes to a file or stdout
func (c *Config) ProcessPL(pl paramList) error {
	data, err := c.marshal(pl)
	if err != nil {
		return err
	}
	if c.OutFile != "" {
		if !c.Overwrite {
			odata, err := ioutil.ReadFile(c.OutFile)
			if err == nil {
				// if we have an empty file just write it out
				if len(odata) > 0 {
					// original param list
					opl := make(paramList, 0)
					err = c.unmarshal(odata, &opl)
					if err != nil {
						return err
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

					if c.RemoveOldParamsNotInTemplate {
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
					data, err = c.marshal(opl)
					if err != nil {
						return err
					}
				}
			}
		}
		data = unescapeBrackets(data)
		err := ioutil.WriteFile(c.OutFile, data, 0644)
		if err != nil {
			return err
		}
	} else {
		data = unescapeBrackets(data)
		fmt.Print(string(data))
	}
	return nil
}

// marshal marshals data in either yaml or json
func (c *Config) marshal(i interface{}) ([]byte, error) {
	if c.OutYaml {
		return yaml.Marshal(i)
	}
	if c.Minimize {
		return json.Marshal(i)
	}
	return json.MarshalIndent(i, "", strings.Repeat(" ", c.Indent))
}

// unmarshal unmarshals data into i
func (c *Config) unmarshal(data []byte, i interface{}) error {
	if c.InYaml {
		return yaml.Unmarshal(data, i)
	}
	return json.Unmarshal(data, i)
}

// unescapeBrackets replaces the unicode points \u003c and \u003e with plaintext brackets - this is designed to be human readable
func unescapeBrackets(data []byte) []byte {
	data = bytes.Replace(data, []byte(`\u003c`), []byte(`<`), -1)
	data = bytes.Replace(data, []byte(`\u003e`), []byte(`>`), -1)
	return data
}

// parameter represents a cloudformation parameter
type parameter struct {
	ParameterKey   string
	ParameterValue string
	Type           string        `json:"-"`
	Description    string        `json:"-"`
	AllowedValues  []interface{} `json:"-"`
	Default        interface{}   `json:"-"`
	AllowedPattern string        `json:"-"`
}

// paramList is a slice of type parameter
type paramList []parameter

// Less implements the sort interface
func (p paramList) Less(i, j int) bool {
	return p[i].ParameterKey < p[j].ParameterKey
}

// Swap implements the sort interface
func (p paramList) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

// Len implements the sort interface
func (p paramList) Len() int {
	return len(p)
}

// RecurseArray converts types to supported types. Specifically interfaces stored in slices or the special case map[interface{}]interface{}
// yaml supports keys of arbitrary types whereas json does not, so we do this conversion to maintain compatibility between types
func RecurseArray(k string, slc []interface{}, container interface{}) {
	nslc := make([]interface{}, 0)
	for _, i := range slc {
		var v interface{}
		switch i := i.(type) {
		case []interface{}:
			RecurseArray(k, i, &nslc)
		case map[interface{}]interface{}:
			m := make(map[string]interface{})
			nslc = append(nslc, m)
			RecurseMapInterface(i, m)
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

// RecurseMapInterface converts types to supported types. Specifically interfaces stored in slices or the special case map[interface{}]interface{}
// yaml supports keys of arbitrary types whereas json does not, so we do this conversion to maintain compatibility between types
func RecurseMapInterface(m map[interface{}]interface{}, newMap map[string]interface{}) {
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
			RecurseMapInterface(t, newMap)
		}

		var nv interface{}
		switch t := v.(type) {
		case map[interface{}]interface{}:
			m := make(map[string]interface{})
			nm[nk] = m
			RecurseMapInterface(t, m)
		case []interface{}:
			RecurseArray(nk, t, nm)
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
