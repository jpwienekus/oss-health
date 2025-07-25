package githubactivity

import (
	"context"
	"fmt"
	"time"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

type GithubClient struct {
	client *githubv4.Client
}

func NewGithubClient(token string, ctx context.Context) *GithubClient {
	src := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	httpClient := oauth2.NewClient(ctx, src)
	return &GithubClient{client: githubv4.NewClient(httpClient)}
}

type CommitEdge struct {
	Node struct {
		CommittedDate time.Time
	}
}

type IssueEdge struct {
	Node struct {
		UpdatedAt time.Time
		State     string
	}
}

type PRNode struct {
	CreatedAt time.Time
	MergedAt  time.Time
}

type RepositoryMetrics struct {
	Name         string
	Stars        int
	Contributers int
	Commits      []time.Time
	OpenIssues   int
	ClosedIssues int
	PullRequests []PRNode
	ReleaseDates []time.Time
}

type repositoryQuery struct {
	Repository struct {
		StargazerCount int
		// Issues         struct {
		// 	TotalCount int
		// } `graphql:"issues(states: OPEN)"`
		// ClosedIssues struct {
		// 	TotalCount int
		// } `graphql:"closedIssues: issues(states: CLOSED)"`
		Issues struct {
			// TotalCount int
			Edges []IssueEdge
		} `graphql:"issues(first: 100, orderBy: {field: UPDATED_AT, direction: DESC})"`
		PullRequests struct {
			Edges []struct {
				Node PRNode
			}
		} `graphql:"pullRequests(states: MERGED, first: 100, orderBy: {field: CREATED_AT, direction: DESC})"`
		DefaultBranchRef struct {
			Target struct {
				Commit struct {
					History struct {
						Edges []CommitEdge
					} `graphql:"history(first: 100, since: $since)"`
				} `graphql:"... on Commit"`
			}
		}
		Releases struct {
			Nodes []struct {
				CreatedAt time.Time
			}
		} `graphql:"releases(last: 100)"`
		MentionableUsers struct {
			TotalCount int
		} `graphql:"mentionableUsers(first:100)"`
	} `graphql:"repository(owner: $owner, name: $name)"`
}

func (c *GithubClient) FetchRepositoryData(ctx context.Context, owner string, name string, monthsBack int) (*RepositoryMetrics, error) {
	var q repositoryQuery
	variables := map[string]any{
		"owner": githubv4.String(owner),
		"name":  githubv4.String(name),
		"since": githubv4.GitTimestamp{Time: time.Now().AddDate(0, -monthsBack, 0)},
	}

	err := c.client.Query(ctx, &q, variables)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch github activity: %v", err)
	}

	since := time.Now().AddDate(0, -monthsBack, 0)

	// Flatten commits
	var commits []time.Time
	for _, edge := range q.Repository.DefaultBranchRef.Target.Commit.History.Edges {
		if edge.Node.CommittedDate.After(since) {
			commits = append(commits, edge.Node.CommittedDate)
		}
	}

	// Flatten PRs
	var prs []PRNode
	for _, edge := range q.Repository.PullRequests.Edges {
		if edge.Node.MergedAt.After(since) {
			prs = append(prs, edge.Node)
		}
	}

	open := 0
	closed := 0

	for _, edge := range q.Repository.Issues.Edges {
		if edge.Node.UpdatedAt.After(since) {
			if edge.Node.State == "OPEN" {
				open++
			} else if edge.Node.State == "CLOSED" {
				closed++
			}
		}
	}

	// Flatten releases
	var releases []time.Time
	for _, r := range q.Repository.Releases.Nodes {
		if r.CreatedAt.After(since) {
			releases = append(releases, r.CreatedAt)
		}
	}

	return &RepositoryMetrics{
		Stars:        q.Repository.StargazerCount,
		OpenIssues:   open,
		ClosedIssues: closed,
		PullRequests: prs,
		Commits:      commits,
		ReleaseDates: releases,
		Contributers: q.Repository.MentionableUsers.TotalCount,
	}, nil

}
