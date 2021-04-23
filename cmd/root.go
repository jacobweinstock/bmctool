package cmd

import (
	"context"
	"flag"
	"strings"

	"github.com/peterbourgon/ff/v3"
	"github.com/peterbourgon/ff/v3/ffcli"
)

const appName = "bmctool"

func getRootCmd(ctx context.Context, cfg *config) *ffcli.Command {
	rootFlagSet := flag.NewFlagSet(appName, flag.ExitOnError)
	rootFlagSet.StringVar(&cfg.IP, "ip", "", "bmc ip")
	rootFlagSet.StringVar(&cfg.User, "user", "", "bmc user")
	rootFlagSet.StringVar(&cfg.Pass, "pass", "", "bmc pass")
	rootFlagSet.IntVar(&cfg.Timeout, "timeout", 30, "timeout (in seconds) for BMC calls")
	rootFlagSet.StringVar(&cfg.Protocol, "preferprotocol", "", "which BMC protocol to try first (ex. redfish)")

	return &ffcli.Command{
		ShortUsage: "bmctool [flags] <subcommand>",
		FlagSet:    rootFlagSet,
		Options:    []ff.Option{ff.WithEnvVarPrefix(strings.ToUpper(appName))},
		Exec: func(ctx context.Context, args []string) error {
			return flag.ErrHelp
		},
	}
}
