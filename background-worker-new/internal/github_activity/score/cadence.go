package score

import githubactivity "github.com/oss-health/background-worker/internal/github_activity"

func ScoreCadence(metrics *githubactivity.RepositoryMetrics, since int) int {
	return calculateCadenceScore(0, 0.0, false, false)
}

func calculateCadenceScore(daysSinceLastRelease int, releasesPerYear float64, usesSemver bool, consistentTagging bool) int {
	score := 0

	if daysSinceLastRelease < 60 {
		score += 3
	}

	if releasesPerYear > 4 {
		score += 4
	}

	if usesSemver {
		score += 2
	}

	if consistentTagging {
		score += 1
	}

	if score > 10 {
		score = 10
	}

	return score
}

