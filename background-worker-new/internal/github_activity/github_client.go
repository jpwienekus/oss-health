package githubactivity

import (
	"context"
	"fmt"
	"strings"
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

type CommitUser struct {
	Login string
}

type CommitAuthor struct {
	User CommitUser
}

type CommitEdge struct {
	Node struct {
		CommittedDate time.Time    `graphql:"committedDate"`
		Author        CommitAuthor `graphql:"author"`
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

type BlobObject struct {
	Text string `graphql:"text"`
}

type Tagger struct {
	Date time.Time `graphql:"date"`
}

type TagRef struct {
	Name   string
	Target struct {
		Tag struct {
			Tagger Tagger `graphql:"tagger"`
		} `graphql:"... on Tag"`
	}
}

type RepositoryMetrics struct {
	Name               string
	Stars              int
	Contributers       int
	ActiveContributors int
	Commits            []time.Time
	OpenIssues         int
	ClosedIssues       int
	PullRequests       []PRNode
	ReleaseDates       []time.Time
	Tags               []TagRef
	IssueActivityRatio float64
	ReleasesPastYear   int
	Readme             string
}

type repositoryQuery struct {
	Repository struct {
		StargazerCount int
		Issues         struct {
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
		} `graphql:"releases(first: 100, orderBy: {field: CREATED_AT, direction: DESC})"`
		MentionableUsers struct {
			TotalCount int
		} `graphql:"mentionableUsers(first:100)"`
		Object struct {
			BlobObject `graphql:"... on Blob"`
		} `graphql:"object(expression: \"HEAD:README.md\")"`
		Refs struct {
			Nodes []TagRef
		} `graphql:"refs(refPrefix: \"refs/tags/\", first: 50, orderBy: {field: TAG_COMMIT_DATE, direction: DESC})"`
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
	commits, activeContributors := commits(q, since)
	prs := pullRequests(q, since)
	open, closed, ratio := openClosedIssues(q, since)
	releases, releasesPastYear := releases(q, since)

	return &RepositoryMetrics{
		Stars:              q.Repository.StargazerCount,
		OpenIssues:         open,
		ClosedIssues:       closed,
		PullRequests:       prs,
		Commits:            commits,
		ReleaseDates:       releases,
		Contributers:       q.Repository.MentionableUsers.TotalCount,
		ActiveContributors: activeContributors,
		IssueActivityRatio: ratio,
		ReleasesPastYear:   releasesPastYear,
		Tags:               q.Repository.Refs.Nodes,
		Readme:             q.Repository.Object.Text,
	}, nil
}

func commits(q repositoryQuery, since time.Time) ([]time.Time, int) {
	seen := make(map[string]struct{})
	var commits []time.Time

	for _, edge := range q.Repository.DefaultBranchRef.Target.Commit.History.Edges {
		if edge.Node.CommittedDate.After(since) {
			commits = append(commits, edge.Node.CommittedDate)

			if edge.Node.Author.User.Login != "" {
				if strings.Contains(edge.Node.Author.User.Login, "bot") || strings.HasSuffix(edge.Node.Author.User.Login, "[bot]") {
					continue
				}

				seen[edge.Node.Author.User.Login] = struct{}{}
			}
		}
	}

	return commits, len(seen)
}

func pullRequests(q repositoryQuery, since time.Time) []PRNode {
	var prs []PRNode
	for _, edge := range q.Repository.PullRequests.Edges {
		if edge.Node.MergedAt.After(since) {
			prs = append(prs, edge.Node)
		}
	}

	return prs
}

func openClosedIssues(q repositoryQuery, since time.Time) (int, int, float64) {
	open := 0
	closed := 0
	active := 0

	for _, edge := range q.Repository.Issues.Edges {
		if edge.Node.UpdatedAt.After(since) {
			active++

			switch edge.Node.State {
			case "OPEN":
				open++
			case "CLOSED":
				closed++
			}
		}
	}

	return open, closed, (float64(active) / float64(len(q.Repository.Issues.Edges)))
}

func releases(q repositoryQuery, since time.Time) ([]time.Time, int) {
	var releases []time.Time
	cutoff := time.Now().AddDate(-1, 0, 0)
	releasesPastYear := 0

	for _, r := range q.Repository.Releases.Nodes {
		if r.CreatedAt.After(since) {
			releases = append(releases, r.CreatedAt)
		}

		if r.CreatedAt.After(cutoff) {
			releasesPastYear++
		}
	}

	return releases, releasesPastYear
}
