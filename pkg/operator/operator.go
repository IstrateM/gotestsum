//Operator package hosts operators needed to interact with different applications or files
package operator

import (
	"context"

	"github.com/google/go-github/github"
	"github.com/istratem/gotestsum/pkg/options"
	"golang.org/x/oauth2"
)

const FailureTag = "test failure"

//FailedTest is used to store information about failed test
type FailedTest struct {
	Title   string
	Issues  string
	IssueNo int
}

//GitOperator is an object that will do basic operations on GitHub
type GitOperator struct {
	client *github.Client
	ctx    context.Context
	owner  string
	repo   string
}

//NewGitOperator provides an implementation of a GitHub Operator with authentication.
func NewGitOperator(owner, repo, token string, ctx context.Context) *GitOperator {
	g := &GitOperator{
		ctx:   ctx,
		owner: owner,
		repo:  repo,
	}
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	g.client = github.NewClient(tc)
	return g
}

//NewUnauthenticatedGitOperator provides an implementation of a GitHub Operator without authentication.
func NewUnauthenticatedGitOperator(owner, repo string, ctx context.Context) *GitOperator {
	return &GitOperator{
		ctx:    ctx,
		owner:  owner,
		repo:   repo,
		client: github.NewClient(nil),
	}
}

//GetTestIssues extracts test related issues from github repository and converts them into a slice of FailedTests.
//Extracted issues are opened and labeled with 'test failure' label.
func (g *GitOperator) GetTestIssues() ([]*FailedTest, error) {
	issues, err := g.GetIssuesByLabel(FailureTag)
	if err != nil {
		return nil, err
	}
	var res []*FailedTest
	for _, i := range issues {
		res = append(res, &FailedTest{
			Title:   *i.Title,
			Issues:  *i.Body,
			IssueNo: *i.Number,
		})
	}
	return res, nil
}

//GetIssuesByLabel extracts all the open issues labeled with the given label.
//It returns a slice of github.Issues and an error.
func (g *GitOperator) GetIssuesByLabel(label string) ([]*github.Issue, error) {
	list, _, err := g.client.Issues.ListByRepo(g.ctx, g.owner, g.repo, &github.IssueListByRepoOptions{
		State:  "open",
		Labels: []string{label},
	})
	if err != nil {
		return nil, err
	}
	return list, err
}

//GetIssuesByLabel pots a new github issue based on a failing test.
func (g *GitOperator) PostNewIssue(f *FailedTest) error {
	_, _, err := g.client.Issues.Create(g.ctx, g.owner, g.repo, g.newIssueForTest(f))
	return err
}

//newIssueForTest creates a new issue from a FailedTest.
func (g *GitOperator) newIssueForTest(f *FailedTest) *github.IssueRequest {
	return &github.IssueRequest{
		Title:  &f.Title,
		Body:   &f.Issues,
		Labels: &[]string{FailureTag},
	}
}

//CloseSolvedIssue closes the open issue based on given test.
func (g *GitOperator) CloseSolvedIssue(f *FailedTest) error {
	req := g.newIssueForTest(f)
	s := "closed"
	req.State = &s
	_, _, err := g.client.Issues.Edit(g.ctx, g.owner, g.repo, f.IssueNo, req)
	return err
}

//JUnitOperator is an object that will do basic operations on JUnit related objects
type JUnitOperator struct {
}

//GetFailedTests extracts FailedTests from a GetFailedTests
func (_ *JUnitOperator) GetFailedTests(opts *options.Options) []*FailedTest {
	var ft []*FailedTest
	for _, s := range opts.JUnitTestSuite.Suites {
		if s.Failures == 0 {
			continue
		}
		for _, t := range s.TestCases {
			if t.Failure != nil {
				ft = append(ft, &FailedTest{
					Title:  s.Name + "/" + t.Classname + "/" + t.Name,
					Issues: t.Failure.Contents,
				})
			}
		}
	}
	return ft
}
