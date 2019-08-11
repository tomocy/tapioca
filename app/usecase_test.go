package app

import (
	"fmt"
	"testing"
	"time"

	"github.com/tomocy/tapioca/domain"
)

func TestSummarizeCommitsOfToday(t *testing.T) {
	repo := newMock()
	uc := NewCommitUsecase(repo)
	expectedCs := repo.cs[:len(repo.cs)-1]
	expected := &domain.Summary{
		Repo: &domain.Repo{
			Owner: "mock",
			Name:  "mock",
		},
		Authors: []string{"alice", "bob", "cris"},
		Commits: expectedCs,
		Diff:    expectedCs.Diff(),
		Date:    today(),
	}
	actual, err := uc.SummarizeCommitsOfToday(expected.Repo.Owner, expected.Repo.Name)
	if err != nil {
		t.Fatalf("%s\n", reportUnexpected("error by SummarizeCommitsOfToday", err, nil))
	}
	if err := assertSummary(actual, expected); err != nil {
		t.Errorf("unexpected summary by SummarizeCommitsOfToday: %s\n", err)
	}
}

func TestSummarizeAuthorCommitsOfToday(t *testing.T) {
	repo := newMock()
	uc := NewCommitUsecase(repo)
	expectedCs := repo.cs[:1]
	expected := &domain.Summary{
		Repo: &domain.Repo{
			Owner: "mock",
			Name:  "mock",
		},
		Authors: []string{"alice"},
		Commits: expectedCs,
		Diff:    expectedCs.Diff(),
		Date:    today(),
	}
	actual, err := uc.SummarizeAuthorCommitsOfToday(expected.Repo.Owner, expected.Repo.Name, expected.Authors[0])
	if err != nil {
		t.Fatalf("%s\n", reportUnexpected("error by SummarizeCommitsOfToday", err, nil))
	}
	if err := assertSummary(actual, expected); err != nil {
		t.Errorf("unexpected summary by SummarizeCommitsOfToday: %s\n", err)
	}
}

func assertSummary(actual, expected *domain.Summary) error {
	if err := assertRepo(actual.Repo, expected.Repo); err != nil {
		return fmt.Errorf("unexpected repo of summary: %s", err)
	}
	if len(actual.Authors) != len(expected.Authors) {
		return reportUnexpected("len of authors of summary", len(actual.Authors), len(expected.Authors))
	}
	for i, expected := range expected.Authors {
		if actual.Authors[i] != expected {
			return reportUnexpected(fmt.Sprintf("authors[%d] of summary", i), actual.Authors[i], expected)
		}
	}
	if err := assertCommits(actual.Commits, expected.Commits); err != nil {
		return fmt.Errorf("unexpected commits of summary: %s", err)
	}
	if err := assertDiff(actual.Diff, expected.Diff); err != nil {
		return fmt.Errorf("unexpected diff of summary: %s", err)
	}
	if !actual.Date.Equal(expected.Date) {
		return reportUnexpected("date of summary", actual.Date, expected.Date)
	}

	return nil
}

func assertRepo(actual, expected *domain.Repo) error {
	if actual.Owner != expected.Owner {
		return reportUnexpected("owner of repo", actual.Owner, expected.Owner)
	}
	if actual.Name != expected.Name {
		return reportUnexpected("name of repo", actual.Name, expected.Name)
	}

	return nil
}

func assertCommits(actuals, expecteds []*domain.Commit) error {
	if len(actuals) != len(expecteds) {
		return reportUnexpected("len of commits", len(actuals), len(expecteds))
	}
	for i, expected := range expecteds {
		if err := assertCommit(actuals[i], expected); err != nil {
			return fmt.Errorf("unexpected commits[%d]: %s", i, err)
		}
	}

	return nil
}

func assertCommit(actual, expected *domain.Commit) error {
	if actual.ID != expected.ID {
		return reportUnexpected("id of commit", actual.ID, expected.ID)
	}
	if err := assertDiff(actual.Diff, expected.Diff); err != nil {
		return fmt.Errorf("unexpected diff of commit: %s", err)
	}

	return nil
}

func assertDiff(actual, expected *domain.Diff) error {
	if actual.Changes != expected.Changes {
		return reportUnexpected("changes of diff", actual.Changes, expected.Changes)
	}
	if actual.Adds != expected.Adds {
		return reportUnexpected("adds of diff", actual.Adds, expected.Adds)
	}
	if actual.Dels != expected.Dels {
		return reportUnexpected("dels of diff", actual.Dels, expected.Dels)
	}

	return nil
}

func reportUnexpected(name string, actual, expected interface{}) error {
	return report("unexpected "+name, actual, expected)
}

func report(name string, actual, expected interface{}) error {
	return fmt.Errorf("%s: got %v, expect %v", name, actual, expected)
}

func newMock() *mock {
	today := today().Add(1 * time.Second)
	return &mock{
		cs: domain.Commits{
			&domain.Commit{
				ID:     "a",
				Author: "alice",
				Diff: &domain.Diff{
					Changes: 3,
					Adds:    3,
				},
				CreatedAt: today,
			},
			&domain.Commit{
				ID:     "b",
				Author: "bob",
				Diff: &domain.Diff{
					Changes: 3,
					Adds:    2,
					Dels:    1,
				},
				CreatedAt: today,
			},
			&domain.Commit{
				ID:     "c",
				Author: "cris",
				Diff: &domain.Diff{
					Changes: 3,
					Adds:    1,
					Dels:    2,
				},
				CreatedAt: today,
			},
			&domain.Commit{
				ID:     "d",
				Author: "dave",
				Diff: &domain.Diff{
					Changes: 3,
					Dels:    3,
				},
			},
		},
	}
}

type mock struct {
	cs   domain.Commits
	date time.Time
}

func (m *mock) FetchCommits(owner, repo string, params domain.Params) (domain.Commits, error) {
	var fetcheds domain.Commits
	for _, c := range m.cs {
		if params.Author != "" && c.Author != params.Author {
			continue
		}
		if !c.CreatedAt.After(params.Since) {
			continue
		}

		fetcheds = append(fetcheds, c)
	}

	return fetcheds, nil
}
