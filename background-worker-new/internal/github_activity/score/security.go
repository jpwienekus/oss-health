package score

import (
	"time"

	"github.com/oss-health/background-worker/internal/vulnerabilities"
)

func ScoreSecurity(vulnerabilities []vulnerabilities.Vulnerability, since int) int {
	if len(vulnerabilities) == 0 {
		return 10
	}

	highCritical := 0
	medium := 0
	low := 0

	for _, v := range vulnerabilities {
		switch v.Severity {
		case "HIGH", "CRITICAL":
			highCritical++
		case "MEDIUM":
			medium++
		case "LOW":
			low++
		}

		// if v.Published.After(mostRecent) {
		// 	mostRecent = v.Published
		// }

	}
	now := time.Now()
	var mostRecent time.Time
	mostRecentVulnOlderThan6Months := mostRecent.Before(now.AddDate(0, since, 0))

	return calculateSecurityScore(highCritical, medium, low, mostRecentVulnOlderThan6Months)
}

func calculateSecurityScore(highCritical int, medium int, low int, mostRecentVulnOlderThan6Months bool) int {
	score := 0

	if highCritical == 0 {
		score += 5
	}

	if medium <= 2 {
		score += 3
	}

	if low <= 3 {
		score += 1
	}

	if mostRecentVulnOlderThan6Months {
		score += 1
	}

	if score > 10 {
		score = 10
	}

	return score
}
