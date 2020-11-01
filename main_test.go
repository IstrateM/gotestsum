package main

import (
	"github.com/istratem/gotestsum/pkg/operator"
	"github.com/istratem/gotestsum/pkg/options"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestTestFailError_Error(t1 *testing.T) {
	type fields struct {
		message string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"expect message", fields{message: "test"}, "test"},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &TestFailError{
				message: tt.fields.message,
			}
			if got := t.Error(); got != tt.want {
				t1.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_contains(t *testing.T) {
	type args struct {
		s []*operator.FailedTest
		e *operator.FailedTest
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"contains", args{
			s: []*operator.FailedTest{{
				Title:  "test",
				Issues: "made up",
			},
			},
			e: &operator.FailedTest{
				Title:  "test",
				Issues: "made up",
			},
		}, true}, {"not_contains", args{

			s: []*operator.FailedTest{{
				Title:   "test",
				Issues:  "made up",
				IssueNo: 0},
			},
			e: &operator.FailedTest{
				Title:   "test2",
				Issues:  "",
				IssueNo: 0},
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := contains(tt.args.s, tt.args.e); got != tt.want {
				t.Errorf("contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_lookEnvWithDefault(t *testing.T) {
	type args struct {
		key      string
		defValue string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"test_env_var", args{
			key:      "test",
			defValue: "test",
		}, "test"},
		{
			"empty_env_var", args{
				key:      "test",
				defValue: "",
			}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := lookEnvWithDefault(tt.args.key, tt.args.defValue); got != tt.want {
				t.Errorf("lookEnvWithDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_run(t *testing.T) {
	type args struct {
		opts *options.Options
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"unauth", args{&options.Options{
			JunitFile: "test.xml",
			Owner:     "astralkn",
			Repo:      "test-bot-playground",
			GitUnAuth: true,
		}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := run(tt.args.opts); (err != nil) != tt.wantErr {
				t.Errorf("run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

//func TestFailExample(t *testing.T) {
//	t.Fatal("This test fails on purpose")
//}
