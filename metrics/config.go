package metrics

import (
	"context"
	"errors"
	"flag"
	"log"
	"time"

	"github.com/graymeta/gmkit/logger"

	"github.com/graymeta/env"
)

// Config contains all the logic for parsing command line flags.
type Config struct {
	flag    *flag.FlagSet
	appName string

	runtimeEnabled  bool
	runtimeInterval time.Duration
	loggingEnabled  bool
	logLevel        string
}

// NewConfig prepares a new configuration for parsing CLI flags.
func NewConfig(l *logger.L, flag *flag.FlagSet, appName string) *Config {
	if flag != nil && flag.Parsed() {
		log.Fatal("flag must be unparsed when calling NewConfig")
	}

	cfg := Config{
		flag:    flag,
		appName: appName,
	}

	flag.BoolVar(&cfg.runtimeEnabled, "runtime-metrics-enabled", env.GetenvBoolWithDefault("gm_runtime_metrics_enabled", false), "Whether or not runtime metrics will be periodically emitted.")
	flag.DurationVar(&cfg.runtimeInterval, "runtime-metrics-interval", env.GetenvDurationWithDefault("gm_runtime_metrics_interval", 10*time.Second), "The interval that runtime metrics will be emitted. Default: 10s")
	flag.BoolVar(&cfg.loggingEnabled, "metrics-logging-enabled", env.GetenvBoolWithDefault("gm_metrics_logging", false), "Whether or not metrics will be emitted to the logs.")
	flag.StringVar(&cfg.logLevel, "metrics-log-level", env.GetenvWithDefault("gm_metrics_log_level", "debug"), "What level to log metrics at. Default: debug")

	return &cfg
}

// Service sets the DefaultStatsd with a given service name
// that should be call just once on the service intialization
// in case the env variable 'gm_service' is not set.
func (cfg *Config) Service(ctx context.Context, l *logger.L) error {
	if cfg.flag != nil && !cfg.flag.Parsed() {
		return errors.New("must parse flags before calling Service")
	}

	DefaultStatsd = Statsd(cfg.appName)

	if cfg.loggingEnabled {
		mc := &MultiClient{}
		mc.Append(DefaultStatsd)
		mc.Append(NewLoggingClient(l, cfg.logLevel))
		DefaultStatsd = mc
	}

	if cfg.runtimeEnabled {
		go func() {
			ticker := time.NewTicker(cfg.runtimeInterval)
			for {
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
					go emitRuntimeStats()
				}
			}
		}()
	}

	return nil
}
