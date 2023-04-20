package printer

import (
	"github.com/logrusorgru/aurora/v4"
)

func FormatGreen(str string) string {
	return aurora.Green(str).String()
}

func FormatYellow(str string) string {
	return aurora.Yellow(str).String()
}

func FormatBold(str string) string {
	return aurora.Bold(str).String()
}

func FormatUnderline(str string) string {
	return aurora.Underline(str).String()
}
