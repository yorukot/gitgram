package githubwebhook

import (
	"errors"
	"testing"

	"github.com/yorukot/gitgram/internal/activity"
)

func TestParseWorkflowRunIgnoresSuccess(t *testing.T) {
	body := []byte(`{
		"action": "completed",
		"repository": {"full_name": "owner/repo"},
		"sender": {"login": "octocat"},
		"workflow_run": {
			"name": "CI",
			"head_branch": "main",
			"conclusion": "success",
			"html_url": "https://github.com/owner/repo/actions/runs/1"
		}
	}`)

	_, err := ParseEvent(activity.EventWorkflowRun, "delivery-1", body)
	if !errors.Is(err, ErrIgnored) {
		t.Fatalf("expected ErrIgnored, got %v", err)
	}
}

func TestParseWorkflowRunFailure(t *testing.T) {
	body := []byte(`{
		"action": "completed",
		"repository": {"full_name": "owner/repo"},
		"sender": {"login": "octocat"},
		"workflow_run": {
			"name": "CI",
			"head_branch": "main",
			"conclusion": "failure",
			"html_url": "https://github.com/owner/repo/actions/runs/2"
		}
	}`)

	got, err := ParseEvent(activity.EventWorkflowRun, "delivery-2", body)
	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}
	if got.Event != activity.EventWorkflowRun {
		t.Fatalf("event = %q, want %q", got.Event, activity.EventWorkflowRun)
	}
	if got.Action != "failed" {
		t.Fatalf("action = %q, want failed", got.Action)
	}
	if got.Repo != "owner/repo" {
		t.Fatalf("repo = %q, want owner/repo", got.Repo)
	}
	if got.Branch != "main" {
		t.Fatalf("branch = %q, want main", got.Branch)
	}
}

func TestParsePushIgnored(t *testing.T) {
	body := []byte(`{
		"ref": "refs/heads/main",
		"compare": "https://github.com/owner/repo/compare/a...b",
		"repository": {"full_name": "owner/repo"},
		"sender": {"login": "octocat"},
		"commits": [
			{
				"id": "abcdef1234567890",
				"message": "Fix login\n\nLong body",
				"url": "https://github.com/owner/repo/commit/abcdef",
				"author": {"name": "Mona"}
			}
		]
	}`)

	_, err := ParseEvent(activity.EventPush, "delivery-3", body)
	if !errors.Is(err, ErrIgnored) {
		t.Fatalf("expected ErrIgnored, got %v", err)
	}
}

func TestParsePullRequestSynchronizeIgnored(t *testing.T) {
	body := []byte(`{
		"action": "synchronize",
		"number": 12,
		"repository": {"full_name": "owner/repo"},
		"sender": {"login": "octocat"},
		"pull_request": {
			"html_url": "https://github.com/owner/repo/pull/12",
			"title": "Add login",
			"number": 12,
			"head": {"ref": "feature/login"},
			"base": {"ref": "main"}
		}
	}`)

	_, err := ParseEvent(activity.EventPullRequest, "delivery-4", body)
	if !errors.Is(err, ErrIgnored) {
		t.Fatalf("expected ErrIgnored, got %v", err)
	}
}

func TestParsePullRequestReviewRequested(t *testing.T) {
	body := []byte(`{
		"action": "review_requested",
		"number": 12,
		"repository": {"full_name": "owner/repo"},
		"sender": {"login": "octocat"},
		"requested_reviewer": {"login": "mona"},
		"pull_request": {
			"html_url": "https://github.com/owner/repo/pull/12",
			"title": "Add login",
			"number": 12,
			"head": {"ref": "feature/login"},
			"base": {"ref": "main"}
		}
	}`)

	got, err := ParseEvent(activity.EventPullRequest, "delivery-5", body)
	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}
	if got.Action != "review requested" {
		t.Fatalf("action = %q, want review requested", got.Action)
	}
	if got.Summary != "review requested from mona" {
		t.Fatalf("summary = %q, want review requested from mona", got.Summary)
	}
}

func TestParseIssuesClosedIgnored(t *testing.T) {
	body := []byte(`{
		"action": "closed",
		"repository": {"full_name": "owner/repo"},
		"sender": {"login": "octocat"},
		"issue": {
			"html_url": "https://github.com/owner/repo/issues/42",
			"title": "Cannot login",
			"number": 42
		}
	}`)

	_, err := ParseEvent(activity.EventIssues, "delivery-6", body)
	if !errors.Is(err, ErrIgnored) {
		t.Fatalf("expected ErrIgnored, got %v", err)
	}
}

func TestParseIssueCommentIgnored(t *testing.T) {
	body := []byte(`{
		"action": "created",
		"repository": {"full_name": "owner/repo"},
		"sender": {"login": "octocat"},
		"issue": {
			"html_url": "https://github.com/owner/repo/issues/42",
			"title": "Cannot login",
			"number": 42
		},
		"comment": {
			"html_url": "https://github.com/owner/repo/issues/42#issuecomment-1",
			"body": "too noisy"
		}
	}`)

	_, err := ParseEvent(activity.EventIssueComment, "delivery-7", body)
	if !errors.Is(err, ErrIgnored) {
		t.Fatalf("expected ErrIgnored, got %v", err)
	}
}

func TestParsePullRequestReviewCommentedIgnored(t *testing.T) {
	body := []byte(`{
		"action": "submitted",
		"repository": {"full_name": "owner/repo"},
		"sender": {"login": "octocat"},
		"pull_request": {
			"html_url": "https://github.com/owner/repo/pull/12",
			"title": "Add login",
			"number": 12
		},
		"review": {
			"html_url": "https://github.com/owner/repo/pull/12#pullrequestreview-1",
			"state": "commented",
			"body": "just a comment"
		}
	}`)

	_, err := ParseEvent(activity.EventPullRequestReview, "delivery-8", body)
	if !errors.Is(err, ErrIgnored) {
		t.Fatalf("expected ErrIgnored, got %v", err)
	}
}

func TestParsePullRequestReviewChangesRequested(t *testing.T) {
	body := []byte(`{
		"action": "submitted",
		"repository": {"full_name": "owner/repo"},
		"sender": {"login": "octocat"},
		"pull_request": {
			"html_url": "https://github.com/owner/repo/pull/12",
			"title": "Add login",
			"number": 12
		},
		"review": {
			"html_url": "https://github.com/owner/repo/pull/12#pullrequestreview-1",
			"state": "changes_requested",
			"body": "Please add tests."
		}
	}`)

	got, err := ParseEvent(activity.EventPullRequestReview, "delivery-9", body)
	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}
	if got.Action != "changes_requested" {
		t.Fatalf("action = %q, want changes_requested", got.Action)
	}
	if got.Summary != "Please add tests." {
		t.Fatalf("summary = %q, want review body", got.Summary)
	}
}
