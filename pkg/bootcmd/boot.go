package bootcmd

import (
	"context"
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/bmc-toolbox/bmclib"
	"github.com/bmc-toolbox/bmclib/bmc"
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
	log := c.rootConfig.Log.WithValues("device", c.Device)
	log.V(0).Info("set boot device start", "device", c.Device)
	ok, metadata, err := c.doBootDevice(ctx)
	if err != nil {
		log.V(0).Info("set boot device complete", "successful", ok, "details", metadata, "error", err.Error())
		return err
	}
	log.V(0).Info("set boot device complete", "successful", ok, "provider", metadata.SuccessfulProvider)
	return err
}

func (c *Config) validateConfig(ctx context.Context) error {
	if err := validator.New().Struct(c); err != nil {
		var errMsg []interface{}
		s := "'%v' not a valid %v, must be %v [%v]"
		for _, msg := range err.(validator.ValidationErrors) {
			errMsg = append(errMsg, msg.Value(), msg.Field(), msg.Tag(), msg.Param())
		}
		return fmt.Errorf(s, errMsg...)
	}
	return nil
}

func (c *Config) doBootDevice(ctx context.Context) (ok bool, metadata bmc.Metadata, err error) {
	client := bmclib.NewClient(c.rootConfig.IP, c.rootConfig.Port, c.rootConfig.User, c.rootConfig.Pass)
	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(c.rootConfig.Timeout))
	defer cancel()
	err = client.Open(ctx)
	if err != nil {
		return false, metadata, errors.Wrap(err, "failed to set boot device")
	}
	defer client.Close(ctx)

	client.Registry.Drivers = client.Registry.PreferProtocol(c.rootConfig.Protocol)
	ok, err = client.SetBootDevice(ctx, c.Device, c.Persistent, false)
	metadata = client.GetMetadata()
	if err != nil {
		return false, metadata, err
	}
	if !ok {
		return false, metadata, errors.New("setting boot device failed with an unknown error")
	}
	return true, metadata, nil
}
