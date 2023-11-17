package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

const (
	EVENTBUS_EVENTIO = "eventio"
	EVENTBUS_GRPC    = "grpc"
	EVENTBUS_NATS    = "nats"
)

// Config structure for user-defined config.yaml
type Config struct {
	GrpcServer struct {
		Addr string `yaml:"addr"`
	} `yaml:"grpc_server"`
	Nats struct {
		Url string `yaml:"url"`
	} `yaml:"nats"`
	EventBus      string        `yaml:"eventbus"`
	EventBusTimer time.Duration `yaml:"eventbus_timer"`
	DataStore     struct {
		Path string `yaml:"path"`
	} `yaml:"datastore"`
	P2p struct {
		Multiaddr string `yaml:"multiaddr"`
		ConnLimit []int  `yaml:"conn_limit"`
	} `yaml:"p2p"`
	Log struct {
		Dev      bool     `yaml:"dev"`
		OutPaths []string `yaml:"out_paths"`
		ErrPaths []string `yaml:"err_paths"`
	} `yaml:"log"`
	Debug bool
}

func ParseConfig(data *[]byte) (*Config, error) {
	var config Config
	config.P2p.ConnLimit = make([]int, 2)

	err := yaml.Unmarshal(*data, &config)
	if err != nil {
		return nil, err
	}

	if len(config.P2p.ConnLimit) != 2 {
		config.P2p.ConnLimit[0] = 100
		config.P2p.ConnLimit[1] = 400
	}

	if config.EventBus == "" {
		config.EventBus = EVENTBUS_EVENTIO
	}

	return &config, nil
}

func ReadConfigYaml(cfgPath string) (*Config, error) {
	data, err := os.ReadFile(cfgPath)
	if err != nil {
		return nil, err
	}
	return ParseConfig(&data)
}
