package service

import (
	"errors"
	"strings"
	"testing"

	"github.com/matteo-gildone/teamtime/internals/storage"
	"github.com/matteo-gildone/teamtime/internals/types"
)

func TestColleagueService_AddColleague(t *testing.T) {
	t.Run("add to empty list", func(t *testing.T) {
		svc, m := setUpTestService(t)
		colleague, err := svc.AddColleague("Alice", "London", "Europe/London")
		if err != nil {
			t.Fatalf("failed to add colleague: %v", err)
		}

		if colleague.Name != "Alice" {
			t.Errorf("got %q ,want name %q", colleague.Name, "Alice")
		}

		assertColleagueCount(t, m, 1)
	})

	t.Run("whitespace is trimmed", func(t *testing.T) {
		svc, m := setUpTestService(t)
		colleague, err := svc.AddColleague("    Alice    ", "London", "Europe/London")
		if err != nil {
			t.Fatalf("failed to add colleague: %v", err)
		}

		if colleague.Name != "Alice" {
			t.Errorf("got %q ,want name %q", colleague.Name, "Alice")
		}

		assertColleagueCount(t, m, 1)
	})

	t.Run("add to existing list", func(t *testing.T) {
		svc, m := setUpTestService(t)
		setupInitialColleagues(t, m, []types.Colleague{
			mustNewColleague(t, "Bob", "NYC", "America/New_York"),
		})

		colleague, err := svc.AddColleague("Matteo", "Tokyo", "Asia/Tokyo")
		if err != nil {
			t.Fatalf("failed to add colleague: %v", err)
		}

		if colleague.Name != "Matteo" {
			t.Errorf("got %q ,want name %q", colleague.Name, "Matteo")
		}

		assertColleagueCount(t, m, 2)
	})

	t.Run("validation error", func(t *testing.T) {
		tests := []struct {
			name        string
			addName     string
			addCity     string
			addTZ       string
			wantErrType error
		}{
			{
				name:        "missing name",
				addName:     "",
				addCity:     "London",
				addTZ:       "Europe/London",
				wantErrType: types.ErrMissingName,
			},
			{
				name:        "missing city",
				addName:     "Alice",
				addCity:     "",
				addTZ:       "Europe/London",
				wantErrType: types.ErrMissingCity,
			},
			{
				name:        "missing timezone",
				addName:     "Alice",
				addCity:     "London",
				addTZ:       "",
				wantErrType: types.ErrMissingTimezone,
			},
			{
				name:    "invalid timezone",
				addName: "Alice",
				addCity: "London",
				addTZ:   "Europe/Timezone",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				svc, m := setUpTestService(t)
				_, err := svc.AddColleague(tt.addName, tt.addCity, tt.addTZ)

				if err == nil {
					t.Fatal("expected error, got nil")
				}

				if tt.wantErrType != nil && !errors.Is(err, tt.wantErrType) {
					t.Errorf("got %v, want error type %v", err, tt.wantErrType)
				}

				assertColleagueCount(t, m, 0)
			})
		}
	})
}

