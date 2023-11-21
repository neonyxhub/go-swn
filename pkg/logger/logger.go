package logger

import (
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	STDOUT = "stdout"
	STDERR = "stderr"
)

type Logger = *zap.Logger

type LoggerCfg struct {
	Name     string
	Dev      bool
	OutPaths []string
	ErrPaths []string
}

func ensureFilepathExists(paths []string) error {
	for _, path := range paths {
		if path == STDOUT || path == STDERR {
			continue
		}

		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
		file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			return err
		}
		if err := file.Close(); err != nil {
			return err
		}
	}

	return nil
}

func validateCfg(cfg *LoggerCfg) error {
	if len(cfg.OutPaths) == 0 {
		cfg.OutPaths = []string{STDOUT}
	}
	if len(cfg.ErrPaths) == 0 {
		cfg.ErrPaths = []string{STDERR}
	}
	if len(cfg.OutPaths) == 2 {
		if err := ensureFilepathExists(cfg.OutPaths); err != nil {
			return err
		}
	}
	if len(cfg.ErrPaths) == 2 {
		if err := ensureFilepathExists(cfg.ErrPaths); err != nil {
			return err
		}
	}

	return nil
}

func New(cfg *LoggerCfg) (*zap.Logger, error) {
	err := validateCfg(cfg)
	if err != nil {
		return nil, err
	}

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	config := zap.Config{
		Level:             zap.NewAtomicLevelAt(zap.InfoLevel),
		Development:       cfg.Dev,
		DisableCaller:     false,
		DisableStacktrace: false,
		Sampling:          nil,
		Encoding:          "json",
		EncoderConfig:     encoderCfg,
		OutputPaths:       cfg.OutPaths,
		ErrorOutputPaths:  cfg.ErrPaths,
		InitialFields: map[string]interface{}{
			"pid":  os.Getpid(),
			"name": cfg.Name,
		},
	}

	log := zap.Must(config.Build())

	return log, nil
}
