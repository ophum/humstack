package main

import (
	"context"
	"flag"
	"net/url"
	"os"
	"os/signal"
	"syscall"

	"github.com/ophum/humstack/v1/pkg/agent"
	"github.com/ophum/humstack/v1/pkg/client"
	"gopkg.in/yaml.v2"
)

var _ yaml.Unmarshaler = &Config{}

type Config struct {
	APIEndpoint           url.URL `yaml:"apiEndpoint"`
	agent.DiskAgentConfig `yaml:",inline"`
}

func (c *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	d := struct {
		APIEndpoint           string `yaml:"apiEndpoint"`
		agent.DiskAgentConfig `yaml:",inline"`
	}{}
	if err := unmarshal(&d); err != nil {
		return err
	}
	endpoint, err := url.Parse(d.APIEndpoint)
	if err != nil {
		return err
	}
	*c = Config{
		APIEndpoint:     *endpoint,
		DiskAgentConfig: d.DiskAgentConfig,
	}
	return nil
}

var (
	config Config
)

func init() {
	configPath := "config.yaml"
	flag.StringVar(&configPath, "config", "config.yaml", "config path")
	flag.Parse()

	f, err := os.Open(configPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if err := yaml.NewDecoder(f).Decode(&config); err != nil {
		panic(err)
	}
}

func main() {
	diskClient := client.NewDiskClient(config.APIEndpoint)
	agent := agent.NewDiskAgent(diskClient, &config.DiskAgentConfig)
	ctx, cancel := context.WithCancel(context.Background())
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sig
		cancel()
	}()

	agent.Start(ctx)
}
