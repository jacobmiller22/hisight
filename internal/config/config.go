package config

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/jacobmiller22/gossentials/omniconfig"
)

const DEFAULT_HISIGHT_CONFIG_FILENAME string = "hisight.json"
const DEFAULT_HISIGHT_LOGFILE_FILENAME string = "hisight.log"

var ErrMissingConfig error = errors.New("missing config")
var ErrTooFewArgs error = errors.New("too few args")
var ErrUnknownLogLevel error = errors.New("unknown log level")

func parseLogLevel(s string) (slog.Leveler, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "debug":
		return slog.LevelDebug, nil
	case "info":
		return slog.LevelWarn, nil
	case "warn", "warning":
		return slog.LevelWarn, nil
	case "fatal", "error":
		return slog.LevelError, nil
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnknownLogLevel, s)
	}
}

func parseLogFilepath(s string) (io.Writer, error) {
	if s == "stderr" {
		return os.Stderr, nil
	}
	if s == "stdout" {
		return os.Stdout, nil
	}

	if err := os.MkdirAll(filepath.Dir(s), 0644); err != nil {
		return nil, err
	}

	fd, err := os.OpenFile(s, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	return fd, nil
}

type LogConfig struct {
	Level string `json:"level"`
	File  string `json:"file"`
}

type DBConfig struct {
	DSN string `json:"dsn"`
}

type HTTPServerConfig struct {
	Port int `json:"port"`
}

type GRPCServerConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type ServerConfig struct {
	HTTP HTTPServerConfig `json:"http"`
	GRPC GRPCServerConfig `json:"grpc"`
}

type Config struct {
	ConfigPath string
	Log        LogConfig    `json:"log"`
	DB         DBConfig     `json:"db"`
	Server     ServerConfig `json:"server"`
}

func (c *Config) LogLeveler() slog.Leveler {
	leveler, err := parseLogLevel(c.Log.Level)
	if err != nil {
		return slog.LevelError
	}
	return leveler
}

func (c *Config) Logger() (*slog.Logger, error) {
	logwriter, err := parseLogFilepath(c.Log.File)
	if err != nil {
		return nil, err
	}

	logger := slog.New(slog.NewJSONHandler(
		logwriter,
		&slog.HandlerOptions{
			AddSource: c.LogLeveler() == slog.LevelDebug,
			Level:     c.LogLeveler(),
		},
	))

	return logger, nil
}

var DefaultConfig *Config = &Config{
	Log: LogConfig{
		Level: "error",
		File:  DefaultLogFilePath(),
	},
	DB: DBConfig{
		DSN: ":memory:",
	},
	Server: ServerConfig{
		HTTPServerConfig{
			Port: 9001,
		},
		GRPCServerConfig{
			Host: "127.0.0.1:9002",
			Port: 9002,
		},
	},
}

func DefaultConfigPath() string {
	base := os.Getenv("XDG_CONFIG_HOME")
	if base == "" {
		base = path.Join(os.Getenv("HOME"), ".config")
	}
	return path.Join(base, DEFAULT_HISIGHT_CONFIG_FILENAME)
}

func DefaultLogFilePath() string {
	base := os.Getenv("XDG_DATA_HOME")
	if base == "" {
		base = path.Join(os.Getenv("HOME"), ".local/share")
	}
	return path.Join(base, DEFAULT_HISIGHT_LOGFILE_FILENAME)
}

func fileExists(p string) bool {
	_, err := os.Stat(p)
	if err != nil {
		return false
	}
	return true
}

var defaultConfigurer omniconfig.StaticConfigurer[Config] = omniconfig.StaticConfigurer[Config]{
	Config: DefaultConfig,
}

type emptyConfigurer[T comparable] struct{}

func (e *emptyConfigurer[T]) Load() (*T, error) {
	return nil, nil
}

func newDefaultJsonConfigurer(path string) omniconfig.Configurer[Config] {
	configurer, err := omniconfig.NewFsIOConfigurer[Config](path)

	if err != nil {
		return &emptyConfigurer[Config]{}
	}

	return configurer.With(func(i *omniconfig.IOConfigurer[Config]) {
		i.Processor = omniconfig.JsonReaderProcessor[Config]
	})
}

func newDefaultFlagConfigurer(args []string) *omniconfig.FlagConfigurer[Config] {
	flagset := flag.NewFlagSet(
		"hisight",
		flag.ContinueOnError,
	)
	var cfg Config
	flagset.IntVar(&cfg.Server.HTTP.Port, "server-http-port", 0, "Port to listen on for HTTP requests")
	flagset.StringVar(&cfg.Server.GRPC.Host, "server-grpc-host", "", "Hostname of the GRPC Server")
	flagset.IntVar(&cfg.Server.GRPC.Port, "server-grpc-port", 0, "Port to listen on for GRPC requests")
	flagset.StringVar(&cfg.DB.DSN, "db-dsn", "", "DSN to database")
	flagset.StringVar(&cfg.Log.Level, "log-level", "", "Log level to use")
	flagset.StringVar(&cfg.Log.File, "log-file", "", "Log file to use")

	// flagset.Usage = func() {
	// 	// Default Usage
	// }

	return omniconfig.NewFlagConfigurer(
		flagset,
		&cfg,
		omniconfig.WithFlagConfigurerArgs[Config](args),
	)
}

var ErrConfigExists error = errors.New("configuration file exists at path")

// InitDefaultConfig creates the default configuration file at the specified path
func InitDefaultConfig() error {
	path := DefaultConfigPath()

	if fileExists(path) {
		return fmt.Errorf("%w: %s", ErrConfigExists, path)
	}

	fd, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	if err := json.NewEncoder(fd).Encode(DefaultConfig); err != nil {
		return err
	}

	return nil
}

func LoadConfig(args []string) *Config {
	flagConfigurer := newDefaultFlagConfigurer(args)
	cfg, err := flagConfigurer.Load()
	if err != nil || cfg == nil {
		return DefaultConfig
	}

	cfg, _, err = omniconfig.MergeConfigurers(
		defaultConfigurer,
		newDefaultJsonConfigurer(cfg.ConfigPath),
		flagConfigurer,
	)

	if err != nil {
		return DefaultConfig
	}

	return cfg
}
