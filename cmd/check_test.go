package cmd

import (
	"strings"
	"testing"
	"time"

	"github.com/matteo-gildone/teamtime/internals/styles"
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
			want: timeOff,
		},
		{
			name: "off too early morning",
			hour: 23,
			want: timeOff,
		},
		{
			name: "work time",
			hour: 9,
			want: timeWork,
		},
		{
			name: "extended work time morning",
			hour: 7,
			want: timeExtended,
		},
		{
			name: "extended work time afternoon",
			hour: 17,
			want: timeExtended,
		},
	}

	for _, tt := range tests {
		got := classifyTimeOfDay(tt.hour)
		if got != tt.want {
			t.Errorf("got %q, want %q, given %v", got, tt.want, tt.hour)
		}
	}
}

func TestGetDisplayTime(t *testing.T) {
	tests := []struct {
		name           string
		hour           int
		noColor        bool
		wantContains   []string
		wantNotContain []string
	}{
		{
			name:           "work hours with color",
			hour:           10,
			noColor:        false,
			wantContains:   []string{"\033[1;36m", "10:00"},
			wantNotContain: []string{"[Off]", "[Extended]"},
		},
		{
			name:           "work hours without color",
			hour:           10,
			noColor:        true,
			wantContains:   []string{"10:00"},
			wantNotContain: []string{"\033[", "[Off]", "[Extended]"},
		},
		{
			name:           "extended morning with color",
			hour:           7,
			noColor:        false,
			wantContains:   []string{"\033[1;33m", "07:00", "[Extended]"},
			wantNotContain: []string{"[Off]"},
		},
		{
			name:           "extended morning without color",
			hour:           7,
			noColor:        true,
			wantContains:   []string{"07:00", "[Extended]"},
			wantNotContain: []string{"\033[", "[Off]"},
		},
		{
			name:           "extended evening with color",
			hour:           18,
			noColor:        false,
			wantContains:   []string{"\033[1;33m", "18:00", "[Extended]"},
			wantNotContain: []string{"[Off]"},
		},
		{
			name:           "extended evening without color",
			hour:           18,
			noColor:        true,
			wantContains:   []string{"18:00", "[Extended]"},
			wantNotContain: []string{"\033[", "[Off]"},
		},
		{
			name:           "off morning with color",
			hour:           23,
			noColor:        false,
			wantContains:   []string{"\033[1;31m", "23:00", "[Off]"},
			wantNotContain: []string{"[Extended]"},
		},
		{
			name:           "off morning without color",
			hour:           23,
			noColor:        true,
			wantContains:   []string{"23:00", "[Off]"},
			wantNotContain: []string{"\033[", "[Extended]"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testTime := time.Date(2025, 12, 6, tt.hour, 0, 0, 0, time.UTC)
			style := styles.NewStylesWithNoColor(tt.noColor)
			result := getDisplayTime(testTime, style)

			for _, want := range tt.wantContains {
				if !strings.Contains(result, want) {
					t.Errorf("want results to contain: %q, got: %q", want, result)
				}
			}

			for _, notWant := range tt.wantNotContain {
				if strings.Contains(result, notWant) {
					t.Errorf("want results not to contain: %q, got: %q", notWant, result)
				}
			}
		})
	}
}
