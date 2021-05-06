package ui

import (
	"strings"
)

func Box(text string) string {

	lines := strings.Split(text, "\n")
	contentWidth := longestLine(lines)
	// contentHeight := len(lines)

	var sb strings.Builder
	sb.WriteString(drawFirstLine(contentWidth + 3))

	for _, line := range lines {
		sb.WriteString(drawLine(line, contentWidth))
	}

	sb.WriteString(drawLastLine(contentWidth + 3))

	return sb.String()
}

func drawFirstLine(length int) string {
	var sb strings.Builder

	sb.WriteRune('┌')
	for sb.Len()/3 < length {
		sb.WriteRune('─')
	}
	sb.WriteRune('┐')
	sb.WriteRune('\n')

	return sb.String()
}

func drawLastLine(length int) string {
	var sb strings.Builder

	sb.WriteRune('└')
	for sb.Len()/3 < length {
		sb.WriteRune('─')
	}
	sb.WriteRune('┘')
	// sb.WriteRune('\n')

	return sb.String()
}

func drawLine(text string, length int) string {
	var sb strings.Builder

	sb.WriteRune('│')
	sb.WriteRune(' ')
	sb.WriteString(text)

	for sb.Len() < length-1 {
		sb.WriteRune(' ')
	}

	sb.WriteRune(' ')
	sb.WriteRune('│')
	sb.WriteRune('\n')

	return sb.String()
}

func longestLine(lines []string) int {
	var result = 0

	for _, line := range lines {
		l := len(line)
		if l > result {
			result = l
		}
	}

	return result
}
