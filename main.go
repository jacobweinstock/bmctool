package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/jacobweinstock/bmctool/cmd"
	"github.com/pkg/errors"
)

func main() {
	exitCode := 0
	defer func() {
		os.Exit(exitCode)
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGHUP, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())

	defer func() {
		signal.Stop(signals)
		cancel()
	}()

	go func() {
		select {
		case <-signals:
			cancel()
		case <-ctx.Done():
		}
	}()

	if err := cmd.Execute(ctx); err != nil {

		///fmt.Printf("%+q\n", err)
		type stackTracer interface {
			StackTrace() errors.StackTrace
		}
		msg := `{"level":"error", "msg":"bmctool failed", "error":%q, "stacktrace":%q}`
		var st string
		e, ok := err.(stackTracer)
		if ok {
			tr := e.StackTrace()
			st = strings.Replace(fmt.Sprintf("%+v", tr), `\n`, "\n", -1)
		}
		fmt.Printf(msg, err, st)
		exitCode = 1
	}
}
