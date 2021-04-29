package rootcmd

import (
	"context"
	"flag"
	"strings"

	"github.com/go-logr/logr"
	"github.com/peterbourgon/ff/v3"
	"github.com/peterbourgon/ff/v3/ffcli"
)

const appName = "bmctool"

type Config struct {
	Auth
	Timeout  int
	Protocol string
	JSON     bool
	AppName  string
	Log      logr.Logger
}

type Auth struct {
	IP   string `validate:"required"`
	User string `validate:"required"`
	Pass string `validate:"required"`
}

func New() (*ffcli.Command, *Config) {
	var cfg Config
	cfg.AppName = appName

	fs := flag.NewFlagSet(appName, flag.ExitOnError)
	cfg.RegisterFlags(fs)

	return &ffcli.Command{
		ShortUsage: "bmctool [flags] <subcommand>",
		FlagSet:    fs,
		Options:    []ff.Option{ff.WithEnvVarPrefix(strings.ToUpper(appName))},
		Exec:       cfg.Exec,
	}, &cfg

}

func (c *Config) RegisterFlags(fs *flag.FlagSet) {
	fs.StringVar(&c.IP, "ip", "", "bmc ip")
	fs.StringVar(&c.User, "user", "", "bmc user")
	fs.StringVar(&c.Pass, "pass", "", "bmc pass")
	fs.IntVar(&c.Timeout, "timeout", 30, "timeout (in seconds) for BMC calls")
	fs.StringVar(&c.Protocol, "protocol", "", "which BMC protocol to try first (ex. redfish)")
	fs.BoolVar(&c.JSON, "json", false, "output logs in json format")
}

// Exec function for this command.
func (c *Config) Exec(context.Context, []string) error {
	// The root command has no meaning, so if it gets executed,
	// display the usage text to the user instead.
	return flag.ErrHelp
}
