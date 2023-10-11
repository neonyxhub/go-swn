package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Config structure for user-defined config.yaml
type Config struct {
	GrpcServer struct {
		Addr string `yaml:"addr"`
	} `yaml:"grpc_server"`
	DataStore struct {
		Path string `yaml:"path"`
	} `yaml:"datastore"`
	P2p struct {
		Multiaddr   string `yaml:"multiaddr"`
		PrivKeyPath string `yaml:"privkey_filepath"`
	} `yaml:"p2p"`
	Log struct {
		Dev      bool     `yaml:"dev"`
		OutPaths []string `yaml:"out_paths"`
		ErrPaths []string `yaml:"err_paths"`
	} `yaml:"log"`
	Debug bool
}

func ParseConfig(data *[]byte) (*Config, error) {
	config := &Config{}
	err := yaml.Unmarshal(*data, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func ReadConfigYaml(cfgPath string) (*Config, error) {
	data, err := os.ReadFile(cfgPath)
	if err != nil {
		return nil, err
	}
	return ParseConfig(&data)
}
