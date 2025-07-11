package pipeline_test

// import (
// 	"errors"
// 	"os"
// 	"strings"
// 	"testing"
//
// 	"github.com/oss-health/background-worker/fetcher"
// )
//
// type mockFetcherChain struct {
// 	fetchers []fetcher.RepositoryFetcher
// 	expectedRepoURL string
// 	called          *bool
// 	returnError     error
// }
//
// func (m *mockFetcherChain) Fetch(repoUrl string, destPath string) error {
// 	if m.called != nil {
// 		*m.called = true
// 	}
// 	if repoUrl != m.expectedRepoURL {
// 		return errors.New("unexpected repo URL")
// 	}
// 	if !strings.Contains(destPath, "repo-clone-") {
// 		return errors.New("unexpected temp dir path")
// 	}
// 	return m.returnError
// }
//
// func TestPipeline_Run_Success(t *testing.T) {
// 	called := false
// 	mockChain := &mockFetcherChain{
// 		expectedRepoURL: "git@github.com:user/repo.git",
// 		called:          &called,
// 		returnError:     nil,
// 	}
//
// 	p := NewPipeline(mockChain)
//
// 	err := p.Run("git@github.com:user/repo.git")
// 	if err != nil {
// 		t.Fatalf("expected no error, got: %v", err)
// 	}
//
// 	if !called {
// 		t.Error("expected FetcherChain.Fetch to be called")
// 	}
// }
//
// func TestPipeline_Run_FetchFails(t *testing.T) {
// 	mockChain := &mockFetcherChain{
// 		expectedRepoURL: "git@github.com:user/repo.git",
// 		returnError:     errors.New("mock fetch failed"),
// 	}
//
// 	p := NewPipeline(mockChain)
//
// 	err := p.Run("git@github.com:user/repo.git")
// 	if err == nil || !strings.Contains(err.Error(), "failed to fetch repository") {
// 		t.Fatalf("expected fetch error, got: %v", err)
// 	}
// }
//
// func TestPipeline_Run_TempDirFails(t *testing.T) {
// 	// Backup original function
// 	defer func() { osMkdirTemp = os.MkdirTemp }()
//
// 	// Replace with mock
// 	osMkdirTemp = func(dir, pattern string) (string, error) {
// 		return "", errors.New("mock temp dir error")
// 	}
//
// 	p := NewPipeline(&mockFetcherChain{})
//
// 	err := p.Run("git@github.com:user/repo.git")
// 	if err == nil || !strings.Contains(err.Error(), "failed to create temp directory") {
// 		t.Fatalf("expected temp dir creation error, got: %v", err)
// 	}
// }
