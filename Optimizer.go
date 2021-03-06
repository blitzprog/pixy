package pixy

import (
	"regexp"
	"strings"
)

var compactCode = regexp.MustCompile(`\n{2,}`)

const (
	writeStringCall = "_b.WriteString("
)

// optimize combines multiple WriteString calls to one.
func optimize(code string) (optimizedCode string, inlined string) {
	lines := strings.Split(code, "\n")
	lastString := strings.Builder{}

	// Count the actual code lines
	lineCount := 0

	for index, line := range lines {
		// Find WriteString call
		pos := strings.Index(line, writeStringCall)

		if pos != -1 {
			if line[pos+len(writeStringCall)] == '"' {
				// Delete this line and save it in a buffer "lastString"
				lastString.WriteString(line[pos+len(writeStringCall)+1 : len(line)-2])
				lines[index] = ""
				continue
			}
		}

		if lastString.Len() > 0 {
			lines[index] = "\t" + writeStringCall + "\"" + lastString.String() + "\")\n" + line
			lastString.Reset()
		}

		lineCount++
	}

	compact := compactCode.ReplaceAllString(strings.Join(lines, "\n"), "\n")
	return compact, inlined
}
