package cmd

import (
	"strings"
	"testing"
	"time"
)

func TestClassifyTimeOfDay(t *testing.T) {
	tests := []struct {
		name string
		hour int
		want timeClassification
	}{
		{
			name: "off too late evening",
			hour: 3,
			want: "off",
		},
		{
			name: "off too early morning",
			hour: 23,
			want: "off",
		},
		{
			name: "work time",
			hour: 9,
			want: "work",
		},
		{
			name: "extended work time morning",
			hour: 7,
			want: "extended",
		},
		{
			name: "extended work time afternoon",
			hour: 17,
			want: "extended",
		},
	}

	for _, tt := range tests {
		got := classifyTimeOfDay(tt.hour)
		if got != tt.want {
			t.Errorf("got %q, want %q, given %v", got, tt.want, tt.hour)
		}
	}
}

func TestGetDisplayTime_WorkHours(t *testing.T) {
	tests := []struct {
		name string
		hour int
	}{
		{
			name: "morning with color",
			hour: 10,
		},
		{
			name: "afternoon no color",
			hour: 14,
		},
		{
			name: "late afternoon",
			hour: 16,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testTime := time.Date(2025, 12, 6, tt.hour, 0, 0, 0, time.UTC)
			result := getDisplayTime(testTime)

			if strings.Contains(result, "[Off]") || strings.Contains(result, "[Extended]") {
				t.Errorf("work hours should not have and indicator, got: %q", result)
			}

			if !strings.Contains(result, "36") {
				t.Errorf("work hours should contain ANSI, got: %q", result)
			}
		})
	}
}

func TestGetDisplayTime_ExtendedHours(t *testing.T) {
	tests := []struct {
		name string
		hour int
	}{
		{
			name: "morning extended",
			hour: 7,
		},
		{
			name: "morning extended",
			hour: 8,
		},
		{
			name: "evening extended",
			hour: 18,
		},
		{
			name: "evening extended",
			hour: 19,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testTime := time.Date(2025, 12, 6, tt.hour, 0, 0, 0, time.UTC)
			result := getDisplayTime(testTime)

			if !strings.Contains(result, "[Extended]") {
				t.Errorf("extended hours should have [Extended], got: %q", result)
			}

			if !strings.Contains(result, "33") {
				t.Errorf("extended hours should contain ANSI, got: %q", result)
			}
		})
	}
}

func TestGetDisplayTime_OffHours(t *testing.T) {
	tests := []struct {
		name string
		hour int
	}{
		{
			name: "morning off",
			hour: 6,
		},
		{
			name: "evening off",
			hour: 20,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testTime := time.Date(2025, 12, 6, tt.hour, 0, 0, 0, time.UTC)
			result := getDisplayTime(testTime)

			if !strings.Contains(result, "[Off]") {
				t.Errorf("off hours should have [Off], got: %q", result)
			}

			if !strings.Contains(result, "31") {
				t.Errorf("off hours should contain ANSI, got: %q", result)
			}
		})
	}
}