func TestColleagueService_RemoveColleague(t *testing.T) {
	t.Run("remove first colleague", func(t *testing.T) {
		svc, m := setUpTestService(t)
		setupInitialColleagues(t, m, []types.Colleague{
			mustNewColleague(t, "Alice", "London", "Europe/London"),
			mustNewColleague(t, "Bob", "NYC", "America/New_York"),
		})

		removed, err := svc.RemoveColleague(1)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if removed.Name != "Alice" {
			t.Errorf("got %q, want %q", removed.Name, "Alice")
		}

		assertColleagueCount(t, m, 1)
	})

	t.Run("remove last colleague", func(t *testing.T) {
		svc, m := setUpTestService(t)
		setupInitialColleagues(t, m, []types.Colleague{
			mustNewColleague(t, "Alice", "London", "Europe/London"),
			mustNewColleague(t, "Bob", "NYC", "America/New_York"),
		})

		removed, err := svc.RemoveColleague(2)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if removed.Name != "Bob" {
			t.Errorf("got %q, want %q", removed.Name, "Bob")
		}

		assertColleagueCount(t, m, 1)
	})

	t.Run("remove only colleague", func(t *testing.T) {
		svc, m := setUpTestService(t)
		setupInitialColleagues(t, m, []types.Colleague{
			mustNewColleague(t, "Alice", "London", "Europe/London"),
		})

		removed, err := svc.RemoveColleague(1)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if removed.Name != "Alice" {
			t.Errorf("got %q, want %q", removed.Name, "Alice")
		}

		assertColleagueCount(t, m, 0)
	})

	t.Run("error cases", func(t *testing.T) {
		tests := []struct {
			name          string
			initial       []types.Colleague
			removeIndex   int
			wantErrorType error
		}{
			{
				name:          "empty list",
				initial:       []types.Colleague{},
				removeIndex:   1,
				wantErrorType: types.ErrEmptyList,
			},
			{
				name: "index zero",
				initial: []types.Colleague{
					mustNewColleague(t, "Alice", "London", "Europe/London"),
				},
				removeIndex:   0,
				wantErrorType: types.ErrorInvalidIndex,
			},
			{
				name: "negative index",
				initial: []types.Colleague{
					mustNewColleague(t, "Alice", "London", "Europe/London"),
				},
				removeIndex:   -1,
				wantErrorType: types.ErrorInvalidIndex,
			},
			{
				name: "index too large",
				initial: []types.Colleague{
					mustNewColleague(t, "Alice", "London", "Europe/London"),
				},
				removeIndex:   5,
				wantErrorType: types.ErrorInvalidIndex,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				svc, m := setUpTestService(t)
				setupInitialColleagues(t, m, tt.initial)
				_, err := svc.RemoveColleague(tt.removeIndex)

				if err == nil {
					t.Fatal("expected error got nil")
				}

				if !errors.Is(err, tt.wantErrorType) {
					t.Errorf("got %v, want %v", err, tt.wantErrorType)
				}
			})
		}
	})
}

func TestColleagueService_AllColleagues(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		svc, _ := setUpTestService(t)

		colleagues, err := svc.AllColleagues()

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(colleagues) != 0 {
			t.Errorf("got: %d, want empty list", len(colleagues))
		}
	})

	t.Run("single colleague", func(t *testing.T) {
		svc, m := setUpTestService(t)

		setupInitialColleagues(t, m, []types.Colleague{
			mustNewColleague(t, "Alice", "London", "Europe/London"),
		})
		colleagues, err := svc.AllColleagues()

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(colleagues) != 1 {
			t.Errorf("got: %d, want 1 colleague", len(colleagues))
		}

		if colleagues[0].Name != "Alice" {
			t.Errorf("got: %q, want %q", colleagues[0].Name, "Alice")
		}
	})

	t.Run("multiple colleagues", func(t *testing.T) {
		svc, m := setUpTestService(t)

		setupInitialColleagues(t, m, []types.Colleague{
			mustNewColleague(t, "Alice", "London", "Europe/London"),
			mustNewColleague(t, "Bob", "NYC", "America/New_York"),
			mustNewColleague(t, "Matteo", "Tokyo", "Asia/Tokyo"),
		})

		colleagues, err := svc.AllColleagues()

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(colleagues) != 3 {
			t.Errorf("got: %d, want 3 colleague", len(colleagues))
		}

		wantNames := []string{"Alice", "Bob", "Matteo"}
		for i, want := range wantNames {
			if colleagues[i].Name != want {
				t.Errorf("got: %q, want %q", colleagues[i].Name, want)
			}
		}
	})
}

