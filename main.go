package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
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
		e, ok := errors.Cause(err).(stackTracer)
		if ok {
			st := e.StackTrace()
			fmt.Printf(`{"level":"error", "msg":"bmctool failed", "error":%v, "stacktrace":"%+v"}`, err, st)
			fmt.Println()
		}
		fmt.Println(ok)
		exitCode = 1
	}
}
