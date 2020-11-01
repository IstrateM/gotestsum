package gotestsum

import (
	"fmt"
	"io"
	"os"

	"github.com/istratem/gotestsum/pkg/junitxml"
	"github.com/istratem/gotestsum/pkg/options"
	"gotest.tools/gotestsum/testjson"
)

type eventHandler struct {
	formatter testjson.EventFormatter
	err       io.Writer
	jsonFile  io.WriteCloser
}

func (h *eventHandler) Err(text string) error {
	_, err := h.err.Write([]byte(text + "\n"))
	return err
}

func (h *eventHandler) Event(event testjson.TestEvent, execution *testjson.Execution) error {
	if h.jsonFile != nil {
		if _, err := h.jsonFile.Write(append(event.Bytes(), '\n')); err != nil {
			return fmt.Errorf("failed to write JSON file %w", err)
		}
	}

	if err := h.formatter.Format(event, execution); err != nil {
		return fmt.Errorf("failed to format event %w", err)
	}
	return nil
}

func (h *eventHandler) Close() error {
	if h.jsonFile != nil {
		if err := h.jsonFile.Close(); err != nil {
			return err
		}
	}
	return nil
}

var _ testjson.EventHandler = &eventHandler{}

func newEventHandler(opts *options.Options, stdout io.Writer, stderr io.Writer) (*eventHandler, error) {
	formatter := testjson.NewEventFormatter(stdout, opts.Format)
	if formatter == nil {
		return nil, fmt.Errorf("unknown format %s", opts.Format)
	}
	handler := eventHandler{
		formatter: formatter,
		err:       stderr,
	}
	var err error
	if opts.JsonFile != "" {
		handler.jsonFile, err = os.Create(opts.JsonFile)
		if err != nil {
			return &handler, fmt.Errorf("failed to open JSON file %w", err)
		}
	}
	return &handler, nil
}

func writeJUnitFile(opts *options.Options, execution *testjson.Execution) error {
	if opts.JunitFile == "" {
		return nil
	}
	junitFile, err := os.Create(opts.JunitFile)
	if err != nil {
		return err
	}
	defer func() {
		err = junitFile.Close()
	}()

	opts.JUnitTestSuite, err = junitxml.Write(junitFile, execution, junitxml.Config{
		FormatTestSuiteName:     opts.JunitTestSuiteNameFormat.Value(),
		FormatTestCaseClassname: opts.JunitTestCaseClassnameFormat.Value(),
	})
	return err
}
