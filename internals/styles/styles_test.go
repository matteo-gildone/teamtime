package styles

import "testing"

func TestStyles_New(t *testing.T) {
	base := NewStyles()

	if len(base.codes) != 0 {
		t.Errorf("expected empty list, got length %d", len(base.codes))
	}
}

func TestStyles_Render(t *testing.T) {
	tests := []struct {
		name    string
		style   Style
		input   string
		noColor bool
		want    string
	}{
		{
			name:    "base style",
			style:   NewStyles(),
			input:   "base style",
			noColor: false,
			want:    "base style",
		},
		{
			name:    "base style without color",
			input:   "base style no terminal color",
			style:   NewStyles(),
			noColor: true,
			want:    "base style no terminal color",
		},
		{
			name:    "bold style without color",
			input:   "bold style no terminal color",
			style:   NewStyles().Bold(),
			noColor: true,
			want:    "bold style no terminal color",
		},
		{
			name:    "bold style",
			input:   "bold style",
			style:   NewStyles().Bold(),
			noColor: false,
			want:    "\033[1mbold style\033[0m",
		},
		{
			name:    "dim style",
			input:   "dim style",
			style:   NewStyles().Dim(),
			noColor: false,
			want:    "\033[2mdim style\033[0m",
		},
		{
			name:    "underline style",
			input:   "underline style",
			style:   NewStyles().Underline(),
			noColor: false,
			want:    "\033[4munderline style\033[0m",
		},
		{
			name:    "italic style",
			input:   "italic style",
			style:   NewStyles().Italic(),
			noColor: false,
			want:    "\033[3mitalic style\033[0m",
		},
		{
			name:    "red style",
			input:   "red style",
			style:   NewStyles().Red(),
			noColor: false,
			want:    "\033[31mred style\033[0m",
		},
		{
			name:    "green style",
			input:   "green style",
			style:   NewStyles().Green(),
			noColor: false,
			want:    "\033[32mgreen style\033[0m",
		},
		{
			name:    "yellow style",
			input:   "yellow style",
			style:   NewStyles().Yellow(),
			noColor: false,
			want:    "\033[33myellow style\033[0m",
		},
		{
			name:    "cyan style",
			input:   "cyan style",
			style:   NewStyles().Cyan(),
			noColor: false,
			want:    "\033[36mcyan style\033[0m",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalNoColor := noColor
			defer func() { noColor = originalNoColor }()

			noColor = tt.noColor
			got := tt.style.Render(tt.input)

			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}

}

func TestStyleChaining(t *testing.T) {
	tests := []struct {
		name  string
		style Style
		input string
		want  string
	}{
		{
			name:  "base style unchanged",
			style: NewStyles(),
			input: "text",
			want:  "text",
		},
		{
			name:  "red style independent",
			style: NewStyles().Red(),
			input: "text",
			want:  "\033[31mtext\033[0m",
		},
		{
			name:  "multiple styles independent",
			style: NewStyles().Red().Bold(),
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
