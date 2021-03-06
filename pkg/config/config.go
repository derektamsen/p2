// Package config provides convenience facilities for Golang-based pods to read their
// configuration files provided either by the environment or a custom path.
package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"sort"

	"gopkg.in/yaml.v2"
)

type Config struct {
	unpacked map[interface{}]interface{}
}

func LoadFromUnpacked(unpacked map[interface{}]interface{}) *Config {
	return &Config{unpacked}
}

func LoadFromEnvironment() (*Config, error) {
	env := os.Getenv("CONFIG_PATH")
	if env == "" {
		return nil, errors.New("No value was found for the environment variable CONFIG_PATH")
	}
	return LoadConfigFile(env)
}

func LoadConfigFile(filepath string) (*Config, error) {
	config := &Config{}
	contents, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	if err = yaml.Unmarshal(contents, &config.unpacked); err != nil {
		return nil, err
	}
	return config, nil
}

func (c *Config) ReadString(key string) (string, error) {
	readVal := c.Read(key)
	if readVal == nil {
		return "", nil
	}
	strVal, ok := readVal.(string)
	if !ok {
		return "", fmt.Errorf("%s is not a string value", key)
	}
	return strVal, nil
}

func (c *Config) ReadBool(key string) (bool, error) {
	readVal := c.Read(key)
	if readVal == nil {
		return false, nil
	}
	boolVal, ok := readVal.(bool)
	if !ok {
		return false, fmt.Errorf("%s is not a bool value", key)
	}
	return boolVal, nil
}

func (c *Config) Read(key string) interface{} {
	return c.unpacked[key]
}

func (c *Config) ReadMap(key string) (*Config, error) {
	readVal := c.Read(key)
	if readVal == nil {
		return &Config{make(map[interface{}]interface{})}, nil
	}
	mapVal, ok := readVal.(map[interface{}]interface{})
	if !ok {
		return &Config{make(map[interface{}]interface{})}, fmt.Errorf("%s is not a map", key)
	}
	return &Config{mapVal}, nil
}

func (c *Config) ReadStringSlice(key string) ([]string, error) {
	readVal := c.Read(key)
	retV := []string{}
	if readVal == nil {
		return retV, nil
	}
	slice, ok := readVal.([]interface{})
	if !ok {
		return retV, fmt.Errorf("%s is not a string slice", key)
	}
	for _, intf := range slice {
		strVal, ok := intf.(string)
		if !ok {
			return []string{}, fmt.Errorf("%v is not a string", intf)
		}
		retV = append(retV, strVal)
	}
	return retV, nil
}

func (c *Config) Keys() []string {
	keys := []string{}
	for intf := range c.unpacked {
		strVal, ok := intf.(string)
		if ok {
			keys = append(keys, strVal)
		}
	}
	sort.Strings(keys)
	return keys
}
