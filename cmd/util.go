package cmd

import (
	"fmt"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/fatih/color"
	"net/url"
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

func dPrintf(format string, a ...any) {
	if debug {
		fmt.Fprintf(os.Stderr, format, a...)
	}
}

func dPrintln(a ...any) {
	if debug {
		fmt.Fprintln(os.Stderr, a...)
	}
}

func ePrintf(format string, a ...any) {
	fmt.Fprintf(os.Stderr, format, a...)
}

func ePrintln(a ...any) {
	fmt.Fprintln(os.Stderr, a...)
}

var SupportedDialects = hashset.New("mysql")

type DatabaseUrl struct {
	Dialect  string
	Host     string
	DbName   string
	User     string
	Password string
}

func (r DatabaseUrl) Dsn() string {
	return fmt.Sprintf("%s%s@(%s)/%s", r.User, optS(r.Password, ":%s"), r.Host, r.DbName)
}

func parseDatabaseUrl(u string) DatabaseUrl {
	p, err := url.Parse(u)
	if err != nil {
		bPrintf("Error: failed to parse URL: %s.\n", u)
		if !strings.Contains(u, "://") {
			ePrintln("URL format: (dialect)://(username)[:(password)]@(hostname)[:(port)][/(database)]")
			ePrintln("Examples:")
			ePrintln("  mysql://root@localhost")
			ePrintln("  mysql://root@localhost:13306/db1")
			ePrintln("  mysql://root:mystrongpassword@mydb.dev.example.com/db1")
		}
		os.Exit(1)
	}
	dialect := p.Scheme
	if !SupportedDialects.Contains(dialect) {
		bPrintf("Error: unsupported dialect: %s\n", dialect)
		ePrintf("Supported dialects: %v\n", SupportedDialects)
		os.Exit(1)
	}

	host := p.Host
	dbName := ""
	if len(p.Path) > 0 {
		dbName = p.Path[1:]
	}
	user := p.User.Username()
	password, _ := p.User.Password()

	return DatabaseUrl{dialect, host, dbName, user, password}
}

func optS(s1 string, s2 string) string {
	if s1 == "" {
		return ""
	}
	return fmt.Sprintf(s2, s1)
}
