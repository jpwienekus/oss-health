package fetcher

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func fakeExecCommand(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestFakeHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

func TestFakeHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return

	}
	args := os.Args
	if strings.Contains(strings.Join(args, " "), "git clone") {
		os.Exit(0)
	}

	os.Exit(1)
}

func TestCanHandle(t *testing.T) {
	g := GithubFetcher{}

	cases := []struct {
		repoUrl string
		expect  bool
	}{
		{"git@github.com:user/repo.git", true},
		{"https://github.com/user/repo.git", true},
		{"github.com/user/repo", true},
		{"", false},
	}

	for _, c := range cases {
		result := g.CanHandle(c.repoUrl)
		if result != c.expect {
			t.Errorf("CanHandle(%q) = %v; want %v", c.repoUrl, result, c.expect)
		}
	}
}

func TestFetch_Success(t *testing.T) {
	g := GithubFetcher{}

	execCommand = fakeExecCommand
	defer func() { execCommand = exec.Command }()

	err := g.Fetch("git@github.com:user/repo.git", "/tmp/repo")
	if err != nil {
		t.Errorf("Fetch() returned error: %v", err)
	}
}
