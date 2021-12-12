package config

import (
	"encoding/json"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/url"
)

type Config struct {
	RawWebUrl    string   `yaml:"web_url"`
	WebUrl       *url.URL `yaml:"-"`
	BrokerUrl    string   `yaml:"broker_url"`
	DBDriver     string   `yaml:"db_driver"`
	DBConnection string   `yaml:"db_connection"`
	MemcachedUrl string   `yaml:"memcached_url"`
	ZipArchPath  string   `yaml:"zip_arch_path"`
	TestsPath    string   `yaml:"tests_path"`
}

func ConfigFromFile(file string) (*Config, error) {
	configBytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var c *Config
	if err := yaml.Unmarshal(configBytes, &c); err != nil {
		return nil, err
	}
	webUrl, err := url.Parse(c.RawWebUrl)
	if err != nil {
		return nil, err
	}
	c.WebUrl = webUrl
	return c, nil
}

func (c *Config) String() string {
	data, err := json.Marshal(c)
	if err != nil {
		return ""
	}
	return string(data)
}
