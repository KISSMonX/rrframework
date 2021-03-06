package rrconfig

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
)

type JsonConfig struct {
	m     map[string]interface{}
	rb    []byte
}

func LoadJsonConfigFromFile(path string) (*JsonConfig, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return LoadJsonConfigFromBytes(b)
}

func LoadJsonConfigFromBytes(b []byte) (*JsonConfig, error) {
	var jm map[string]interface{}
	if err := json.Unmarshal(b, &jm); err != nil {
		return nil, err
	}
	s := &JsonConfig{
		m:     jm,
		rb:    b,
	}
	return s, nil
}

func (s *JsonConfig) Dump() (string, error) {
	var rj bytes.Buffer
	if err := json.Indent(&rj, s.rb, "", "\t"); err != nil {
		return "", err
	}
	return string(rj.Bytes()), nil
}

// Get("a.b.c")
func (s *JsonConfig) Get(key string) (interface{}, error) {
	nodes := strings.Split(key, ".")
	m := s.m
	for i := 0; i < len(nodes); i++ {
		if v, ok := m[nodes[i]]; ok {
			if vv, okk := v.(map[string]interface{}); okk {
				// not end
				m = vv
			} else {
				return v, nil
			}
		} else {
			return nil, fmt.Errorf("no value for key %s", key)
		}
	}
	return m, nil
}

func (s *JsonConfig) GetStringSlice(key string) ([]string, error) {
	empty := []string{}
	f, err := s.Get(key)
	if err != nil {
		return empty, err
	}
	if _, ok := f.([]interface{}); !ok {
		return empty, fmt.Errorf("value for key %s is not slice", key)
	}
	sf := f.([]interface{})
	ss := make([]string, len(sf))
	for i, v := range sf {
		if vv, ok := v.(string); ok {
			ss[i] = vv
		} else {
			return empty, fmt.Errorf("%s[%d] is not a string", key, i)
		}
	}
	return ss, nil
}

func (s *JsonConfig) GetString(key string) (string, error) {
	f, err := s.Get(key)
	if err != nil {
		return "", err
	}
	if _, ok := f.(string); !ok {
		return "", fmt.Errorf("value for key %s is not string", key)
	}
	return f.(string), nil
}

func (s *JsonConfig) GetInt(key string) (int, error) {
	f, err := s.Get(key)
	if err != nil {
		return 0, err
	}
	if _, ok := f.(float64); !ok {
		return 0, fmt.Errorf("value for key %s is not int", key)
	}
	return int(f.(float64)), nil
}

func (s *JsonConfig) GetFloat64(key string) (float64, error) {
	f, err := s.Get(key)
	if err != nil {
		return 0.0, err
	}
	if _, ok := f.(float64); !ok {
		return 0.0, fmt.Errorf("value for key %s is not float64", key)
	}
	return f.(float64), nil
}

func (s *JsonConfig) GetInterfaceSlice(key string) ([]interface{}, error) {
	f, err := s.Get(key)
	if err != nil {
		return nil, err
	}
	if _, ok := f.([]interface{}); !ok {
		return nil, fmt.Errorf("value for key %s is not []interface{}", key)
	}
	return f.([]interface{}), nil
}

