package operator

import (
	"context"
	"github.com/google/go-github/github"
	"github.com/istratem/gotestsum/pkg/options"
	"reflect"
	"testing"
)

func TestGitOperator_CloseSolvedIssue(t *testing.T) {
	type fields struct {
		client *github.Client
		ctx    context.Context
		owner  string
		repo   string
	}
	type args struct {
		f *FailedTest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GitOperator{
				client: tt.fields.client,
				ctx:    tt.fields.ctx,
				owner:  tt.fields.owner,
				repo:   tt.fields.repo,
			}
			if err := g.CloseSolvedIssue(tt.args.f); (err != nil) != tt.wantErr {
				t.Errorf("CloseSolvedIssue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGitOperator_GetIssuesByLabel(t *testing.T) {
	type fields struct {
		client *github.Client
		ctx    context.Context
		owner  string
		repo   string
	}
	type args struct {
		label string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*github.Issue
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GitOperator{
				client: tt.fields.client,
				ctx:    tt.fields.ctx,
				owner:  tt.fields.owner,
				repo:   tt.fields.repo,
			}
			got, err := g.GetIssuesByLabel(tt.args.label)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetIssuesByLabel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetIssuesByLabel() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGitOperator_GetTestIssues(t *testing.T) {
	type fields struct {
		client *github.Client
		ctx    context.Context
		owner  string
		repo   string
	}
	tests := []struct {
		name    string
		fields  fields
		want    []*FailedTest
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GitOperator{
				client: tt.fields.client,
				ctx:    tt.fields.ctx,
				owner:  tt.fields.owner,
				repo:   tt.fields.repo,
			}
			got, err := g.GetTestIssues()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTestIssues() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTestIssues() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGitOperator_PostNewIssue(t *testing.T) {
	type fields struct {
		client *github.Client
		ctx    context.Context
		owner  string
		repo   string
	}
	type args struct {
		f *FailedTest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GitOperator{
				client: tt.fields.client,
				ctx:    tt.fields.ctx,
				owner:  tt.fields.owner,
				repo:   tt.fields.repo,
			}
			if err := g.PostNewIssue(tt.args.f); (err != nil) != tt.wantErr {
				t.Errorf("PostNewIssue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGitOperator_newIssueForTest(t *testing.T) {
	type fields struct {
		client *github.Client
		ctx    context.Context
		owner  string
		repo   string
	}
	type args struct {
		f *FailedTest
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *github.IssueRequest
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GitOperator{
				client: tt.fields.client,
				ctx:    tt.fields.ctx,
				owner:  tt.fields.owner,
				repo:   tt.fields.repo,
			}
			if got := g.newIssueForTest(tt.args.f); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newIssueForTest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJUnitOperator_GetFailedTests(t *testing.T) {
	type args struct {
		opts *options.Options
	}
	tests := []struct {
		name string
		args args
		want []*FailedTest
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ju := &JUnitOperator{}
			if got := ju.GetFailedTests(tt.args.opts); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetFailedTests() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewGitOperator(t *testing.T) {
	type args struct {
		owner string
		repo  string
		token string
		ctx   context.Context
	}
	tests := []struct {
		name string
		args args
		want *GitOperator
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewGitOperator(tt.args.owner, tt.args.repo, tt.args.token, tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewGitOperator() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewUnauthenticatedGitOperator(t *testing.T) {
	type args struct {
		owner string
		repo  string
		ctx   context.Context
	}
	tests := []struct {
		name string
		args args
		want *GitOperator
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUnauthenticatedGitOperator(tt.args.owner, tt.args.repo, tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUnauthenticatedGitOperator() = %v, want %v", got, tt.want)
			}
		})
	}
}
