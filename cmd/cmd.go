package cmd

import (
	"context"
	"os"

	"github.com/go-playground/validator"
	"github.com/peterbourgon/ff/v3/ffcli"
)

func Execute(ctx context.Context) error {
	cfg := config{}

	root := getRootCmd(ctx, &cfg)
	root.Subcommands = []*ffcli.Command{getPowerCmd(ctx, &cfg), getBootDeviceCmd(ctx, &cfg)}

	if err := root.Parse(os.Args[1:]); err != nil {
		return err
	}
	if err := validator.New().Struct(cfg.auth); err != nil {
		return err
	}
	if err := root.Run(context.Background()); err != nil {
		return err
	}

	return nil
}
