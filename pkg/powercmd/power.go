package powercmd

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

const powerCmd = "power"

// Config for the create subcommand, including a reference to the API client.
type Config struct {
	rootConfig *rootcmd.Config
	Action     string `validate:"oneof=on off cycle reset"`
}

func New(rootConfig *rootcmd.Config) *ffcli.Command {
	cfg := Config{
		rootConfig: rootConfig,
	}

	fs := flag.NewFlagSet(powerCmd, flag.ExitOnError)
	cfg.RegisterFlags(fs)

	return &ffcli.Command{
		Name:       powerCmd,
		ShortUsage: "bmctool power -action [on|off|cycle|reset|status]",
		ShortHelp:  "Power actions for a machine.",
		FlagSet:    fs,
		Options:    []ff.Option{ff.WithEnvVarPrefix(strings.ToUpper(cfg.rootConfig.AppName))},
		Exec:       cfg.Exec,
	}

}

func (c *Config) RegisterFlags(fs *flag.FlagSet) {
	fs.StringVar(&c.Action, "action", "", "power action")
}

// Exec function for this command.
func (c *Config) Exec(ctx context.Context, args []string) error {
	if err := c.validateConfig(ctx); err != nil {
		return err
	}
	log := c.rootConfig.Log.WithValues("action", c.Action)
	log.V(0).Info("power action start")
	ok, metadata, err := c.doPower(ctx)
	if err != nil {
		log.V(0).Info("power action complete", "successful", ok, "details", metadata, "error", err.Error())
		return err
	}
	log.V(0).Info("power action complete", "successful", ok, "provider", metadata.SuccessfulProvider)
	return err
}

func (c *Config) validateConfig(ctx context.Context) error {
	if err := validator.New().StructPartialCtx(ctx, c, "Action"); err != nil {
		var errMsg []interface{}
		s := "'%v' not a valid %v, must be %v [%v]"
		for _, msg := range err.(validator.ValidationErrors) {
			errMsg = append(errMsg, msg.Value(), msg.Field(), msg.Tag(), msg.Param())
		}
		return fmt.Errorf(s, errMsg...)
	}
	return nil
}

func (c *Config) doPower(ctx context.Context) (ok bool, metadata bmc.Metadata, err error) {
	client := bmclib.NewClient(c.rootConfig.IP, "623", c.rootConfig.User, c.rootConfig.Pass)
	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(c.rootConfig.Timeout))
	defer cancel()
	err = client.Open(ctx)
	if err != nil {
		return false, metadata, errors.Wrapf(err, "failed to set power state: %+v", client.GetMetadata())
	}
	defer client.Close(ctx)

	client.Registry.Drivers = client.Registry.PreferProtocol(c.rootConfig.Protocol)
	ok, err = client.SetPowerState(ctx, c.Action)
	metadata = client.GetMetadata()
	if err != nil {
		return false, metadata, err
	}
	if !ok {
		return false, metadata, errors.New("an unknown error occured")
	}
	return true, metadata, nil
}
