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

	"gopkg.in/yaml.v2"
)

var (
	ErrMissingInFile     = errors.New("Missing required argument -f")
	ErrMissingParameters = errors.New("Parameters not found in file")
)

type Config struct {
	InFile                       string
	OutFile                      string
	Minimize                     bool
	Indent                       int
	Overwrite                    bool
	RemoveOldParamsNotInTemplate bool
	OutYaml                      bool
	InYaml                       bool
}

func Run(c *Config) error {
	flag.Parse()

	if c.InFile == "" {
		return ErrMissingInFile
	}

	m := make(map[string]interface{})
	data, err := ioutil.ReadFile(c.InFile)
	if err != nil {
		return err
	}
	err = unmarshal(data, &m, c)
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

	pl := make(paramList, 0)

	for k := range params {
		p := parameter{ParameterKey: k}
		pl = append(pl, p)
	}

	sort.Sort(pl)
	data, err = marshal(pl, c)
	if err != nil {
		return err
	}
	if c.OutFile != "" {
		if !c.Overwrite {
			odata, err := ioutil.ReadFile(c.OutFile)
			if err == nil {
				// original param list
				opl := make(paramList, 0)
				err = unmarshal(odata, &opl, c)
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
				data, err = marshal(opl, c)
				if err != nil {
					return err
				}
			}
		}

		err = ioutil.WriteFile(c.OutFile, data, 0644)
		if err != nil {
			return err
		}
	} else {
		fmt.Print(string(data))
	}
	return nil
}

func marshal(i interface{}, c *Config) ([]byte, error) {
	if c.OutYaml {
		return yaml.Marshal(i)
	}
	if c.Minimize {
		return json.Marshal(i)
	}
	return json.MarshalIndent(i, "", strings.Repeat(" ", c.Indent))
}

func unmarshal(data []byte, i interface{}, c *Config) error {
	if c.InYaml {
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
