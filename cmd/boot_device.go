package cmd

import (
	"context"
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/bmc-toolbox/bmclib"
	"github.com/go-playground/validator"
	"github.com/peterbourgon/ff/v3"
	"github.com/peterbourgon/ff/v3/ffcli"
	"github.com/pkg/errors"
)

const bootDeviceCmd = "boot"

func getBootDeviceCmd(ctx context.Context, cfg *config) *ffcli.Command {
	bootDeviceFlagSet := flag.NewFlagSet(bootDeviceCmd, flag.ExitOnError)
	bootDeviceFlagSet.StringVar(&cfg.Device, "device", "", "name of device")
	bootDeviceFlagSet.BoolVar(&cfg.Persistent, "persistent", false, "persist the boot device across reboots")
	return &ffcli.Command{
		Name:       bootDeviceCmd,
		ShortUsage: "bmctool bootdevice -device",
		ShortHelp:  "Set a boot device.",
		FlagSet:    bootDeviceFlagSet,
		Options:    []ff.Option{ff.WithEnvVarPrefix(strings.ToUpper(appName))},
		Exec: func(ctx context.Context, args []string) error {
			if err := validator.New().Struct(cfg.bootDevice); err != nil {
				fields := make([]interface{}, 2)
				for _, msg := range err.(validator.ValidationErrors) {
					fields[0] = msg.Value()
					fields[1] = msg.Param()
				}
				return fmt.Errorf("'%v' not a valid device, must be one of [%v]", fields[0], fields[1])
			}
			client := bmclib.NewClient(cfg.IP, "623", cfg.User, cfg.Pass)
			var err error
			ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(cfg.Timeout))
			defer cancel()
			client.Registry.Drivers, err = client.Open(ctx)
			if err != nil {
				return errors.Wrap(err, "failed to set boot device")
			}
			defer client.Close(ctx)

			client.Registry.Drivers = client.Registry.PreferProtocol(cfg.Protocol)
			ok, err := client.SetBootDevice(ctx, cfg.Device, cfg.Persistent, false)
			if err != nil {
				return errors.Wrap(err, "failed to set boot device")
			}
			if !ok {
				return errors.New("an unknown error occurred")
			}
			return nil
		},
	}
}
