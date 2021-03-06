package main

import (
	"context"
	"os"
	"os/exec"

	"github.com/fatih/color"
	"github.com/istratem/gotestsum/pkg/gotestsum"
	"github.com/istratem/gotestsum/pkg/operator"
	"github.com/istratem/gotestsum/pkg/options"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"gotest.tools/gotestsum/testjson"
)

var version = "master"

type TestFailError struct {
	message string
}

func (t *TestFailError) Error() string {
	return t.message
}

func main() {
	name := os.Args[0]
	flags, opts := setupFlags(name)
	switch err := flags.Parse(os.Args[1:]); {
	case err == pflag.ErrHelp:
		os.Exit(0)
	case err != nil:
		log.Error(err.Error())
		flags.Usage()
		os.Exit(1)
	}
	opts.Args = flags.Args()
	setupLogging(opts)
	if opts.Version {
		log.Printf("gotestmng version %s\n", version)
		os.Exit(0)
	}
	log.Println(opts)
	if opts.Post == true {
		if opts.GitUnAuth == true {
			log.Println(os.Args)
			log.Error("can not post issues without authentication")
			os.Exit(1)
		}
	}
	err := run(opts)
	switch err.(type) {
	case nil:
	case *exec.ExitError:
		os.Exit(0)
	case *TestFailError:
		log.Error(name + ": Error :" + err.Error())
		os.Exit(1)
	default:
		log.Error(name + ": Error :" + err.Error())
		os.Exit(3)
	}
}

func setupFlags(name string) (*pflag.FlagSet, *options.Options) {
	opts := options.Options{
		NoSummary:                    options.NewNoSummaryValue(),
		JunitTestCaseClassnameFormat: &options.JunitFieldFormatValue{},
		JunitTestSuiteNameFormat:     &options.JunitFieldFormatValue{},
	}
	flags := pflag.NewFlagSet(name, pflag.ContinueOnError)
	flags.SetInterspersed(false)
	flags.Usage = func() {
		log.Errorf(`Usage:
    %s [flags] [--] [go test flags]

Flags:
`, name)
		flags.PrintDefaults()
		log.Error(`
Formats:
    dots                    print a character for each test
    dots-v2                 experimental dots format, one package per line
    pkgname                 print a line for each package
    pkgname-and-test-fails  print a line for each package and failed test output
    testname                print a line for each test and package
    standard-quiet          standard go test format
    standard-verbose        standard go test -v format
`)
	}
	flags.BoolVar(&opts.Debug, "debug", false, "enabled debug")
	flags.StringVar(&opts.Token, "token", "", "set remote github auth token")
	flags.StringVar(&opts.Owner, "owner", "", "set remote github repository owner ")
	flags.StringVar(&opts.Repo, "repo", "", "set remote github repository")
	flags.BoolVar(&opts.Post, "post", false, "post found failures on github")
	flags.StringVarP(&opts.Format, "format", "f", lookEnvWithDefault("GOTESTSUM_FORMAT", "short"), "print format of test input")
	flags.BoolVar(&opts.RawCommand, "raw-command", false, "don't prepend 'go test -json' to the 'go test' command")
	flags.StringVar(&opts.JsonFile, "jsonfile", lookEnvWithDefault("GOTESTSUM_JSONFILE", ""), "write all TestEvents to file")
	flags.StringVar(&opts.JunitFile, "junitfile", lookEnvWithDefault("GOTESTSUM_JUNITFILE", ""), "write a JUnit XML file")
	flags.BoolVar(&opts.NoColor, "no-color", color.NoColor, "disable color output")
	flags.Var(opts.NoSummary, "no-summary", "do not print summary of: "+testjson.SummarizeAll.String())
	flags.Var(opts.JunitTestSuiteNameFormat, "junitfile-testsuite-name", "format the testsuite name field as: "+options.JunitFieldFormatValues)
	flags.Var(opts.JunitTestCaseClassnameFormat, "junitfile-testcase-classname", "format the testcase classname field as: "+options.JunitFieldFormatValues)
	flags.BoolVar(&opts.Version, "version", false, "show version and exit")
	flags.BoolVar(&opts.GitUnAuth, "unauth", false, "use unauthenticated git operator")
	return flags, &opts
}

func run(opts *options.Options) error {
	err := gotestsum.GoTestSum(opts, log.StandardLogger())
	switch e := err.(type) {
	case *exec.ExitError:
		if e.ExitCode() != 1 {
			return err
		}
	}
	var junitOperator operator.JUnitOperator
	failedTests := junitOperator.GetFailedTests(opts)
	var gitOperator *operator.GitOperator
	if opts.GitUnAuth {
		gitOperator = operator.NewUnauthenticatedGitOperator(opts.Owner, opts.Repo, context.Background())

	} else {
		gitOperator = operator.NewGitOperator(opts.Owner, opts.Repo, opts.Token, context.Background())
	}
	knownIssues, err := gitOperator.GetTestIssues()
	if err != nil {
		return err
	}

	var newIssues []operator.FailedTest
	var solvedIssues []operator.FailedTest

	for _, t := range failedTests {
		if !contains(knownIssues, t) {
			newIssues = append(newIssues, *t)
		}
	}

	for _, t := range knownIssues {
		if !contains(failedTests, t) {
			solvedIssues = append(solvedIssues, *t)
		}
	}
	for _, t := range newIssues {
		log.Printf("NEW Failed Test found : %v\n", t)
	}

	if opts.Post {
		for _, i := range newIssues {
			err = gitOperator.PostNewIssue(&i)
			if err != nil {
				return err
			}
			log.Println("New issue created on git", i)
		}
		for _, i := range solvedIssues {
			err = gitOperator.CloseSolvedIssue(&i)
			if err != nil {
				return err
			}
			log.Println("Issue closed on git", i)
		}
	}

	if len(newIssues) > 0 {
		return &TestFailError{message: "New failing test(s) found"}
	}

	return nil
}

func lookEnvWithDefault(key, defValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	} else {
		return defValue
	}
}

func contains(s []*operator.FailedTest, e *operator.FailedTest) bool {
	for _, a := range s {
		if a.Title == e.Title /*&& a.Issues == e.Issues*/ {
			return true
		}
	}
	return false
}

func setupLogging(opts *options.Options) {
	if opts.Debug {
		log.SetLevel(log.DebugLevel)
	}
}
