package types

import (
	"errors"
	"testing"
)

func TestColleagueList_Add(t *testing.T) {
	t.Run("add to empty list", func(t *testing.T) {
		cl := ColleagueList{}
		colleagueName := "Maurizio"
		colleagueCity := "Bari"
		colleagueTZ := "Europe/Rome"
		err := cl.Add(colleagueName, colleagueCity, colleagueTZ)

		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if len(cl) != 1 {
			t.Errorf("expected %d, got %d instead.", 1, len(cl))
		}

		if cl[0].Name != colleagueName {
			t.Errorf("expected %q, got %q instead.", colleagueName, cl[0].Name)
		}

		if cl[0].City != colleagueCity {
			t.Errorf("expected %q, got %q instead.", colleagueCity, cl[0].City)
		}

		if cl[0].Timezone != colleagueTZ {
			t.Errorf("expected %q, got %q instead.", colleagueTZ, cl[0].Timezone)
		}
	})

	t.Run("add multiple colleague", func(t *testing.T) {
		colleagues := [][]string{
			{"Alice", "London", "Europe/London"},
			{"Bob", "Manchester", "Europe/London"},
			{"Sabina", "Verona", "Europe/Rome"},
		}
		cl := ColleagueList{}

		for _, colleague := range colleagues {
			err := cl.Add(colleague[0], colleague[1], colleague[2])
			if err != nil {
				t.Errorf("%v", err)
			}
		}

		if len(cl) != 3 {
			t.Errorf("expected list length %d, got %d instead.", 3, len(cl))
		}

		for i, colleague := range colleagues {
			if cl[i].Name != colleague[0] {
				t.Errorf("expected %q, got %q instead.", colleague[0], cl[i].Name)
			}

			if cl[i].City != colleague[1] {
				t.Errorf("expected %q, got %q instead.", colleague[1], cl[i].City)
			}

			if cl[i].Timezone != colleague[2] {
				t.Errorf("expected %q, got %q instead.", colleague[2], cl[i].Timezone)
			}
		}

	})
}

