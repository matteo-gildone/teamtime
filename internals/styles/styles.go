package styles

import (
	"os"
	"strings"
)

type Style struct {
	codes   []string
	noColor bool
}

// NewStyles creates a new style with auto-detect colour support
// Color is disabled if NO_COLOR env var is set or if TERM "dumb" or empty
func NewStyles() Style {
	return Style{
		codes: []string{},
		noColor: os.Getenv("NO_COLOR") != "" ||
			os.Getenv("TERM") == "dumb" ||
			os.Getenv("TERM") == "",
	}
}

// NewStylesWithNoColor creates a style with explicit color control.
// This is primarily useful for testing
func NewStylesWithNoColor(noColor bool) Style {
	return Style{
		codes:   []string{},
		noColor: noColor,
	}
}

func (s Style) NoColor() bool {
	return s.noColor
}

// addCode return a new Style with an additional ANSI code
func (s Style) addCode(code string) Style {
	codes := append([]string(nil), s.codes...)
	codes = append(codes, code)
	return Style{
		codes:   codes,
		noColor: s.noColor,
	}
}

// Text styles

// Bold returns a style with bold text
func (s Style) Bold() Style {
	return s.addCode("1")
}

// Dim returns a style with dim text
func (s Style) Dim() Style {
	return s.addCode("2")
}

// Italic returns a style with italic text
func (s Style) Italic() Style {
	return s.addCode("3")
}

// Underline returns a style with underline text
func (s Style) Underline() Style {
	return s.addCode("4")
}

// Foreground colors

// Red returns a style with red text
func (s Style) Red() Style {
	return s.addCode("31")
}

// Green returns a style with green text
func (s Style) Green() Style {
	return s.addCode("32")
}

// Yellow returns a style with yellow text
func (s Style) Yellow() Style {
	return s.addCode("33")
}

// Cyan returns a style with cyan text
func (s Style) Cyan() Style {
	return s.addCode("36")
}

// Render applies the style to given text
func (s Style) Render(text string) string {
	if len(s.codes) == 0 || s.noColor {
		return text
	}

	var sb strings.Builder
	sb.WriteString("\033[")
	sb.WriteString(s.codes[0])

	for _, code := range s.codes[1:] {
		sb.WriteString(";")
		sb.WriteString(code)
	}

	sb.WriteString("m")
	sb.WriteString(text)
	sb.WriteString("\033[0m")

	return sb.String()
}
