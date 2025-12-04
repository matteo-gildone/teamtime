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

type Styles struct {
	codes []string
}

func NewStyles() Styles {
	return Styles{
		codes: []string{},
	}
}

func (s Styles) addCode(code string) Styles {
	newCodes := make([]string, len(s.codes), len(s.codes)+1)
	copy(newCodes, s.codes)
	s.codes = append(newCodes, code)
	return s
}

// Text styles

func (s Styles) Bold() Styles {
	return s.addCode("1")
}

func (s Styles) Dim() Styles {
	return s.addCode("2")
}

func (s Styles) Italic() Styles {
	return s.addCode("3")
}

func (s Styles) Underline() Styles {
	return s.addCode("4")
}

// Foreground colors

func (s Styles) Red() Styles {
	return s.addCode("31")
}

func (s Styles) Cyan() Styles {
	return s.addCode("36")
}

func (s Styles) Render(text string) string {
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
