package styles

import (
	"os"
	"strings"
)

var noColor bool

func init() {
	noColor = os.Getenv("NO_COLOR") != "" ||
		os.Getenv("TERM") == "dumb" ||
		os.Getenv("TERM") == ""
}

type Style struct {
	codes []string
}

func NewStyles() Style {
	return Style{
		codes: []string{},
	}
}

func (s Style) addCode(code string) Style {
	codes := append([]string(nil), s.codes...)
	codes = append(codes, code)
	return Style{codes: codes}
}

// Text styles

func (s Style) Bold() Style {
	return s.addCode("1")
}

func (s Style) Dim() Style {
	return s.addCode("2")
}

func (s Style) Italic() Style {
	return s.addCode("3")
}

func (s Style) Underline() Style {
	return s.addCode("4")
}

// Foreground colors

func (s Style) Red() Style {
	return s.addCode("31")
}

func (s Style) Green() Style {
	return s.addCode("32")
}

func (s Style) Yellow() Style {
	return s.addCode("33")
}

func (s Style) Cyan() Style {
	return s.addCode("36")
}

func (s Style) Render(text string) string {
	if len(s.codes) == 0 || noColor {
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
