package config

import (
	"time"

	arg "github.com/alexflint/go-arg"
)

// AppConfig holds the configuration for the application
type AppConfig struct {
	PingDuration time.Duration
	NodePort     int
	InitialNodes []string
}

// Get will parse the cli args and create an AppConfig
func Get() *AppConfig {
	var args struct {
		Nodes    []string `arg:"positional"`
		Port     int      `arg:"-P,required"`
		PingTime int
	}
	args.PingTime = 30

	arg.MustParse(&args)

	return &AppConfig{
		PingDuration: time.Duration(args.PingTime) * time.Second,
		NodePort:     args.Port,
		InitialNodes: args.Nodes,
	}
}
