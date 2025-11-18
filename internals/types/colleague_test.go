package types

import (
	"errors"
	"testing"
)

func TestColleagueList_Add(t *testing.T) {
	t.Run("Add to empty list", func(t *testing.T) {
		cl := ColleagueList{}
		colleagueName := "Maurizio"
		colleagueCity := "Bari"
		colleagueTZ := "Europe/Rome"
		cl.Add(colleagueName, colleagueCity, colleagueTZ)

		if len(cl) != 1 {
			t.Errorf("Expected %d, got %d instead.", 1, len(cl))
		}

		if cl[0].Name != colleagueName {
			t.Errorf("Expected %q, got %q instead.", colleagueName, cl[0].Name)
		}

		if cl[0].City != colleagueCity {
			t.Errorf("Expected %q, got %q instead.", colleagueCity, cl[0].City)
		}

		if cl[0].Timezone != colleagueTZ {
			t.Errorf("Expected %q, got %q instead.", colleagueTZ, cl[0].Timezone)
		}
	})

	t.Run("Add multiple colleague", func(t *testing.T) {
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
			t.Errorf("Expected list length %d, got %d instead.", 3, len(cl))
		}

		for i, colleague := range colleagues {
			if cl[i].Name != colleague[0] {
				t.Errorf("Expected %q, got %q instead.", colleague[0], cl[i].Name)
			}

			if cl[i].City != colleague[1] {
				t.Errorf("Expected %q, got %q instead.", colleague[1], cl[i].City)
			}

			if cl[i].Timezone != colleague[2] {
				t.Errorf("Expected %q, got %q instead.", colleague[2], cl[i].Timezone)
			}
		}

	})

	t.Run("Add colleague without a name", func(t *testing.T) {
		cl := ColleagueList{}
		colleagueName := ""
		colleagueCity := "Bari"
		colleagueTZ := "Europe/Rome"
		err := cl.Add(colleagueName, colleagueCity, colleagueTZ)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		if !errors.Is(err, ErrMissingName) {
			t.Errorf("expected ErrMissingName, got %v", err)
		}
	})

	t.Run("Add colleague without a city", func(t *testing.T) {
		cl := ColleagueList{}
		colleagueName := "Giulio"
		colleagueCity := ""
		colleagueTZ := "Europe/Rome"
		err := cl.Add(colleagueName, colleagueCity, colleagueTZ)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		if !errors.Is(err, ErrMissingCity) {
			t.Errorf("expected ErrMissingCity, got %v", err)
		}
	})

	t.Run("Add colleague without timezone", func(t *testing.T) {
		cl := ColleagueList{}
		colleagueName := "Giulio"
		colleagueCity := "Bari"
		colleagueTZ := ""
		err := cl.Add(colleagueName, colleagueCity, colleagueTZ)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		if !errors.Is(err, ErrMissingTimezone) {
			t.Errorf("expected ErrMissingTimezone, got %v", err)
		}
	})

	t.Run("Add colleague with unknown timezone", func(t *testing.T) {
		cl := ColleagueList{}
		colleagueName := "Giulio"
		colleagueCity := "Bari"
		colleagueTZ := "Asia/Tokio"
		err := cl.Add(colleagueName, colleagueCity, colleagueTZ)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		if err.Error() != "unknown time zone Asia/Tokio" {
			t.Errorf("%v", err)
		}
	})
}

func TestColleagueList_Delete(t *testing.T) {
	t.Run("Delete middle colleague", func(t *testing.T) {
		cl := ColleagueList{
			{Name: "Alice", City: "London", Timezone: "Europe/London"},
			{Name: "Bob", City: "NYC", Timezone: "America/New_York"},
			{Name: "Daisuke", City: "Tokio", Timezone: "Asia/Tokio"},
		}

		if err := cl.Delete(2); err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if len(cl) != 2 {
			t.Errorf("Expected list length %d, got %d instead.", 2, len(cl))
		}

		if cl[0].Name != "Alice" {
			t.Errorf("Expected %q, got %q instead.", "Alice", cl[0].Name)
		}

		if cl[1].Name != "Daisuke" {
			t.Errorf("Expected %q, got %q instead.", "Daisuke", cl[1].Name)
		}
	})

	t.Run("Delete first colleague", func(t *testing.T) {
		cl := ColleagueList{
			{Name: "Alice", City: "London", Timezone: "Europe/London"},
			{Name: "Bob", City: "NYC", Timezone: "America/New_York"},
		}

		if err := cl.Delete(1); err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if len(cl) != 1 {
			t.Errorf("Expected list length %d, got %d instead.", 1, len(cl))
		}

		if cl[0].Name != "Bob" {
			t.Errorf("Expected Bob, got %s instead.", cl[0].Name)
		}
	})

	t.Run("Delete last colleague", func(t *testing.T) {
		cl := ColleagueList{
			{Name: "Alice", City: "London", Timezone: "Europe/London"},
			{Name: "Bob", City: "NYC", Timezone: "America/New_York"},
		}

		if err := cl.Delete(2); err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if len(cl) != 1 {
			t.Errorf("Expected list length %d, got %d instead.", 1, len(cl))
		}

		if cl[0].Name != "Alice" {
			t.Errorf("Expected Alice, got %s instead.", cl[0].Name)
		}
	})

	t.Run("Delete with index 0", func(t *testing.T) {
		cl := ColleagueList{
			{Name: "Alice", City: "London", Timezone: "Europe/London"},
			{Name: "Bob", City: "NYC", Timezone: "America/New_York"},
		}

		err := cl.Delete(0)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		if !errors.Is(err, ErrorInvalidIndex) {
			t.Errorf("expected ErrorInvalidIndex, got %v", err)
		}
	})

	t.Run("Delete with negative index", func(t *testing.T) {
		cl := ColleagueList{
			{Name: "Alice", City: "London", Timezone: "Europe/London"},
			{Name: "Bob", City: "NYC", Timezone: "America/New_York"},
		}

		err := cl.Delete(-1)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		if !errors.Is(err, ErrorInvalidIndex) {
			t.Errorf("expected ErrorInvalidIndex, got %v", err)
		}
	})

	t.Run("Delete with index too large", func(t *testing.T) {
		cl := ColleagueList{
			{Name: "Alice", City: "London", Timezone: "Europe/London"},
			{Name: "Bob", City: "NYC", Timezone: "America/New_York"},
		}

		err := cl.Delete(-1)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		if !errors.Is(err, ErrorInvalidIndex) {
			t.Errorf("expected ErrorInvalidIndex, got %v", err)
		}
	})

	t.Run("Delete from empty list", func(t *testing.T) {
		cl := ColleagueList{}

		err := cl.Delete(1)

		if err == nil {
			t.Fatal("Expected error, got nil")
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
