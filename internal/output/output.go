package output

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

var (
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Red    = "\033[31m"
	Cyan   = "\033[36m"
	Bold   = "\033[1m"
	Dim    = "\033[2m"
	Reset  = "\033[0m"
)

func init() {
	if os.Getenv("NO_COLOR") != "" || !isTerminal() {
		Green, Yellow, Red, Cyan, Bold, Dim, Reset = "", "", "", "", "", "", ""
	}
}

func isTerminal() bool {
	fi, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}

func Success(msg string) { fmt.Printf("%s✓%s %s\n", Green, Reset, msg) }
func Error(msg string)   { fmt.Fprintf(os.Stderr, "%s✗%s %s\n", Red, Reset, msg) }
func Info(msg string)    { fmt.Printf("%s→%s %s\n", Cyan, Reset, msg) }

func Header(title string) {
	fmt.Printf("\n%s%s%s\n%s\n", Bold, title, Reset, strings.Repeat("─", len(title)))
}

func Prompt(label string) (string, error) {
	fmt.Printf("%s%s%s: ", Bold, label, Reset)
	var s string
	fmt.Scanln(&s)
	return s, nil
}

func JSON(v json.RawMessage) {
	var obj interface{}
	if err := json.Unmarshal(v, &obj); err != nil {
		fmt.Println(string(v))
		return
	}
	data, _ := json.MarshalIndent(obj, "", "  ")
	fmt.Println(colorJSON(string(data)))
}

func colorJSON(s string) string {
	if Bold == "" {
		return s
	}
	lines := strings.Split(s, "\n")
	var b strings.Builder
	for _, line := range lines {
		trimmed := strings.TrimLeft(line, " ")
		indent := line[:len(line)-len(trimmed)]
		if strings.HasPrefix(trimmed, `"`) && strings.Contains(trimmed, `":`) {
			i := strings.Index(trimmed, `":`)
			key := trimmed[:i+1]
			rest := trimmed[i+1:]
			b.WriteString(indent + Cyan + key + Reset + ":" + colorValue(rest) + "\n")
		} else {
			b.WriteString(line + "\n")
		}
	}
	return strings.TrimRight(b.String(), "\n")
}

func colorValue(s string) string {
	t := strings.TrimSpace(s)
	comma := ""
	if strings.HasSuffix(t, ",") {
		comma = ","
		t = t[:len(t)-1]
	}
	var c string
	switch {
	case t == "true" || t == "false":
		c = Yellow + t + Reset
	case t == "null":
		c = Dim + t + Reset
	case strings.HasPrefix(t, `"`):
		c = Green + t + Reset
	default:
		c = t
	}
	return " " + c + comma
}

func Table(headers []string, rows [][]string) {
	if len(rows) == 0 {
		fmt.Println(Dim + "(no results)" + Reset)
		return
	}
	widths := make([]int, len(headers))
	for i, h := range headers {
		widths[i] = len(h)
	}
	for _, row := range rows {
		for i, cell := range row {
			if i < len(widths) && len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}
	}
	fmt.Print(Bold)
	for i, h := range headers {
		fmt.Printf("%-*s", widths[i]+2, h)
	}
	fmt.Println(Reset)
	for _, w := range widths {
		fmt.Print(strings.Repeat("─", w+2))
	}
	fmt.Println()
	for _, row := range rows {
		for i, cell := range row {
			if i < len(widths) {
				fmt.Printf("%-*s", widths[i]+2, cell)
			}
		}
		fmt.Println()
	}
}
