package main

import (
	"log"
	"os"
	"path/filepath"
	"strconv"

	"gopkg.in/yaml.v3"
)

type ConfigLocalAddress struct {
	Name    string `yaml:"name"`
	Address string `yaml:"address"`
}

type Config struct {
	ListenAddr       string               `yaml:"listen"`
	UpstreamDNS      []string             `yaml:"upstream"`
	BlacklistSources []string             `yaml:"blacklist"`
	Whitelist        []string             `yaml:"whitelist"`
	LocalAddresses   []ConfigLocalAddress `yaml:"local"`
}

var ConfigInstance *Config = &Config{}

func GetConfig() *Config {
	return ConfigInstance
}

func (c *Config) ReadConfig() {
	configPath, err := os.Getwd()
	if (err != nil) || (configPath == "") {
		log.Fatalln("could neither get system config dir nor current working dir")
	}
	configPath = filepath.Join(configPath, "config.yaml")
	data, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("could not read config yaml from %s\n", configPath)
	}
	if err := yaml.Unmarshal(data, &c); err != nil {
		log.Fatalf("could not parse config yaml: %s\n", err.Error())
	}
	c.ReadEnv()
}

func (c *Config) ReadEnv() {
	listenAddr := c.getEnv("LISTEN_ADDR", "")
	if listenAddr != "" {
		c.ListenAddr = listenAddr
	}
	for i := 1; i <= 10; i++ {
		server := c.getEnv("UPSTREAM_DNS_"+strconv.Itoa(i), "")
		if server != "" {
			c.UpstreamDNS = append(c.UpstreamDNS, server)
		}
	}
}

func (c *Config) Print() {
	s, _ := yaml.Marshal(c)
	log.Println("Using config:\n" + string(s))
}

func (c *Config) getEnv(key, defaultValue string) string {
	res := os.Getenv(key)
	if res == "" {
		return defaultValue
	}
	return res
}
