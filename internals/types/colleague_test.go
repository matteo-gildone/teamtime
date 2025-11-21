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
		cl.Add(colleagueName, colleagueCity, colleagueTZ)

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

func TestColleagueList_Delete(t *testing.T) {
	t.Run("delete middle colleague", func(t *testing.T) {
		cl := ColleagueList{
			{Name: "Alice", City: "London", Timezone: "Europe/London"},
			{Name: "Bob", City: "NYC", Timezone: "America/New_York"},
			{Name: "Daisuke", City: "Tokio", Timezone: "Asia/Tokio"},
		}

		if err := cl.Delete(2); err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if len(cl) != 2 {
			t.Errorf("expected list length %d, got %d instead.", 2, len(cl))
		}

		if cl[0].Name != "Alice" {
			t.Errorf("expected %q, got %q instead.", "Alice", cl[0].Name)
		}

		if cl[1].Name != "Daisuke" {
			t.Errorf("expected %q, got %q instead.", "Daisuke", cl[1].Name)
		}
	})

	t.Run("delete first colleague", func(t *testing.T) {
		cl := ColleagueList{
			{Name: "Alice", City: "London", Timezone: "Europe/London"},
			{Name: "Bob", City: "NYC", Timezone: "America/New_York"},
		}

		if err := cl.Delete(1); err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if len(cl) != 1 {
			t.Errorf("expected list length %d, got %d instead.", 1, len(cl))
		}

		if cl[0].Name != "Bob" {
			t.Errorf("expected Bob, got %s instead.", cl[0].Name)
		}
	})

	t.Run("delete last colleague", func(t *testing.T) {
		cl := ColleagueList{
			{Name: "Alice", City: "London", Timezone: "Europe/London"},
			{Name: "Bob", City: "NYC", Timezone: "America/New_York"},
		}

		if err := cl.Delete(2); err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if len(cl) != 1 {
			t.Errorf("expected list length %d, got %d instead.", 1, len(cl))
		}

		if cl[0].Name != "Alice" {
			t.Errorf("expected Alice, got %s instead.", cl[0].Name)
		}
	})

	t.Run("delete with index 0", func(t *testing.T) {
		cl := ColleagueList{
			{Name: "Alice", City: "London", Timezone: "Europe/London"},
			{Name: "Bob", City: "NYC", Timezone: "America/New_York"},
		}

		err := cl.Delete(0)

		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !errors.Is(err, ErrorInvalidIndex) {
			t.Errorf("expected ErrorInvalidIndex, got %v", err)
		}
	})

	t.Run("delete with negative index", func(t *testing.T) {
		cl := ColleagueList{
			{Name: "Alice", City: "London", Timezone: "Europe/London"},
			{Name: "Bob", City: "NYC", Timezone: "America/New_York"},
		}

		err := cl.Delete(-1)

		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !errors.Is(err, ErrorInvalidIndex) {
			t.Errorf("expected ErrorInvalidIndex, got %v", err)
		}
	})

	t.Run("delete with index too large", func(t *testing.T) {
		cl := ColleagueList{
			{Name: "Alice", City: "London", Timezone: "Europe/London"},
			{Name: "Bob", City: "NYC", Timezone: "America/New_York"},
		}

		err := cl.Delete(5)

		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !errors.Is(err, ErrorInvalidIndex) {
			t.Errorf("expected ErrorInvalidIndex, got %v", err)
		}
	})

	t.Run("delete from empty list", func(t *testing.T) {
		cl := ColleagueList{}

		err := cl.Delete(1)

		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !errors.Is(err, ErrEmptyList) {
			t.Errorf("expected ErrEmptyList, got %v", err)
		}
	})
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
