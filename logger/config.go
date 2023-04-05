package logger

import (
	"errors"
	"flag"
	"io"
	"log"
	"os"

	"github.com/graymeta/env"
)

// Config is used to parse CLI flags related to logging
type Config struct {
	flag *flag.FlagSet

	appName      string
	level        string
	processLevel string
}

// NewConfig prepares a new configuration for parsing CLI flags
func NewConfig(flag *flag.FlagSet, appName string) *Config {
	if flag != nil && flag.Parsed() {
		log.Fatal("flag must be unparsed when calling NewConfig")
	}

	cfg := Config{
		flag:    flag,
		appName: appName,
	}

	flag.StringVar(&cfg.level, "log", env.GetenvWithDefault("log_level", "info"), "log level (all|err|warn|info|debug)")
	flag.StringVar(&cfg.processLevel, "process-level-log", os.Getenv("log_level_"+appName), "process level override (all|err|warn|info|debug)")

	return &cfg
}

// Logger gets the logger
func (cfg *Config) Logger(w io.Writer, keyvals ...any) (*L, error) {
	if cfg.flag != nil && !cfg.flag.Parsed() {
		return nil, errors.New("must parse flags before calling Logger")
	}

	l := New(w, cfg.appName, cfg.Level(), keyvals...)
	return l, nil
}

// Level exposes the configured log level as a string. If a process level logging
// configuration has been specified it will override the generic "gm_log_level"
func (cfg *Config) Level() string {
	if cfg.processLevel != "" {
		return cfg.processLevel
	}
	return cfg.level
}
