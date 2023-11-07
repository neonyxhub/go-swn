package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	neo_swn "go.neonyx.io/go-swn/pkg/swn"
	neo_cfg "go.neonyx.io/go-swn/pkg/swn/config"
)

func main() {
	cfgPath := flag.String("config", "config.yaml", "path to the config file")
	debug := flag.Bool("debug", false, "debug mode and save peer info to debug.yml")
	flag.Parse()

	cfg, err := neo_cfg.ReadConfigYaml(*cfgPath)
	if err != nil {
		log.Fatalf("failed to read and parse config yaml file: %v", err)
	}

	cfg.Debug = *debug

	swn, err := neo_swn.New(cfg)
	if err != nil {
		log.Fatalf("failed to spawn a swn: %v", err)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	err = swn.Run()
	if err != nil {
		log.Fatalf("failed to run swn: %v", err)
	}

	<-sigCh
	err = swn.Stop()
	if err != nil {
		log.Fatalf("failed to stop swn: %v", err)
	}
}
