package score

import (
	"log"

	githubactivity "github.com/oss-health/background-worker/internal/github_activity"
)

func ScoreCommunity(metrics *githubactivity.RepositoryMetrics, since int) int {
	score := calculateCommunityScore(metrics.Stars, metrics.ActiveContributors, metrics.IssueActivityRatio, false)
	log.Printf("Stars: %d, Active Contributors: %d, Issue Activity Ratio: %.2f, Has Forum: %t", metrics.Stars, metrics.ActiveContributors, metrics.IssueActivityRatio, false)
	log.Printf("Community Score: %d", score)
	return score
}

func calculateCommunityScore(stars int, activeContributors int, issueActivitityRatio float64, hasForum bool) int {
	score := 0

	if stars > 1000 {
		score += 2
	}

	if activeContributors > 5 {
		score += 4
	}

	if issueActivitityRatio > 0.5 {
		score += 2
	}

	if hasForum {
		score += 2
	}

	if score > 10 {
		score = 10
	}

	return score
}
