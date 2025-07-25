package score

import (
	"log"
	"time"

	githubactivity "github.com/oss-health/background-worker/internal/github_activity"
)

func ScoreMaintenance(metrics *githubactivity.RepositoryMetrics, since int) int {
	lastCommitAgeInDays := float64(-1)
	commitsPerMonth := float64(-1)

	if len(metrics.Commits) > 0 {
		lastCommitAgeInDays = float64(time.Since(metrics.Commits[0]).Hours() / 24)
		commitsPerMonth = float64(len(metrics.Commits)) / float64(since)
	}

	issueClosureRate := float64(metrics.ClosedIssues) / float64(metrics.ClosedIssues+metrics.OpenIssues)
	avgPRMergeTimeInDays := float64(-1)

	if len(metrics.PullRequests) > 0 {
		var totalPRMergeDays float64
		for _, pr := range metrics.PullRequests {
			totalPRMergeDays += float64(pr.MergedAt.Sub(pr.CreatedAt).Hours() / 24)
		}
		avgPRMergeTimeInDays = totalPRMergeDays / float64(len(metrics.PullRequests))
	}

	score := calculateScore(lastCommitAgeInDays, commitsPerMonth, issueClosureRate, avgPRMergeTimeInDays)
	log.Printf("Last Commit Age (days): %.2f, Commits Per Month: %.2f, Issue Closure Rate: %.2f, Average PR Merge Time (days): %.2f", lastCommitAgeInDays, commitsPerMonth, issueClosureRate, avgPRMergeTimeInDays)
	log.Printf("Maintenance Score: %d", calculateScore(lastCommitAgeInDays, commitsPerMonth, issueClosureRate, avgPRMergeTimeInDays))

	return score
}

func calculateScore(lastCommitAgeInDays float64, commitsPerMonth float64, issueClosureRate float64, avgPRMergeDays float64) int {
	score := 0

	if lastCommitAgeInDays > -1 && lastCommitAgeInDays < 30 {
		score += 3
	}

	if commitsPerMonth > 10 {
		score += 2
	}

	if issueClosureRate > 0.7 {
		score += 2
	}

	if avgPRMergeDays > -1 && avgPRMergeDays < 7 {
		score += 3
	}

	// This should never happen, but cap just incase
	if score > 10 {
		score = 10
	}

	return score
}
