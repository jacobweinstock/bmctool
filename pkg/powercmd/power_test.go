package powercmd

import (
	"context"
	"flag"
	"testing"

	"github.com/go-logr/logr"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jacobweinstock/bmctool/pkg/rootcmd"
	"github.com/peterbourgon/ff/v3/ffcli"
	"github.com/pkg/errors"
)

func TestNew(t *testing.T) {
	action := ""
	fs := flag.NewFlagSet(powerCmd, flag.ExitOnError)
	fs.StringVar(&action, "action", "", "power action")

	expectedNew := &ffcli.Command{
		Name:       powerCmd,
		ShortUsage: "bmctool power -action [on|off|cycle|reset|status]",
		ShortHelp:  "Power actions for a machine.",
		FlagSet:    fs,
	}
	rootConfig := &rootcmd.Config{}
	got := New(rootConfig)
	opts := []cmp.Option{
		cmpopts.IgnoreUnexported(ffcli.Command{}),
		cmpopts.IgnoreFields(ffcli.Command{}, "Exec", "Options"),
		cmp.AllowUnexported(flag.FlagSet{}),
		cmpopts.IgnoreFields(flag.FlagSet{}, "Usage"),
	}
	if diff := cmp.Diff(got, expectedNew, opts...); diff != "" {
		t.Fatal(diff)
	}
}

func TestValidateConfig(t *testing.T) {
	log := logr.Discard()
	tests := map[string]struct {
		input *Config
		want  error
	}{
		"successful validation": {input: &Config{rootConfig: &rootcmd.Config{Log: log}, Action: "on"}, want: nil},
		"failed validation":     {input: &Config{rootConfig: &rootcmd.Config{Log: log}, Action: "blah"}, want: errors.New("'blah' not a valid Action, must be oneof [on off cycle reset]")},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.input.validateConfig(context.Background())
			if tc.want != nil {
				if diff := cmp.Diff(tc.want.Error(), got.Error()); diff != "" {
					t.Fatalf(diff)
				}
			} else {
				if diff := cmp.Diff(tc.want, got); diff != "" {
					t.Fatalf(diff)
				}
			}
		})
	}
}
