package server

import (
	"gopkg.in/yaml.v2"
)

type Config struct {
	Listen    string
	Cassandra struct {
		Keyspace  string
		Addresses []string
	}
	Sources []struct {
		Name        string
		Timekey     string
		Consistency string
	}
}

func NewConfigFromYAML(input []byte) (*Config, error) {
	c := &Config{}
	return c, yaml.Unmarshal(input, c)
}
