module github.com/adodon2go/gotestsum

go 1.13

require (
	github.com/fatih/color v1.10.0
	github.com/google/go-github v17.0.0+incompatible
	github.com/google/go-querystring v1.0.0 // indirect
	github.com/sirupsen/logrus v1.7.0
	github.com/spf13/pflag v1.0.5
	golang.org/x/oauth2 v0.0.0-20200902213428-5d25da1a8d43
	gotest.tools v2.2.0+incompatible
	gotest.tools/gotestsum v0.0.0-00010101000000-000000000000
)

replace gotest.tools/gotestsum => gotest.tools/gotestsum v0.5.0
