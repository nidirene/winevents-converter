package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var (
	version string
)

type Options struct {
	Verbose       bool `yaml:"verbose"`
	Logger        *log.Logger
	LogFile       string `yaml:"log-file"`
	EventFilePath string `yaml:"ev-file"`
	version       bool
}

func init() {
	if version == "" {
		version = "unknown"
	}
}

func NewOptions() *Options {
	return &Options{
		Verbose:       false,
		version:       false,
		Logger:        log.New(os.Stderr, "[nfcapture] ", log.Ldate|log.Ltime),
		LogFile:       "",
		EventFilePath: "tests",
	}
}

func GetOptions() *Options {
	opts := NewOptions()

	if opts.Verbose {
		opts.Logger.Printf("the full logging enabled")
		opts.Logger.SetFlags(log.LstdFlags | log.Lshortfile)
	}

	if opts.LogFile != "" {
		f, err := os.OpenFile(opts.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			opts.Logger.Println(err)
		} else {
			opts.Logger.SetOutput(f)
		}
	}
	return opts
}

func (opts *Options) wincFlagSet() {
	// global options
	flag.BoolVar(&opts.Verbose, "verbose", opts.Verbose, "enable/disable verbose logging")
	flag.BoolVar(&opts.version, "version", opts.version, "show version")
	flag.StringVar(&opts.LogFile, "log-file", opts.LogFile, "log file name")
	flag.StringVar(&opts.EventFilePath, "ev-file", opts.EventFilePath, "event file name")

	flag.Usage = func() {
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, `
Example:
# enable verbose logging
  winc -verbose=true`)

	}
}
