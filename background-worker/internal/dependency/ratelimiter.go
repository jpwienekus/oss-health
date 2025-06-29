package dependency

import (
	"golang.org/x/time/rate"

	"time"

	"github.com/oss-health/background-worker/internal/utils"
)

const (
	NpmRps    = 10
	NpmBurst  = 10
	PypiRps   = 1
	PypiBurst = 1
)

func InitRateLimiters(rateLimiter *utils.DefaultRateLimiter) {
	registerRateLimter(rateLimiter, "npm", NpmRps, NpmBurst)
	registerRateLimter(rateLimiter, "pypi", PypiRps, PypiBurst)
}

func registerRateLimter(rateLimiter *utils.DefaultRateLimiter, registry string, rps int, burst int) {
	periodPerRequest := time.Second / time.Duration(rps)
	rateLimiter.RegisterLimiter(registry, rate.Every(periodPerRequest), burst)
}
