package bootcmd

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
	device := ""
	persistent := false
	fs := flag.NewFlagSet(bootDeviceCmd, flag.ExitOnError)
	fs.StringVar(&device, "device", "", "name of device")
	fs.BoolVar(&persistent, "persistent", false, "persist the boot device across reboots")

	expectedNew := &ffcli.Command{
		Name:       bootDeviceCmd,
		ShortUsage: "bmctool bootdevice -device",
		ShortHelp:  "Set a boot device.",
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
		"successful validation": {input: &Config{rootConfig: &rootcmd.Config{Log: log}, Device: "pxe", Persistent: false}, want: nil},
		"failed validation":     {input: &Config{rootConfig: &rootcmd.Config{Log: log}, Device: "blah", Persistent: false}, want: errors.New("'blah' not a valid Device, must be oneof [pxe disk]")},
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

/*
func TestDoBootDevice(t *testing.T) {
	s := ipmi.NewSimulator(net.UDPAddr{Port: 0})
	err := s.Run()
	if err != nil {
		t.Fatal(err)
	}
	defer s.Stop()
	s.SetHandler(ipmi.NetworkFunctionChassis, ipmi.CommandSetSystemBootOptions, func(*ipmi.Message) ipmi.Response {
		return nil
	})
	c := s.NewConnection()
	addr := s.LocalAddr()
	t.Log("ip", addr.IP.String())
	t.Log("port", strconv.Itoa(addr.Port))
	t.Log("user", c.Username)
	t.Log("pass", c.Password)
	config := &Config{
		rootConfig: &rootcmd.Config{
			Auth: rootcmd.Auth{
				IP:   addr.IP.String(),
				Port: strconv.Itoa(addr.Port),
				User: c.Username,
				Pass: c.Password,
			},
			Log:     logr.Discard(),
			Timeout: 20,
		},
		Device:     "pxe",
		Persistent: false,
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second)*30)
	defer cancel()
	ok, err := config.doBootDevice(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Fatal(ok)
}
*/
