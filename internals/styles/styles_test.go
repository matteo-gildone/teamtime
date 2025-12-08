package styles

import "testing"

func TestStyles_New(t *testing.T) {
	base := NewStyles()

	if len(base.codes) != 0 {
		t.Errorf("expected empty list, got length %d", len(base.codes))
	}
}

func TestStylesWithColor_Render(t *testing.T) {
	tests := []struct {
		name  string
		style Style
		input string
		want  string
	}{
		{
			name:  "base style",
			style: NewStylesWithNoColor(false),
			input: "base style",
			want:  "base style",
		},
		{
			name:  "bold style",
			input: "bold style",
			style: NewStylesWithNoColor(false).Bold(),
			want:  "\033[1mbold style\033[0m",
		},
		{
			name:  "dim style",
			input: "dim style",
			style: NewStylesWithNoColor(false).Dim(),
			want:  "\033[2mdim style\033[0m",
		},
		{
			name:  "underline style",
			input: "underline style",
			style: NewStylesWithNoColor(false).Underline(),
			want:  "\033[4munderline style\033[0m",
		},
		{
			name:  "italic style",
			input: "italic style",
			style: NewStylesWithNoColor(false).Italic(),
			want:  "\033[3mitalic style\033[0m",
		},
		{
			name:  "red style",
			input: "red style",
			style: NewStylesWithNoColor(false).Red(),
			want:  "\033[31mred style\033[0m",
		},
		{
			name:  "green style",
			input: "green style",
			style: NewStylesWithNoColor(false).Green(),
			want:  "\033[32mgreen style\033[0m",
		},
		{
			name:  "yellow style",
			input: "yellow style",
			style: NewStylesWithNoColor(false).Yellow(),
			want:  "\033[33myellow style\033[0m",
		},
		{
			name:  "cyan style",
			input: "cyan style",
			style: NewStylesWithNoColor(false).Cyan(),
			want:  "\033[36mcyan style\033[0m",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.style.Render(tt.input)

			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}

}

func TestStylesWithNoColor_Render(t *testing.T) {
	tests := []struct {
		name  string
		style Style
		input string
		want  string
	}{
		{
			name:  "base style",
			style: NewStylesWithNoColor(true),
			input: "base style",
			want:  "base style",
		},
		{
			name:  "bold style",
			input: "bold style",
			style: NewStylesWithNoColor(true).Bold(),
			want:  "bold style",
		},
		{
			name:  "dim style",
			input: "dim style",
			style: NewStylesWithNoColor(true).Dim(),
			want:  "dim style",
		},
		{
			name:  "underline style",
			input: "underline style",
			style: NewStylesWithNoColor(true).Underline(),
			want:  "underline style",
		},
		{
			name:  "italic style",
			input: "italic style",
			style: NewStylesWithNoColor(true).Italic(),
			want:  "italic style",
		},
		{
			name:  "red style",
			input: "red style",
			style: NewStylesWithNoColor(true).Red(),
			want:  "red style",
		},
		{
			name:  "green style",
			input: "green style",
			style: NewStylesWithNoColor(true).Green(),
			want:  "green style",
		},
		{
			name:  "yellow style",
			input: "yellow style",
			style: NewStylesWithNoColor(true).Yellow(),
			want:  "yellow style",
		},
		{
			name:  "cyan style",
			input: "cyan style",
			style: NewStylesWithNoColor(true).Cyan(),
			want:  "cyan style",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.style.Render(tt.input)

			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}

}

func TestStyleChaining(t *testing.T) {
	tests := []struct {
		name    string
		style   Style
		input   string
		noColor bool
		want    string
	}{
		{
			name:  "base style unchanged",
			style: NewStylesWithNoColor(false),
			input: "text",
			want:  "text",
		},
		{
			name:  "red style independent",
			style: NewStylesWithNoColor(false).Red(),
			input: "text",
			want:  "\033[31mtext\033[0m",
		},
		{
			name:  "multiple styles independent",
			style: NewStylesWithNoColor(false).Red().Bold(),
			input: "text",
			want:  "\033[31;1mtext\033[0m",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.style.Render(tt.input)
			if got != tt.want {
				t.Errorf("Render() = %q, want %q", got, tt.want)
			}
		})
	}
}
