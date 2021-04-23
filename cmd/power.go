package cmd

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/bmc-toolbox/bmclib"
	"github.com/go-playground/validator"
	"github.com/peterbourgon/ff/v3"
	"github.com/peterbourgon/ff/v3/ffcli"
)

const powerCmd = "power"

func getPowerCmd(ctx context.Context, cfg *config) *ffcli.Command {
	powerFlagSet := flag.NewFlagSet(powerCmd, flag.ExitOnError)
	powerFlagSet.StringVar(&cfg.Action, "action", "", "power action")
	return &ffcli.Command{
		Name:       powerCmd,
		ShortUsage: "bmctool power [-n times] <arg>",
		ShortHelp:  "Power actions for a machine.",
		FlagSet:    powerFlagSet,
		Options:    []ff.Option{ff.WithEnvVarPrefix(strings.ToUpper(appName))},
		Exec: func(ctx context.Context, args []string) error {
			if err := validator.New().Struct(cfg.power); err != nil {
				fields := make([]interface{}, 2)
				for _, msg := range err.(validator.ValidationErrors) {
					fields[0] = msg.Value()
					fields[1] = msg.Param()
				}
				return fmt.Errorf("got '%v', power action must be one of [%v]", fields[0], fields[1])
			}
			client := bmclib.NewClient(cfg.IP, "623", cfg.User, cfg.Pass)
			var err error
			ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(cfg.Timeout))
			defer cancel()
			client.Registry.Drivers, err = client.Open(ctx)
			if err != nil {
				return err
			}
			defer client.Close(ctx)

			ok, err := client.SetPowerState(ctx, cfg.Action)
			if err != nil {
				return err
			}
			if !ok {
				return errors.New("an unknown error occured")
			}
			return nil
		},
	}
}