func TestColleague_Add_Validation(t *testing.T) {
	tests := []struct {
		name      string
		inputName string
		inputCity string
		inputTZ   string
		wantErr   error
	}{
		{
			name:      "missing name",
			inputName: "",
			inputCity: "Bari",
			inputTZ:   "Europe/Rome",
			wantErr:   ErrMissingName,
		},
		{
			name:      "missing city",
			inputName: "Mariolino",
			inputCity: "",
			inputTZ:   "Europe/Rome",
			wantErr:   ErrMissingCity,
		},
		{
			name:      "missing timezone",
			inputName: "Mariolino",
			inputCity: "Aprilia",
			inputTZ:   "",
			wantErr:   ErrMissingTimezone,
		},
		{
			name:      "whitespace only name",
			inputName: "   ",
			inputCity: "Bari",
			inputTZ:   "Europe/Rome",
			wantErr:   ErrMissingName,
		},
		{
			name:      "whitespace only city",
			inputName: "Gregorio",
			inputCity: "   ",
			inputTZ:   "Europe/Rome",
			wantErr:   ErrMissingCity,
		},
		{
			name:      "whitespace only timezone",
			inputName: "Gregorio",
			inputCity: "Alba",
			inputTZ:   "   ",
			wantErr:   ErrMissingTimezone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cl := ColleagueList{}
			err := cl.Add(tt.inputName, tt.inputCity, tt.inputTZ)

			if err == nil {
				t.Fatal("expected error, got nil")
			}

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("expected %v, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestColleague_Add_InvalidTimezones(t *testing.T) {
	tests := []struct {
		name    string
		inputTZ string
	}{
		{
			name:    "typo in timezone",
			inputTZ: "Asia/Tokio",
		},
		{
			name:    "completely invalid",
			inputTZ: "NotaTimeZone",
		},
		{
			name:    "partial timezone",
			inputTZ: "Europe/",
		},
		{
			name:    "numbers only",
			inputTZ: "12345",
		},
		{
			name:    "special characters",
			inputTZ: "Europe/@£$",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cl := ColleagueList{}
			err := cl.Add("Test", "City", tt.inputTZ)

			if err == nil {
				t.Fatalf("expected error for timezone %s, got nil", tt.inputTZ)
			}

			if err.Error() == "" {
				t.Error("expected non-empty error message")
			}
		})
	}
}

func TestColleagueList_Delete_Success(t *testing.T) {
	tests := []struct {
		name              string
		initial           ColleagueList
		deleteIndex       int
		expectedDeleted   string
		expectedLength    int
		expectedRemaining []string
	}{
		{
			name: "delete first colleague",
			initial: ColleagueList{
				{Name: "Alice", City: "London", Timezone: "Europe/London"},
				{Name: "Bob", City: "NYC", Timezone: "America/New_York"},
			},
			deleteIndex:       1,
			expectedDeleted:   "Alice",
			expectedLength:    1,
			expectedRemaining: []string{"Bob"},
		},
		{
			name: "delete middle colleague",
			initial: ColleagueList{
				{Name: "Alice", City: "London", Timezone: "Europe/London"},
				{Name: "Bob", City: "NYC", Timezone: "America/New_York"},
				{Name: "Daisuke", City: "Tokyo", Timezone: "Asia/Tokyo"},
			},
			deleteIndex:       2,
			expectedDeleted:   "Bob",
			expectedLength:    2,
			expectedRemaining: []string{"Alice", "Daisuke"},
		},
		{
			name: "delete last colleague",
			initial: ColleagueList{
				{Name: "Alice", City: "London", Timezone: "Europe/London"},
				{Name: "Bob", City: "NYC", Timezone: "America/New_York"},
			},
			deleteIndex:       2,
			expectedDeleted:   "Bob",
			expectedLength:    1,
			expectedRemaining: []string{"Alice"},
		},
		{
			name: "delete only colleague",
			initial: ColleagueList{
				{Name: "Alice", City: "London", Timezone: "Europe/London"},
			},
			deleteIndex:       1,
			expectedDeleted:   "Alice",
			expectedLength:    0,
			expectedRemaining: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cl := make(ColleagueList, len(tt.initial))
			copy(cl, tt.initial)

			deleted, err := cl.Delete(tt.deleteIndex)
			if err != nil {
				t.Fatalf("expected no error, got: %v", err)
			}

			if deleted.Name != tt.expectedDeleted {
				t.Errorf("expected to delete %s, got: %s", tt.expectedDeleted, deleted.Name)
			}

			if len(cl) != tt.expectedLength {
				t.Errorf("expected list of length %d, got %d", tt.expectedLength, len(cl))
			}

			for i, expectedName := range tt.expectedRemaining {
				if i >= len(cl) {
					t.Errorf("expected colleague %q at index %d, but list is too short", expectedName, i)
					continue
				}

				if cl[i].Name != expectedName {
					t.Errorf("expected %q at index %d, got %q", expectedName, i, cl[i].Name)
				}
			}
		})
	}
}

func TestColleagueList_Delete_Errors(t *testing.T) {
	tests := []struct {
		name        string
		initial     ColleagueList
		deleteIndex int
		wantErr     error
	}{
		{
			name:        "remove from empty list",
			initial:     ColleagueList{},
			deleteIndex: 1,
			wantErr:     ErrEmptyList,
		},
		{
			name:        "index 0",
			initial:     ColleagueList{{Name: "Alice", City: "London", Timezone: "Europe/London"}},
			deleteIndex: 0,
			wantErr:     ErrorInvalidIndex,
		},
		{
			name:        "negative index",
			initial:     ColleagueList{{Name: "Alice", City: "London", Timezone: "Europe/London"}},
			deleteIndex: -1,
			wantErr:     ErrorInvalidIndex,
		},
		{
			name: "index too large",
			initial: ColleagueList{
				{Name: "Alice", City: "London", Timezone: "Europe/London"},
				{Name: "Bob", City: "NYC", Timezone: "America/New_York"},
			},
			deleteIndex: 5,
			wantErr:     ErrorInvalidIndex,
		},
		{
			name: "index one past end",
			initial: ColleagueList{
				{Name: "Alice", City: "London", Timezone: "Europe/London"},
			},
			deleteIndex: 2,
			wantErr:     ErrorInvalidIndex,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cl := make(ColleagueList, len(tt.initial))
			copy(cl, tt.initial)

			_, err := cl.Delete(tt.deleteIndex)
			if err == nil {
				t.Fatal("expected error, got nil")
			}

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("expected error %v, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestNewColleagues(t *testing.T) {
	cl := NewColleagues()

	if cl == nil {
		t.Fatal("expected non-nil ColleagueList, got nil")
	}

	if len(*cl) != 0 {
		t.Errorf("expected empty list, got length %d", len(*cl))
	}
}
