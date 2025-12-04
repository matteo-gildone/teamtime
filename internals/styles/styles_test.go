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
		style   Styles
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
			name:    "dim style without color",
			input:   "dim style no terminal color",
			style:   NewStyles().Dim(),
			noColor: true,
			want:    "dim style no terminal color",
		},
		{
			name:    "dim style",
			input:   "dim style",
			style:   NewStyles().Dim(),
			noColor: false,
			want:    "\033[2mdim style\033[0m",
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
