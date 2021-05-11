package rootcmd

import (
	"flag"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/peterbourgon/ff/v3/ffcli"
)

func TestNew(t *testing.T) {
	c := Config{AppName: appName}
	fs := flag.NewFlagSet(appName, flag.ExitOnError)
	fs.StringVar(&c.IP, "ip", "", "bmc ip")
	fs.StringVar(&c.Port, "port", "623", "bmc port")
	fs.StringVar(&c.User, "user", "", "bmc user")
	fs.StringVar(&c.Pass, "pass", "", "bmc pass")
	fs.IntVar(&c.Timeout, "timeout", 30, "timeout (in seconds) for BMC calls")
	fs.StringVar(&c.Protocol, "protocol", "", "which BMC protocol to try first (ex. redfish)")
	fs.BoolVar(&c.JSON, "json", false, "output logs in json format")

	expectedNew := &ffcli.Command{
		Name:       "bmctool",
		ShortUsage: "bmctool [flags] <subcommand>",
		FlagSet:    fs,
	}
	got, config := New()
	opts := []cmp.Option{
		cmpopts.IgnoreUnexported(ffcli.Command{}),
		cmpopts.IgnoreFields(ffcli.Command{}, "Exec", "Options"),
		cmp.AllowUnexported(flag.FlagSet{}),
		cmpopts.IgnoreFields(flag.FlagSet{}, "Usage"),
	}
	if diff := cmp.Diff(got, expectedNew, opts...); diff != "" {
		t.Fatal(diff)
	}
	if diff := cmp.Diff(config, &c, opts...); diff != "" {
		t.Fatal(diff)
	}
}
