package bootcmd

import (
	"context"
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/bmc-toolbox/bmclib"
	"github.com/go-playground/validator"
	"github.com/jacobweinstock/bmctool/pkg/rootcmd"
	"github.com/peterbourgon/ff/v3"
	"github.com/peterbourgon/ff/v3/ffcli"
	"github.com/pkg/errors"
)

const bootDeviceCmd = "boot"

// Config for the create subcommand, including a reference to the API client.
type Config struct {
	rootConfig *rootcmd.Config
	Device     string `validate:"oneof=pxe disk"`
	Persistent bool
}

func New(rootConfig *rootcmd.Config) *ffcli.Command {
	cfg := Config{
		rootConfig: rootConfig,
	}

	fs := flag.NewFlagSet(bootDeviceCmd, flag.ExitOnError)
	cfg.RegisterFlags(fs)

	return &ffcli.Command{
		Name:       bootDeviceCmd,
		ShortUsage: "bmctool bootdevice -device",
		ShortHelp:  "Set a boot device.",
		FlagSet:    fs,
		Options:    []ff.Option{ff.WithEnvVarPrefix(strings.ToUpper(cfg.rootConfig.AppName))},
		Exec:       cfg.Exec,
	}
}

func (c *Config) RegisterFlags(fs *flag.FlagSet) {
	fs.StringVar(&c.Device, "device", "", "name of device")
	fs.BoolVar(&c.Persistent, "persistent", false, "persist the boot device across reboots")
}

// Exec function for this command.
func (c *Config) Exec(ctx context.Context, args []string) error {
	if err := c.validateConfig(ctx); err != nil {
		return err
	}
	c.rootConfig.Log.V(0).Info("boot device start", "device", c.Device)
	ok, err := c.doBootDevice(ctx)
	if err != nil {
		c.rootConfig.Log.V(0).Info("boot device complete", "device", c.Device, "successful", ok, "error", err.Error())
		return err
	}
	c.rootConfig.Log.V(0).Info("boot device complete", "device", c.Device, "successful", ok)
	return err
}

func (c *Config) validateConfig(ctx context.Context) error {
	if err := validator.New().Struct(c); err != nil {
		fields := make([]interface{}, 2)
		for _, msg := range err.(validator.ValidationErrors) {
			fields[0] = msg.Value()
			fields[1] = msg.Param()
		}
		return fmt.Errorf("'%v' not a valid device, must be one of [%v]", fields[0], fields[1])
	}
	return nil
}

func (c *Config) doBootDevice(ctx context.Context) (bool, error) {
	client := bmclib.NewClient(c.rootConfig.IP, "623", c.rootConfig.User, c.rootConfig.Pass)
	var err error
	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(c.rootConfig.Timeout))
	defer cancel()
	err = client.Open(ctx)
	if err != nil {
		return false, errors.Wrap(err, "failed to set boot device")
	}
	defer client.Close(ctx)

	client.Registry.Drivers = client.Registry.PreferProtocol(c.rootConfig.Protocol)
	ok, err := client.SetBootDevice(ctx, c.Device, c.Persistent, false)
	if err != nil {
		return false, errors.Wrap(err, "failed to set boot device")
	}
	if !ok {
		return false, errors.New("an unknown error occurred")
	}
	return true, nil
}
