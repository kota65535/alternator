package cmd

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"log"
	"os"
	"strings"
)

var (
	Bold   = color.New(color.Bold)
	Green  = color.New(color.FgGreen)
	Red    = color.New(color.FgRed)
	Yellow = color.New(color.FgYellow)
	Cyan   = color.New(color.FgCyan)
)

func bPrint(a ...any) {
	Bold.Fprint(os.Stderr, a...)
}

func bPrintf(format string, a ...any) {
	Bold.Fprintf(os.Stderr, format, a...)
}

func bPrintln(a ...any) {
	Bold.Fprintln(os.Stderr, a...)
}

func printlnDiff(str string) {
	split := strings.Split(str, "\n")
	for _, s := range split {
		if strings.HasPrefix(s, "+") {
			Green.Println(s)
		} else if strings.HasPrefix(s, "-") {
			Red.Println(s)
		} else if strings.HasPrefix(s, "~") {
			Yellow.Println(s)
		} else if strings.HasPrefix(s, "@") {
			Cyan.Println(s)
		} else {
			fmt.Println(s)
		}
	}
	fmt.Println()
}

func ePrintf(format string, a ...any) {
	fmt.Fprintf(os.Stderr, format, a...)
}

func ePrintln(a ...any) {
	fmt.Fprintln(os.Stderr, a...)
}

func optS(s1 string, s2 string) string {
	if s1 == "" {
		return ""
	}
	return fmt.Sprintf(s2, s1)
}

func confirm(s string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		bPrintf("%s [y/n]: ", s)

		response, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" {
			return true
		} else if response == "n" || response == "no" {
			return false
		}
	}
}