func TestColleagueService_FindColleague(t *testing.T) {
	t.Run("single match", func(t *testing.T) {
		svc, m := setUpTestService(t)
		setupInitialColleagues(t, m, []types.Colleague{
			mustNewColleague(t, "Alice", "London", "Europe/London"),
			mustNewColleague(t, "Bob", "NYC", "America/New_York"),
		})

		results, err := svc.FindColleague("Alice")

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(results) != 1 {
			t.Errorf("got: %d, want 1", len(results))
		}

		if results[0].Name != "Alice" {
			t.Errorf("got: %q, want %q", results[0].Name, "Alice")
		}
	})

	t.Run("case insensitive", func(t *testing.T) {
		svc, m := setUpTestService(t)
		setupInitialColleagues(t, m, []types.Colleague{
			mustNewColleague(t, "Alice", "London", "Europe/London"),
		})

		results, err := svc.FindColleague("alice")

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(results) != 1 {
			t.Errorf("got: %d, want 1", len(results))
		}

		if results[0].Name != "Alice" {
			t.Errorf("got: %q, want %q", results[0].Name, "Alice")
		}
	})

	t.Run("no matches", func(t *testing.T) {
		svc, m := setUpTestService(t)
		setupInitialColleagues(t, m, []types.Colleague{
			mustNewColleague(t, "Alice", "London", "Europe/London"),
			mustNewColleague(t, "Bob", "NYC", "America/New_York"),
		})

		colleagues, _ := svc.FindColleague("Matteo")

		if len(colleagues) != 0 {
			t.Fatalf("expected empty list got: %d", len(colleagues))
		}
	})

	t.Run("multiple matches", func(t *testing.T) {
		svc, m := setUpTestService(t)
		setupInitialColleagues(t, m, []types.Colleague{
			mustNewColleague(t, "Alice", "London", "Europe/London"),
			mustNewColleague(t, "Alia", "NYC", "America/New_York"),
		})

		results, err := svc.FindColleague("Ali")

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(results) != 2 {
			t.Errorf("got: %d, want 0", len(results))
		}

		for _, result := range results {
			if !strings.Contains(result.Name, "Ali") {
				t.Errorf("got: %q, want %q", result.Name, "Ali")
			}
		}
	})
}

func TestColleagueService_Integration(t *testing.T) {
	svc, m := setUpTestService(t)

	svc.AddColleague("Alice", "London", "Europe/London")
	svc.AddColleague("Bob", "NYC", "America/New_York")
	assertColleagueCount(t, m, 2)

	results, _ := svc.FindColleague("Alice")
	if len(results) != 1 {
		t.Fatalf("got: %d want 1", len(results))
	}

	svc.RemoveColleague(1)
	assertColleagueCount(t, m, 1)

	all, _ := svc.AllColleagues()

	if all[0].Name != "Bob" {
		t.Errorf("got: %q want: Bob", all[0].Name)
	}
}

func setUpTestService(t *testing.T) (*ColleagueService, *storage.Manager) {
	t.Helper()
	tempDir := t.TempDir()
	m, err := storage.NewManager(tempDir)

	if err != nil {
		t.Fatalf("failed to setup test service: %v", err)
	}

	if err := m.EnsureFolder(); err != nil {
		t.Fatalf("failed to ensure folder: %v", err)
	}

	return NewColleagueService(m), m
}

func setupInitialColleagues(t *testing.T, m *storage.Manager, colleagues []types.Colleague) {
	t.Helper()
	if len(colleagues) == 0 {
		return
	}

	cl := types.NewColleagues()
	for _, c := range colleagues {
		cl.Add(c)
	}

	if err := m.Save(cl); err != nil {
		t.Fatalf("failed to save initial state: %v", err)
	}
}

func assertColleagueCount(t *testing.T, m *storage.Manager, want int) {
	t.Helper()
	loaded, err := m.Load()
	if err != nil {
		t.Fatalf("failed to load colleagues: %v", err)
	}

	if got := len(*loaded); got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}

func mustNewColleague(t *testing.T, name, city, tz string) types.Colleague {
	t.Helper()
	colleague, err := types.NewColleague(name, city, tz)
	if err != nil {
		t.Fatalf("failed to create test colleague: %v", err)
	}
	return colleague
}
