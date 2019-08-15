package main

import (
	"regexp"
	"github.com/spf13/viper"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	Port     int
	Hostname string
	Chunk    int
	Interval int
	Tcponly  bool
	SynRst   bool
	Monitor struct {
		Enabled bool
		Listen  string
		Timeout int
	}
	Http     struct {
		Method   string
		Query    string
		Version  string
		Validate struct {
			Enabled bool
			Status  int
			Body    string
			regex   *regexp.Regexp
		}
	}
	Tries    struct {
		Up      int
		Down    int
		History int
	}
	Timeout  struct {
		Connect int
		Check   int
	}
}

func array_contains(haystack []string, needle string) (bool) {
	for _, v := range haystack {
		if v == needle {
			return true
		}
	}

	return false
}

func LoadConfig() {
	viper.SetConfigFile("config.yaml")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Could not read config.yaml: %s ", err)
	}

	viper.Unmarshal(&config)

	if !array_contains([]string {"1.0", "1.1"}, config.Http.Version) {
		log.Fatalf("Allowed HTTP versions are 1.0 and 1.1. Given %s", config.Http.Version)
	}

	if config.Tries.History < config.Tries.Up || config.Tries.History < config.Tries.Down {
		log.Fatalf("Tries history must be greater or equal to Up tries and Down tries.")
	}

	config.Http.Validate.regex = regexp.MustCompile(config.Http.Validate.Body)
}
