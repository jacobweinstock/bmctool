package cmd

import (
	"context"
	"os"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"github.com/go-playground/validator"
	"github.com/jacobweinstock/bmctool/pkg/bootcmd"
	"github.com/jacobweinstock/bmctool/pkg/powercmd"
	"github.com/jacobweinstock/bmctool/pkg/rootcmd"
	"github.com/peterbourgon/ff/v3/ffcli"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Execute(ctx context.Context) error {
	var (
		rootCommand, rootConfig = rootcmd.New()
		powerCommand            = powercmd.New(rootConfig)
		bootCommand             = bootcmd.New(rootConfig)
	)

	rootCommand.Subcommands = []*ffcli.Command{
		powerCommand,
		bootCommand,
	}

	if err := rootCommand.Parse(os.Args[1:]); err != nil {
		return err
	}

	if err := validator.New().Struct(rootConfig.Auth); err != nil {
		return err
	}
	rootConfig.Log = defaultLogger(rootConfig.JSON)

	if err := rootCommand.Run(context.Background()); err != nil {
		return err
	}
	return nil
}

func defaultLogger(jsonLogs bool) logr.Logger {
	encoding := "console"
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	encoderConfig.EncodeCaller = nil
	encoderConfig.TimeKey = ""
	level := zap.NewAtomicLevelAt(zap.InfoLevel)
	if jsonLogs {
		encoding = "json"
		encoderConfig = zap.NewProductionEncoderConfig()
	}
	config := zap.Config{
		Level:            level,
		Encoding:         encoding,
		EncoderConfig:    encoderConfig,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
	logger, err := config.Build()
	if err != nil {
		panic(err)
	}

	return zapr.NewLogger(logger)
}
