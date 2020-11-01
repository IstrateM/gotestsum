//Package options stores the necessary data for running that application and implements flag operations
package options

import (
	"github.com/istratem/gotestsum/pkg/junitxml"
)

//Options stores flags value
type Options struct {
	Args                         []string
	Format                       string
	Debug                        bool
	RawCommand                   bool
	JsonFile                     string
	JunitFile                    string
	NoColor                      bool
	NoSummary                    *NoSummaryValue
	JunitTestSuiteNameFormat     *JunitFieldFormatValue
	JunitTestCaseClassnameFormat *JunitFieldFormatValue
	Version                      bool
	Post                         bool
	Token                        string
	Owner                        string
	Repo                         string
	JUnitTestSuite               junitxml.JUnitTestSuites
	GitUnAuth                    bool
}
