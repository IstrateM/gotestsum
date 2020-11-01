//Package gotestsum is a slightly modified version of github.com/gotestyourself/gotestsum package that allows integration
// of the package within another application.
package gotestsum

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/istratem/gotestsum/pkg/options"
	"gotest.tools/gotestsum/testjson"
)

//TODO: ADD Tests

type log interface {
	Debugf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Printf(format string, args ...interface{})
}

var logger log

//GoTestSum runs go test command on a package then analyzes the json output of tests and generates a junit xml report
//based on the test results.
func GoTestSum(opts *options.Options, log log) error {
	logger = log
	logger.Debugf("running gotestsum")
	err := run(opts)
	return err
}

func lookEnvWithDefault(key, defValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defValue
}

func run(opts *options.Options) error {
	ctx := context.Background()
	goTestProc, err := startGoTest(ctx, goTestCmdArgs(opts))
	if err != nil {
		logger.Errorf("failed to run %s %s: %w",
			goTestProc.cmd.Path,
			strings.Join(goTestProc.cmd.Args, " "), err)
		return err
	}
	defer goTestProc.cancel()

	out := os.Stdout
	handler, err := newEventHandler(opts, out, os.Stderr)
	if err != nil {
		return err
	}
	defer handler.Close() // nolint: errcheck
	execution, err := testjson.ScanTestOutput(testjson.ScanConfig{
		Stdout:  goTestProc.stdout,
		Stderr:  goTestProc.stderr,
		Handler: handler,
	})
	if err != nil {
		return err
	}
	testjson.PrintSummary(out, execution, opts.NoSummary.Value)
	if err := writeJUnitFile(opts, execution); err != nil {
		return err
	}
	return goTestProc.cmd.Wait()
}

func goTestCmdArgs(opts *options.Options) []string {
	logger.Debugf("setting command line args")
	args := opts.Args
	defaultArgs := []string{"go", "test"}
	switch {
	case opts.RawCommand:
		logger.Debugf("args set:%v", args)
		return args
	case len(args) == 0:
		a := append(defaultArgs, "-json", pathFromEnv("./..."))
		logger.Debugf("args set:%v", a)
		return a
	case !hasJSONArg(args):
		defaultArgs = append(defaultArgs, "-json")
	}
	if testPath := pathFromEnv(""); testPath != "" {
		args = append(args, testPath)
	}
	a := append(defaultArgs, args...)
	logger.Debugf("args set:%v", a)
	return a
}

func pathFromEnv(defaultPath string) string {
	return lookEnvWithDefault("TEST_DIRECTORY", defaultPath)
}

func hasJSONArg(args []string) bool {
	for _, arg := range args {
		if arg == "-json" || arg == "--json" {
			return true
		}
	}
	return false
}

type proc struct {
	cmd    *exec.Cmd
	stdout io.Reader
	stderr io.Reader
	cancel func()
}

func startGoTest(ctx context.Context, args []string) (proc, error) {
	if len(args) == 0 {
		return proc{}, fmt.Errorf("missing command to run")
	}

	ctx, cancel := context.WithCancel(ctx)
	p := proc{
		cmd:    exec.CommandContext(ctx, args[0], args[1:]...),
		cancel: cancel,
	}
	logger.Debugf("exec: %s", p.cmd.Args)
	var err error
	p.stdout, err = p.cmd.StdoutPipe()
	if err != nil {
		return p, err
	}
	p.stderr, err = p.cmd.StderrPipe()
	if err != nil {
		return p, err
	}
	logger.Debugf("executing args")
	err = p.cmd.Start()
	if err == nil {
		logger.Debugf("go test pid: %d", p.cmd.Process.Pid)
	}
	return p, err
}
